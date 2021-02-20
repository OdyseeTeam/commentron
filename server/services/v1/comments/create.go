package comments

import (
	"database/sql"
	"net/http"

	"github.com/lbryio/commentron/commentapi"
	"github.com/lbryio/commentron/flags"
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
	channel, err := m.Channels(m.ChannelWhere.ClaimID.EQ(null.StringFromPtr(args.ChannelID).String)).OneG()
	if errors.Is(err, sql.ErrNoRows) {
		channel = &m.Channel{
			ClaimID: null.StringFromPtr(args.ChannelID).String,
			Name:    null.StringFromPtr(args.ChannelName).String,
		}
		err = nil
		err := channel.InsertG(boil.Infer())
		if err != nil {
			return errors.Err(err)
		}
	}
	blockedEntry, err := m.BlockedEntries(m.BlockedEntryWhere.UniversallyBlocked.EQ(null.BoolFrom(true)), m.BlockedEntryWhere.BlockedChannelID.EQ(null.StringFromPtr(args.ChannelID))).OneG()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return errors.Err(err)
	}
	if blockedEntry != nil {
		return api.StatusError{Err: errors.Err("channel is not allowed to post comments"), Status: http.StatusBadRequest}
	}

	err = lbry.ValidateSignature(util.StrFromPtr(args.ChannelID), util.StrFromPtr(args.Signature), util.StrFromPtr(args.SigningTS), args.CommentText)
	if err != nil {
		return errors.Prefix("could not authenticate channel signature:", err)
	}

	commentID, timestamp, err := createCommentID(args.CommentText, null.StringFromPtr(args.ChannelID).String)
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
	signingChannel, err := lbry.GetSigningChannelForClaim(args.ClaimID)
	if err != nil {
		return errors.Err(err)
	}
	if signingChannel != nil {
		blockedEntry, err := m.BlockedEntries(m.BlockedEntryWhere.BlockedByChannelID.EQ(null.StringFrom(signingChannel.ClaimID)), m.BlockedEntryWhere.BlockedChannelID.EQ(null.StringFromPtr(args.ChannelID))).OneG()
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return errors.Err(err)
		}
		if blockedEntry != nil {
			return api.StatusError{Err: errors.Err("channel %s is blocked by publisher %s", args.ChannelID, signingChannel.Name)}
		}
	}

	comment = &m.Comment{
		CommentID:   commentID,
		LbryClaimID: args.ClaimID,
		ChannelID:   null.StringFromPtr(args.ChannelID),
		Body:        args.CommentText,
		ParentID:    null.StringFromPtr(args.ParentID),
		Signature:   null.StringFromPtr(args.Signature),
		Signingts:   null.StringFromPtr(args.SigningTS),
		Timestamp:   int(timestamp),
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

	go lbry.Notify(lbry.NotifyOptions{
		ActionType: "C",
		CommentID:  item.CommentID,
		ChannelID:  &item.ChannelID,
		ParentID:   &item.ParentID,
		Comment:    &item.Comment,
		ClaimID:    item.ClaimID,
	})
	return nil
}
