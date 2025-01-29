package comments

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/OdyseeTeam/commentron/commentapi"
	"github.com/OdyseeTeam/commentron/config"
	"github.com/OdyseeTeam/commentron/db"
	"github.com/OdyseeTeam/commentron/flags"
	"github.com/OdyseeTeam/commentron/helper"
	m "github.com/OdyseeTeam/commentron/model"
	"github.com/OdyseeTeam/commentron/server/lbry"
	"github.com/OdyseeTeam/commentron/server/websocket"
	"github.com/OdyseeTeam/commentron/sockety"

	"github.com/OdyseeTeam/sockety/socketyapi"
	"github.com/lbryio/lbry.go/v2/extras/api"
	"github.com/lbryio/lbry.go/v2/extras/errors"
	"github.com/lbryio/lbry.go/v2/extras/jsonrpc"
	"github.com/lbryio/lbry.go/v2/extras/util"
	v "github.com/lbryio/ozzo-validation"

	"github.com/Avalanche-io/counter"
	"github.com/btcsuite/btcutil"
	"github.com/hbakhtiyor/strsim"
	"github.com/karlseguin/ccache/v2"
	"github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/paymentintent"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/machinebox/graphql"
)

var specialLogFile *os.File

// Temp variable to allow testing
var useOldTipAmountChecks bool

func init() {
	var err error
	specialLogFile, err = os.OpenFile("special.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logrus.Fatal(err)
	}
}
func create(_ *http.Request, args *commentapi.CreateArgs, reply *commentapi.CreateResponse) error {
	err := v.ValidateStruct(args,
		v.Field(&args.ClaimID, v.Required))
	if err != nil {
		return api.StatusError{Err: errors.Err(err), Status: http.StatusBadRequest}
	}
	//log what this special commenter is doing to find the bug
	if args.ChannelID == "ccf4e035d8164d8a6540d96d1a689a4f068b6bc7" {
		stuffToLog, err := json.Marshal(args)
		stuffToLog = append(stuffToLog, '\n')
		if err == nil {
			_, _ = specialLogFile.Write(stuffToLog)
		}

	}
	channel, err := helper.FindOrCreateChannel(args.ChannelID, args.ChannelName)
	if err != nil {
		return errors.Err(err)
	}
	request := &createRequest{args: args}
	err = checkAllowedAndValidate(request)
	if err != nil {
		return err
	}

	err = createComment(request)
	if err != nil {
		return err
	}

	// Temp to allow testing
	useOldTipAmountChecks = args.Amount == nil

	var frequencyCheck = checkFrequency
	if args.SupportTxID != nil || args.PaymentIntentID != nil || args.Currency != nil {
		if args.DryRun {
			if args.Amount != nil {
				if args.PaymentIntentID != nil {
					cents := uint64(*args.Amount * 100)
					request.comment.Amount.SetValid(cents)
				} else if args.SupportTxID != nil {
					lbc, err := btcutil.NewAmount(*args.Amount)
					if err != nil {
						return errors.Err(err)
					}
					request.comment.Amount.SetValid(uint64(lbc.ToUnit(btcutil.AmountSatoshi)))
				} else if *args.Currency == "USDC" {
					cents := uint64(*args.Amount * 100)
					request.comment.Amount.SetValid(cents)
					request.comment.Currency.SetValid(*args.Currency)
				}
			}
		} else {
			err := updateSupportInfo(request)
			if err != nil {
				return err
			}
		}
		// ignore the frequency if its a tipped comment
		frequencyCheck = ignoreFrequency
	}

	// This is strategically placed, nothing can be done before this using the comment id or timestamp
	commentID, timestamp, err := createCommentID(request.args.CommentText, null.StringFrom(request.args.ChannelID).String, frequencyCheck)
	if err != nil {
		return errors.Err(err)
	}
	request.comment.CommentID = commentID
	request.comment.Timestamp = int(timestamp)
	request.comment.IsProtected = args.IsProtected

	item := populateItem(request.comment, channel, 0)

	err = applyModStatus(&item, args.ChannelID, args.ClaimID)
	if err != nil {
		return err
	}
	if !item.IsModerator {
		err = blockedByCreator(request)
		if err != nil {
			return err
		}
	}

	if !(args.Sticker && (args.SupportTxID != nil || args.PaymentIntentID != nil)) {
		flags.CheckComment(request.comment)
	}

	err = EnsureClaimToChannelExists(request.comment.LbryClaimID)
	if err != nil {
		return err
	}

	if args.DryRun {
		reply.CommentItem = &item
		return nil
	}

	err = request.comment.Insert(db.RW, boil.Infer())
	if err != nil {
		return err
	}

	reply.CommentItem = &item
	if !request.comment.IsFlagged {
		pushClaimID := args.ClaimID
		if item.IsProtected {
			pushClaimID = helper.ReverseString(args.ClaimID)
		}
		go pushItem(item, pushClaimID, args.MentionedChannels)
		amount, err := btcutil.NewAmount(item.SupportAmount)
		if err != nil {
			return errors.Err(err)
		}
		go lbry.API.Notify(lbry.NotifyOptions{
			ActionType:  "C",
			CommentID:   item.CommentID,
			ChannelID:   &item.ChannelID,
			ParentID:    &item.ParentID,
			Comment:     &item.Comment,
			ClaimID:     item.ClaimID,
			Amount:      uint64(amount),
			IsFiat:      item.IsFiat,
			Currency:    util.PtrToString(item.Currency),
			IsProtected: item.IsProtected,
		})
	}

	return nil
}

