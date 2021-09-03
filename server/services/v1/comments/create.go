package comments

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/lbryio/commentron/commentapi"
	"github.com/lbryio/commentron/config"
	"github.com/lbryio/commentron/db"
	"github.com/lbryio/commentron/flags"
	"github.com/lbryio/commentron/helper"
	m "github.com/lbryio/commentron/model"
	"github.com/lbryio/commentron/server/lbry"
	"github.com/lbryio/commentron/server/websocket"
	"github.com/lbryio/commentron/sockety"

	"github.com/lbryio/lbry.go/v2/extras/api"
	"github.com/lbryio/lbry.go/v2/extras/errors"
	"github.com/lbryio/lbry.go/v2/extras/jsonrpc"
	"github.com/lbryio/lbry.go/v2/extras/util"
	v "github.com/lbryio/ozzo-validation"
	"github.com/lbryio/sockety/socketyapi"

	"github.com/Avalanche-io/counter"
	"github.com/btcsuite/btcutil"
	"github.com/hbakhtiyor/strsim"
	"github.com/karlseguin/ccache"
	"github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/paymentintent"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

func create(_ *http.Request, args *commentapi.CreateArgs, reply *commentapi.CreateResponse) error {
	err := v.ValidateStruct(args,
		v.Field(&args.ClaimID, v.Required))
	if err != nil {
		return api.StatusError{Err: errors.Err(err), Status: http.StatusBadRequest}
	}
	channel, err := m.Channels(m.ChannelWhere.ClaimID.EQ(null.StringFrom(args.ChannelID).String)).One(db.RO)
	if errors.Is(err, sql.ErrNoRows) {
		channel = &m.Channel{
			ClaimID: null.StringFrom(args.ChannelID).String,
			Name:    null.StringFrom(args.ChannelName).String,
		}
		err = nil
		err := channel.Insert(db.RW, boil.Infer())
		if err != nil {
			return errors.Err(err)
		}
	}
	if err != nil {
		return errors.Err(err)
	}
	request := &createRequest{args: args}
	err = checkAllowedAndValidate(args)
	if err != nil {
		return err
	}

	err = createComment(request)
	if err != nil {
		return err
	}

	if args.SupportTxID != nil || args.PaymentIntentID != nil {
		err := updateSupportInfo(request)
		if err != nil {
			return err
		}
	}

	err = blockedByCreator(request)
	if err != nil {
		return err
	}

	err = flags.CheckComment(request.comment)
	if err != nil {
		return err
	}

	err = request.comment.Insert(db.RW, boil.Infer())
	if err != nil {
		return err
	}

	item := populateItem(request.comment, channel, 0)

	err = applyModStatus(&item, args.ChannelID, args.ClaimID)
	if err != nil {
		return err
	}

	reply.CommentItem = &item
	if !request.comment.IsFlagged {
		go pushItem(item, args.ClaimID)
		amount, err := btcutil.NewAmount(item.SupportAmount)
		if err != nil {
			return errors.Err(err)
		}
		go lbry.API.Notify(lbry.NotifyOptions{
			ActionType: "C",
			CommentID:  item.CommentID,
			ChannelID:  &item.ChannelID,
			ParentID:   &item.ParentID,
			Comment:    &item.Comment,
			ClaimID:    item.ClaimID,
			Amount:     uint64(amount),
			IsFiat:     item.IsFiat,
			Currency:   util.PtrToString(item.Currency),
		})
	}

	return nil
}

func createComment(request *createRequest) error {
	commentID, timestamp, err := createCommentID(request.args.CommentText, null.StringFrom(request.args.ChannelID).String)
	if err != nil {
		return errors.Err(err)
	}

	request.comment = &m.Comment{
		CommentID:   commentID,
		LbryClaimID: request.args.ClaimID,
		ChannelID:   null.StringFrom(request.args.ChannelID),
		Body:        request.args.CommentText,
		ParentID:    null.StringFromPtr(request.args.ParentID),
		Signature:   null.StringFrom(request.args.Signature),
		Signingts:   null.StringFrom(request.args.SigningTS),
		Timestamp:   int(timestamp),
	}
	return nil
}

