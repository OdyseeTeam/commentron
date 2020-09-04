package comments

import (
	"net/http"

	"github.com/lbryio/commentron/commentapi"

	"github.com/lbryio/commentron/model"

	"github.com/lbryio/lbry.go/extras/api"
	"github.com/lbryio/lbry.go/v2/extras/errors"
)

func byID(_ *http.Request, args *commentapi.ByIDArgs) (commentapi.CommentItem, error) {
	comment, err := model.Comments(model.CommentWhere.CommentID.EQ(args.CommentID)).OneG()
	if err != nil {
		return commentapi.CommentItem{}, errors.Err(err)
	}
	if comment == nil {
		return commentapi.CommentItem{}, api.StatusError{Err: errors.Err("comment for id %s could not be found", args.CommentID), Status: http.StatusBadRequest}
	}
	return populateItem(comment, nil), nil
}