func createComment(request *createRequest) error {

	request.comment = &m.Comment{
		LbryClaimID: request.args.ClaimID,
		ChannelID:   null.StringFrom(request.args.ChannelID),
		Body:        request.args.CommentText,
		ParentID:    null.StringFromPtr(request.args.ParentID),
		Signature:   null.StringFrom(request.args.Signature),
		Signingts:   null.StringFrom(request.args.SigningTS),
	}
	return nil
}

func checkAllowedAndValidate(request *createRequest) error {
	blockedEntry, err := m.BlockedEntries(m.BlockedEntryWhere.UniversallyBlocked.EQ(null.BoolFrom(true)), m.BlockedEntryWhere.BlockedChannelID.EQ(null.StringFrom(request.args.ChannelID))).One(db.RO)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return errors.Err(err)
	}

	if blockedEntry != nil {
		return api.StatusError{Err: errors.Err("channel is not allowed to post comments"), Status: http.StatusBadRequest}
	}

	if request.args.ParentID != nil {
		contentCreatorChannel, err := lbry.SDK.GetSigningChannelForClaim(request.args.ClaimID)
		if err != nil {
			return errors.Err(err)
		}
		err = helper.AllowedToRespond(util.StrFromPtr(request.args.ParentID), request.args.ChannelID, contentCreatorChannel)
		if err != nil {
			return err
		}
	}

	err = lbry.ValidateSignatureAndTS(request.args.ChannelID, request.args.Signature, request.args.SigningTS, request.args.CommentText)
	if err != nil {
		return errors.Prefix("could not authenticate channel signature:", err)
	}
	matches := commentapi.StickerRE.FindStringSubmatch(request.args.CommentText)
	if len(matches) > 0 && !request.args.Sticker {
		return errors.Err("a sticker cannot be passed with the sticker flag true")
	}
	if request.args.Sticker {
		if len(matches) != 2 {
			return errors.Err("invalid sticker code")
		}
		paid, ok := allowedStickers[matches[1]]
		if !ok {
			return errors.Err("%s is not an authorized Odysee sticker", matches[1])
		}
		if paid && request.args.PaymentIntentID == nil && request.args.SupportTxID == nil {
			return errors.Err("%s requires a support to post", matches[1])
		}
	}

	isProtected, err := IsProtectedContent(request.args.ClaimID)
	if err != nil {
		return err
	}
	isLivestream, err := IsLivestreamClaim(request.args.ClaimID)
	if err != nil {
		return err
	}
	request.args.IsProtected = isProtected
	request.isLivestream = isLivestream

	if isProtected {
		hasAccess, err := HasAccessToProtectedContent(request.args.ClaimID, request.args.ChannelID)
		if err != nil {
			return err
		}
		if !hasAccess {
			return api.StatusError{Err: errors.Err("channel does not have permissions to comment on this claim"), Status: http.StatusForbidden}
		}
	}

	return nil
}

// IsProtectedContent resolves a claim and checks if it's a protected claim which would require authentication
func IsProtectedContent(claimID string) (bool, error) {
	claim, err := lbry.SDK.GetClaim(claimID)
	if err != nil {
		return false, err
	}

	for _, t := range claim.Value.GetTags() {
		if t == "c:members-only" {
			return true, nil
		}
	}
	return false, nil
}

