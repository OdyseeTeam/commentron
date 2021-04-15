package comments

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/lbryio/commentron/server/websocket"

	"github.com/lbryio/commentron/commentapi"
	"github.com/lbryio/commentron/flags"
	"github.com/lbryio/commentron/helper"
	m "github.com/lbryio/commentron/model"
	"github.com/lbryio/commentron/server/lbry"

	"github.com/lbryio/lbry.go/extras/api"
	"github.com/lbryio/lbry.go/extras/util"
	"github.com/lbryio/lbry.go/v2/extras/errors"
	v "github.com/lbryio/ozzo-validation"

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
		comment.TXID.SetValid(util.StrFromPtr(args.SupportTxID))
	}

	err = flags.CheckComment(comment)
	if err != nil {
		return err
	}

	err = errors.Err(comment.InsertG(boil.Infer()))
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
		return api.StatusError{Err: errors.Err("channel is blocked by publisher")}
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
