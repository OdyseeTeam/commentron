package comments

import (
	"database/sql"
	"math"
	"net/http"
	"time"

	"github.com/OdyseeTeam/commentron/commentapi"
	"github.com/OdyseeTeam/commentron/db"
	"github.com/OdyseeTeam/commentron/helper"
	m "github.com/OdyseeTeam/commentron/model"
	"github.com/OdyseeTeam/commentron/server/lbry"

	"github.com/lbryio/lbry.go/v2/extras/api"
	"github.com/lbryio/lbry.go/v2/extras/errors"
	"github.com/lbryio/lbry.go/v2/extras/util"

	"github.com/karlseguin/ccache/v2"
	"github.com/sirupsen/logrus"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"golang.org/x/sync/singleflight"
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
			return api.StatusError{Err: errors.Err("mismatch in is_protected"), Status: http.StatusBadRequest}
		}
	} else {
		if args.RequestorChannelID == nil {
			return api.StatusError{Err: errors.Err("requestor channel id is required to list own comments"), Status: http.StatusBadRequest}
		}
		ownerChannel, err := helper.FindOrCreateChannel(*args.RequestorChannelID, args.RequestorChannelName)
		if err != nil {
			return err
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
		return api.StatusError{Err: errors.Err("requestor channel id is required to list protected comments"), Status: http.StatusBadRequest}
	}
	if args.IsProtected && args.ClaimID != nil && args.RequestorChannelID != nil {
		hasAccess, err := HasAccessToProtectedContent(*args.ClaimID, *args.RequestorChannelID)
		if err != nil {
			return err
		}
		if !hasAccess {
			return api.StatusError{Err: errors.Err("channel does not have permissions to comment on this claim"), Status: http.StatusForbidden}
		}
		commenterChannel, err := helper.FindOrCreateChannel(*args.RequestorChannelID, args.RequestorChannelName)
		if err != nil {
			return err
		}

		err = lbry.ValidateSignatureAndTS(commenterChannel.ClaimID, args.Signature, args.SigningTS, args.RequestorChannelName)
		if err != nil {
			return err
		}
		if commenterChannel.ClaimID != *args.RequestorChannelID {
			return api.StatusError{Err: errors.Err("channel mismatch, someone trying to spoof"), Status: http.StatusForbidden}
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

var repliesCountCache = ccache.New(ccache.Configure().MaxSize(100000))
var sf singleflight.Group

func getItems(comments m.CommentSlice, creatorChannel *m.Channel, skipBlocked bool) ([]commentapi.CommentItem, int64, error) {
	var items []commentapi.CommentItem
	var blockedCommentCnt int64
	var alreadyInSet = make(map[string]bool)

comments:
	for _, comment := range comments {
		if comment.R == nil || comment.R.Channel == nil || comment.R.Channel.R == nil {
			continue
		}

		blockedFrom := comment.R.Channel.R.BlockedChannelBlockedEntries
		if len(blockedFrom) > 0 && !skipBlocked {
			channel, err := lbry.SDK.GetSigningChannelForClaim(comment.LbryClaimID)
			if err != nil {
				// Cannot find claim commented on in SDK, ignore, nil channel by default
			} else if channel != nil {
				for _, entry := range blockedFrom {
					if creatorChannel != nil && creatorChannel.BlockedListID.Valid && creatorChannel.BlockedListID == entry.BlockedListID {
						if !entry.Expiry.Valid || entry.Expiry.Time.After(time.Now()) {
							blockedCommentCnt++
							continue comments
						}
					}
					if entry.UniversallyBlocked.Bool || entry.CreatorChannelID.String == channel.ClaimID {
						if !entry.Expiry.Valid || entry.Expiry.Time.After(time.Now()) {
							blockedCommentCnt++
							continue comments
						}
					}
				}
			}
		}

		channel := comment.R.Channel
		if channel == nil || channel.Name == "" {
			continue
		}

		if alreadyInSet[comment.CommentID] {
			continue
		}

		repliesCachedCount := repliesCountCache.Get(comment.CommentID)
		if repliesCachedCount != nil && !repliesCachedCount.Expired() {
			alreadyInSet[comment.CommentID] = true
			item := populateItem(comment, channel, int(repliesCachedCount.Value().(int64)))
			err := applyModStatus(&item, comment.ChannelID.String, comment.LbryClaimID)
			if err != nil {
				return items, blockedCommentCnt, err
			}
			items = append(items, item)
			continue
		}

		val, err, _ := sf.Do(comment.CommentID, func() (interface{}, error) {
			replies, err := comment.ParentComments().Count(db.RO)
			if err != nil && errors.Is(err, sql.ErrNoRows) {
				return nil, errors.Err(err)
			}
			repliesCountCache.Set(comment.CommentID, replies, 30*time.Second)
			return replies, nil
		})
		if err != nil {
			return items, blockedCommentCnt, err
		}
		replies := val.(int64)
		item := populateItem(comment, channel, int(replies))
		err := applyModStatus(&item, comment.ChannelID.String, comment.LbryClaimID)
		if err != nil {
			return items, blockedCommentCnt, err
		}
		items = append(items, item)
	}

	return items, blockedCommentCnt, nil
}
