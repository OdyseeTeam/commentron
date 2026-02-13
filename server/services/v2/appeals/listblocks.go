package appeals

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/OdyseeTeam/commentron/commentapi"
	"github.com/OdyseeTeam/commentron/db"
	"github.com/OdyseeTeam/commentron/model"
	"github.com/OdyseeTeam/commentron/server/auth"
	"github.com/OdyseeTeam/commentron/server/services/v2/blockedlists"

	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/lbryio/lbry.go/v2/extras/errors"
)

func listBlocks(r *http.Request, args *commentapi.AppealBlockListArgs, reply *commentapi.AppealBlockListResponse) error {
	ownerChannel, _, err := auth.Authenticate(r, &args.Authorization)
	if err != nil {
		return err
	}
	blockedEntries, err := ownerChannel.BlockedChannelBlockedEntries(
		qm.Load(model.BlockedEntryRels.BlockedList),
		qm.Load(model.BlockedEntryRels.CreatorChannel),
		qm.Load(model.BlockedEntryRels.BlockedListAppeals),
	).All(db.RO)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return errors.Err(err)
	}

	var blocks []commentapi.Appeal
	for _, be := range blockedEntries {
		sbl := commentapi.SharedBlockedList{}
		if be.R.BlockedList != nil {
			err := blockedlists.PopulateSharedBlockedList(&sbl, be.R.BlockedList)
			if err != nil {
				return errors.Err(err)
			}
		}
		var blockedFor time.Duration
		var blockRemaining time.Duration
		if be.Expiry.Valid {
			blockedFor = be.Expiry.Time.Sub(be.CreatedAt)
			if be.Expiry.Time.After(time.Now()) {
				blockRemaining = time.Until(be.Expiry.Time)
			}
		}
		appeal := commentapi.AppealRequest{}
		if len(be.R.BlockedListAppeals) > 0 {
			a := be.R.BlockedListAppeals[0]
			appeal.AppealStatus = getAppealStatus(a)
			appeal.TxID = a.TXID.String
			appeal.AppealMessage = a.Appeal
			appeal.ResponseMessage = a.Response
		}
		block := commentapi.Appeal{
			BlockedList: sbl,
			BlockedChannel: commentapi.BlockedChannel{
				BlockedChannelID:     ownerChannel.ClaimID,
				BlockedChannelName:   ownerChannel.Name,
				BlockedByChannelID:   be.R.CreatorChannel.ClaimID,
				BlockedByChannelName: be.R.CreatorChannel.Name,
				BlockedAt:            be.CreatedAt,
				BlockedFor:           blockedFor,
				BlcokRemaining:       blockRemaining,
			},
			AppealRequest: appeal,
		}
		blocks = append(blocks, block)
	}
	reply.Blocks = blocks
	return nil
}

func getAppealStatus(appeal *model.BlockedListAppeal) commentapi.AppealStatus {
	if !appeal.Approved.Valid {
		return commentapi.AppealPending
	} else if appeal.Approved.Valid && !appeal.Approved.Bool {
		if appeal.Escalated.Valid {
			return commentapi.AppealEscalated
		}
		return commentapi.AppealRejected
	}
	return commentapi.AppealAccepted
}
