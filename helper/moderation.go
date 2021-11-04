package helper

import (
	"database/sql"
	"net/http"

	"github.com/lbryio/commentron/db"
	"github.com/lbryio/commentron/model"
	m "github.com/lbryio/commentron/model"

	"github.com/lbryio/lbry.go/v2/extras/api"
	"github.com/lbryio/lbry.go/v2/extras/errors"
	"github.com/volatiletech/sqlboiler/queries/qm"

	"github.com/volatiletech/null"
)

// AllowedToRespond checks if the creator of the comment will allow a response from the respondent
func AllowedToRespond(parentCommentID, commenterClaimID string) error {
	parentComment, err := m.Comments(m.CommentWhere.CommentID.EQ(parentCommentID)).One(db.RO)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return errors.Err(err)
	}
	if parentComment != nil {
		parentChannel, err := parentComment.Channel().One(db.RO)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return errors.Err(err)
		}
		if parentChannel != nil {

			blockedEntry, err := m.BlockedEntries(
				m.BlockedEntryWhere.CreatorChannelID.EQ(null.StringFrom(parentChannel.ClaimID)),
				m.BlockedEntryWhere.BlockedChannelID.EQ(null.StringFrom(commenterClaimID))).One(db.RO)
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				return errors.Err(err)
			}
			if blockedEntry != nil {
				return api.StatusError{Err: errors.Err("'%s' has blocked you from replying to their comments", parentChannel.Name), Status: http.StatusBadRequest}
			}
		}
	}
	return nil
}

// GetModerator returns the validated moderator and the creator which delegated the moderator. If a creator is not passed
// the moderator will be returned as the creator and will be equal.
func GetModerator(modChannelID, modChannelName, creatorChannelID, creatorChannelName string) (*model.Channel, *model.Channel, error) {
	modChannel, err := FindOrCreateChannel(modChannelID, modChannelName)
	if err != nil {
		return nil, nil, errors.Err(err)
	}
	var creatorChannel = modChannel
	if creatorChannelID != "" && creatorChannelName != "" {
		creatorChannel, err = FindOrCreateChannel(creatorChannelID, creatorChannelName)
		if err != nil {
			return nil, nil, errors.Err(err)
		}
		dmRels := model.DelegatedModeratorRels
		dmWhere := model.DelegatedModeratorWhere
		loadCreatorChannels := qm.Load(dmRels.CreatorChannel, dmWhere.CreatorChannelID.EQ(creatorChannelID))
		exists, err := modChannel.ModChannelDelegatedModerators(loadCreatorChannels, dmWhere.CreatorChannelID.EQ(creatorChannelID)).Exists(db.RO)
		if err != nil {
			return nil, nil, errors.Err(err)
		}
		if !exists {
			return nil, nil, errors.Err("%s is not delegated by %s to be a moderator", modChannel.Name, creatorChannel.Name)
		}
	}
	return modChannel, creatorChannel, nil
}
