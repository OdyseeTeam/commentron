package appeals

import (
	"database/sql"
	"net/http"

	"github.com/OdyseeTeam/commentron/commentapi"
	"github.com/OdyseeTeam/commentron/db"
	"github.com/OdyseeTeam/commentron/model"
	"github.com/OdyseeTeam/commentron/server/auth"

	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/lbryio/lbry.go/v2/extras/errors"
)

func listAppeals(r *http.Request, args *commentapi.AppealListArgs, reply *commentapi.AppealListResponse) error {
	modChannel, creatorChannel, _, err := auth.ModAuthenticate(r, &args.ModAuthorization)
	if err != nil {
		return err
	}

	ownerappeals, err := getAppeals(modChannel.ClaimID)
	if err != nil {
		return errors.Err(err)
	}
	moderatedAppeals, err := getModeratedAppeals(modChannel, creatorChannel)
	if err != nil {
		return errors.Err(err)
	}

	reply.Appeals = ownerappeals
	reply.ModeratedAppeals = moderatedAppeals

	return nil
}

func getAppeals(channelID string) ([]commentapi.Appeal, error) {
	entries, err := model.BlockedEntries(
		model.BlockedEntryWhere.CreatorChannelID.EQ(null.StringFrom(channelID)),
		qm.Load(model.BlockedEntryRels.BlockedListAppeals),
		qm.Load(model.BlockedEntryRels.BlockedList),
	).All(db.RO)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Err(err)
	}
	for _, entry := range entries {
		_ = entry
	}

	return nil, nil
}

func getModeratedAppeals(moderator *model.Channel, creator *model.Channel) ([]commentapi.Appeal, error) {
	if creator.ClaimID == moderator.ClaimID {
		return nil, nil
	}

	return getAppeals(creator.ClaimID)
}