// IsLivestreamClaim resolves a claim and checks if it has a source
func IsLivestreamClaim(claimID string) (bool, error) {
	claim, err := lbry.SDK.GetClaim(claimID)
	if err != nil {
		return true, err
	}

	if claim.Value.GetStream() != nil && claim.Value.GetStream().GetSource() == nil {
		return true, nil
	}

	return false, nil
}

var claimToChannelExistsCache = ccache.New(ccache.Configure().MaxSize(10000))

// EnsureClaimToChannelExists ensures that a claim to channel exists for the given claim id
func EnsureClaimToChannelExists(claimID string) error {
	// check cache first. it's only storing a boolean but it lets us do db upserts less.
	_, err := claimToChannelExistsCache.Fetch(claimID, 24*time.Hour, func() (interface{}, error) {

		// SDK calls have their own cache and the GetClaim call is done multiple times in create,
		// so this doesn't add much overhead.
		claim, err := lbry.SDK.GetClaim(claimID)
		if err != nil {
			return true, err
		}

		// It may be an anonymous channel.
		channel := claim.SigningChannel
		if channel == nil {
			return true, nil
		}

		// Create the claim to channel.
		cl2ch := &m.ClaimToChannel{
			ClaimID:   claimID,
			ChannelID: channel.ClaimID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Upsert it.
		err = cl2ch.Upsert(db.RW, boil.None(), boil.Infer())
		if err != nil {
			return false, err
		}

		return true, nil
	})

	return err
}

// HasAccessToProtectedContent checks if a channel has access to a protected claim
func HasAccessToProtectedContent(claimID, channelID string) (bool, error) {
	contentType := "Exclusive content"
	isLivestream, err := IsLivestreamClaim(claimID)
	if err != nil {
		return true, err
	}
	if isLivestream {
		contentType = "Exclusive livestreams"
	}

	hasAccess, err := lbry.API.CheckPerk(lbry.CheckPerkOptions{
		ChannelClaimID: channelID,
		ClaimID:        claimID,
		Type:           contentType,
	})
	if err != nil {
		return false, err
	}
	return hasAccess, nil
}

// HasAccessToProtectedChat checks if a channel has access to chat perk (members only mode)
func HasAccessToProtectedChat(claimID, channelID string) (bool, error) {
	hasAccess, err := lbry.API.CheckPerk(lbry.CheckPerkOptions{
		ChannelClaimID: channelID,
		ClaimID:        claimID,
		Type:           "Members-only chat",
	})
	if err != nil {
		return false, err
	}
	return hasAccess, nil
}

type modStatus struct {
	IsGlobalMod bool
	IsCreator   bool
	IsModerator bool
}

var modStatusCache = ccache.New(ccache.Configure().MaxSize(100000))

func applyModStatus(item *commentapi.CommentItem, channelID, claimID string) error {
	// Define a unique key for the cache based on channelID and claimID
	cacheKey := channelID + ":" + claimID

	// Attempt to retrieve the cached result
	cachedStatus := modStatusCache.Get(cacheKey)
	if cachedStatus != nil {
		// If cache hit, use the cached result
		if status, ok := cachedStatus.Value().(*modStatus); ok {
			item.IsGlobalMod = status.IsGlobalMod
			item.IsCreator = status.IsCreator
			item.IsModerator = status.IsModerator
			return nil
		}
	}

	// Cache miss, proceed to check mod status
	var isCreator bool
	var isModerator bool
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
		isCreator = channelID == signingChannel.ClaimID
		item.IsCreator = isCreator
		filterCreator := m.DelegatedModeratorWhere.CreatorChannelID.EQ(signingChannel.ClaimID)
		filterCommenter := m.DelegatedModeratorWhere.ModChannelID.EQ(channelID)
		isModerator, err = m.DelegatedModerators(filterCreator, filterCommenter).Exists(db.RO)
		if err != nil {
			return errors.Err(err)
		}
		item.IsModerator = isModerator
	}

	// Cache the moderation status
	modStatus := &modStatus{
		IsGlobalMod: isGlobalMod,
		IsCreator:   isCreator,
		IsModerator: isModerator,
	}
	modStatusCache.Set(cacheKey, modStatus, time.Minute*10)

	return nil
}

