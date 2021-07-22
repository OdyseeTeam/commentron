package moderation

import (
	"database/sql"
	"net/http"

	"github.com/lbryio/commentron/commentapi"
	"github.com/lbryio/commentron/helper"
	"github.com/lbryio/commentron/model"
	"github.com/lbryio/commentron/server/lbry"

	"github.com/lbryio/lbry.go/extras/api"
	"github.com/lbryio/lbry.go/extras/util"
	"github.com/lbryio/lbry.go/v2/extras/errors"

	"github.com/volatiletech/null"
)

func unBlock(_ *http.Request, args *commentapi.UnBlockArgs, reply *commentapi.UnBlockResponse) error {

	modChannel, creatorChannel, err := getModerator(args.ModChannelID, args.ModChannelName, args.CreatorChannelID, args.CreatorChannelName)
	if err != nil {
		return err
	}
	err = lbry.ValidateSignature(modChannel.ClaimID, args.Signature, args.SigningTS, args.ModChannelName)
	if err != nil {
		return err
	}

	bannedChannel, err := helper.FindOrCreateChannel(args.UnBlockedChannelID, args.UnBlockedChannelName)
	if err != nil {
		return errors.Err(err)
	}

	entries, err := bannedChannel.BlockedChannelBlockedEntries(model.BlockedEntryWhere.CreatorChannelID.EQ(null.StringFrom(creatorChannel.ClaimID))).AllG()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return errors.Err(err)
	}

	isMod, err := modChannel.ModChannelModerators().ExistsG()
	if err != nil {
		return errors.Err(err)
	}

	if !isMod && args.GlobalUnBlock {
		return api.StatusError{Err: errors.Err("you must be a global moderator to take global action"), Status: http.StatusBadRequest}

	}

	if args.GlobalUnBlock {
		entries, err := bannedChannel.BlockedChannelBlockedEntries(model.BlockedEntryWhere.UniversallyBlocked.EQ(null.BoolFrom(true))).AllG()
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return errors.Err(err)
		}
		err = entries.DeleteAllG()
		if err != nil {
			return errors.Err(err)
		}
	} else {
		if len(entries) > 0 {
			for _, be := range entries {
				if be.CreatorChannelID.String == creatorChannel.ClaimID {
					err := be.DeleteG()
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
