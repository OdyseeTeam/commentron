package moderation

import (
	"database/sql"
	"net/http"

	"github.com/OdyseeTeam/commentron/commentapi"
	"github.com/OdyseeTeam/commentron/db"
	"github.com/OdyseeTeam/commentron/helper"
	"github.com/OdyseeTeam/commentron/model"
	"github.com/OdyseeTeam/commentron/server/auth"

	"github.com/lbryio/lbry.go/v2/extras/api"
	"github.com/lbryio/lbry.go/v2/extras/errors"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type delegatedModLevel int

const defaultLevel = delegatedModLevel(0)

func addDelegate(r *http.Request, args *commentapi.AddDelegateArgs, reply *commentapi.ListDelegateResponse) error {
	creatorChannel, _, err := auth.Authenticate(r, &args.Authorization)
	if err != nil {
		return err
	}

	modChannel, err := helper.FindOrCreateChannel(args.ModChannelID, args.ModChannelName)
	if modChannel != nil && creatorChannel != nil && modChannel.ClaimID == creatorChannel.ClaimID {
		return api.StatusError{Err: errors.Err("you are the creator, one cannot simply delegate to themselves"), Status: http.StatusBadRequest}
	}
	if err != nil {
		return errors.Err(err)
	}
	exists, err := creatorChannel.CreatorChannelDelegatedModerators(
		model.DelegatedModeratorWhere.ModChannelID.EQ(modChannel.ClaimID),
		model.DelegatedModeratorWhere.CreatorChannelID.EQ(creatorChannel.ClaimID)).Exists(db.RO)
	if err != nil {
		return errors.Err(err)
	}
	if exists {
		return errors.Err("channel %s already is a moderator for %s", args.ModChannelName, args.CreatorChannelName)
	}

	delegatedModerator := &model.DelegatedModerator{
		ModChannelID: modChannel.ClaimID,
		Permissons:   uint64(defaultLevel),
	}

	err = creatorChannel.AddCreatorChannelDelegatedModerators(db.RW, true, delegatedModerator)
	if err != nil {
		return errors.Err(err)
	}

	reply.Delegates = append(reply.Delegates, commentapi.Delegate{
		ChannelID:   modChannel.ClaimID,
		ChannelName: modChannel.Name,
	})

	return nil
}

func removeDelegate(r *http.Request, args *commentapi.RemoveDelegateArgs, reply *commentapi.ListDelegateResponse) error {
	creatorChannel, _, err := auth.Authenticate(r, &args.Authorization)
	if err != nil {
		return err
	}

	modChannel, err := helper.FindOrCreateChannel(args.ModChannelID, args.ModChannelName)
	modEntry, err := creatorChannel.CreatorChannelDelegatedModerators(model.DelegatedModeratorWhere.ModChannelID.EQ(modChannel.ClaimID)).One(db.RO)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return errors.Err(err)
	}

	if modEntry == nil {
		return errors.Err("Mod channel %s is not a moderator for channel %s", args.ModChannelName, args.CreatorChannelName)
	}

	err = modEntry.Delete(db.RW)
	if err != nil {
		return errors.Err(err)
	}

	reply.Delegates = append(reply.Delegates, commentapi.Delegate{
		ChannelID:   modChannel.ClaimID,
		ChannelName: modChannel.Name,
	})

	return nil
}

func listDelegates(r *http.Request, args *commentapi.ListDelegatesArgs, reply *commentapi.ListDelegateResponse) error {
	creatorChannel, _, err := auth.Authenticate(r, &args.Authorization)
	if err != nil {
		return err
	}

	delegatedModEntries, err := creatorChannel.CreatorChannelDelegatedModerators(qm.Load(model.DelegatedModeratorRels.ModChannel)).All(db.RO)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return errors.Err(err)
	}

	for _, m := range delegatedModEntries {
		reply.Delegates = append(reply.Delegates, commentapi.Delegate{
			ChannelID:   m.R.ModChannel.ClaimID,
			ChannelName: m.R.ModChannel.Name,
		})
	}

	return nil
}