func pushItem(item commentapi.CommentItem, claimID string, mentionedChannels []commentapi.MentionedChannel) {
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

	for _, mc := range mentionedChannels {
		go sockety.SendNotification(socketyapi.SendNotificationArgs{
			Service: socketyapi.Commentron,
			Type:    "mention",
			IDs:     []string{mc.ChannelID, "mentions"},
			Data:    map[string]interface{}{"comment": item, "channel": mc.ChannelName, "channel_id": mc.ChannelID},
		})
	}

}

func checkForDuplicate(commentID string) error {
	// ignore checking for soft delete in this context
	comment, err := m.Comments(
		m.CommentWhere.CommentID.EQ(commentID),
		qm.WithDeleted(),
	).One(db.RO)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return errors.Err(err)
	}

	if comment != nil {
		return api.StatusError{Err: errors.Err("duplicate comment!"), Status: http.StatusBadRequest}
	}
	return nil
}

func checkForDuplicateTxID(txID string) error {
	// ignore checking for soft delete in this context
	comment, err := m.Comments(
		m.CommentWhere.TXID.EQ(null.StringFrom(txID)),
		qm.WithDeleted(),
	).One(db.RO)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return errors.Err(err)
	}

	if comment != nil {
		return api.StatusError{Err: errors.Err("txid already associated with a different comment!"), Status: http.StatusBadRequest}
	}
	return nil
}

var slowModeCache = ccache.New(ccache.Configure().MaxSize(10000))

type createRequest struct {
	args             *commentapi.CreateArgs
	comment          *m.Comment
	creatorChannel   *m.Channel
	commenterChannel *m.Channel
	signingChannel   *jsonrpc.Claim
	currency         string
	isFiat           bool
	isLivestream     bool
}

const maxSimilaryScoreToCreatorName = 0.8

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
			expiresIn := blockedListEntry.Expiry.Time.Sub(time.Now())
			timeLeft := helper.FormatDur(expiresIn)
			message := fmt.Sprintf("channel %s added you to the shared block list %s and you will not be able to comment with %s remaining.", blockedByChannel, blockedListName, timeLeft)
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

const maxSimilaryScoreToBlockedWord = 0.6

func checkMinTipAmountComment(settings *m.CreatorSetting, request *createRequest) error {
	if settings.MinTipAmountComment.IsZero() {
		return nil
	}
	if request.args.PaymentIntentID != nil || request.comment.Amount.IsZero() {
		return api.StatusError{Err: errors.Err("you must include LBC tip in order to comment as required by creator"), Status: http.StatusBadRequest}
	}
	if request.comment.Amount.Uint64 < settings.MinTipAmountComment.Uint64 {
		return api.StatusError{Err: errors.Err("you must tip at least %.2f LBC with this comment as required by %s", btcutil.Amount(settings.MinTipAmountComment.Uint64).ToBTC(), request.creatorChannel.Name), Status: http.StatusBadRequest}
	}
	return nil
}

func checkMinUsdcTipAmountComment(settings *m.CreatorSetting, request *createRequest) error {
	if settings.MinUsdcTipAmountComment.IsZero() {
		return nil
	}
	if (request.args.PaymentIntentID == nil && (request.args.Currency == nil || *request.args.Currency != "USDC")) || request.comment.Amount.IsZero() {
		return api.StatusError{Err: errors.Err("you must include USDC tip in order to comment as required by creator"), Status: http.StatusBadRequest}
	}
	if request.comment.Amount.Uint64 < settings.MinUsdcTipAmountComment.Uint64 {
		return api.StatusError{Err: errors.Err("you must tip at least %.2f USDC with this comment as required by %s", float64(settings.MinUsdcTipAmountComment.Uint64)/float64(100), request.creatorChannel.Name), Status: http.StatusBadRequest}
	}
	return nil
}

func checkMinTipAmountSuperChat(settings *m.CreatorSetting, request *createRequest) error {
	if settings.MinTipAmountSuperChat.IsZero() {
		return nil
	}
	if request.args.PaymentIntentID != nil || request.comment.Amount.Uint64 < settings.MinTipAmountSuperChat.Uint64 {
		return api.StatusError{Err: errors.Err("a min tip of %.2f LBC is required to hyperchat", btcutil.Amount(settings.MinTipAmountSuperChat.Uint64).ToBTC()), Status: http.StatusBadRequest}
	}
	return nil
}