func checkAllowedAndValidate(args *commentapi.CreateArgs) error {
	blockedEntry, err := m.BlockedEntries(m.BlockedEntryWhere.UniversallyBlocked.EQ(null.BoolFrom(true)), m.BlockedEntryWhere.BlockedChannelID.EQ(null.StringFrom(args.ChannelID))).One(db.RO)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return errors.Err(err)
	}

	if blockedEntry != nil {
		return api.StatusError{Err: errors.Err("channel is not allowed to post comments"), Status: http.StatusBadRequest}
	}

	if args.ParentID != nil {
		err = helper.AllowedToRespond(util.StrFromPtr(args.ParentID), args.ChannelID)
		if err != nil {
			return err
		}
	}

	err = lbry.ValidateSignature(args.ChannelID, args.Signature, args.SigningTS, args.CommentText)
	if err != nil {
		return errors.Prefix("could not authenticate channel signature:", err)
	}

	return nil
}

func applyModStatus(item *commentapi.CommentItem, channelID, claimID string) error {
	isGlobalMod, err := m.Moderators(m.ModeratorWhere.ModChannelID.EQ(null.StringFrom(channelID))).Exists(db.RO)
	if err != nil {
		return errors.Err(err)
	}
	item.IsGlobalMod = isGlobalMod

	signingChannel, err := lbry.SDK.GetSigningChannelForClaim(claimID)
	if err != nil {
		return errors.Err(err)
	}
	if signingChannel != nil {
		item.IsCreator = channelID == signingChannel.ClaimID
		filterCreator := m.DelegatedModeratorWhere.CreatorChannelID.EQ(signingChannel.ClaimID)
		filterCommenter := m.DelegatedModeratorWhere.ModChannelID.EQ(channelID)
		isMod, err := m.DelegatedModerators(filterCreator, filterCommenter).Exists(db.RO)
		if err != nil {
			return errors.Err(err)
		}
		item.IsModerator = isMod
	}
	return nil
}

func pushItem(item commentapi.CommentItem, claimID string) {
	websocket.PushTo(&websocket.PushNotification{
		Type: "delta",
		Data: map[string]interface{}{"comment": item},
	}, claimID)

	go sockety.SendNotification(socketyapi.SendNotificationArgs{
		Service: socketyapi.Commentron,
		Type:    "delta",
		IDs:     []string{claimID, "comments"},
		Data:    map[string]interface{}{"comment": item},
	})

}

func checkForDuplicate(commentID string) error {
	comment, err := m.Comments(m.CommentWhere.CommentID.EQ(commentID)).One(db.RO)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return errors.Err(err)
	}

	if comment != nil {
		return api.StatusError{Err: errors.Err("duplicate comment!"), Status: http.StatusBadRequest}
	}
	return nil
}

var slowModeCache = ccache.New(ccache.Configure().MaxSize(10000))

type createRequest struct {
	args           *commentapi.CreateArgs
	comment        *m.Comment
	creatorChannel *m.Channel
	signingChannel *jsonrpc.Claim
	currency       string
	isFiat         bool
}

const maxSimilaryScoreToCreatorName = 0.6

