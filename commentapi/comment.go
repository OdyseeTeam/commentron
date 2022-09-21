package commentapi

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/OdyseeTeam/commentron/validator"

	"github.com/lbryio/lbry.go/v2/extras/api"
	"github.com/lbryio/lbry.go/v2/extras/errors"

	v "github.com/lbryio/ozzo-validation"
)

// StickerRE is the regex for a valid sticker as a comment.
var StickerRE = regexp.MustCompile(`^<stkr>:(?P<sticker>[a-zA-Z0-9_]+):<stkr>$`)

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
	IsCreator     bool    `json:"is_creator,omitempty"`
	IsModerator   bool    `json:"is_moderator,omitempty"`
	IsGlobalMod   bool    `json:"is_global_mod,omitempty"`
	IsHidden      bool    `json:"is_hidden"`
	IsPinned      bool    `json:"is_pinned"`
	IsFiat        bool    `json:"is_fiat"`
	IsProtected   bool    `json:"is_protected"`
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

// MentionedChannel channels mentioned in comment
type MentionedChannel struct {
	ChannelName string `json:"channel_name"`
	ChannelID   string `json:"channel_id"`
}

// CreateArgs arguments for the comment.Create rpc call
type CreateArgs struct {
	CommentText       string             `json:"comment"`
	ClaimID           string             `json:"claim_id"`
	ParentID          *string            `json:"parent_id"`
	ChannelID         string             `json:"channel_id"`
	ChannelName       string             `json:"channel_name"`
	Sticker           bool               `json:"sticker"`
	SupportTxID       *string            `json:"support_tx_id"`
	SupportVout       *uint64            `json:"support_vout"`
	PaymentIntentID   *string            `json:"payment_intent_id"`
	Environment       *string            `json:"environment"`
	Signature         string             `json:"signature"`
	SigningTS         string             `json:"signing_ts"`
	MentionedChannels []MentionedChannel `json:"mentioned_channels"`
	IsProtected       bool               `json:"is_protected"`
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
	Authorization

	CommentID string `json:"comment_id"`
	Remove    bool   `json:"remove"`
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
	// NewestNoPins sorts the comments by newest first but removes the presort for pinned comments
	NewestNoPins
)

// ListArgs arguments for the comment.List rpc call
type ListArgs struct {
	Authorization

	RequestorChannelName string  `json:"requestor_channel_name"` // Used for Author ID filter authorization Only your comments!
	RequestorChannelID   *string `json:"requestor_channel_id"`   // Used for Author ID filter authorization
	ClaimID              *string `json:"claim_id"`               // claim id of claim being commented on
	AuthorClaimID        *string `json:"author_claim_id"`        // filters comments to just this author
	ParentID             *string `json:"parent_id"`              // filters comments to those under this thread
	Page                 int     `json:"page"`                   // pagination: which page of results
	PageSize             int     `json:"page_size"`              // pagination: nr of comments to show in a page (max 200)
	TopLevel             bool    `json:"top_level"`              // filters to only top level comments
	Hidden               bool    `json:"hidden"`                 // if true will show hidden comments as well
	SortBy               Sort    `json:"sort_by"`                // can be popularity, controversy, default is time (newest)
	IsProtected          bool    `json:"is_protected"`           // if true, only return protected when authorized
}

// Key returns the hash of the list args struct for caching
func (c ListArgs) Key() (string, error) {
	//this is a value receiver, so we can delete a bunch of fields without impacting the original struct
	c.ChannelName = ""
	c.ChannelID = ""
	c.Signature = ""
	c.SigningTS = ""
	c.RequestorChannelName = ""
	c.RequestorChannelID = nil

	a, err := json.Marshal(c)
	if err != nil {
		return "", errors.Prefix("could not marshall args: ", err)
	}
	sha256 := sha256.New()
	_, err = sha256.Write(a)
	if err != nil {
		return "", errors.Prefix("could not hash json form of list args: ", err)
	}
	return hex.EncodeToString(sha256.Sum(nil)), nil
}

// AbandonArgs are the arguments passed to comment.Abandon RPC call. If creator args are passed
// the signing channel of the content of the comment is checked these args and signature
// verification happens against the creators public key for authorization.
type AbandonArgs struct {
	ModAuthorization
	CommentID string `json:"comment_id"`
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
	Page                 int           `json:"page"`
	PageSize             int           `json:"page_size"`
	TotalPages           int           `json:"total_pages"`
	TotalItems           int64         `json:"total_items"`
	TotalFilteredItems   int64         `json:"total_filtered_items"`
	Items                []CommentItem `json:"items,omitempty"`
	HasHiddenComments    bool          `json:"has_hidden_comments"`
	HasProtectedComments bool          `json:"has_protected_comments"`
}

// Validate validates the data in the list args
func (c ListArgs) Validate() api.StatusError {
	err := v.ValidateStruct(&c,
		v.Field(&c.ChannelID, validator.ClaimID),
		v.Field(&c.ChannelName),
	)
	if err != nil {
		return api.StatusError{Err: errors.Err(err), Status: http.StatusBadRequest}
	}
	if c.ClaimID == nil && c.AuthorClaimID == nil {
		return api.StatusError{Err: errors.Err("you must pass either claim_id or author_claim_id"), Status: http.StatusBadRequest}
	}
	return api.StatusError{}
}
