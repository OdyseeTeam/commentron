package comments

import (
	"database/sql"
	"math"
	"net/http"
	"time"

	"github.com/lbryio/commentron/commentapi"
	"github.com/lbryio/commentron/db"
	"github.com/lbryio/commentron/helper"
	m "github.com/lbryio/commentron/model"
	"github.com/lbryio/commentron/server/lbry"

	"github.com/lbryio/lbry.go/v2/extras/api"
	"github.com/lbryio/lbry.go/v2/extras/errors"
	"github.com/lbryio/lbry.go/v2/extras/util"

	"github.com/karlseguin/ccache/v2"
	"github.com/sirupsen/logrus"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func list(_ *http.Request, args *commentapi.ListArgs, reply *commentapi.ListResponse) error {
	args.ApplyDefaults()
	creatorChannel, err := checkCommentsEnabled(args.ChannelName, args.ChannelID)
	if err != nil {
		return err
	}

	isListingOwnComments := args.AuthorClaimID != nil
	actualIsProtected := args.IsProtected
	if !isListingOwnComments {
		actualIsProtected, err := IsProtectedContent(*args.ClaimID)
		if err != nil {
			return err
		}
		if actualIsProtected != args.IsProtected {
			return errors.Err("mismatch in is_protected")
		}
	} else {
		if args.RequestorChannelID == nil {
			return errors.Err("requestor channel id is required to list own comments")
		}
		ownerChannel, err := helper.FindOrCreateChannel(*args.RequestorChannelID, args.RequestorChannelName)
		if err != nil {
			return errors.Err(err)
		}
		err = lbry.ValidateSignatureAndTS(ownerChannel.ClaimID, args.Signature, args.SigningTS, args.RequestorChannelName)
		if err != nil {
			return err
		}
		if ownerChannel.ClaimID != *args.AuthorClaimID {
			return api.StatusError{Err: errors.Err("you can only view your comments, not others"), Status: http.StatusBadRequest}
		}
	}
	loadChannels := qm.Load("Channel.BlockedChannelBlockedEntries")
	filterIsHidden := m.CommentWhere.IsHidden.EQ(null.BoolFrom(true))
	filterIsProtected := m.CommentWhere.IsProtected.EQ(true)
	filterClaimID := m.CommentWhere.LbryClaimID.EQ(util.StrFromPtr(args.ClaimID))
	filterAuthorClaimID := m.CommentWhere.ChannelID.EQ(null.StringFromPtr(args.AuthorClaimID))
	filterTopLevel := m.CommentWhere.ParentID.IsNull()
	filterParent := m.CommentWhere.ParentID.EQ(null.StringFrom(util.StrFromPtr(args.ParentID)))
	filterFlaggedComments := m.CommentWhere.IsFlagged.EQ(false)

	totalFilteredCommentsQuery := make([]qm.QueryMod, 0)
	totalCommentsQuery := make([]qm.QueryMod, 0)
	offset := (args.Page - 1) * args.PageSize
	getCommentsQuery := applySorting(args.SortBy, []qm.QueryMod{loadChannels, qm.Offset(offset), qm.Limit(args.PageSize)})
	hasHiddenCommentsQuery := []qm.QueryMod{filterIsHidden, qm.Limit(1)}
	HasProtectedCommentsQuery := []qm.QueryMod{filterIsProtected, qm.Limit(1)}

	if isListingOwnComments {
		getCommentsQuery = append(getCommentsQuery, filterAuthorClaimID)
		hasHiddenCommentsQuery = append(hasHiddenCommentsQuery, filterAuthorClaimID)
		HasProtectedCommentsQuery = append(HasProtectedCommentsQuery, filterAuthorClaimID)
		totalFilteredCommentsQuery = append(totalFilteredCommentsQuery, filterAuthorClaimID)
		totalCommentsQuery = append(totalCommentsQuery, filterAuthorClaimID)
	}

	if args.ClaimID != nil {
		getCommentsQuery = append(getCommentsQuery, filterClaimID)
		hasHiddenCommentsQuery = append(hasHiddenCommentsQuery, filterClaimID)
		HasProtectedCommentsQuery = append(HasProtectedCommentsQuery, filterClaimID)
		totalFilteredCommentsQuery = append(totalFilteredCommentsQuery, filterClaimID)
		totalCommentsQuery = append(totalCommentsQuery, filterClaimID)
	}

	if args.TopLevel {
		getCommentsQuery = append(getCommentsQuery, filterTopLevel)
		hasHiddenCommentsQuery = append(hasHiddenCommentsQuery, filterTopLevel)
		HasProtectedCommentsQuery = append(HasProtectedCommentsQuery, filterTopLevel)
		totalFilteredCommentsQuery = append(totalFilteredCommentsQuery, filterTopLevel)
	}

	if args.ParentID != nil {
		getCommentsQuery = append(getCommentsQuery, filterParent)
		hasHiddenCommentsQuery = append(hasHiddenCommentsQuery, filterParent)
		HasProtectedCommentsQuery = append(HasProtectedCommentsQuery, filterParent)
		totalFilteredCommentsQuery = append(totalFilteredCommentsQuery, filterParent)
		totalCommentsQuery = append(totalCommentsQuery, filterParent)
	}
	if !isListingOwnComments {
		totalCommentsQuery = append(totalCommentsQuery, m.CommentWhere.IsProtected.EQ(actualIsProtected))
	}
	getCommentsQuery = append(getCommentsQuery, filterFlaggedComments)

	totalFilteredItems, err := m.Comments(totalFilteredCommentsQuery...).Count(db.RO)
	if err != nil {
		return errors.Err(err)
	}

	totalItems, err := m.Comments(totalCommentsQuery...).Count(db.RO)
	if err != nil {
		return errors.Err(err)
	}

	hasHiddenComments, err := m.Comments(hasHiddenCommentsQuery...).Exists(db.RO)
	if err != nil {
		return errors.Err(err)
	}

	HasProtectedComments, err := m.Comments(HasProtectedCommentsQuery...).Exists(db.RO)
	if err != nil {
		return errors.Err(err)
	}

	comments, err := m.Comments(getCommentsQuery...).All(db.RO)
	if err != nil {
		return errors.Err(err)
	}

	// if listing own comments, show all including blocked ones
	skipBlocked := isListingOwnComments
	items, blockedCommentCnt, err := getItems(comments, creatorChannel, skipBlocked)
	if err != nil {
		logrus.Error(errors.FullTrace(err))
	}

	totalFilteredItems = totalFilteredItems - blockedCommentCnt
	reply.Items = items
	reply.Page = args.Page
	reply.PageSize = args.PageSize
	reply.TotalFilteredItems = totalFilteredItems
	reply.TotalItems = totalItems
	reply.TotalPages = int(math.Ceil(float64(totalFilteredItems) / float64(args.PageSize)))
	reply.HasHiddenComments = hasHiddenComments
	reply.HasProtectedComments = HasProtectedComments

	return nil
}

var commentListCache = ccache.New(ccache.Configure().GetsPerPromote(1).MaxSize(100000))

func getCachedList(r *http.Request, args *commentapi.ListArgs, reply *commentapi.ListResponse) error {
	listingOwnComments := args.AuthorClaimID != nil

	if args.IsProtected && args.RequestorChannelID == nil {
		return errors.Err("requestor channel id is required to list protected comments")
	}
	if args.IsProtected && args.ClaimID != nil && args.RequestorChannelID != nil {
		hasAccess, err := HasAccessToProtectedContent(*args.ClaimID, *args.RequestorChannelID)
		if err != nil {
			return err
		}
		if !hasAccess {
			return errors.Err("channel does not have permissions to comment on this claim")
		}
	}

	var cachedReply *commentapi.ListResponse
	if listingOwnComments {
		err := list(r, args, reply)
		if err != nil {
			return err
		}
		cachedReply = reply
	} else {
		key, err := args.Key()
		if err != nil {
			return err
		}
		item, err := commentListCache.Fetch(key, 15*time.Second, func() (interface{}, error) {
			err := list(r, args, reply)
			if err != nil {
				return nil, err
			}
			return reply, nil
		})
		if err != nil {
			return err
		}
		var ok bool
		cachedReply, ok = item.Value().(*commentapi.ListResponse)
		if !ok {
			return errors.Prefix("could not convert item to ListResponse: ", err)
		}
	}

	reply.PageSize = cachedReply.PageSize
	reply.Page = cachedReply.Page
	reply.Items = cachedReply.Items
	reply.TotalItems = cachedReply.TotalItems
	reply.HasHiddenComments = cachedReply.HasHiddenComments
	reply.HasProtectedComments = cachedReply.HasProtectedComments
	reply.TotalFilteredItems = cachedReply.TotalFilteredItems
	reply.TotalPages = cachedReply.TotalPages
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
		} else if sort == commentapi.NewestNoPins {
			queryMods = append(queryMods, qm.OrderBy(m.CommentColumns.Timestamp+" DESC"))
		}
	} else {
		queryMods = append(queryMods, qm.OrderBy(m.CommentColumns.IsPinned+" DESC, "+m.CommentColumns.Timestamp+" DESC"))
	}

	return queryMods
}

