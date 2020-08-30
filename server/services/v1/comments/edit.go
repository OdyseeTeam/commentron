package comments

// EditArgs arguments for the comment.Edit rpc call
type EditArgs struct {
	CommentID string `json:"comment_id"`
}

// EditResponse response for the comment.Edit rpc call
type EditResponse struct {
	*CommentItem
}