func checkMinUsdcTipAmountSuperChat(settings *m.CreatorSetting, request *createRequest) error {
	if settings.MinUsdcTipAmountSuperChat.IsZero() {
		return nil
	}
	if (request.args.PaymentIntentID == nil && (request.args.Currency == nil || *request.args.Currency != "USDC")) || request.comment.Amount.Uint64 < settings.MinUsdcTipAmountSuperChat.Uint64 {
		return api.StatusError{Err: errors.Err("a min tip of %.2f USDC is required to hyperchat", (float64(settings.MinUsdcTipAmountSuperChat.Uint64) / float64(100))), Status: http.StatusBadRequest}
	}
	return nil
}

func checkSettings(settings *m.CreatorSetting, request *createRequest) error {
	isMod, err := m.DelegatedModerators(m.DelegatedModeratorWhere.ModChannelID.EQ(request.args.ChannelID), m.DelegatedModeratorWhere.CreatorChannelID.EQ(request.signingChannel.ClaimID)).Exists(db.RO)
	if err != nil {
		return errors.Err(err)
	}
	if !isMod && request.args.ChannelID != request.creatorChannel.ClaimID {
		if useOldTipAmountChecks {
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
		} else {
			if !request.comment.Amount.IsZero() {
				if request.args.PaymentIntentID == nil && request.args.Currency == nil {
					err = checkMinTipAmountSuperChat(settings, request)
				} else {
					err = checkMinUsdcTipAmountSuperChat(settings, request)
				}
				if err != nil {
					return err
				}
			}
			if !settings.MinTipAmountComment.IsZero() || !settings.MinUsdcTipAmountComment.IsZero() {
				if request.comment.Amount.IsZero() {
					return api.StatusError{Err: errors.Err("you must include tip in order to comment as required by creator"), Status: http.StatusBadRequest}
				}
				if request.args.PaymentIntentID == nil && request.args.Currency == nil {
					err = checkMinTipAmountComment(settings, request)
				} else {
					err = checkMinUsdcTipAmountComment(settings, request)
				}
				if err != nil {
					return err
				}
			}
		}

		if !settings.SlowModeMinGap.IsZero() {
			err := checkMinGap(request.args.ChannelID+request.creatorChannel.ClaimID, time.Duration(settings.SlowModeMinGap.Uint64)*time.Second, request.args.DryRun)
			if err != nil {
				return err
			}
		}
		if !settings.MutedWords.IsZero() {
			blockedWords := strings.Split(settings.MutedWords.String, ",")
			lowerComment := strings.ToLower(request.args.CommentText)
			for _, blockedWord := range blockedWords {
				lowerBlockedWord := strings.ToLower(blockedWord)
				if strings.Contains(lowerComment, lowerBlockedWord) {
					return api.StatusError{Err: errors.Err("the comment contents are blocked by %s", request.signingChannel.Name)}
				} else if strsim.Compare(lowerComment, lowerBlockedWord) > maxSimilaryScoreToBlockedWord {
					return api.StatusError{Err: errors.Err("the comment contents are blocked (by %s)", request.signingChannel.Name)}
				}
				if settings.BlockedWordsFuzzinessMatch.Valid {
					for _, commentWord := range strings.Split(lowerComment, " ") {
						if strsim.Compare(commentWord, lowerBlockedWord) > float64(settings.BlockedWordsFuzzinessMatch.Int64)/100.0 {
							return api.StatusError{Err: errors.Err("the comment contents are blocked [by %s]", request.signingChannel.Name)}
						}
					}
				}
			}
		}
		if request.isLivestream {
			if settings.LivestreamChatMembersOnly {
				hasAccess, err := HasAccessToProtectedChat(request.args.ClaimID, request.args.ChannelID)
				if err != nil {
					return err
				}
				if !hasAccess {
					return api.StatusError{Err: errors.Err("livestream chats are set to members only by the creator"), Status: http.StatusForbidden}
				}
			}
		} else {
			if settings.CommentsMembersOnly {
				hasAccess, err := HasAccessToProtectedChat(request.args.ClaimID, request.args.ChannelID)
				if err != nil {
					return err
				}
				if !hasAccess {
					return api.StatusError{Err: errors.Err("comments are set to members only by the creator"), Status: http.StatusForbidden}
				}
			}
		}
	}

	if !settings.CommentsEnabled.Bool {
		return api.StatusError{Err: errors.Err("comments are disabled by the creator"), Status: http.StatusForbidden}
	}

	if settings.TimeSinceFirstComment.Valid {
		request.commenterChannel, err = helper.FindOrCreateChannel(request.args.ChannelID, request.args.ChannelName)
		if err != nil {
			return errors.Err(err)
		}
		if time.Since(request.commenterChannel.CreatedAt) < time.Duration(settings.TimeSinceFirstComment.Int64)*time.Minute {
			return api.StatusError{Err: errors.Err(fmt.Sprintf("this creator has set minimum account age requirements that are not currently met: %d minutes", settings.TimeSinceFirstComment.Int64)), Status: http.StatusBadRequest}
		}
	}
	return nil
}

