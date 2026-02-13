package comments

import (
	"math"
	"net/http"
	"time"

	"github.com/OdyseeTeam/commentron/commentapi"
	"github.com/OdyseeTeam/commentron/db"
	"github.com/OdyseeTeam/commentron/helper"
	m "github.com/OdyseeTeam/commentron/model"
	"github.com/OdyseeTeam/commentron/server/lbry"

	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/btcsuite/btcutil"
	"github.com/karlseguin/ccache/v2"
	"github.com/lbryio/lbry.go/v2/extras/api"
	"github.com/lbryio/lbry.go/v2/extras/errors"
	"github.com/lbryio/lbry.go/v2/extras/util"
)

func superChatList(_ *http.Request, args *commentapi.SuperListArgs, reply *commentapi.SuperListResponse) error {
	args.ApplyDefaults()
	isListingOwnSuperChats := args.AuthorClaimID != nil
	actualIsProtected := args.IsProtected
	if !isListingOwnSuperChats {
		actualIsProtected, err := IsProtectedContent(*args.ClaimID)
		if err != nil {
			return err
		}
		if actualIsProtected != args.IsProtected {
			return errors.Err("mismatch in is_protected")
		}
	} else {
		if args.RequestorChannelID == nil {
			return errors.Err("requestor channel id is required to list own superchats")
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
			return api.StatusError{Err: errors.Err("you can only view your superchats, not others"), Status: http.StatusBadRequest}
		}
	}

	loadChannels := qm.Load("Channel.BlockedChannelBlockedEntries")
	filterIsHidden := m.CommentWhere.IsHidden.EQ(null.BoolFrom(true))
	filterIsProtected := m.CommentWhere.IsProtected.EQ(true)
	filterClaimID := m.CommentWhere.LbryClaimID.EQ(util.StrFromPtr(args.ClaimID))
	filterAuthorClaimID := m.CommentWhere.ChannelID.EQ(null.StringFromPtr(args.AuthorClaimID))
	filterTopLevel := m.CommentWhere.ParentID.IsNull()
	filterParent := m.CommentWhere.ParentID.EQ(null.StringFrom(util.StrFromPtr(args.ParentID)))
	filterSuperChats := m.CommentWhere.Amount.GTE(null.Uint64From(uint64(args.SuperChatsAmount)))

	totalCommentsQuery := make([]qm.QueryMod, 0)
	totalSuperChatAmountQuery := []qm.QueryMod{qm.Select(`SUM(` + m.CommentColumns.Amount + `)`)}
	offset := (args.Page - 1) * args.PageSize
	getCommentsQuery := []qm.QueryMod{loadChannels, qm.Offset(offset), qm.Limit(args.PageSize), qm.OrderBy(m.CommentColumns.IsFiat + " DESC, " + m.CommentColumns.Amount + " DESC, " + m.CommentColumns.Timestamp + " DESC")}
	hasHiddenCommentsQuery := []qm.QueryMod{filterIsHidden, qm.Limit(1)}
	HasProtectedCommentsQuery := []qm.QueryMod{filterIsProtected, qm.Limit(1)}

	if isListingOwnSuperChats {
		getCommentsQuery = append(getCommentsQuery, filterAuthorClaimID)
		hasHiddenCommentsQuery = append(hasHiddenCommentsQuery, filterAuthorClaimID)
		HasProtectedCommentsQuery = append(HasProtectedCommentsQuery, filterAuthorClaimID)
		totalCommentsQuery = append(totalCommentsQuery, filterAuthorClaimID)
		totalSuperChatAmountQuery = append(totalSuperChatAmountQuery, filterAuthorClaimID)
	}

	if args.ClaimID != nil {
		getCommentsQuery = append(getCommentsQuery, filterClaimID)
		hasHiddenCommentsQuery = append(hasHiddenCommentsQuery, filterClaimID)
		HasProtectedCommentsQuery = append(HasProtectedCommentsQuery, filterClaimID)
		totalCommentsQuery = append(totalCommentsQuery, filterClaimID)
		totalSuperChatAmountQuery = append(totalSuperChatAmountQuery, filterClaimID)
	}

	if args.TopLevel {
		getCommentsQuery = append(getCommentsQuery, filterTopLevel)
		hasHiddenCommentsQuery = append(hasHiddenCommentsQuery, filterTopLevel)
		HasProtectedCommentsQuery = append(HasProtectedCommentsQuery, filterTopLevel)
		totalCommentsQuery = append(totalCommentsQuery, filterTopLevel)
		totalSuperChatAmountQuery = append(totalSuperChatAmountQuery, filterTopLevel)
	}

	if args.ParentID != nil {
		getCommentsQuery = append(getCommentsQuery, filterParent)
		hasHiddenCommentsQuery = append(hasHiddenCommentsQuery, filterParent)
		HasProtectedCommentsQuery = append(HasProtectedCommentsQuery, filterParent)
		totalCommentsQuery = append(totalCommentsQuery, filterParent)
		totalSuperChatAmountQuery = append(totalSuperChatAmountQuery, filterParent)
	}

	if args.SuperChatsAmount > 0 {
		getCommentsQuery = append(getCommentsQuery, filterSuperChats)
		hasHiddenCommentsQuery = append(hasHiddenCommentsQuery, filterSuperChats)
		HasProtectedCommentsQuery = append(HasProtectedCommentsQuery, filterSuperChats)
		totalCommentsQuery = append(totalCommentsQuery, filterSuperChats)
		totalSuperChatAmountQuery = append(totalSuperChatAmountQuery, filterSuperChats)
	}
	if !isListingOwnSuperChats {
		totalCommentsQuery = append(totalCommentsQuery, m.CommentWhere.IsProtected.EQ(actualIsProtected))
	}

	var superChatAmount null.Uint64
	result := m.Comments(totalSuperChatAmountQuery...).QueryRow(db.RO)
	err := result.Scan(&superChatAmount)
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

	channelClaim, err := lbry.SDK.GetSigningChannelForClaim(util.StrFromPtr(args.ClaimID))
	if err != nil {
		return errors.Err(err)
	}
	var creatorChannel *m.Channel
	if channelClaim != nil {
		creatorChannel, err = m.Channels(m.ChannelWhere.ClaimID.EQ(channelClaim.ClaimID)).One(db.RO)
		if err != nil {
			return errors.Err(err)
		}
	}

	// if listing own comments, show all including blocked ones
	skipBlocked := isListingOwnSuperChats

	items, blockedCommentCnt, err := getItems(comments, creatorChannel, skipBlocked)
	if err != nil {
		return errors.Err(err)
	}

	totalItems = totalItems - blockedCommentCnt
	reply.Items = items
	reply.Page = args.Page
	reply.PageSize = args.PageSize
	reply.TotalItems = totalItems
	reply.TotalPages = int(math.Ceil(float64(totalItems) / float64(args.PageSize)))
	reply.HasHiddenComments = hasHiddenComments
	reply.HasProtectedComments = HasProtectedComments
	reply.TotalAmount = btcutil.Amount(superChatAmount.Uint64).ToBTC()

	return nil
}

var superChatListCache = ccache.New(ccache.Configure().GetsPerPromote(1).MaxSize(100000))

func getCachedSuperChatList(r *http.Request, args *commentapi.SuperListArgs, reply *commentapi.SuperListResponse) error {
	listingOwnSuperChats := args.AuthorClaimID != nil
	if args.IsProtected && args.RequestorChannelID == nil {
		return errors.Err("requestor channel id is required to list protected superchats")
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

	var cachedReply *commentapi.SuperListResponse
	if listingOwnSuperChats {
		err := superChatList(r, args, reply)
		if err != nil {
			return err
		}
		cachedReply = reply
	} else {
		key, err := args.Key()
		if err != nil {
			return err
		}
		item, err := superChatListCache.Fetch(key, 1*time.Minute, func() (interface{}, error) {
			err := superChatList(r, args, reply)
			if err != nil {
				return nil, err
			}
			return reply, nil
		})
		if err != nil {
			return err
		}
		var ok bool
		cachedReply, ok = item.Value().(*commentapi.SuperListResponse)
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
	reply.TotalPages = cachedReply.TotalPages
	reply.TotalAmount = cachedReply.TotalAmount
	return nil
}
