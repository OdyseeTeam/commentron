package moderation

import (
	"net/http"
	"strings"
	"time"

	"github.com/OdyseeTeam/commentron/commentapi"
	"github.com/OdyseeTeam/commentron/db"
	"github.com/OdyseeTeam/commentron/model"
	"github.com/OdyseeTeam/commentron/server/auth"

	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/lbryio/lbry.go/v2/extras/api"
	"github.com/lbryio/lbry.go/v2/extras/errors"
)

const reactionDeleteBatchSize = 250
const reactionDeleteMaxAttempts = 4
const reactionDeleteInitialRetryDelay = 250 * time.Millisecond

func isRetryableReactionDeleteError(err error) bool {
	lower := strings.ToLower(err.Error())
	return strings.Contains(lower, "lock wait timeout") ||
		strings.Contains(lower, "deadlock") ||
		strings.Contains(lower, "error 1205") ||
		strings.Contains(lower, "error 1213")
}

func deleteReactionBatch(ids []uint64) error {
	if len(ids) == 0 {
		return nil
	}

	args := make([]interface{}, 0, len(ids))
	for _, id := range ids {
		args = append(args, id)
	}

	delay := reactionDeleteInitialRetryDelay
	var err error
	for attempt := 0; attempt < reactionDeleteMaxAttempts; attempt++ {
		err = model.Reactions(
			qm.WhereIn(model.ReactionColumns.ID+" IN ?", args...),
		).DeleteAll(db.RW)
		if err == nil {
			return nil
		}
		if !isRetryableReactionDeleteError(err) || attempt == reactionDeleteMaxAttempts-1 {
			return err
		}
		time.Sleep(delay)
		delay *= 2
	}

	return err
}

func wipeReactions(r *http.Request, args *commentapi.WipeReactionsArgs, reply *commentapi.WipeReactionsResponse) error {
	modChannel, _, _, err := auth.ModAuthenticate(r, &args.ModAuthorization)
	if err != nil {
		return err
	}

	isMod, err := modChannel.ModChannelModerators().Exists(db.RO)
	if err != nil {
		return errors.Err(err)
	}
	if !isMod {
		return api.StatusError{Err: errors.Err("cannot wipe reactions without admin privileges"), Status: http.StatusForbidden}
	}

	reactions, err := model.Reactions(
		model.ReactionWhere.ChannelID.EQ(null.StringFrom(args.TargetChannelID)),
		qm.OrderBy(model.ReactionColumns.ID+" ASC"),
	).All(db.RO)
	if err != nil {
		return errors.Err(err)
	}

	reply.TargetChannelID = args.TargetChannelID
	reply.DeletedReactionCount = uint64(len(reactions))
	if len(reactions) == 0 {
		return nil
	}

	reactionIDs := make([]uint64, 0, len(reactions))
	for _, reaction := range reactions {
		reactionIDs = append(reactionIDs, reaction.ID)
	}

	for start := 0; start < len(reactionIDs); start += reactionDeleteBatchSize {
		end := start + reactionDeleteBatchSize
		if end > len(reactionIDs) {
			end = len(reactionIDs)
		}
		err = deleteReactionBatch(reactionIDs[start:end])
		if err != nil {
			return errors.Err(err)
		}
	}

	return nil
}
