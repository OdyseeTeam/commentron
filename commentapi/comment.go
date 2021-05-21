package commentapi

// CommentItem is the data structure of a comment returned from commentron
type CommentItem struct {
	Comment       string  `json:"comment"`
	CommentID     string  `json:"comment_id"`
	ClaimID       string  `json:"claim_id"`
	Timestamp     int     `json:"timestamp"`
	ParentID      string  `json:"parent_id,omitempty"`
	Signature     string  `json:"signature,omitempty"`
	SigningTs     string  `json:"signing_ts,omitempty"`
	IsHidden      bool    `json:"is_hidden"`
	IsPinned      bool    `json:"is_pinned"`
	ChannelID     string  `json:"channel_id,omitempty"`
	ChannelName   string  `json:"channel_name,omitempty"`
	ChannelURL    string  `json:"channel_url,omitempty"`
	Replies       int     `json:"replies,omitempty"`
	SupportAmount float64 `json:"support_amount"`
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
	CommentText string  `json:"comment"`
	ClaimID     string  `json:"claim_id"`
	ParentID    *string `json:"parent_id"`
	ChannelID   string  `json:"channel_id"`
	ChannelName string  `json:"channel_name"`
	SupportTxID *string `json:"support_tx_id"`
	SupportVout *uint64 `json:"support_vout"`
	Signature   string  `json:"signature"`
	SigningTS   string  `json:"signing_ts"`
}

// CreateResponse response for the comment.Create rpc call
type CreateResponse struct {
	*CommentItem
}

// ByIDArgs arguments for the comment.List rpc call
type ByIDArgs struct {
	CommentID string `json:"comment_id"`
}

// ByIDResponse response for the comment.ByID rpc call
type ByIDResponse struct {
	Item CommentItem `json:"items,omitempty"`
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

// ListArgs arguments for the comment.List rpc call
type ListArgs struct {
	ClaimID       *string `json:"claim_id"`
	AuthorClaimID *string `json:"author_claim_id"`
	ParentID      *string `json:"parent_id"`
	Page          int     `json:"page"`
	PageSize      int     `json:"page_size"`
	TopLevel      bool    `json:"top_level"`
	Hidden        bool    `json:"hidden"`
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
	if c.PageSize > 500 {
		c.PageSize = 500
	}
}

// ListResponse response for the comment.List rpc call
type ListResponse struct {
	Page              int           `json:"page"`
	PageSize          int           `json:"page_size"`
	TotalPages        int           `json:"total_pages"`
	TotalItems        int64         `json:"total_items"`
	Items             []CommentItem `json:"items,omitempty"`
	HasHiddenComments bool          `json:"has_hidden_comments"`
}
