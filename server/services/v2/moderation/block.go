package moderation

import (
	"database/sql"
	"encoding/hex"
	"net/http"

	"github.com/lbryio/commentron/commentapi"
	"github.com/lbryio/commentron/model"
	"github.com/lbryio/commentron/server/lbry"
	"github.com/lbryio/commentron/util"

	"github.com/lbryio/lbry.go/extras/api"
	"github.com/lbryio/lbry.go/v2/extras/errors"

	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
)

func block(r *http.Request, args *commentapi.BlockArgs, reply *commentapi.BlockResponse) error {
	modChannel, err := util.FindOrCreateChannel(args.ModChannelID, args.ModChannelName)
	if err != nil {
		return errors.Err(err)
	}
	err = lbry.ValidateSignature(modChannel.ClaimID, args.Signature, args.SigningTS, hex.EncodeToString([]byte(args.ModChannelName)))
	if err != nil {
		return err
	}

	bannedChannel, err := util.FindOrCreateChannel(args.BannedChannelID, args.BannedChannelName)
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

	err = blockedEntry.UpsertG(boil.Infer(), boil.Infer())
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
