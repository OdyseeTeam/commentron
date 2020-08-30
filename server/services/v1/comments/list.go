package comments

// ListArgs arguments for the comment.List rpc call
type ListArgs struct {
	ClaimID  *string `json:"claim_id"`
	ParentID *string `json:"parent_id"`
	Page     int     `json:"page"`
	PageSize int     `json:"page_size"`
	TopLevel bool    `json:"top_level"`
}

// ApplyDefaults applies the default values for arguments passed that are different from normal defaults.
func (c *ListArgs) ApplyDefaults() {
	if c.Page == 0 {
		c.Page = 1
	}

	if c.PageSize == 0 {
		c.PageSize = 50
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
