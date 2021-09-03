package moderation

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/lbryio/commentron/commentapi"
	"github.com/lbryio/commentron/db"
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

	// Only get the block list they were invited to.
	participatingBlockedList, err := creatorChannel.BlockedListInvite().One(db.RO)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return errors.Err(err)
	}

	bannedChannel, err := helper.FindOrCreateChannel(args.BlockedChannelID, args.BlockedChannelName)
	if err != nil {
		return errors.Err(err)
	}
	blockedEntry, err := model.BlockedEntries(
		model.BlockedEntryWhere.BlockedChannelID.EQ(null.StringFrom(args.BlockedChannelID)),
		model.BlockedEntryWhere.CreatorChannelID.EQ(null.StringFrom(creatorChannel.ClaimID))).One(db.RO)
	if err != nil && err != sql.ErrNoRows {
		return errors.Err(err)
	}
	strikes := 0
	insert := false
	if blockedEntry == nil {
		blocklistID := null.Uint64{}
		if participatingBlockedList != nil {
			blocklistID.SetValid(participatingBlockedList.ID)
		}
		blockedEntry = &model.BlockedEntry{
			BlockedChannelID: null.StringFrom(bannedChannel.ClaimID),
			CreatorChannelID: null.StringFrom(creatorChannel.ClaimID),
			BlockedListID:    blocklistID,
		}
		insert = true
	} else {
		blockedEntry.Strikes.SetValid(blockedEntry.Strikes.Int + 1)
		strikes = blockedEntry.Strikes.Int
	}
	if participatingBlockedList != nil && args.TimeOut > 0 {
		return api.StatusError{Err: errors.Err("the block list rules you are participating have their time out hours settings per strike. You must stop participating in the shared blocked list to customize timeouts"), Status: http.StatusBadRequest}
	} else if participatingBlockedList != nil {
		blockedEntry.Expiry.SetValid(time.Now().Add(getStrikeDuration(strikes, participatingBlockedList)))
	} else if args.TimeOut > 0 {
		blockedEntry.Expiry.SetValid(time.Now().Add(time.Duration(args.TimeOut) * time.Second))
	} else if participatingBlockedList == nil { // Only reset expiry if not participating in shared blockedlist, this should never exist from check above!
		blockedEntry.Expiry.Valid = false
		blockedEntry.Expiry.Time = time.Time{}
	}

	isMod, err := modChannel.ModChannelModerators().Exists(db.RO)
	if err != nil {
		return errors.Err(err)
	}
	if args.BlockAll {
		if !isMod {
			return api.StatusError{Err: errors.Err("cannot block universally without admin privileges"), Status: http.StatusForbidden}
		}
		blockedEntry.CreatorChannelID.SetValid(creatorChannel.ClaimID)
		blockedEntry.UniversallyBlocked.SetValid(true)
		reply.AllBlocked = true
	} else {
		reply.BannedFrom = &creatorChannel.ClaimID
	}

	if modChannel.ClaimID != creatorChannel.ClaimID {
		blockedEntry.DelegatedModeratorChannelID = null.StringFrom(modChannel.ClaimID)
	}

	if insert {
		err := blockedEntry.Insert(db.RW, boil.Infer())
		if err != nil {
			return errors.Err(err)
		}
	} else {
		err = blockedEntry.Update(db.RW, boil.Infer())
		if err != nil {
			return errors.Err(err)
		}
	}

	var deletedCommentIDs []string
	if args.DeleteAll {
		if !isMod {
			return api.StatusError{Err: errors.Err("cannot delete all comments of user without admin priviledges"), Status: http.StatusForbidden}
		}

		comments, err := model.Comments(model.CommentWhere.ChannelID.EQ(null.StringFrom(bannedChannel.ClaimID))).All(db.RO)
		if err != nil {
			return errors.Err(err)
		}
		err = comments.DeleteAll(db.RW)
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

const defaultStrikeTimeout = 4 * time.Hour

func getStrikeDuration(strike int, list *model.BlockedList) time.Duration {
	if list.StrikeThree.Valid && strike == 3 {
		return time.Duration(list.StrikeThree.Uint64) * time.Hour
	} else if list.StrikeTwo.Valid && strike == 2 {
		return time.Duration(list.StrikeTwo.Uint64) * time.Hour
	} else if list.StrikeOne.Valid && strike == 1 {
		return time.Duration(list.StrikeOne.Uint64) * time.Hour
	} else {
		return defaultStrikeTimeout
	}
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
		exists, err := modChannel.ModChannelDelegatedModerators(loadCreatorChannels).Exists(db.RO)
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
	modChannel, _, err := getModerator(args.ModChannelID, args.ModChannelName, args.CreatorChannelID, args.CreatorChannelName)
	if err != nil {
		return err
	}
	err = lbry.ValidateSignature(modChannel.ClaimID, args.Signature, args.SigningTS, args.ModChannelName)
	if err != nil {
		return err
	}

	isMod, err := modChannel.ModChannelModerators().Exists(db.RO)
	if err != nil {
		return errors.Err(err)
	}

	var blockedByMod model.BlockedEntrySlice
	var blockedByCreator model.BlockedEntrySlice
	var blockedGlobally model.BlockedEntrySlice

	blockedByMod, err = modChannel.CreatorChannelBlockedEntries(qm.Load(model.BlockedEntryRels.BlockedChannel), model.BlockedEntryWhere.UniversallyBlocked.EQ(null.BoolFrom(false))).All(db.RO)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return errors.Err(err)
	}

	blockedByCreator, err = getDelegatedEntries(modChannel)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return errors.Err(err)
	}

	if isMod {
		blockedGlobally, err = model.BlockedEntries(qm.Load(model.BlockedEntryRels.BlockedChannel), qm.Load(model.BlockedEntryRels.CreatorChannel), model.BlockedEntryWhere.UniversallyBlocked.EQ(null.BoolFrom(true))).All(db.RO)
		if err != nil && errors.Is(err, sql.ErrNoRows) {
			return errors.Err(err)
		}
	}

	reply.BlockedChannels = populateBlockedChannelsReply(modChannel, filterBlocks(blockedByMod))
	reply.DelegatedBlockedChannels = populateBlockedChannelsReply(nil, filterBlocks(blockedByCreator))
	reply.GloballyBlockedChannels = populateBlockedChannelsReply(modChannel, filterBlocks(blockedGlobally))

	return nil
}

