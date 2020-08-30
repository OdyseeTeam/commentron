package comments

type CreateArgs struct {
	CommentText string  `json:"comment"`
	ClaimID     string  `json:"claim_id"`
	ParentID    *string `json:"parent_id"`
	ChannelID   *string `json:"channel_id"`
	ChannelName *string `json:"channel_name"`
	Signature   *string `json:"signature"`
	SigningTS   *string `json:"signing_ts"`
}

type CreateResponse struct {
	*CommentItem
}
