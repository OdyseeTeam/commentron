package blockedlists

import (
	"database/sql"
	"net/http"

	"github.com/lbryio/commentron/helper"
	"github.com/lbryio/commentron/server/lbry"

	"github.com/lbryio/commentron/commentapi"
	"github.com/lbryio/commentron/db"
	"github.com/lbryio/commentron/model"

	"github.com/lbryio/lbry.go/extras/api"
	"github.com/lbryio/lbry.go/v2/extras/errors"

	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

func get(_ *http.Request, args *commentapi.SharedBlockedListGetArgs, reply *commentapi.SharedBlockedListGetResponse) error {
	var list *model.BlockedList
	var err error
	var ownerChannel *model.Channel
	if args.SharedBlockedListID != 0 {
		list, err = model.BlockedLists(model.BlockedListWhere.ID.EQ(args.SharedBlockedListID)).One(db.RO)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return errors.Err(err)
		}
	} else {
		ownerChannel, err = helper.FindOrCreateChannel(args.ChannelID, args.ChannelName)
		if err != nil {
			return errors.Err(err)
		}
		err = lbry.ValidateSignature(ownerChannel.ClaimID, args.Signature, args.SigningTS, args.ChannelName)
		if err != nil {
			return err
		}

		list, err = model.BlockedLists(model.BlockedListWhere.ChannelID.EQ(ownerChannel.ClaimID)).One(db.RO)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return errors.Err(err)
		}
	}

	if list == nil {
		return api.StatusError{Err: errors.Err("blocked list not found"), Status: http.StatusNotFound}
	}

	err = populateSharedBlockedList(&reply.BlockedList, list)
	if err != nil {
		return err
	}
	var invitedMembers []commentapi.SharedBlockedListInvitedMember
	if args.Status != commentapi.None && ownerChannel != nil {
		invitesFilters := []qm.QueryMod{qm.Load(model.BlockedListInviteRels.InvitedChannel), qm.Load(model.BlockedListInviteRels.InviterChannel)}
		if args.Status == commentapi.Pending {
			invitesFilters = append(invitesFilters, model.BlockedListInviteWhere.Accepted.EQ(null.Bool{}))
		} else if args.Status == commentapi.Accepted {
			invitesFilters = append(invitesFilters, model.BlockedListInviteWhere.Accepted.EQ(null.BoolFrom(true)))
		} else if args.Status == commentapi.Rejected {
			invitesFilters = append(invitesFilters, model.BlockedListInviteWhere.Accepted.EQ(null.BoolFrom(true)))
		}
		invites, err := list.BlockedListInvites(invitesFilters...).All(db.RO)
		if err != nil {
			return errors.Err(err)
		}

		for _, invite := range invites {
			if invite.R != nil && invite.R.InvitedChannel != nil {
				member := commentapi.SharedBlockedListInvitedMember{
					InvitedByChannelName: invite.R.InviterChannel.Name,
					InvitedByChannelID:   invite.R.InviterChannel.ClaimID,
					InvitedChannelName:   invite.R.InvitedChannel.Name,
					InvitedChannelID:     invite.R.InvitedChannel.ClaimID,
					Status:               commentapi.InviteMemberStatusFrom(invite.Accepted, invite.CreatedAt, list.InviteExpiration),
					InviteMessage:        invite.Message,
				}
				invitedMembers = append(invitedMembers, member)
			}
		}
		reply.InvitedMembers = invitedMembers
	}

	return nil
}