func filterBlocks(list model.BlockedEntrySlice) model.BlockedEntrySlice {
	var out model.BlockedEntrySlice
	for _, l := range list {
		if l.Expiry.Valid && l.Expiry.Time.Before(time.Now()) {
			continue
		}
		out = append(out, l)
	}
	return out
}

func getDelegatedEntries(modChannel *model.Channel) (model.BlockedEntrySlice, error) {
	var blockedByCreator model.BlockedEntrySlice
	moderations, err := modChannel.ModChannelDelegatedModerators(qm.Load(model.DelegatedModeratorRels.CreatorChannel)).All(db.RO)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Err(err)
	}
	var creatorIDs []interface{}
	for _, m := range moderations {
		creatorIDs = append(creatorIDs, m.CreatorChannelID)
	}
	blockedByCreator, err = model.BlockedEntries(qm.WhereIn(model.BlockedEntryColumns.CreatorChannelID+" IN ?", creatorIDs...), qm.Load(model.BlockedEntryRels.BlockedChannel), qm.Load(model.BlockedEntryRels.CreatorChannel)).All(db.RO)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Err(err)
	}
	return blockedByCreator, nil
}

func populateBlockedChannelsReply(blockedBy *model.Channel, blocked model.BlockedEntrySlice) []commentapi.BlockedChannel {
	var blockedChannels []commentapi.BlockedChannel
	for _, b := range blocked {
		blockedByChannel := blockedBy
		if b.R != nil && b.R.BlockedChannel != nil {
			if b.R.CreatorChannel != nil && blockedBy == nil {
				blockedByChannel = b.R.CreatorChannel
			}
			var blockedFor time.Duration
			var blockRemaining time.Duration
			if b.Expiry.Valid {
				blockedFor = b.Expiry.Time.Sub(b.CreatedAt)
				if b.Expiry.Time.After(time.Now()) {
					blockRemaining = b.Expiry.Time.Sub(time.Now())
				}
			}
			blockedChannels = append(blockedChannels, commentapi.BlockedChannel{
				BlockedChannelID:     b.R.BlockedChannel.ClaimID,
				BlockedChannelName:   b.R.BlockedChannel.Name,
				BlockedByChannelID:   blockedByChannel.ClaimID,
				BlockedByChannelName: blockedByChannel.Name,
				BlockedAt:            b.CreatedAt,
				BlockedFor:           blockedFor,
				BlcokRemaining:       blockRemaining,
			})
		}
	}
	return blockedChannels
}
