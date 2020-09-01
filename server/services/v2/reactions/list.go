package reactions

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/lbryio/commentron/server/lbry"
	"github.com/volatiletech/sqlboiler/boil"

	"github.com/lbryio/commentron/model"
	"github.com/lbryio/errors.go"
	"github.com/lbryio/lbry.go/v2/extras/util"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

// ListArgs are the arguments passed to comment.Abandon RPC call
type ListArgs struct {
	CommentIDs  string `json:"comment_ids"`
	Signature   string `json:"signature"`
	SigningTS   string `json:"signing_ts"`
	Types       *string
	ChannelID   *string `json:"channel_id"`
	ChannelName *string `json:"channel_name"`
}

// ListResponse the response to the abandon call
type ListResponse struct {
	MyReactions     reactions `json:"my_reactions,omitempty"`
	OthersReactions reactions `json:"others_reactions"`
}

func list(_ *http.Request, args *ListArgs, reply *ListResponse) error {
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
	var myfilters = []qm.QueryMod{qm.WhereIn(model.ReactionColumns.ChannelID+" IN ?", commentIDs...),
		qm.Load("ReactionType"),
		qm.Load("Comment")}
	var allfilters = []qm.QueryMod{qm.WhereIn(model.ReactionColumns.ClaimID+" IN ?", commentIDs...),
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
	chanErr := lbry.ValidateSignature(util.StrFromPtr(args.ChannelID), args.Signature, args.SigningTS, util.StrFromPtr(args.ChannelName))
	var userReactions reactions
	if chanErr == nil {
		allfilters = append(allfilters, qm.Where(model.ReactionColumns.ChannelID+" != ?", channel.ClaimID))
		reactionlist, err := channel.Reactions(myfilters...).AllG()
		if err != nil {
			return errors.Err(err)
		}
		userReactions = newReactions(strings.Split(args.CommentIDs, ","), args.Types)
		for _, r := range reactionlist {
			userReactions[r.R.Channel.ClaimID].Add(r.R.ReactionType.Name)
		}
	}

	reactionlist, err := model.Reactions(allfilters...).AllG()
	if err != nil {
		return errors.Err(err)
	}
	var othersReactions = newReactions(strings.Split(args.CommentIDs, ","), args.Types)
	for _, r := range reactionlist {
		othersReactions[r.R.Channel.ClaimID].Add(r.R.ReactionType.Name)
	}
	reply.MyReactions = userReactions
	reply.OthersReactions = othersReactions
	return nil
}

type reactions map[string]commentReaction

type commentReaction map[string]int

func (c commentReaction) Add(reactionType string) {
	curr, ok := c[reactionType]
	if !ok {
		c[reactionType] = 0
	}
	c[reactionType] = curr + 1
}

func newReactions(commentIDs []string, types *string) reactions {
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
	r := make(map[string]commentReaction, len(commentIDs))
	for _, c := range commentIDs {
		r[c] = make(map[string]int)
		for _, t := range reactionTypes {
			r[c][t] = 0
		}
	}
	return r
}
