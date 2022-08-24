package appeals

import (
	"database/sql"
	"net/http"

	"github.com/lbryio/commentron/commentapi"
	"github.com/lbryio/commentron/db"
	"github.com/lbryio/commentron/model"
	"github.com/lbryio/commentron/server/auth"

	"github.com/lbryio/lbry.go/v2/extras/errors"

	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
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
	//var appeals []commentapi.Appeal
	for _, entry := range entries {
		if len(entry.R.BlockedListAppeals) > 0 {

		}
	}

	return nil, nil
}

func getModeratedAppeals(moderator *model.Channel, creator *model.Channel) ([]commentapi.Appeal, error) {
	if creator.ClaimID == moderator.ClaimID {
		return nil, nil
	}

	return getAppeals(creator.ClaimID)
}
