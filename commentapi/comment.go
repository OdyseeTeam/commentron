package commentapi

import (
	"net/http"

	"github.com/lbryio/commentron/validator"
	"github.com/lbryio/lbry.go/v2/extras/api"
	"github.com/lbryio/lbry.go/v2/extras/errors"
	v "github.com/lbryio/ozzo-validation"
)

// CommentItem is the data structure of a comment returned from commentron
type CommentItem struct {
	Comment       string  `json:"comment"`
	CommentID     string  `json:"comment_id"`
	ClaimID       string  `json:"claim_id"`
	Timestamp     int     `json:"timestamp"`
	ParentID      string  `json:"parent_id,omitempty"`
	Signature     string  `json:"signature,omitempty"`
	SigningTs     string  `json:"signing_ts,omitempty"`
	ChannelID     string  `json:"channel_id,omitempty"`
	ChannelName   string  `json:"channel_name,omitempty"`
	ChannelURL    string  `json:"channel_url,omitempty"`
	Currency      string  `json:"currency"`
	Replies       int     `json:"replies,omitempty"`
	SupportAmount float64 `json:"support_amount"`
	IsHidden      bool    `json:"is_hidden"`
	IsPinned      bool    `json:"is_pinned"`
	IsFiat        bool    `json:"is_fiat"`
}

// ChannelArgs arguments to the comment.GetChannelForCommentID call
type ChannelArgs struct {
	CommentID string `json:"comment_id"`
}

// ChannelResponse response to the comment.GetChannelForCommentID call
type ChannelResponse struct {
	ChannelID   string `json:"channel_id"`
	ChannelName string `json:"channel_name"`
}

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

// CreateArgs arguments for the comment.Create rpc call
type CreateArgs struct {
	CommentText     string  `json:"comment"`
	ClaimID         string  `json:"claim_id"`
	ParentID        *string `json:"parent_id"`
	ChannelID       string  `json:"channel_id"`
	ChannelName     string  `json:"channel_name"`
	SupportTxID     *string `json:"support_tx_id"`
	SupportVout     *uint64 `json:"support_vout"`
	PaymentIntentID *string `json:"payment_intent_id"`
	Environment     *string `json:"environment"`
	Signature       string  `json:"signature"`
	SigningTS       string  `json:"signing_ts"`
}

// CreateResponse response for the comment.Create rpc call
type CreateResponse struct {
	*CommentItem
}

// ByIDArgs arguments for the comment.List rpc call
type ByIDArgs struct {
	CommentID     string `json:"comment_id"`
	WithAncestors bool   `json:"with_ancestors"`
}

// ByIDResponse response for the comment.ByID rpc call
type ByIDResponse struct {
	Item      CommentItem   `json:"items,omitempty"`
	Ancestors []CommentItem `json:"ancestors,omitempty"`
}

// PinArgs arguments for the comment.Pin rpc call. The comment id must be signed with a timestamp for authentication.
type PinArgs struct {
	CommentID   string `json:"comment_id"`
	ChannelID   string `json:"channel_id"`
	ChannelName string `json:"channel_name"` //Technical debt? probably dont need this since we can get the channel from the comment claim
	Remove      bool   `json:"remove"`
	Signature   string `json:"signature"`
	SigningTS   string `json:"signing_ts"`
}

// PinResponse response for the comment.Pin rpc call
type PinResponse struct {
	Item CommentItem `json:"items,omitempty"`
}

// Sort defines the type of sort for the comment.List api
type Sort int

const (
	// Newest sorts the comments by newest first
	Newest Sort = iota
	// Oldest sorts the comments from first to last
	Oldest
	// Controversy sorts the comments by controversy
	Controversy
	// Popularity sorts the comments by how popular it is
	Popularity
)

// ListArgs arguments for the comment.List rpc call
type ListArgs struct {
	ChannelName   *string `json:"channel_name"`    // signing channel name of claim
	ChannelID     *string `json:"channel_id"`      // signing channel claim id of claim
	ClaimID       *string `json:"claim_id"`        // claim id of claim being commented on
	AuthorClaimID *string `json:"author_claim_id"` // filters comments to just this author
	ParentID      *string `json:"parent_id"`       // filters comments to those under this thread
	Page          int     `json:"page"`            // pagination: which page of results
	PageSize      int     `json:"page_size"`       // pagination: nr of comments to show in a page (max 200)
	TopLevel      bool    `json:"top_level"`       // filters to only top level comments
	Hidden        bool    `json:"hidden"`          // if true will show hidden comments as well
	SortBy        Sort    `json:"sort_by"`         // can be popularity, controversy, default is time (newest)
}

// AbandonArgs are the arguments passed to comment.Abandon RPC call. If creator args are passed
// the signing channel of the content of the comment is checked these args and signature
// verification happens against the creators public key for authorization.
type AbandonArgs struct {
	CommentID          string  `json:"comment_id"`
	Signature          string  `json:"signature"`
	SigningTS          string  `json:"signing_ts"`
	CreatorChannelID   *string `json:"creator_channel_id"`
	CreatorChannelName *string `json:"creator_channel_name"`
}

// AbandonResponse the response to the abandon call
type AbandonResponse struct {
	*CommentItem
	Abandoned bool `json:"abandoned"`
}

// ApplyDefaults applies the default values for arguments passed that are different from normal defaults.
func (c *ListArgs) ApplyDefaults() {
	if c.Page == 0 {
		c.Page = 1
	}

	if c.PageSize == 0 {
		c.PageSize = 50
	}
	if c.PageSize > 600 {
		c.PageSize = 600
	}
}

// ListResponse response for the comment.List rpc call
type ListResponse struct {
	Page               int           `json:"page"`
	PageSize           int           `json:"page_size"`
	TotalPages         int           `json:"total_pages"`
	TotalItems         int64         `json:"total_items"`
	TotalFilteredItems int64         `json:"total_filtered_items"`
	Items              []CommentItem `json:"items,omitempty"`
	HasHiddenComments  bool          `json:"has_hidden_comments"`
}

// Validate validates the data in the list args
func (c ListArgs) Validate() api.StatusError {
	err := v.ValidateStruct(&c,
		v.Field(&c.ChannelID, validator.ClaimID, v.Required),
		v.Field(&c.ChannelName, v.Required),
	)
	if err != nil {
		return api.StatusError{Err: errors.Err(err), Status: http.StatusBadRequest}
	}
	if c.ClaimID == nil && c.AuthorClaimID == nil {
		return api.StatusError{Err: errors.Err("you must pass either claim_id or author_claim_id"), Status: http.StatusBadRequest}
	}
	return api.StatusError{}
}
