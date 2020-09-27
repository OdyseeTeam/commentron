package reactions

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/lbryio/commentron/commentapi"
	"github.com/lbryio/commentron/model"
	"github.com/lbryio/commentron/server/lbry"

	"github.com/lbryio/errors.go"
	"github.com/lbryio/lbry.go/v2/extras/util"

	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

func list(_ *http.Request, args *commentapi.ReactionListArgs, reply *commentapi.ReactionListResponse) error {
	comments, err := model.Comments(qm.WhereIn(model.CommentColumns.CommentID+" IN ?", util.StringSplitArg(args.CommentIDs, ",")...)).AllG()
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
		types, err := model.ReactionTypes(qm.WhereIn(model.ReactionTypeColumns.Name+" IN ?", typeNames...)).AllG()
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
	channel, err := model.Channels(model.ChannelWhere.ClaimID.EQ(util.StrFromPtr(args.ChannelID))).OneG()
	if errors.Is(err, sql.ErrNoRows) {
		channel = &model.Channel{
			ClaimID: util.StrFromPtr(args.ChannelID),
			Name:    util.StrFromPtr(args.ChannelName),
		}
		err = nil
		err := channel.InsertG(boil.Infer())
		if err != nil {
			return errors.Err(err)
		}
	}
	var userReactions commentapi.Reactions
	if args.ChannelName != nil {
		chanErr := lbry.ValidateSignature(util.StrFromPtr(args.ChannelID), args.Signature, args.SigningTS, util.StrFromPtr(args.ChannelName))
		if chanErr == nil {
			allfilters = append(allfilters, qm.Where(model.ReactionColumns.ChannelID+" != ?", channel.ClaimID))
			reactionlist, err := channel.Reactions(myfilters...).AllG()
			if err != nil {
				return errors.Err(err)
			}
			userReactions = newReactions(strings.Split(args.CommentIDs, ","), args.Types)
			for _, r := range reactionlist {
				addTo(userReactions[r.CommentID], r.R.ReactionType.Name)
			}
		}
	}

	reactionlist, err := model.Reactions(allfilters...).AllG()
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
		rts, err := model.ReactionTypes().AllG()
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
