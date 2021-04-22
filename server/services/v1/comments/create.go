package comments

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/lbryio/commentron/config"

	"github.com/lbryio/sockety/socketyapi"

	"github.com/lbryio/commentron/commentapi"
	"github.com/lbryio/commentron/flags"
	"github.com/lbryio/commentron/helper"
	m "github.com/lbryio/commentron/model"
	"github.com/lbryio/commentron/server/lbry"
	"github.com/lbryio/commentron/server/websocket"

	"github.com/lbryio/lbry.go/v2/extras/api"
	"github.com/lbryio/lbry.go/v2/extras/errors"
	"github.com/lbryio/lbry.go/v2/extras/jsonrpc"
	"github.com/lbryio/lbry.go/v2/extras/util"
	v "github.com/lbryio/ozzo-validation"

	"github.com/btcsuite/btcutil"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
)

func create(_ *http.Request, args *commentapi.CreateArgs, reply *commentapi.CreateResponse) error {
	err := v.ValidateStruct(args,
		v.Field(&args.ClaimID, v.Required))
	if err != nil {
		return api.StatusError{Err: errors.Err(err), Status: http.StatusBadRequest}
	}
	channel, err := m.Channels(m.ChannelWhere.ClaimID.EQ(null.StringFrom(args.ChannelID).String)).OneG()
	if errors.Is(err, sql.ErrNoRows) {
		channel = &m.Channel{
			ClaimID: null.StringFrom(args.ChannelID).String,
			Name:    null.StringFrom(args.ChannelName).String,
		}
		err = nil
		err := channel.InsertG(boil.Infer())
		if err != nil {
			return errors.Err(err)
		}
	}
	blockedEntry, err := m.BlockedEntries(m.BlockedEntryWhere.UniversallyBlocked.EQ(null.BoolFrom(true)), m.BlockedEntryWhere.BlockedChannelID.EQ(null.StringFrom(args.ChannelID))).OneG()
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

	commentID, timestamp, err := createCommentID(args.CommentText, null.StringFrom(args.ChannelID).String)
	if err != nil {
		return errors.Err(err)
	}

	comment, err := m.Comments(m.CommentWhere.CommentID.EQ(commentID)).OneG()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return errors.Err(err)
	}

	if comment != nil {
		return api.StatusError{Err: errors.Err("duplicate comment!"), Status: http.StatusBadRequest}
	}
	err = blockedByCreator(args.ClaimID, args.ChannelID, args.CommentText)
	if err != nil {
		return errors.Err(err)
	}

	comment = &m.Comment{
		CommentID:   commentID,
		LbryClaimID: args.ClaimID,
		ChannelID:   null.StringFrom(args.ChannelID),
		Body:        args.CommentText,
		ParentID:    null.StringFromPtr(args.ParentID),
		Signature:   null.StringFrom(args.Signature),
		Signingts:   null.StringFrom(args.SigningTS),
		Timestamp:   int(timestamp),
	}

	if args.SupportTxID != nil {
		err := updateSupportInfo(channel.ClaimID, comment, args.SupportTxID, args.SupportVout)
		if err != nil {
			return errors.Err(err)
		}
	}

	err = flags.CheckComment(comment)
	if err != nil {
		return err
	}

	err = comment.InsertG(boil.Infer())
	if err != nil {
		return errors.Err(err)
	}
	item := populateItem(comment, channel, 0)
	reply.CommentItem = &item

	go pushItem(item, args.ClaimID)
	go lbry.API.Notify(lbry.NotifyOptions{
		ActionType: "C",
		CommentID:  item.CommentID,
		ChannelID:  &item.ChannelID,
		ParentID:   &item.ParentID,
		Comment:    &item.Comment,
		ClaimID:    item.ClaimID,
	})
	return nil
}

func pushItem(item commentapi.CommentItem, claimID string) {
	websocket.PushTo(&websocket.PushNotification{
		Type: "delta",
		Data: map[string]interface{}{"comment": item},
	}, claimID)

	go sendMessage(item, claimID)

}

