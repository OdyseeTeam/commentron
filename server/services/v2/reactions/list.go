package reactions

import (
	"database/sql"
	"net/http"
	"strings"
	"time"

	"github.com/lbryio/commentron/commentapi"
	"github.com/lbryio/commentron/db"
	"github.com/lbryio/commentron/model"
	"github.com/lbryio/commentron/server/lbry"

	"github.com/lbryio/errors.go"
	"github.com/lbryio/lbry.go/v2/extras/util"

	"github.com/karlseguin/ccache"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

func list(_ *http.Request, args *commentapi.ReactionListArgs, reply *commentapi.ReactionListResponse) error {
	comments, err := model.Comments(qm.WhereIn(model.CommentColumns.CommentID+" IN ?", util.StringSplitArg(args.CommentIDs, ",")...)).All(db.RO)
	if err != nil {
		return errors.Err(err)
	}
	if len(comments) == 0 {
		return errors.Err("could not find comment(s)")
	}
	var commentIDs []interface{}
	for _, p := range comments {
		commentIDs = append(commentIDs, p.CommentID)
	}
	var myfilters = []qm.QueryMod{qm.WhereIn(model.ReactionColumns.CommentID+" IN ?", commentIDs...),
		qm.Load("ReactionType"),
		qm.Load("Comment")}
	var allfilters = []qm.QueryMod{qm.WhereIn(model.ReactionColumns.CommentID+" IN ?", commentIDs...),
		qm.Load("ReactionType"),
		qm.Load("Comment")}
	if args.Types != nil {
		typeNames := util.StringSplitArg(util.StrFromPtr(args.Types), ",")
		types, err := model.ReactionTypes(qm.WhereIn(model.ReactionTypeColumns.Name+" IN ?", typeNames...)).All(db.RO)
		if err != nil {
			return errors.Err(err)
		}
		var typeIDs []interface{}
		for _, t := range types {
			typeIDs = append(typeIDs, t.ID)
		}
		if len(typeIDs) == 0 {
			return errors.Err("none of the types %s are in use in commentron", util.StrFromPtr(args.Types))
		}
		myfilters = append(myfilters, qm.WhereIn(model.ReactionColumns.ReactionTypeID+" IN ?", typeIDs...))
		allfilters = append(allfilters, qm.WhereIn(model.ReactionColumns.ReactionTypeID+" IN ?", typeIDs...))
	}
	channel, err := model.Channels(model.ChannelWhere.ClaimID.EQ(util.StrFromPtr(args.ChannelID))).One(db.RO)
	if errors.Is(err, sql.ErrNoRows) {
		channel = &model.Channel{
			ClaimID: util.StrFromPtr(args.ChannelID),
			Name:    util.StrFromPtr(args.ChannelName),
		}
		err = nil
		err := channel.Insert(db.RW, boil.Infer())
		if err != nil {
			return errors.Err(err)
		}
	}
	if err != nil {
		return errors.Err(err)
	}
	var userReactions commentapi.Reactions
	if args.ChannelName != nil {
		chanErr := lbry.ValidateSignature(util.StrFromPtr(args.ChannelID), args.Signature, args.SigningTS, util.StrFromPtr(args.ChannelName))
		if chanErr == nil {
			allfilters = append(allfilters, qm.Where(model.ReactionColumns.ChannelID+" != ?", channel.ClaimID))
			reactionlist, err := channel.Reactions(myfilters...).All(db.RO)
			if err != nil {
				return errors.Err(err)
			}
			userReactions = newReactions(strings.Split(args.CommentIDs, ","), args.Types)
			for _, r := range reactionlist {
				addTo(userReactions[r.CommentID], r.R.ReactionType.Name)
			}
		}
	}

	reactionlist, err := model.Reactions(allfilters...).All(db.RO)
	if err != nil {
		return errors.Err(err)
	}
	var othersReactions = newReactions(strings.Split(args.CommentIDs, ","), args.Types)
	for _, r := range reactionlist {
		addTo(othersReactions[r.CommentID], r.R.ReactionType.Name)
	}
	reply.MyReactions = userReactions
	reply.OthersReactions = othersReactions
	return nil
}

func newReactions(commentIDs []string, types *string) commentapi.Reactions {
	var reactionTypes []string
	if types == nil {
		rts, err := getReactionTypes()
		if err == nil {
			for _, r := range rts {
				reactionTypes = append(reactionTypes, r.Name)
			}
		}
	} else {
		reactionTypes = strings.Split(*types, ",")
	}
	r := make(map[string]commentapi.CommentReaction, len(commentIDs))
	for _, c := range commentIDs {
		r[c] = make(map[string]int)
		for _, t := range reactionTypes {
			r[c][t] = 0
		}
	}
	return r
}

func addTo(c commentapi.CommentReaction, reactionType string) {
	curr, ok := c[reactionType]
	if !ok {
		c[reactionType] = 0
	}
	c[reactionType] = curr + 1
}

var reactionTypeCache = ccache.New(ccache.Configure().MaxSize(100))

func getReactionTypes() (model.ReactionTypeSlice, error) {
	v, err := reactionTypeCache.Fetch("all", 30*time.Minute, func() (interface{}, error) {
		rts, err := model.ReactionTypes().All(db.RO)
		if err != nil {
			return nil, errors.Err(err)
		}
		return rts, nil
	})
	if err != nil {
		return nil, err
	}
	slice, ok := v.Value().(model.ReactionTypeSlice)
	if !ok {
		return nil, errors.Err("could not convert cached value to ReactionTypeSlice")
	}
	return slice, nil
}
