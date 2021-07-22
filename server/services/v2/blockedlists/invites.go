package blockedlists

import (
	"database/sql"
	"net/http"

	"github.com/lbryio/commentron/commentapi"
	"github.com/lbryio/commentron/db"
	"github.com/lbryio/commentron/helper"
	"github.com/lbryio/commentron/model"
	"github.com/lbryio/commentron/server/lbry"

	"github.com/lbryio/lbry.go/extras/api"
	"github.com/lbryio/lbry.go/v2/extras/errors"

	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
)

func invite(_ *http.Request, args *commentapi.SharedBlockedListInviteArgs, reply *commentapi.SharedBlockedListInviteResponse) error {
	err := lbry.ValidateSignature(args.ChannelID, args.Signature, args.SigningTS, args.ChannelName)
	if err != nil {
		return err
	}
	blockedList, err := model.BlockedLists(model.BlockedListWhere.ID.EQ(args.SharedBlockedListID)).One(db.RO)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return api.StatusError{Err: errors.Err("there is no shared block list with id %d", args.SharedBlockedListID), Status: http.StatusBadRequest}
		}
		return errors.Err(err)
	}

	inviter, err := helper.FindOrCreateChannel(args.ChannelID, args.ChannelName)
	if err != nil {
		return errors.Err(err)
	}
	if (blockedList.MemberInviteEnabled.Valid && !blockedList.MemberInviteEnabled.Bool) && inviter.ClaimID != blockedList.ChannelID {
		return api.StatusError{Err: errors.Err("shared blocked list %s does not have member inviting enabled", blockedList.Name)}
	}
	if !inviter.BlockedListInviteID.Valid {
		return api.StatusError{Err: errors.Err("channel %s is not authorized member of the shared blocked list %s", inviter.Name, blockedList.Name), Status: http.StatusUnauthorized}
	}
	if inviter.BlockedListInviteID.Uint64 != blockedList.ID {
		return api.StatusError{Err: errors.Err("channel %s is not a member of the shared blocked list %s", inviter.Name, blockedList.Name), Status: http.StatusBadRequest}
	}

	invitee, err := helper.FindOrCreateChannel(args.InviteeChannelID, args.InviteeChannelName)
	if err != nil {
		return errors.Err(err)
	}
	if invitee.BlockedListInviteID.Valid && invitee.BlockedListInviteID.Uint64 == blockedList.ID {
		return api.StatusError{Err: errors.Err("channel %s is already a member of the shared blocked list %s", invitee.Name, blockedList.Name), Status: http.StatusBadRequest}
	}
	where := model.BlockedListInviteWhere
	invite, err := model.BlockedListInvites(where.BlockedListID.EQ(args.SharedBlockedListID), where.InvitedChannelID.EQ(args.InviteeChannelID)).One(db.RO)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return errors.Err(err)
	}
	if invite != nil {
		if !invite.Accepted.Valid {
			return api.StatusError{Err: errors.Err("channel %s already has an invite pending", invitee.Name), Status: http.StatusBadRequest}
		} else if !invite.Accepted.Bool {
			return api.StatusError{Err: errors.Err("channel %s already an invite and has rejected joining the shared blocked list %s", invitee.Name, blockedList.Name)}
		}
	}

	invite = &model.BlockedListInvite{
		BlockedListID:    blockedList.ID,
		InviterChannelID: inviter.ClaimID,
		InvitedChannelID: invitee.ClaimID,
		Message:          args.Message,
	}
	err = invite.Insert(db.RW, boil.Infer())
	if err != nil {
		return errors.Err(err)
	}

	return nil
}

func accept(_ *http.Request, args *commentapi.SharedBlockedListInviteAcceptArgs, _ *commentapi.SharedBlockedListInviteResponse) error {
	err := lbry.ValidateSignature(args.ChannelID, args.Signature, args.SigningTS, args.ChannelName)
	if err != nil {
		return err
	}

	channel, err := helper.FindOrCreateChannel(args.ChannelID, args.ChannelName)
	if err != nil {
		return err
	}

	where := model.BlockedListInviteWhere
	invite, err := model.BlockedListInvites(where.BlockedListID.EQ(args.SharedBlockedListID), where.InvitedChannelID.EQ(channel.ClaimID)).One(db.RO)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return errors.Err(err)
	}

	if invite == nil {
		return api.StatusError{Err: errors.Err("channel %s does not have an invite for the shared block list %s to accept", args.ChannelName)}
	}

	blockedList, err := model.BlockedLists(model.BlockedListWhere.ID.EQ(args.SharedBlockedListID)).One(db.RO)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return api.StatusError{Err: errors.Err("there is no shared block list with id %d", args.SharedBlockedListID), Status: http.StatusBadRequest}
		}
		return errors.Err(err)
	}

	var blockedListID = null.Uint64{}
	if args.Accepted {
		blockedListID = null.Uint64From(blockedList.ID)
	}

	blockedListCol := map[string]interface{}{model.BlockedEntryColumns.BlockedListID: blockedListID}
	err = channel.CreatorChannelBlockedEntries().UpdateAll(db.RW, blockedListCol)
	if err != nil {
		return errors.Err(err)
	}

	channel.BlockedListID = blockedListID
	channel.BlockedListInviteID = blockedListID
	err = channel.Update(db.RW, boil.Whitelist(model.ChannelColumns.BlockedListID, model.ChannelColumns.BlockedListInviteID))
	if err != nil {
		return errors.Err(err)
	}
	return nil
}
