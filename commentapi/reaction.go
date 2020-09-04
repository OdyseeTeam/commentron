package commentapi

// ReactArgs are the arguments passed to comment.Abandon RPC call
type ReactArgs struct {
	CommentIDs  string  `json:"comment_ids"`
	Signature   string  `json:"signature"`
	SigningTS   string  `json:"signing_ts"`
	Remove      bool    `json:"remove"`
	ClearTypes  string  `json:"clear_types"`
	Type        string  `json:"type"`
	ChannelID   *string `json:"channel_id"`
	ChannelName *string `json:"channel_name"`
}

// ReactResponse the response to the abandon call
type ReactResponse struct {
	Reactions
}

// ListArgs are the arguments passed to comment.Abandon RPC call
type ReactionListArgs struct {
	CommentIDs  string `json:"comment_ids"`
	Signature   string `json:"signature"`
	SigningTS   string `json:"signing_ts"`
	Types       *string
	ChannelID   *string `json:"channel_id"`
	ChannelName *string `json:"channel_name"`
}

// ListResponse the response to the abandon call
type ReactionListResponse struct {
	MyReactions     Reactions `json:"my_reactions,omitempty"`
	OthersReactions Reactions `json:"others_reactions"`
}

type Reactions map[string]CommentReaction

type CommentReaction map[string]int
