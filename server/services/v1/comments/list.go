package comments

import (
	"database/sql"
	"math"
	"net/http"

	"github.com/lbryio/commentron/commentapi"
	m "github.com/lbryio/commentron/model"
	"github.com/lbryio/commentron/server/lbry"

	"github.com/lbryio/lbry.go/extras/util"
	"github.com/lbryio/lbry.go/v2/extras/errors"

	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

func list(_ *http.Request, args *commentapi.ListArgs, reply *commentapi.ListResponse) error {
	args.ApplyDefaults()
	loadChannels := qm.Load("Channel.BlockedChannelBlockedEntries")
	filterIsHidden := m.CommentWhere.IsHidden.EQ(null.BoolFrom(true))
	filterClaimID := m.CommentWhere.LbryClaimID.EQ(util.StrFromPtr(args.ClaimID))
	filterAuthorClaimID := m.CommentWhere.ChannelID.EQ(null.StringFromPtr(args.AuthorClaimID))
	filterTopLevel := m.CommentWhere.ParentID.IsNull()
	filterParent := m.CommentWhere.ParentID.EQ(null.StringFrom(util.StrFromPtr(args.ParentID)))

	totalCommentsQuery := make([]qm.QueryMod, 0)
	offset := (args.Page - 1) * args.PageSize
	getCommentsQuery := []qm.QueryMod{loadChannels, qm.Offset(offset), qm.Limit(args.PageSize), qm.OrderBy(m.CommentColumns.Timestamp + " DESC")}
	hasHiddenCommentsQuery := []qm.QueryMod{filterIsHidden, qm.Limit(1)}

	if args.AuthorClaimID != nil {
		getCommentsQuery = append(getCommentsQuery, filterAuthorClaimID)
		hasHiddenCommentsQuery = append(hasHiddenCommentsQuery, filterAuthorClaimID)
		totalCommentsQuery = append(totalCommentsQuery, filterAuthorClaimID)
	}

	if args.ClaimID != nil {
		getCommentsQuery = append(getCommentsQuery, filterClaimID)
		hasHiddenCommentsQuery = append(hasHiddenCommentsQuery, filterClaimID)
		totalCommentsQuery = append(totalCommentsQuery, filterClaimID)
	}

	if args.TopLevel {
		getCommentsQuery = append(getCommentsQuery, filterTopLevel)
		hasHiddenCommentsQuery = append(hasHiddenCommentsQuery, filterTopLevel)
		totalCommentsQuery = append(totalCommentsQuery, filterTopLevel)
	}

	if args.ParentID != nil {
		getCommentsQuery = append(getCommentsQuery, filterParent)
		hasHiddenCommentsQuery = append(hasHiddenCommentsQuery, filterParent)
		totalCommentsQuery = append(totalCommentsQuery, filterParent)
	}

	totalItems, err := m.Comments(totalCommentsQuery...).CountG()
	if err != nil {
		return errors.Err(err)
	}

	hasHiddenComments, err := m.Comments(hasHiddenCommentsQuery...).ExistsG()
	if err != nil {
		return errors.Err(err)
	}

	comments, err := m.Comments(getCommentsQuery...).AllG()
	if err != nil {
		return errors.Err(err)
	}

	items, blockedCommentCnt, err := getItems(comments)

	totalItems = totalItems - blockedCommentCnt
	reply.Items = items
	reply.Page = args.Page
	reply.PageSize = args.PageSize
	reply.TotalItems = totalItems
	reply.TotalPages = int(math.Ceil(float64(totalItems) / float64(args.PageSize)))
	reply.HasHiddenComments = hasHiddenComments

	return nil
}

func getItems(comments m.CommentSlice) ([]commentapi.CommentItem, int64, error) {
	var items []commentapi.CommentItem
	var blockedCommentCnt int64
Comments:
	for _, comment := range comments {
		if comment.R != nil && comment.R.Channel != nil && comment.R.Channel.R != nil {
			blockedFrom := comment.R.Channel.R.BlockedChannelBlockedEntries
			if len(blockedFrom) > 0 {
				channel, err := lbry.SDK.GetSigningChannelForClaim(comment.LbryClaimID)
				if err != nil {
					return items, blockedCommentCnt, errors.Err(err)
				}
				if channel != nil {
					for _, entry := range blockedFrom {
						if entry.UniversallyBlocked.Bool || entry.BlockedByChannelID.String == channel.ClaimID {
							blockedCommentCnt++
							continue Comments
						}
					}
				}
			}
		}
		var channel *m.Channel
		if comment.R != nil {
			channel = comment.R.Channel
			if channel != nil && channel.Name != "" {
				replies, err := comment.ParentComments().CountG()
				if err != nil && errors.Is(err, sql.ErrNoRows) {
					return items, blockedCommentCnt, errors.Err(err)
				}
				items = append(items, populateItem(comment, channel, int(replies)))
			}
		}
	}
	return items, blockedCommentCnt, nil
}