func sendMessage(item commentapi.CommentItem, claimID string) {
	resp, err := socketyapi.NewClient("https://sockety.lbry.com", config.SocketyToken).SendNotification(socketyapi.SendNotificationArgs{
		Service: socketyapi.Commentron,
		Type:    "delta",
		IDs:     []string{claimID},
		Data:    map[string]interface{}{"comment": item},
	})
	if err != nil {
		logrus.Error(errors.Prefix("Sockety SendTo: ", err))
	}
	if resp != nil && resp.Error != nil {
		logrus.Error(errors.Prefix("Sockety SendToResp: ", errors.Base(*resp.Error)))
	}
}

func blockedByCreator(contentClaimID, commenterChannelID, comment string) error {
	signingChannel, err := lbry.SDK.GetSigningChannelForClaim(contentClaimID)
	if err != nil {
		return errors.Err(err)
	}
	if signingChannel == nil {
		return nil
	}

	blockedEntry, err := m.BlockedEntries(m.BlockedEntryWhere.BlockedByChannelID.EQ(null.StringFrom(signingChannel.ClaimID)), m.BlockedEntryWhere.BlockedChannelID.EQ(null.StringFrom(commenterChannelID))).OneG()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return errors.Err(err)
	}

	if blockedEntry != nil {
		return api.StatusError{Err: errors.Err("channel is blocked by publisher"), Status: http.StatusBadRequest}
	}

	creatorChannel, err := helper.FindOrCreateChannel(signingChannel.ClaimID, signingChannel.Name)
	if err != nil {
		return err
	}
	settings, err := creatorChannel.CreatorChannelCreatorSettings().OneG()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return errors.Err(err)
	}
	if settings != nil && !settings.MutedWords.IsZero() {
		blockedWords := strings.Split(settings.MutedWords.String, ",")
		for _, blockedWord := range blockedWords {
			if strings.Contains(comment, blockedWord) {
				return api.StatusError{Err: errors.Err("the comment contents are blocked by %s", signingChannel.Name)}
			}
		}
	}
	return nil
}

func allowedToPostReply(parentID, commenterClaimID string) error {
	parentComment, err := m.Comments(m.CommentWhere.CommentID.EQ(parentID)).OneG()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return errors.Err(err)
	}
	if parentComment != nil {
		parentChannel, err := parentComment.Channel().OneG()
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return errors.Err(err)
		}
		if parentChannel != nil {

			blockedEntry, err := m.BlockedEntries(
				m.BlockedEntryWhere.BlockedByChannelID.EQ(null.StringFrom(parentChannel.ClaimID)),
				m.BlockedEntryWhere.BlockedChannelID.EQ(null.StringFrom(commenterClaimID))).OneG()
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				return errors.Err(err)
			}
			if blockedEntry != nil {
				return api.StatusError{Err: errors.Err("'%s' has blocked you from replying to their comments", parentChannel.Name), Status: http.StatusBadRequest}
			}
		}
	}
	return nil
}

func updateSupportInfo(channelID string, comment *m.Comment, supportTxID *string, supportVout *uint64) error {
	comment.TXID.SetValid(util.StrFromPtr(supportTxID))
	txSummary, err := lbry.SDK.GetTx(comment.TXID.String)
	if err != nil {
		return errors.Err(err)
	}
	if txSummary == nil {
		return errors.Err("transaction not found for txid %s", comment.TXID.String)
	}
	var vout uint64
	if supportVout != nil {
		vout = *supportVout
	}
	amount, err := getVoutAmount(channelID, txSummary, vout)
	if err != nil {
		return errors.Err(err)
	}
	comment.Amount.SetValid(amount)
	return nil
}

func getVoutAmount(channelID string, summary *jsonrpc.TransactionSummary, vout uint64) (uint64, error) {
	if summary == nil {
		return 0, errors.Err("transaction summary missing")
	}

	if len(summary.Outputs) <= int(vout) {
		return 0, errors.Err("there are not enough outputs on the transaction to for position %d", vout)
	}
	output := summary.Outputs[int(vout)]

	if output.SigningChannel == nil {
		return 0, errors.Err("Expected signed support for %s in transaction %s", channelID, summary.Txid)
	}

	if output.SigningChannel.ChannelID != channelID {
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
