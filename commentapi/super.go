package commentapi

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"github.com/lbryio/lbry.go/v2/extras/errors"
)

// SuperListArgs arguments for the comment.List rpc call
type SuperListArgs struct {
	ClaimID       *string `json:"claim_id"`
	AuthorClaimID *string `json:"author_claim_id"`
	ParentID      *string `json:"parent_id"`
	Page          int     `json:"page"`
	PageSize      int     `json:"page_size"`
	TopLevel      bool    `json:"top_level"`
	Hidden        bool    `json:"hidden"`
	// Satoshi amount to filter below >= x
	SuperChatsAmount int `json:"super_chat"`
}

// SuperListResponse response for the comment.List rpc call
type SuperListResponse struct {
	Page              int           `json:"page"`
	PageSize          int           `json:"page_size"`
	TotalPages        int           `json:"total_pages"`
	TotalItems        int64         `json:"total_items"`
	TotalAmount       float64       `json:"total_amount"`
	Items             []CommentItem `json:"items,omitempty"`
	HasHiddenComments bool          `json:"has_hidden_comments"`
}

// ApplyDefaults applies the default values for arguments passed that are different from normal defaults.
func (c *SuperListArgs) ApplyDefaults() {
	if c.Page == 0 {
		c.Page = 1
	}

	if c.PageSize == 0 {
		c.PageSize = 100
	}

	if c.SuperChatsAmount == 0 {
		c.SuperChatsAmount = 1
	}

	if c.PageSize > 100 {
		c.PageSize = 100
	}
}

//Key returns the hash of the list args struct for caching
func (c *SuperListArgs) Key() (string, error) {
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
