package comments

// ChannelArgs arguments to the comment.GetChannelForCommentID call
type ChannelArgs struct {
	CommentID string `json:"comment_id"`
}

// ChannelResponse response to the comment.GetChannelForCommentID call
type ChannelResponse struct {
	ChannelID   string `json:"channel_id"`
	ChannelName string `json:"channel_name"`
}
