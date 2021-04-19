package comments

import (
	"math"
	"net/http"
	"sort"

	"github.com/lbryio/commentron/commentapi"
	m "github.com/lbryio/commentron/model"
	"github.com/lbryio/lbry.go/extras/util"
	"github.com/lbryio/lbry.go/v2/extras/errors"
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

	if args.SuperChatsAmount > 0 {
		getCommentsQuery = append(getCommentsQuery, filterSuperChats)
		hasHiddenCommentsQuery = append(hasHiddenCommentsQuery, filterSuperChats)
		totalCommentsQuery = append(totalCommentsQuery, filterSuperChats)
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

	sort.SliceStable(items, func(i, j int) bool {
		return items[j].SupportAmount <= items[i].SupportAmount
	})

	totalItems = totalItems - blockedCommentCnt
	reply.Items = items
	reply.Page = args.Page
	reply.PageSize = args.PageSize
	reply.TotalItems = totalItems
	reply.TotalPages = int(math.Ceil(float64(totalItems) / float64(args.PageSize)))
	reply.HasHiddenComments = hasHiddenComments

	return nil
}
