package comments

import (
	"database/sql"
	"net/http"

	"github.com/lbryio/commentron/commentapi"
	"github.com/lbryio/commentron/db"
	"github.com/lbryio/commentron/model"

	"github.com/lbryio/lbry.go/extras/api"
	"github.com/lbryio/lbry.go/v2/extras/errors"

	"github.com/volatiletech/sqlboiler/queries/qm"
)

func byID(_ *http.Request, args *commentapi.ByIDArgs) (commentapi.CommentItem, []commentapi.CommentItem, error) {
	comment, err := model.Comments(model.CommentWhere.CommentID.EQ(args.CommentID), qm.Load(model.CommentRels.Channel)).One(db.RO)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return commentapi.CommentItem{}, nil, errors.Err(err)
	}
	if comment == nil {
		return commentapi.CommentItem{}, nil, api.StatusError{Err: errors.Err("comment for id %s could not be found", args.CommentID), Status: http.StatusBadRequest}
	}
	var channel *model.Channel
	if comment.R != nil && comment.R.Channel != nil {
		channel = comment.R.Channel
	}
	replies, err := comment.ParentComments().Count(db.RO)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return commentapi.CommentItem{}, nil, errors.Err(err)
	}
	var ancestors []commentapi.CommentItem
	if args.WithAncestors {
		lastcomment := comment
		for !lastcomment.ParentID.IsZero() {
			parentComment, err := lastcomment.Parent(qm.Load(model.CommentRels.Channel)).One(db.RO)
			if err != nil {
				return commentapi.CommentItem{}, nil, errors.Err(err)
			}
			var parentChannel *model.Channel
			if parentComment.R != nil && parentComment.R.Channel != nil {
				parentChannel = parentComment.R.Channel
			}
			parentReplies, err := parentComment.ParentComments().Count(db.RO)
			if err != nil && errors.Is(err, sql.ErrNoRows) {
				return commentapi.CommentItem{}, nil, errors.Err(err)
			}
			ancestors = append(ancestors, populateItem(parentComment, parentChannel, int(parentReplies)))
			lastcomment = parentComment
		}
	}

	return populateItem(comment, channel, int(replies)), ancestors, nil
}
