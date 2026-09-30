package moderation

import (
	"database/sql"
	"net/http"

	"github.com/OdyseeTeam/commentron/commentapi"
	"github.com/OdyseeTeam/commentron/db"
	"github.com/OdyseeTeam/commentron/helper"
	"github.com/OdyseeTeam/commentron/model"
	"github.com/OdyseeTeam/commentron/server/auth"

	"github.com/aarondl/null/v8"
	"github.com/lbryio/lbry.go/v2/extras/api"
	"github.com/lbryio/lbry.go/v2/extras/errors"
	"github.com/lbryio/lbry.go/v2/extras/util"
)

func unBlock(r *http.Request, args *commentapi.UnBlockArgs, reply *commentapi.UnBlockResponse) error {
	modChannel, creatorChannel, _, err := auth.ModAuthenticate(r, &args.ModAuthorization)
	if err != nil {
		return err
	}

	bannedChannel, err := helper.FindOrCreateChannel(args.UnBlockedChannelID, args.UnBlockedChannelName)
	if err != nil {
		return errors.Err(err)
	}

	entries, err := bannedChannel.BlockedChannelBlockedEntries(model.BlockedEntryWhere.CreatorChannelID.EQ(null.StringFrom(creatorChannel.ClaimID))).All(db.RO)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return errors.Err(err)
	}

	isMod, err := modChannel.ModChannelModerators().Exists(db.RO)
	if err != nil {
		return errors.Err(err)
	}

	if !isMod && args.GlobalUnBlock {
		return api.StatusError{Err: errors.Err("you must be a global moderator to take global action"), Status: http.StatusBadRequest}
	}

	if args.GlobalUnBlock {
		entries, err := bannedChannel.BlockedChannelBlockedEntries(model.BlockedEntryWhere.UniversallyBlocked.EQ(null.BoolFrom(true))).All(db.RO)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return errors.Err(err)
		}
		err = entries.DeleteAll(db.RW)
		if err != nil {
			return errors.Err(err)
		}
	} else {
		if len(entries) > 0 {
			for _, be := range entries {
				if be.CreatorChannelID.String == creatorChannel.ClaimID {
					err := be.Delete(db.RW)
					if err != nil {
						return errors.Err(err)
					}
					reply.UnBlockedFrom = util.PtrToString(creatorChannel.ClaimID)
				}
			}
		}
	}

	reply.GlobalUnBlock = args.GlobalUnBlock
	reply.UnBlockedChannelID = bannedChannel.ClaimID

	return nil
}
