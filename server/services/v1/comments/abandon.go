package comments

import (
	"github.com/lbryio/commentron/commentapi"
	"github.com/lbryio/commentron/model"
	"github.com/lbryio/commentron/server/lbry"

	"github.com/lbryio/lbry.go/v2/extras/errors"
)

func abandon(args *commentapi.AbandonArgs) (*commentapi.CommentItem, error) {
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
	item := populateItem(comment, channel, 0)
	err = comment.DeleteG()
	if err != nil {
		return nil, errors.Err(err)
	}
	return &item, nil

}
