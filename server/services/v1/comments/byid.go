package comments

import (
	"net/http"

	"github.com/lbryio/commentron/model"

	"github.com/lbryio/lbry.go/extras/api"
	"github.com/lbryio/lbry.go/v2/extras/errors"
)

// ByIDArgs arguments for the comment.List rpc call
type ByIDArgs struct {
	CommentID string `json:"comment_id"`
}

// ByIDResponse response for the comment.ByID rpc call
type ByIDResponse struct {
	Item CommentItem `json:"items,omitempty"`
}

func byID(_ *http.Request, args *ByIDArgs) (CommentItem, error) {
	comment, err := model.Comments(model.CommentWhere.CommentID.EQ(args.CommentID)).OneG()
	if err != nil {
		return CommentItem{}, errors.Err(err)
	}
	if comment == nil {
		return CommentItem{}, api.StatusError{Err: errors.Err("comment for id %s could not be found", args.CommentID), Status: http.StatusBadRequest}
	}
	return populateItem(comment, nil), nil
}
