package helper

import (
	"database/sql"
	"net/http"

	"github.com/lbryio/commentron/db"
	m "github.com/lbryio/commentron/model"

	"github.com/lbryio/lbry.go/extras/api"
	"github.com/lbryio/lbry.go/v2/extras/errors"

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