func blockedByCreator(request *createRequest) error {
	var err error
	request.signingChannel, err = lbry.SDK.GetSigningChannelForClaim(request.args.ClaimID)
	if err != nil {
		return errors.Err(err)
	}
	if request.signingChannel == nil {
		return nil
	}
	request.creatorChannel, err = helper.FindOrCreateChannel(request.signingChannel.ClaimID, request.signingChannel.Name)
	if err != nil {
		return errors.Err(err)
	}
	//Make sure commenter is not commenting from a channel that is "like" the creator.
	similarity := strsim.Compare(request.creatorChannel.Name, request.args.ChannelName)
	if request.args.ChannelID != request.signingChannel.ClaimID && similarity > maxSimilaryScoreToCreatorName {
		return errors.Err("your user name %s is too close to the creator's user name %s and may cause confusion. Please use another identity.", request.args.ChannelName, request.creatorChannel.Name)
	}

	creatorFilter := m.BlockedEntryWhere.CreatorChannelID.EQ(null.StringFrom(request.signingChannel.ClaimID))
	userFilter := m.BlockedEntryWhere.BlockedChannelID.EQ(null.StringFrom(request.args.ChannelID))
	blockedListFilter := m.BlockedEntryWhere.BlockedListID.EQ(request.creatorChannel.BlockedListID)
	blockedEntry, err := m.BlockedEntries(creatorFilter, userFilter).One(db.RO)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return errors.Err(err)
	}

	if blockedEntry != nil && !blockedEntry.Expiry.Valid {
		return api.StatusError{Err: errors.Err("channel is blocked by publisher"), Status: http.StatusBadRequest}
	} else if blockedEntry != nil && blockedEntry.Expiry.Valid && time.Since(blockedEntry.Expiry.Time) < time.Duration(0) {
		timeLeft := helper.FormatDur(blockedEntry.Expiry.Time.Sub(time.Now()))
		message := fmt.Sprintf("publisher %s has given you a temporary ban with %s remaining.", request.creatorChannel.Name, timeLeft)
		return api.StatusError{Err: errors.Err(message), Status: http.StatusBadRequest}
	}

	blockedListEntry, err := m.BlockedEntries(blockedListFilter, userFilter, qm.Load(m.BlockedEntryRels.BlockedList), qm.Load(m.BlockedEntryRels.CreatorChannel)).One(db.RO)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return errors.Err(err)
	}
	if blockedListEntry != nil && blockedListEntry.R != nil && blockedListEntry.R.BlockedList != nil {
		blockedByChannel := "UNKNOWN"
		blockedListName := blockedListEntry.R.BlockedList.Name
		if blockedListEntry.R.CreatorChannel != nil {
			blockedByChannel = blockedListEntry.R.CreatorChannel.Name
		}
		if blockedListEntry.Expiry.Valid && time.Since(blockedListEntry.Expiry.Time) < time.Duration(0) {
			timeLeft := time.Since(blockedListEntry.Expiry.Time)
			message := fmt.Sprintf("channel %s added you to the shared block list %s and you will not be able to comment until %g hrs has passed.", blockedByChannel, blockedListName, timeLeft.Hours())
			return api.StatusError{Err: errors.Err(message), Status: http.StatusBadRequest}
		}
	}

	settings, err := request.creatorChannel.CreatorChannelCreatorSettings().One(db.RO)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return errors.Err(err)
	}
	if settings != nil {
		return checkSettings(settings, request)
	}
	return nil
}

func checkSettings(settings *m.CreatorSetting, request *createRequest) error {
	isMod, err := m.DelegatedModerators(m.DelegatedModeratorWhere.ModChannelID.EQ(request.args.ChannelID), m.DelegatedModeratorWhere.CreatorChannelID.EQ(request.signingChannel.ClaimID)).Exists(db.RO)
	if err != nil {
		return errors.Err(err)
	}
	if !isMod && request.args.ChannelID != request.creatorChannel.ClaimID {
		if !settings.MinTipAmountSuperChat.IsZero() && !request.comment.Amount.IsZero() && request.args.PaymentIntentID == nil {
			if request.comment.Amount.Uint64 < settings.MinTipAmountSuperChat.Uint64 {
				return api.StatusError{Err: errors.Err("a min tip of %d LBC is required to hyperchat", settings.MinTipAmountSuperChat.Uint64), Status: http.StatusBadRequest}
			}
		}
		if !settings.MinTipAmountComment.IsZero() {
			if request.comment.Amount.IsZero() {
				return api.StatusError{Err: errors.Err("you must include tip in order to comment as required by creator"), Status: http.StatusBadRequest}
			}
			if request.comment.Amount.Uint64 < settings.MinTipAmountComment.Uint64 {
				return api.StatusError{Err: errors.Err("you must tip at least %d with this comment as required by %s", settings.MinTipAmountComment.Uint64, request.creatorChannel.Name), Status: http.StatusBadRequest}
			}
		}
		if !settings.SlowModeMinGap.IsZero() {
			err := checkMinGap(request.args.ChannelID+request.creatorChannel.ClaimID, time.Duration(settings.SlowModeMinGap.Uint64)*time.Second)
			if err != nil {
				return err
			}
		}
		if !settings.MutedWords.IsZero() {
			blockedWords := strings.Split(settings.MutedWords.String, ",")
			for _, blockedWord := range blockedWords {
				if strings.Contains(request.args.CommentText, blockedWord) {
					return api.StatusError{Err: errors.Err("the comment contents are blocked by %s", request.signingChannel.Name)}
				}
			}
		}
	}
	if !settings.CommentsEnabled.Valid {
		for _, tag := range request.signingChannel.Value.Tags {
			if tag == "comments-disabled" {
				settings.CommentsEnabled.SetValid(false)
				err := settings.Update(db.RW, boil.Whitelist(m.CreatorSettingColumns.CommentsEnabled))
				if err != nil {
					return errors.Err(err)
				}
			}
		}
	}
	if !settings.CommentsEnabled.Bool {
		return api.StatusError{Err: errors.Err("comments are disabled by the creator"), Status: http.StatusBadRequest}
	}
	return nil
}

