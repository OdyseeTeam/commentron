package moderation

import (
	"database/sql"
	"net/http"

	"github.com/lbryio/commentron/commentapi"
	"github.com/lbryio/commentron/helper"
	"github.com/lbryio/commentron/model"
	"github.com/lbryio/commentron/server/lbry"

	"github.com/lbryio/lbry.go/extras/api"
	"github.com/lbryio/lbry.go/v2/extras/errors"

	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

func block(_ *http.Request, args *commentapi.BlockArgs, reply *commentapi.BlockResponse) error {
	modChannel, err := helper.FindOrCreateChannel(args.ModChannelID, args.ModChannelName)
	if err != nil {
		return errors.Err(err)
	}
	err = lbry.ValidateSignature(modChannel.ClaimID, args.Signature, args.SigningTS, args.ModChannelName)
	if err != nil {
		return err
	}

	bannedChannel, err := helper.FindOrCreateChannel(args.BannedChannelID, args.BannedChannelName)
	if err != nil {
		return errors.Err(err)
	}
	blockedEntry, err := model.BlockedEntries(model.BlockedEntryWhere.BlockedChannelID.EQ(null.StringFrom(args.BannedChannelID))).OneG()
	if err != nil && err != sql.ErrNoRows {
		return errors.Err(err)
	}
	if blockedEntry == nil {
		blockedEntry = &model.BlockedEntry{
			BlockedChannelID:   null.StringFrom(bannedChannel.ClaimID),
			BlockedByChannelID: null.StringFrom(modChannel.ClaimID),
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
		reply.BannedFrom = &modChannel.ClaimID
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
			})
		}
	}
	return nil
}
