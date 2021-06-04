package moderation

import (
	"database/sql"
	"net/http"

	"github.com/lbryio/commentron/commentapi"
	"github.com/lbryio/commentron/helper"
	"github.com/lbryio/commentron/model"
	"github.com/lbryio/commentron/server/lbry"
	"github.com/lbryio/commentron/validator"

	"github.com/lbryio/lbry.go/extras/api"
	"github.com/lbryio/lbry.go/v2/extras/errors"
	v "github.com/lbryio/ozzo-validation"

	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

func block(_ *http.Request, args *commentapi.BlockArgs, reply *commentapi.BlockResponse) error {
	err := v.ValidateStruct(args,
		v.Field(&args.BlockedChannelID, validator.ClaimID, v.Required),
		v.Field(&args.BlockedChannelName, v.Required),
		v.Field(&args.ModChannelID, validator.ClaimID, v.Required),
		v.Field(&args.ModChannelName, v.Required),
	)
	if err != nil {
		return api.StatusError{Err: errors.Err(err), Status: http.StatusBadRequest}
	}
	modChannel, creatorChannel, err := getModerator(args.ModChannelID, args.ModChannelName, args.CreatorChannelID, args.CreatorChannelName)
	if err != nil {
		return err
	}
	err = lbry.ValidateSignature(modChannel.ClaimID, args.Signature, args.SigningTS, args.ModChannelName)
	if err != nil {
		return err
	}

	bannedChannel, err := helper.FindOrCreateChannel(args.BlockedChannelID, args.BlockedChannelName)
	if err != nil {
		return errors.Err(err)
	}
	blockedEntry, err := model.BlockedEntries(
		model.BlockedEntryWhere.BlockedChannelID.EQ(null.StringFrom(args.BlockedChannelID)),
		model.BlockedEntryWhere.BlockedByChannelID.EQ(null.StringFrom(creatorChannel.ClaimID))).OneG()
	if err != nil && err != sql.ErrNoRows {
		return errors.Err(err)
	}
	if blockedEntry == nil {
		blockedEntry = &model.BlockedEntry{
			BlockedChannelID:   null.StringFrom(bannedChannel.ClaimID),
			BlockedByChannelID: null.StringFrom(creatorChannel.ClaimID),
		}
		err := blockedEntry.InsertG(boil.Infer())
		if err != nil {
			return errors.Err(err)
		}
	}
	isMod, err := modChannel.ModChannelModerators().ExistsG()
	if err != nil {
		return errors.Err(err)
	}
	if args.BlockAll {
		if !isMod {
			return api.StatusError{Err: errors.Err("cannot block universally without admin privileges"), Status: http.StatusForbidden}
		}
		blockedEntry.UniversallyBlocked.SetValid(true)
		reply.AllBlocked = true
	} else {
		reply.BannedFrom = &creatorChannel.ClaimID
	}

	err = blockedEntry.UpdateG(boil.Infer())
	if err != nil {
		return errors.Err(err)
	}
	var deletedCommentIDs []string
	if args.DeleteAll {
		if !isMod {
			return api.StatusError{Err: errors.Err("cannot delete all comments of user without admin priviledges"), Status: http.StatusForbidden}
		}

		comments, err := model.Comments(model.CommentWhere.ChannelID.EQ(null.StringFrom(bannedChannel.ClaimID))).AllG()
		if err != nil {
			return errors.Err(err)
		}
		err = comments.DeleteAllG()
		if err != nil {
			return errors.Err(err)
		}
		for _, c := range comments {
			deletedCommentIDs = append(deletedCommentIDs, c.CommentID)
		}
		reply.DeletedCommentIDs = deletedCommentIDs
	}

	reply.BannedChannelID = bannedChannel.ClaimID

	return nil
}

func getModerator(modChannelID, modChannelName, creatorChannelID, creatorChannelName string) (*model.Channel, *model.Channel, error) {
	modChannel, err := helper.FindOrCreateChannel(modChannelID, modChannelName)
	if err != nil {
		return nil, nil, errors.Err(err)
	}
	var creatorChannel = modChannel
	if creatorChannelID != "" && creatorChannelName != "" {
		creatorChannel, err = helper.FindOrCreateChannel(creatorChannelID, creatorChannelName)
		if err != nil {
			return nil, nil, errors.Err(err)
		}
		dmRels := model.DelegatedModeratorRels
		dmWhere := model.DelegatedModeratorWhere
		loadCreatorChannels := qm.Load(dmRels.CreatorChannel, dmWhere.CreatorChannelID.EQ(creatorChannelID))
		exists, err := modChannel.ModChannelDelegatedModerators(loadCreatorChannels).ExistsG()
		if err != nil {
			return nil, nil, errors.Err(err)
		}
		if !exists {
			return nil, nil, errors.Err("%s is not delegated by %s to be a moderator", modChannel.Name, creatorChannel.Name)
		}
	}
	return modChannel, creatorChannel, nil
}

func blockedList(_ *http.Request, args *commentapi.BlockedListArgs, reply *commentapi.BlockedListResponse) error {
	modChannel, err := helper.FindOrCreateChannel(args.ModChannelID, args.ModChannelName)
	if err != nil {
		return errors.Err(err)
	}
	err = lbry.ValidateSignature(modChannel.ClaimID, args.Signature, args.SigningTS, args.ModChannelName)
	if err != nil {
		return err
	}

	blockedByMod, err := modChannel.BlockedByChannelBlockedEntries(qm.Load(model.BlockedEntryRels.BlockedChannel)).AllG()
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return errors.Err(err)
	}

	for _, b := range blockedByMod {
		if b.R != nil && b.R.BlockedChannel != nil {
			reply.BlockedChannels = append(reply.BlockedChannels, commentapi.BlockedChannel{
				BlockedChannelID:     b.R.BlockedChannel.ClaimID,
				BlockedChannelName:   b.R.BlockedChannel.Name,
				BlockedByChannelID:   modChannel.ClaimID,
				BlockedByChannelName: modChannel.Name,
				BlockedAt:            b.CreatedAt,
			})
		}
	}
	return nil
}
