package helper

import (
	"database/sql"
	"net/http"

	m "github.com/lbryio/commentron/model"

	"github.com/lbryio/lbry.go/extras/api"
	"github.com/lbryio/lbry.go/v2/extras/errors"

	"github.com/volatiletech/null"
)

func AllowedToRespond(parentCommentID, commenterClaimID string) error {
	parentComment, err := m.Comments(m.CommentWhere.CommentID.EQ(parentCommentID)).OneG()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return errors.Err(err)
	}
	if parentComment != nil {
		parentChannel, err := parentComment.Channel().OneG()
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return errors.Err(err)
		}
		if parentChannel != nil {

			blockedEntry, err := m.BlockedEntries(
				m.BlockedEntryWhere.BlockedByChannelID.EQ(null.StringFrom(parentChannel.ClaimID)),
				m.BlockedEntryWhere.BlockedChannelID.EQ(null.StringFrom(commenterClaimID))).OneG()
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
