package moderation

import (
	"net/http"

	"github.com/OdyseeTeam/commentron/commentapi"
	"github.com/OdyseeTeam/commentron/db"
	"github.com/OdyseeTeam/commentron/jobs/commentclassification"
	"github.com/OdyseeTeam/commentron/model"

	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/lbryio/lbry.go/v2/extras/errors"
)

func actOnClassification(r *http.Request, args *commentapi.ActOnClassificationArgs, reply *commentapi.ActOnClassificationResponse) error {
	if err := commentclassification.IsAuthenticated(r); err != nil {
		return err
	}

	// Lookup the comment classification by claim id.
	// There is a foreign key constraint on the comment id, so this will fail if the
	// comment doesn't exist assuming soft deletes.
	commentClassification, err := model.CommentClassifications(
		model.CommentClassificationWhere.CommentID.EQ(args.CommentID),
		qm.Load(model.CommentClassificationRels.Comment),
	).One(db.RO)
	if err != nil {
		return errors.Err(err)
	}

	err = db.WithTx(db.RW, nil, func(tx boil.Transactor) error {
		// Delete the comment.
		if args.DoDelete {
			err = commentClassification.R.Comment.Delete(tx, false)
			if err != nil {
				return err
			}
		}

		// Update the comment classification.
		commentClassification.IsReviewed = null.BoolFrom(true)
		commentClassification.ReviewerApproved = null.BoolFrom(args.Confirm)

		return commentClassification.Update(tx, boil.Infer())
	})

	if err != nil {
		return errors.Err(err)
	}

	reply.Status = "ok"

	return nil
}
