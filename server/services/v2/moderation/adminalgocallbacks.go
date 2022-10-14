package moderation

import (
	"net/http"
	"time"

	"github.com/OdyseeTeam/commentron/commentapi"
	"github.com/OdyseeTeam/commentron/db"
	"github.com/OdyseeTeam/commentron/jobs/commentclassification"
	"github.com/OdyseeTeam/commentron/model"
	"github.com/lbryio/lbry.go/v2/extras/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

// adminAlgoCallbacks adds or deletes rows from the ChannelAlgoCallbacks table
//
// Ideally this would be restful with a PUT and DELETE, but this is such a narrow
// API surface area that I'd prefer to keep it in one place.
func adminAlgoCallbacks(r *http.Request, args *commentapi.AdminAlgoCallbacksArgs, reply *commentapi.AdminAlgoCallbacksResponse) error {
	if err := commentclassification.IsAuthenticated(r); err != nil {
		return err
	}

	// Verify that the channel exists.
	exists, err := model.ChannelExists(db.RO, args.ChannelID)
	if err != nil {
		return errors.Err(err)
	} else if !exists {
		return errors.Err("channel `%s` does not exist in commentron", args.ChannelID)
	}

	if args.Add {
		err := db.WithTx(db.RW, nil, func(tx boil.Transactor) error {
			// Check if it already exists.
			exists, err := model.ChannelAlgoCallbacks(
				model.ChannelAlgoCallbackWhere.ChannelID.EQ(args.ChannelID),
				model.ChannelAlgoCallbackWhere.WatcherID.EQ(args.WatcherID),
			).Exists(tx)

			if err != nil {
				return err
			} else if exists {
				return nil
			}

			// Insert the channel algo callback.
			callback := &model.ChannelAlgoCallback{
				ChannelID: args.ChannelID,
				WatcherID: args.WatcherID,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			if err = callback.Insert(tx, boil.Infer()); err != nil {
				return errors.Err(err)
			}
			return nil
		})
		if err != nil {
			return errors.Err(err)
		}
		reply.Status = "added"
	} else {
		// Delete the channel algo callback.
		err := db.WithTx(db.RW, nil, func(tx boil.Transactor) error {
			callback, err := model.ChannelAlgoCallbacks(
				model.ChannelAlgoCallbackWhere.ChannelID.EQ(args.ChannelID),
				model.ChannelAlgoCallbackWhere.WatcherID.EQ(args.WatcherID),
			).One(db.RO)

			if err != nil {
				return err
			}

			return callback.Delete(tx)
		})
		if err != nil {
			return errors.Err(err)
		}
		reply.Status = "deleted"
	}

	return nil
}
