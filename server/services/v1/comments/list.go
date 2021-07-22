package comments

import (
	"database/sql"
	"math"
	"net/http"

	"github.com/lbryio/commentron/commentapi"
	"github.com/lbryio/commentron/helper"
	m "github.com/lbryio/commentron/model"
	"github.com/lbryio/commentron/server/lbry"

	"github.com/lbryio/lbry.go/extras/util"
	"github.com/lbryio/lbry.go/v2/extras/api"
	"github.com/lbryio/lbry.go/v2/extras/errors"

	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

func list(_ *http.Request, args *commentapi.ListArgs, reply *commentapi.ListResponse) error {
	args.ApplyDefaults()
	creatorChannel, err := checkCommentsEnabled(null.StringFromPtr(args.ChannelName), null.StringFromPtr(args.ChannelID))
	if err != nil {
		return err
	}
	loadChannels := qm.Load("Channel.BlockedChannelBlockedEntries")
	filterIsHidden := m.CommentWhere.IsHidden.EQ(null.BoolFrom(true))
	filterClaimID := m.CommentWhere.LbryClaimID.EQ(util.StrFromPtr(args.ClaimID))
	filterAuthorClaimID := m.CommentWhere.ChannelID.EQ(null.StringFromPtr(args.AuthorClaimID))
	filterTopLevel := m.CommentWhere.ParentID.IsNull()
	filterParent := m.CommentWhere.ParentID.EQ(null.StringFrom(util.StrFromPtr(args.ParentID)))

	totalFilteredCommentsQuery := make([]qm.QueryMod, 0)
	totalCommentsQuery := make([]qm.QueryMod, 0)
	offset := (args.Page - 1) * args.PageSize
	getCommentsQuery := applySorting(args.SortBy, []qm.QueryMod{loadChannels, qm.Offset(offset), qm.Limit(args.PageSize)})
	hasHiddenCommentsQuery := []qm.QueryMod{filterIsHidden, qm.Limit(1)}

	if args.AuthorClaimID != nil {
		getCommentsQuery = append(getCommentsQuery, filterAuthorClaimID)
		hasHiddenCommentsQuery = append(hasHiddenCommentsQuery, filterAuthorClaimID)
		totalFilteredCommentsQuery = append(totalFilteredCommentsQuery, filterAuthorClaimID)
		totalCommentsQuery = append(totalCommentsQuery, filterAuthorClaimID)
	}

	if args.ClaimID != nil {
		getCommentsQuery = append(getCommentsQuery, filterClaimID)
		hasHiddenCommentsQuery = append(hasHiddenCommentsQuery, filterClaimID)
		totalFilteredCommentsQuery = append(totalFilteredCommentsQuery, filterClaimID)
		totalCommentsQuery = append(totalCommentsQuery, filterClaimID)
	}

	if args.TopLevel {
		getCommentsQuery = append(getCommentsQuery, filterTopLevel)
		hasHiddenCommentsQuery = append(hasHiddenCommentsQuery, filterTopLevel)
		totalFilteredCommentsQuery = append(totalFilteredCommentsQuery, filterTopLevel)
	}

	if args.ParentID != nil {
		getCommentsQuery = append(getCommentsQuery, filterParent)
		hasHiddenCommentsQuery = append(hasHiddenCommentsQuery, filterParent)
		totalFilteredCommentsQuery = append(totalFilteredCommentsQuery, filterParent)
		totalCommentsQuery = append(totalCommentsQuery, filterParent)
	}

	totalFilteredItems, err := m.Comments(totalFilteredCommentsQuery...).CountG()
	if err != nil {
		return errors.Err(err)
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

	items, blockedCommentCnt, err := getItems(comments, creatorChannel)

	totalFilteredItems = totalFilteredItems - blockedCommentCnt
	reply.Items = items
	reply.Page = args.Page
	reply.PageSize = args.PageSize
	reply.TotalFilteredItems = totalFilteredItems
	reply.TotalItems = totalItems
	reply.TotalPages = int(math.Ceil(float64(totalFilteredItems) / float64(args.PageSize)))
	reply.HasHiddenComments = hasHiddenComments

	return nil
}

func applySorting(sort commentapi.Sort, queryMods []qm.QueryMod) []qm.QueryMod {
	if sort != commentapi.Newest {
		if sort == commentapi.Popularity {
			queryMods = append(queryMods, qm.OrderBy(m.CommentColumns.IsPinned+" DESC, "+m.CommentColumns.PopularityScore+" DESC, "+m.CommentColumns.Timestamp+" DESC"))
		} else if sort == commentapi.Controversy {
			queryMods = append(queryMods, qm.OrderBy(m.CommentColumns.IsPinned+" DESC, "+m.CommentColumns.ControversyScore+" DESC, "+m.CommentColumns.Timestamp+" DESC"))
		} else if sort == commentapi.Oldest {
			queryMods = append(queryMods, qm.OrderBy(m.CommentColumns.IsPinned+" DESC, "+m.CommentColumns.Timestamp+" ASC"))
		}
	} else {
		queryMods = append(queryMods, qm.OrderBy(m.CommentColumns.IsPinned+" DESC, "+m.CommentColumns.Timestamp+" DESC"))
	}

	return queryMods
}

func checkCommentsEnabled(channelName, ChannelID null.String) (*m.Channel, error) {
	if channelName.Valid && ChannelID.Valid {
		creatorChannel, err := helper.FindOrCreateChannel(ChannelID.String, channelName.String)
		if err != nil {
			return nil, err
		}
		settings, err := creatorChannel.CreatorChannelCreatorSettings().OneG()
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Err(err)
		}
		if settings != nil {
			if !settings.CommentsEnabled.Bool {
				return nil, api.StatusError{Err: errors.Err("comments are disabled by the creator"), Status: http.StatusBadRequest}
			}
		}
		return creatorChannel, nil
	}
	return nil, nil
}

func getItems(comments m.CommentSlice, creatorChannel *m.Channel) ([]commentapi.CommentItem, int64, error) {
	var items []commentapi.CommentItem
	var blockedCommentCnt int64
	var alreadyInSet = map[string]bool{}
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
						if creatorChannel != nil && creatorChannel.BlockedListID.Valid {
							if creatorChannel.BlockedListID == entry.BlockedListID {
								blockedCommentCnt++
								continue Comments
							}
						}
						if entry.UniversallyBlocked.Bool || entry.CreatorChannelID.String == channel.ClaimID {
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
				if !alreadyInSet[comment.CommentID] {
					replies, err := comment.ParentComments().CountG()
					if err != nil && errors.Is(err, sql.ErrNoRows) {
						return items, blockedCommentCnt, errors.Err(err)
					}
					alreadyInSet[comment.CommentID] = true
					items = append(items, populateItem(comment, channel, int(replies)))
				}
			}
		}
	}
	return items, blockedCommentCnt, nil
}
