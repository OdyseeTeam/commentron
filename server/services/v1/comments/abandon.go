package comments

import (
	"github.com/lbryio/commentron/commentapi"
	"github.com/lbryio/commentron/helper"
	"github.com/lbryio/commentron/model"
	"github.com/lbryio/commentron/server/lbry"

	"github.com/lbryio/lbry.go/extras/api"
	"github.com/lbryio/lbry.go/extras/util"
	"github.com/lbryio/lbry.go/v2/extras/errors"
)

func abandon(args *commentapi.AbandonArgs) (*commentapi.CommentItem, error) {
	comment, err := model.Comments(model.CommentWhere.CommentID.EQ(args.CommentID)).OneG()
	if err != nil {
		return nil, errors.Err(err)
	}
	var channel *model.Channel
	if args.CreatorChannelID != nil && args.CreatorChannelName != nil {
		channel, err = helper.FindOrCreateChannel(util.StrFromPtr(args.CreatorChannelID), util.StrFromPtr(args.CreatorChannelName))
		if err != nil {
			return nil, errors.Err(err)
		}
		content, err := lbry.GetClaim(comment.LbryClaimID)
		if err != nil {
			return nil, errors.Err(err)
		}
		signingChannelClaimID := content.ClaimID
		if content.SigningChannel != nil {
			signingChannelClaimID = content.SigningChannel.ClaimID
		}
		if signingChannelClaimID != channel.ClaimID {
			return nil, api.StatusError{Err: errors.Err("you do not have creator authorizations to remove this comment on %s", comment.LbryClaimID)}
		}
	} else {
		channel, err = model.Channels(model.ChannelWhere.ClaimID.EQ(comment.ChannelID.String)).OneG()
		if err != nil {
			return nil, errors.Err(err)
		}
	}

	err = lbry.ValidateSignature(channel.ClaimID, args.Signature, args.SigningTS, args.CommentID)
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
