package comments

import (
	"fmt"

	"github.com/lbryio/commentron/commentapi"
	m "github.com/lbryio/commentron/model"

	"github.com/btcsuite/btcutil"
	"github.com/volatiletech/null"
)

func populateItem(comment *m.Comment, channel *m.Channel, replies int) commentapi.CommentItem {
	var channelName null.String
	var channelURL null.String
	if channel != nil {
		channelName = null.StringFrom(channel.Name)
		channelURL = null.StringFrom(fmt.Sprintf("lbry://%s#%s", channelName.String, comment.ChannelID.String))
	}

	item := commentapi.CommentItem{
		Comment:       comment.Body,
		CommentID:     comment.CommentID,
		ClaimID:       comment.LbryClaimID,
		Timestamp:     comment.Timestamp,
		ParentID:      comment.ParentID.String,
		Signature:     comment.Signature.String,
		SigningTs:     comment.Signingts.String,
		IsHidden:      comment.IsHidden.Bool,
		IsPinned:      comment.IsPinned,
		ChannelID:     comment.ChannelID.String,
		ChannelName:   channelName.String,
		ChannelURL:    channelURL.String,
		Replies:       replies,
		SupportAmount: btcutil.Amount(comment.Amount.Uint64).ToBTC(),
	}

	return item
}