func checkCommentsEnabled(channelName, ChannelID string) (*m.Channel, error) {
	if channelName != "" && ChannelID != "" {
		creatorChannel, err := helper.FindOrCreateChannel(ChannelID, channelName)
		if err != nil {
			return nil, err
		}
		settings, err := creatorChannel.CreatorChannelCreatorSettings().One(db.RO)
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

func getItems(comments m.CommentSlice, creatorChannel *m.Channel, skipBlocked bool) ([]commentapi.CommentItem, int64, error) {
	var items []commentapi.CommentItem
	var blockedCommentCnt int64
	var alreadyInSet = map[string]bool{}
Comments:
	for _, comment := range comments {
		if comment.R != nil && comment.R.Channel != nil && comment.R.Channel.R != nil {
			blockedFrom := comment.R.Channel.R.BlockedChannelBlockedEntries
			if len(blockedFrom) > 0 && !skipBlocked {
				channel, err := lbry.SDK.GetSigningChannelForClaim(comment.LbryClaimID)
				if err != nil {
					//cannot find claim commented on in SDK, ignore, nil channel by default
				}
				if channel != nil {
					for _, entry := range blockedFrom {
						if creatorChannel != nil && creatorChannel.BlockedListID.Valid {
							if creatorChannel.BlockedListID == entry.BlockedListID {
								if !entry.Expiry.Valid || (entry.Expiry.Valid && time.Since(entry.Expiry.Time) < time.Duration(0)) {
									blockedCommentCnt++
									continue Comments
								}
							}
						}
						if entry.UniversallyBlocked.Bool || entry.CreatorChannelID.String == channel.ClaimID {
							if !entry.Expiry.Valid || (entry.Expiry.Valid && time.Since(entry.Expiry.Time) < time.Duration(0)) {
								blockedCommentCnt++
								continue Comments
							}
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
					replies, err := comment.ParentComments().Count(db.RO)
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
