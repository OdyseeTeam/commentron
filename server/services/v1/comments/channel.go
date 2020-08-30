package comments

type ChannelArgs struct {
	CommentID string `json:"comment_id"`
}

type ChannelResponse struct {
	ChannelID   string `json:"channel_id"`
	ChannelName string `json:"channel_name"`
}
