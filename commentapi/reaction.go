package commentapi

import (
	"net/http"
	"strings"

	"github.com/lbryio/lbry.go/v2/extras/errors"

	"github.com/lbryio/lbry.go/v2/extras/api"
)

// ReactArgs are the arguments passed to comment.Abandon RPC call
type ReactArgs struct {
	Authorization

	CommentIDs string `json:"comment_ids"`
	Remove     bool   `json:"remove"`
	ClearTypes string `json:"clear_types"`
	Type       string `json:"type"`
}

// ReactResponse the response to the abandon call
type ReactResponse struct {
	Reactions
}

// ReactionListArgs are the arguments passed to comment.Abandon RPC call
type ReactionListArgs struct {
	Authorization
	CommentIDs string `json:"comment_ids"`
	Types      *string
}

// Validate validates the data in the list args
func (rl ReactionListArgs) Validate() api.StatusError {
	if len(strings.Split(rl.CommentIDs, ",")) > 50 {
		return api.StatusError{Err: errors.Err("too many comment ids passed"), Status: http.StatusBadRequest}
	}
	return api.StatusError{}
}

// ReactionListResponse the response to the abandon call
type ReactionListResponse struct {
	MyReactions     Reactions `json:"my_reactions,omitempty"`
	OthersReactions Reactions `json:"others_reactions"`
}

// Reactions a map structure where the key is the comment_id and the value is a CommentReaction
type Reactions map[string]CommentReaction

// CommentReaction is a map for representing the reaction and its quantity for a comment
type CommentReaction map[string]int
