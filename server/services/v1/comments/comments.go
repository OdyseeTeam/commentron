package comments

import (
	"fmt"

	"github.com/lbryio/commentron/commentapi"
	m "github.com/lbryio/commentron/model"

	"github.com/btcsuite/btcutil"
	"github.com/volatiletech/null/v8"
)

var currencyMap = map[string]uint64{"USD": 100}

func populateItem(comment *m.Comment, channel *m.Channel, replies int) commentapi.CommentItem {
	var channelName null.String
	var channelURL null.String
	if channel != nil {
		channelName = null.StringFrom(channel.Name)
		channelURL = null.StringFrom(fmt.Sprintf("lbry://%s#%s", channelName.String, comment.ChannelID.String))
	}
	supportAmount := btcutil.Amount(comment.Amount.Uint64).ToBTC()
	if comment.IsFiat {
		divisor := currencyMap[comment.Currency.String]
		if divisor == 0 {
			divisor = 100
		}
		supportAmount = float64(comment.Amount.Uint64) / float64(divisor)
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
		SupportAmount: supportAmount,
		IsFiat:        comment.IsFiat,
		Currency:      comment.Currency.String,
		IsProtected:   comment.IsProtected,
	}

	return item
}
