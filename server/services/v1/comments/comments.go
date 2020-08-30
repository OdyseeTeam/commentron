package comments

import (
	"fmt"

	m "github.com/lbryio/commentron/model"

	"github.com/volatiletech/null"
)

// CommentItem is the data structure of a comment returned from commentron
type CommentItem struct {
	Comment     string `json:"comment"`
	CommentID   string `json:"comment_id"`
	ClaimID     string `json:"claim_id"`
	Timestamp   int    `json:"timestamp"`
	ParentID    string `json:"parent_id,omitempty"`
	Signature   string `json:"signature,omitempty"`
	SigningTs   string `json:"signing_ts,omitempty"`
	IsHidden    bool   `json:"is_hidden"`
	ChannelID   string `json:"channel_id,omitempty"`
	ChannelName string `json:"channel_name,omitempty"`
	ChannelURL  string `json:"channel_url,omitempty"`
}

func populateItem(comment *m.Comment, channel *m.Channel) CommentItem {
	var channelName null.String
	var channelURL null.String
	if channel != nil {
		channelName = null.StringFrom(channel.Name)
		channelURL = null.StringFrom(fmt.Sprintf("lbry://%s#%s", channelName.String, comment.ChannelID.String))
	}

	item := CommentItem{
		Comment:     comment.Body,
		CommentID:   comment.CommentID,
		ClaimID:     comment.LbryClaimID,
		Timestamp:   comment.Timestamp,
		ParentID:    comment.ParentID.String,
		Signature:   comment.Signature.String,
		SigningTs:   comment.Signingts.String,
		IsHidden:    comment.IsHidden.Bool,
		ChannelID:   comment.ChannelID.String,
		ChannelName: channelName.String,
		ChannelURL:  channelURL.String,
	}

	return item
}
