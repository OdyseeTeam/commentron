package comments

type EditArgs struct {
	CommentID string `json:"comment_id"`
}

type EditResponse struct {
	*CommentItem
}
