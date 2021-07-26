package comments

import (
	"math"
	"net/http"

	"github.com/lbryio/commentron/commentapi"
	"github.com/lbryio/commentron/db"
	m "github.com/lbryio/commentron/model"
	"github.com/lbryio/commentron/server/lbry"

	"github.com/lbryio/lbry.go/extras/util"
	"github.com/lbryio/lbry.go/v2/extras/errors"

	"github.com/btcsuite/btcutil"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

func superChatList(_ *http.Request, args *commentapi.SuperListArgs, reply *commentapi.SuperListResponse) error {
	args.ApplyDefaults()
	loadChannels := qm.Load("Channel.BlockedChannelBlockedEntries")
	filterIsHidden := m.CommentWhere.IsHidden.EQ(null.BoolFrom(true))
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

	if args.AuthorClaimID != nil {
		getCommentsQuery = append(getCommentsQuery, filterAuthorClaimID)
		hasHiddenCommentsQuery = append(hasHiddenCommentsQuery, filterAuthorClaimID)
		totalCommentsQuery = append(totalCommentsQuery, filterAuthorClaimID)
		totalSuperChatAmountQuery = append(totalSuperChatAmountQuery, filterAuthorClaimID)
	}

	if args.ClaimID != nil {
		getCommentsQuery = append(getCommentsQuery, filterClaimID)
		hasHiddenCommentsQuery = append(hasHiddenCommentsQuery, filterClaimID)
		totalCommentsQuery = append(totalCommentsQuery, filterClaimID)
		totalSuperChatAmountQuery = append(totalSuperChatAmountQuery, filterClaimID)
	}

	if args.TopLevel {
		getCommentsQuery = append(getCommentsQuery, filterTopLevel)
		hasHiddenCommentsQuery = append(hasHiddenCommentsQuery, filterTopLevel)
		totalCommentsQuery = append(totalCommentsQuery, filterTopLevel)
		totalSuperChatAmountQuery = append(totalSuperChatAmountQuery, filterTopLevel)
	}

	if args.ParentID != nil {
		getCommentsQuery = append(getCommentsQuery, filterParent)
		hasHiddenCommentsQuery = append(hasHiddenCommentsQuery, filterParent)
		totalCommentsQuery = append(totalCommentsQuery, filterParent)
		totalSuperChatAmountQuery = append(totalSuperChatAmountQuery, filterParent)
	}

	if args.SuperChatsAmount > 0 {
		getCommentsQuery = append(getCommentsQuery, filterSuperChats)
		hasHiddenCommentsQuery = append(hasHiddenCommentsQuery, filterSuperChats)
		totalCommentsQuery = append(totalCommentsQuery, filterSuperChats)
		totalSuperChatAmountQuery = append(totalSuperChatAmountQuery, filterSuperChats)
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

	items, blockedCommentCnt, err := getItems(comments, creatorChannel)

	totalItems = totalItems - blockedCommentCnt
	reply.Items = items
	reply.Page = args.Page
	reply.PageSize = args.PageSize
	reply.TotalItems = totalItems
	reply.TotalPages = int(math.Ceil(float64(totalItems) / float64(args.PageSize)))
	reply.HasHiddenComments = hasHiddenComments
	reply.TotalAmount = btcutil.Amount(superChatAmount.Uint64).ToBTC()

	return nil
}
