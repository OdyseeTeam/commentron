package comments

import (
	"time"

	"github.com/lbryio/commentron/model"
	"github.com/lbryio/commentron/server/lbry"
	"github.com/lbryio/lbry.go/v2/extras/errors"
	"github.com/volatiletech/sqlboiler/boil"
)

// EditArgs arguments for the comment.Edit rpc call
type EditArgs struct {
	Comment   string `json:"comment"`
	CommentID string `json:"comment_id"`
	Signature string `json:"signature"`
	SigningTS string `json:"signing_ts"`
}

// EditResponse response for the comment.Edit rpc call
type EditResponse struct {
	*CommentItem
}

func edit(args *EditArgs) (*CommentItem, error) {
	comment, err := model.Comments(model.CommentWhere.CommentID.EQ(args.CommentID)).OneG()
	if err != nil {
		return nil, errors.Err(err)
	}
	channel, err := model.Channels(model.ChannelWhere.ClaimID.EQ(comment.ChannelID.String)).OneG()
	if err != nil {
		return nil, errors.Err(err)
	}
	err = lbry.ValidateSignature(comment.ChannelID.String, args.Signature, args.SigningTS, args.Comment)
	if err != nil {
		return nil, err
	}

	comment.Body = args.Comment
	comment.Signature.SetValid(args.Signature)
	comment.Signingts.SetValid(args.SigningTS)
	comment.Timestamp = int(time.Now().Unix())
	err = comment.UpdateG(boil.Infer())
	if err != nil {
		return nil, errors.Err(err)
	}
	item := populateItem(comment, channel)
	err = comment.DeleteG()
	if err != nil {
		return nil, errors.Err(err)
	}
	return &item, nil
}