func checkMinGap(key string, expiration time.Duration) error {
	creatorCounter, err := getCounter(key, expiration)
	if err != nil {
		return err
	}
	if creatorCounter.Get() > 0 {
		minGapViolated := fmt.Sprintf("Slow mode is on. Please wait at most %d seconds before commenting again.", int(expiration.Seconds()))
		return api.StatusError{Err: errors.Err(minGapViolated), Status: http.StatusBadRequest}
	}
	creatorCounter.Add(1)

	return nil
}

func getCounter(key string, expiration time.Duration) (*counter.Counter, error) {
	result, err := slowModeCache.Fetch(key, expiration, func() (interface{}, error) {
		return counter.New(), nil
	})
	if err != nil {
		return nil, errors.Err(err)
	}
	creatorCounter, ok := result.Value().(*counter.Counter)
	if !ok {
		return nil, errors.Err("could not convert counter from cache!")
	}
	return creatorCounter, nil
}

func updateSupportInfo(request *createRequest) error {
	triesLeft := 3
	for {
		triesLeft--
		err := updateSupportInfoAttempt(request)
		if err == nil {
			return nil
		}
		if triesLeft == 0 {
			return err
		}
		time.Sleep(1 * time.Second)
	}
}

func updateSupportInfoAttempt(request *createRequest) error {
	if request.args.PaymentIntentID != nil {
		env := ""
		if request.args.Environment != nil {
			env = *request.args.Environment
		}
		paymentintentClient := &paymentintent.Client{B: stripe.GetBackend(stripe.APIBackend), Key: config.ConnectAPIKey(config.From(env))}
		pi, err := paymentintentClient.Get(*request.args.PaymentIntentID, &stripe.PaymentIntentParams{})
		if err != nil {
			logrus.Error(errors.Prefix("could not get payment intent %s", *request.args.PaymentIntentID))
			return errors.Err("could not validate tip")
		}
		request.comment.Amount.SetValid(uint64(pi.Amount))
		request.comment.IsFiat = true
		request.comment.Currency.SetValid(pi.Currency)
		return nil

	}
	request.comment.TXID.SetValid(util.StrFromPtr(request.args.SupportTxID))
	txSummary, err := lbry.SDK.GetTx(request.comment.TXID.String)
	if err != nil {
		return errors.Err(err)
	}
	if txSummary == nil {
		return errors.Err("transaction not found for txid %s", request.comment.TXID.String)
	}
	var vout uint64
	if request.args.SupportVout != nil {
		vout = *request.args.SupportVout
	}
	amount, err := getVoutAmount(request.args.ChannelID, txSummary, vout)
	if err != nil {
		return errors.Err(err)
	}
	request.comment.Amount.SetValid(amount)
	return nil
}

func getVoutAmount(channelID string, summary *jsonrpc.TransactionSummary, vout uint64) (uint64, error) {
	if summary == nil {
		return 0, errors.Err("transaction summary missing")
	}

	if len(summary.Outputs) <= int(vout) {
		return 0, errors.Err("there are not enough outputs on the transaction for position %d", vout)
	}
	output := summary.Outputs[int(vout)]

	if output.SigningChannel == nil {
		return 0, errors.Err("Expected signed support for %s in transaction %s", channelID, summary.Txid)
	}

	if output.SigningChannel.ChannelID != channelID && !config.IsTestMode {
		return 0, errors.Err("The support was not signed by %s, but was instead signed by channel %s", channelID, output.SigningChannel.ChannelID)
	}
	amountStr := output.Amount
	amountFloat, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return 0, errors.Err(err)
	}
	amount, err := btcutil.NewAmount(amountFloat)
	if err != nil {
		return 0, errors.Err(err)
	}
	return uint64(amount), nil
}
