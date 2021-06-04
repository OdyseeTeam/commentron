package moderation

import (
	"database/sql"
	"net/http"

	"github.com/volatiletech/sqlboiler/queries/qm"

	"github.com/lbryio/commentron/helper"
	"github.com/lbryio/commentron/model"
	"github.com/lbryio/commentron/server/lbry"
	"github.com/lbryio/lbry.go/v2/extras/errors"

	"github.com/lbryio/commentron/commentapi"
)

type delegatedModLevel int

const defaultLevel = delegatedModLevel(0)

func addDelegate(r *http.Request, args *commentapi.AddDelegateArgs, reply *commentapi.ListDelegateResponse) error {
	creatorChannel, err := helper.FindOrCreateChannel(args.CreatorChannelID, args.CreatorChannelName)
	if err != nil {
		return errors.Err(err)
	}

	err = lbry.ValidateSignature(creatorChannel.ClaimID, args.Signature, args.SigningTS, args.CreatorChannelName)
	if err != nil {
		return err
	}

	modChannel, err := helper.FindOrCreateChannel(args.ModChannelID, args.ModChannelName)
	if err != nil {
		return errors.Err(err)
	}
	exists, err := creatorChannel.CreatorChannelDelegatedModerators(
		model.DelegatedModeratorWhere.ModChannelID.EQ(modChannel.ClaimID),
		model.DelegatedModeratorWhere.CreatorChannelID.EQ(creatorChannel.ClaimID)).ExistsG()
	if err != nil {
		return errors.Err(err)
	}
	if exists {
		return errors.Err("channel %s already is a moderation for %s", args.ModChannelName, args.CreatorChannelName)
	}

	delegatedModerator := &model.DelegatedModerator{
		ModChannelID: modChannel.ClaimID,
		Permissons:   uint64(defaultLevel),
	}

	err = creatorChannel.AddCreatorChannelDelegatedModeratorsG(true, delegatedModerator)
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
	creatorChannel, err := helper.FindOrCreateChannel(args.CreatorChannelID, args.CreatorChannelName)
	if err != nil {
		return errors.Err(err)
	}

	err = lbry.ValidateSignature(creatorChannel.ClaimID, args.Signature, args.SigningTS, args.CreatorChannelName)
	if err != nil {
		return err
	}

	modChannel, err := helper.FindOrCreateChannel(args.ModChannelID, args.ModChannelName)
	if err != nil {
		return errors.Err(err)
	}

	modEntry, err := creatorChannel.CreatorChannelDelegatedModerators(model.DelegatedModeratorWhere.ModChannelID.EQ(modChannel.ClaimID)).OneG()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return errors.Err(err)
	}

	if modEntry == nil {
		return errors.Err("Mod channel %s is not a moderator for channel %s", args.ModChannelName, args.CreatorChannelName)
	}

	err = modEntry.DeleteG()
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
	creatorChannel, err := helper.FindOrCreateChannel(args.CreatorChannelID, args.CreatorChannelName)
	if err != nil {
		return errors.Err(err)
	}

	err = lbry.ValidateSignature(creatorChannel.ClaimID, args.Signature, args.SigningTS, args.CreatorChannelName)
	if err != nil {
		return err
	}

	delegatedModEntries, err := creatorChannel.CreatorChannelDelegatedModerators(qm.Load(model.DelegatedModeratorRels.ModChannel)).AllG()
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