func checkMinGap(key string, expiration time.Duration, dryRun bool) error {
	creatorCounter, err := getCounter(key, expiration)
	if err != nil {
		return err
	}
	if creatorCounter.Get() > 0 {
		minGapViolated := fmt.Sprintf("Slow mode is on. Please wait at most %d seconds before commenting again.", int(expiration.Seconds()))
		return api.StatusError{Err: errors.Err(minGapViolated), Status: http.StatusBadRequest}
	}
	if !dryRun {
		creatorCounter.Add(1)
	}

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

func handleUsdcTip(request *createRequest) {
	client := graphql.NewClient("https://arweave.net/graphql")

	query := `
		query($ids: [ID!]) {
			transactions(
				ids: $ids
			) {
				edges {
					node {
						id
						recipient
						owner {
                    		address
                		}
						tags {
							name
							value
                		}
					}
				}
			}
		}	
	`

	req := graphql.NewRequest(query)

	if request.args.SupportTxID == nil {
		logrus.Error(fmt.Sprintf("Can't verify the tip, support txid not given. CommentID: %s", request.comment.CommentID))
		return
	}
	defaultErrorInfo := fmt.Sprintf("TxID: %s CommentID: %s", *request.args.SupportTxID, request.comment.CommentID)

	err := checkForDuplicateTxID(*request.args.SupportTxID)
	if err != nil {
		logrus.Error(fmt.Sprintf("%v %s", err.Error(), defaultErrorInfo))
	}

	req.Var("ids", []string{*request.args.SupportTxID})

	var respData struct {
		Transactions struct {
			Edges []struct {
				Node struct {
					ID        string `json:"id"`
					Recipient string `json:"recipient"`
					Owner     struct {
						Address string `json:"address"`
					} `json:"owner"`
					Tags []struct {
						Name  string `json:"name"`
						Value string `json:"value"`
					} `json:"Tags"`
				} `json:"node"`
			} `json:"edges"`
		} `json:"transactions"`
	}

	amount := request.comment.Amount.Uint64
	triesLeft := 2
	for triesLeft > 0 {
		triesLeft--
		ctx := context.Background()
		if err := client.Run(ctx, req, &respData); err != nil {
			logrus.Error(fmt.Sprintf("failed to execute query: %v", err.Error()))
		}

		if triesLeft == 0 {
			if len(respData.Transactions.Edges) == 0 {
				logrus.Error(fmt.Sprintf("Tx for id %s not found", *request.args.SupportTxID))
			}
			// If tx can't be found let it pass (assume a delay with gateway indexing)
			request.comment.Amount.SetValid(amount)
			request.comment.Currency.SetValid(*request.args.Currency)
			return
		}
		time.Sleep(3 * time.Second)
	}

	if len(respData.Transactions.Edges[0].Node.Tags) == 0 {
		logrus.Error(fmt.Sprintf("No tags found. %s", defaultErrorInfo))
	}

	wUsdcPid := "7zH9dlMNoxprab9loshv3Y7WG45DOny_Vrq9KrXObdQ"
	if respData.Transactions.Edges[0].Node.Recipient != wUsdcPid {
		logrus.Error(fmt.Sprintf("Expected recipient %s, got %s. %s", wUsdcPid, respData.Transactions.Edges[0].Node.Recipient, defaultErrorInfo))
	}

	if request.args.Owner != nil {
		foundOwner := respData.Transactions.Edges[0].Node.Owner.Address
		if foundOwner != *request.args.Owner {
			logrus.Error(fmt.Sprintf("Expected Owner %s, got %s. %s", *request.args.Owner, foundOwner, defaultErrorInfo))
		}
	} else {
		logrus.Error(fmt.Sprintf("Given nil Owner. %s", defaultErrorInfo))
	}

	tagsLeftToCheck := []string{"Action", "Quantity", "Recipient", "TimeStamp", "Signature", "SignatureTS"}

	var signature string
	var signatureTS string
	for i := 0; i < len(respData.Transactions.Edges[0].Node.Tags); i++ {
		tag := respData.Transactions.Edges[0].Node.Tags[i]
		switch tag.Name {
		case "Action":
			if tag.Value != "Transfer" {
				logrus.Error(fmt.Sprintf("Action not Transfer. %s", defaultErrorInfo))
			}
		case "Quantity":
			quantity, err := strconv.Atoi(tag.Value)
			if err != nil {
				logrus.Error(fmt.Sprintf("Failed to parse Quantity %v. %s", err.Error(), defaultErrorInfo))
			}
			quantityCents := uint64(quantity / 10000)
			amount = quantityCents
			if quantityCents != request.comment.Amount.Uint64 {
				logrus.Error(fmt.Sprintf("Quantity amount mismatch. Expected %d, got %d. %s", request.comment.Amount.Uint64, quantityCents, defaultErrorInfo))
			}
		case "Recipient":
			if request.args.Recipient != nil {
				if tag.Value != *request.args.Recipient {
					logrus.Error(fmt.Sprintf("Transfer recipient mismatch. Expected %s, got %s. %s", *request.args.Recipient, tag.Value, defaultErrorInfo))
				}
			} else {
				logrus.Error(fmt.Sprintf("Given nil Transfer Recipient. %s", defaultErrorInfo))
			}
		case "TimeStamp":
			allowedTimeDifference, _ := time.ParseDuration("5m")
			timeStamp, err := strconv.Atoi(tag.Value)
			if err != nil {
				logrus.Error(fmt.Sprintf("Failed to parse timestamp %v. %s", err.Error(), defaultErrorInfo))
			}
			deltaMs := math.Abs(float64(int64(timeStamp) - time.Now().UnixMilli()))
			if int64(deltaMs) > allowedTimeDifference.Milliseconds() {
				parsedDelta, _ := time.ParseDuration(fmt.Sprintf("%dms", int64(deltaMs)))
				logrus.Error(fmt.Sprintf("Timestamp %d over allowed difference of %dm, difference %dm. %s", timeStamp, int64(allowedTimeDifference.Minutes()), int64(parsedDelta.Minutes()), defaultErrorInfo))
			}
		case "Signature":
			signature = tag.Value
		case "SignatureTS":
			signatureTS = tag.Value
		}

		err := lbry.ValidateSignatureAndTS(request.args.ChannelID, signature, signatureTS, request.args.ChannelName)
		if err != nil {
			logrus.Error(fmt.Sprintf("%v %s", errors.Prefix("could not authenticate channel signature:", err), defaultErrorInfo))
		}

		for i, v := range tagsLeftToCheck {
			if v == tag.Name {
				tagsLeftToCheck = append(tagsLeftToCheck[:i], tagsLeftToCheck[i+1:]...)
			}
		}
	}

	if len(tagsLeftToCheck) != 0 {
		logrus.Error(fmt.Sprintf("Didn't found tags %v from the tx", tagsLeftToCheck))
	}

	request.comment.Amount.SetValid(amount)
	request.comment.Currency.SetValid(*request.args.Currency)

}

func updateSupportInfo(request *createRequest) error {
	if request.args.Currency != nil && *request.args.Currency == "USDC" {
		handleUsdcTip(request)
		return nil
	}
	triesLeft := 3
	for {
		triesLeft--
		err := updateSupportInfoAttempt(request, true)
		if err == nil {
			return nil
		}
		if triesLeft == 0 {
			return err
		}
		time.Sleep(1 * time.Second)
	}
}

func updateSupportInfoAttempt(request *createRequest, retry bool) error {
	if request.args.PaymentIntentID != nil {
		env := ""
		if request.args.Environment != nil {
			env = *request.args.Environment
		}
		paymentintentClient := &paymentintent.Client{B: stripe.GetBackend(stripe.APIBackend), Key: config.ConnectAPIKey(config.From(env))}
		pi, err := paymentintentClient.Get(*request.args.PaymentIntentID, &stripe.PaymentIntentParams{})
		if err != nil {
			if !retry {
				logrus.Error(errors.Prefix("could not get payment intent %s", *request.args.PaymentIntentID))
				return errors.Err("could not validate tip")
			}
			// in the rare event that the payment intent is not found, wait a bit and try again once
			time.Sleep(5 * time.Second)
			return updateSupportInfoAttempt(request, false)
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
