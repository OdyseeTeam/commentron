package blockedlists

import (
	"database/sql"
	"net/http"

	"github.com/lbryio/commentron/commentapi"
	"github.com/lbryio/commentron/db"
	"github.com/lbryio/commentron/model"

	"github.com/lbryio/lbry.go/extras/api"
	"github.com/lbryio/lbry.go/v2/extras/errors"

	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

func get(_ *http.Request, args *commentapi.SharedBlockedListGetArgs, reply *commentapi.SharedBlockedListGetResponse) error {
	list, err := model.BlockedLists(model.BlockedListWhere.ID.EQ(args.SharedBlockedListID)).One(db.RO)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return errors.Err(err)
	}

	if list == nil {
		return api.StatusError{Err: errors.Err("blocked list not found"), Status: http.StatusNotFound}
	}

	var acceptedFilter qm.QueryMod
	if args.Status == commentapi.Pending {
		acceptedFilter = model.BlockedListInviteWhere.Accepted.EQ(null.Bool{})
	} else if args.Status == commentapi.Accepted {
		acceptedFilter = model.BlockedListInviteWhere.Accepted.EQ(null.BoolFrom(true))
	} else if args.Status == commentapi.Rejected {
		acceptedFilter = model.BlockedListInviteWhere.Accepted.EQ(null.BoolFrom(false))
	}

	err = populateSharedBlockedList(&reply.BlockedList, list)
	if err != nil {
		return err
	}
	var invitedMembers []commentapi.SharedBlockedListInvitedMember
	if args.Status != commentapi.None {
		invites, err := list.BlockedListInvites(acceptedFilter,
			qm.Load(model.BlockedListInviteRels.InvitedChannel),
			qm.Load(model.BlockedListInviteRels.InviterChannel)).All(db.RO)
		if err != nil {
			return errors.Err(err)
		}

		for _, invite := range invites {
			if invite.R != nil && invite.R.InvitedChannel != nil {
				member := commentapi.SharedBlockedListInvitedMember{
					InviterChannelName: invite.R.InvitedChannel.Name,
					InviterChannelID:   invite.R.InvitedChannel.ClaimID,
					Status:             commentapi.InviteMemberStatusFrom(invite.Accepted),
					InviteMessage:      invite.Message,
				}
				if invite.R.InviterChannel != nil {
					member.InviterChannelName = invite.R.InviterChannel.Name
					member.InviterChannelID = invite.R.InviterChannel.ClaimID
				}
				invitedMembers = append(invitedMembers, member)
			}
		}
		reply.InvitedMembers = invitedMembers
	}

	return nil
}
