package comments

import (
	"github.com/lbryio/commentron/server/lbry"
	"github.com/lbryio/lbry.go/v2/extras/errors"

	"github.com/lbryio/commentron/model"
)

// AbandonArgs are the arguments passed to comment.Abandon RPC call
type AbandonArgs struct {
	CommentID string `json:"comment_id"`
	Signature string `json:"signature"`
	SigningTS string `json:"signing_ts"`
}

// AbandonResponse the response to the abandon call
type AbandonResponse struct {
	*CommentItem
}

func abandon(args *AbandonArgs) (*CommentItem, error) {
	comment, err := model.Comments(model.CommentWhere.CommentID.EQ(args.CommentID)).OneG()
	if err != nil {
		return nil, errors.Err(err)
	}
	channel, err := model.Channels(model.ChannelWhere.ClaimID.EQ(comment.ChannelID.String)).OneG()
	if err != nil {
		return nil, errors.Err(err)
	}
	err = lbry.ValidateSignature(comment.ChannelID.String, args.Signature, args.SigningTS, args.CommentID)
	if err != nil {
		return nil, err
	}
	item := populateItem(comment, channel)
	err = comment.DeleteG()
	if err != nil {
		return nil, errors.Err(err)
	}
	return &item, nil

}
