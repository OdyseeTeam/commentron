package reactions

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/volatiletech/null"

	"github.com/lbryio/commentron/server/lbry"

	"github.com/lbryio/commentron/db"
	"github.com/lbryio/commentron/model"

	"github.com/lbryio/errors.go"
	"github.com/lbryio/lbry.go/extras/api"
	"github.com/lbryio/lbry.go/v2/extras/util"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

// ReactArgs are the arguments passed to comment.Abandon RPC call
type ReactArgs struct {
	CommentIDs  string  `json:"comment_ids"`
	Signature   string  `json:"signature"`
	SigningTS   string  `json:"signing_ts"`
	Remove      bool    `json:"remove"`
	ClearTypes  string  `json:"clear_types"`
	Type        string  `json:"type"`
	ChannelID   *string `json:"channel_id"`
	ChannelName *string `json:"channel_name"`
}

// ReactResponse the response to the abandon call
type ReactResponse struct {
	reactions
}

// React creates/updates a reaction to a comment
func react(_ *http.Request, args *ReactArgs, reply *ReactResponse) error {

	comments, err := model.Comments(qm.WhereIn(model.CommentColumns.CommentID+" IN ?", util.StringSplitArg(args.CommentIDs, ",")...)).AllG()
	if err != nil {
		return errors.Err(err)
	}
	if len(comments) == 0 {
		return errors.Err("could not find comments(s)")
	}
	var commentIDs []interface{}
	for _, p := range comments {
		commentIDs = append(commentIDs, p.CommentID)
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
	err = lbry.ValidateSignature(util.StrFromPtr(args.ChannelID), args.Signature, args.SigningTS, util.StrFromPtr(args.ChannelName))
	if err != nil {
		return errors.Prefix("could not authenticate channel signature: %s", err)
	}
	modifiedReactions, err := updateReactions(channel, args, commentIDs, comments)
	if err != nil {
		return errors.Err(err)
	}
	reply.reactions = modifiedReactions
	return nil
}
func updateReactions(channel *model.Channel, args *ReactArgs, commentIDs []interface{}, comments model.CommentSlice) (reactions, error) {
	var modifiedReactions = newReactions(strings.Split(args.CommentIDs, ","), &args.Type)
	err := db.WithTx(nil, func(tx boil.Transactor) error {
		if len(args.ClearTypes) > 0 {
			typeNames := util.StringSplitArg(args.ClearTypes, ",")
			reactionTypes, err := model.ReactionTypes(qm.WhereIn(model.ReactionTypeColumns.Name+" IN ?", typeNames...)).All(tx)
			if err != nil {
				return errors.Err(err)
			}
			if len(reactionTypes) > 0 {
				var typesToClear []interface{}
				for _, rt := range reactionTypes {
					typesToClear = append(typesToClear, rt.ID)
				}
				err = channel.Reactions(
					qm.Where(model.ReactionColumns.ChannelID+"=?", channel.ClaimID),
					qm.WhereIn(model.ReactionColumns.ReactionTypeID+" IN ?", typesToClear...)).DeleteAll(tx)
				if err != nil {
					return errors.Err(err)
				}
			}
		}

		reactionType, err := model.ReactionTypes(model.ReactionTypeWhere.Name.EQ(args.Type)).One(tx)
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
			reactionType := &model.ReactionType{Name: args.Type}
			err = reactionType.Insert(tx, boil.Infer())
		}
		if err != nil {
			return errors.Err(err)
		}
		if args.Remove {
			existingReactions, err := channel.Reactions(
				qm.WhereIn(model.ReactionColumns.CommentID+"=?", commentIDs...),
				qm.Where(model.ReactionColumns.ReactionTypeID+"=?", reactionType.ID),
				qm.Load("Comment")).All(tx)
			if err != nil {
				return errors.Err(err)
			}
			if len(existingReactions) == 0 {
				return api.StatusError{Err: errors.Err("there are no reactions for the claim(s) to remove"), Status: http.StatusBadRequest}
			}
			for _, r := range existingReactions {
				modifiedReactions[r.R.Comment.CommentID].Add(args.Type)
			}
			err = existingReactions.DeleteAll(tx)
			return errors.Err(err)
		}
		for _, p := range comments {
			newReaction := &model.Reaction{ChannelID: null.StringFrom(channel.ClaimID), CommentID: p.CommentID, ReactionTypeID: reactionType.ID}
			err := newReaction.Insert(tx, boil.Infer())
			if err != nil {
				if strings.Contains(err.Error(), "Duplicate entry") {
					return api.StatusError{Err: errors.Err("reaction already acknowledged!"), Status: http.StatusBadRequest}
				}
				return errors.Err(err)
			}
			modifiedReactions[p.CommentID].Add(reactionType.Name)
		}
		return nil
	})
	if err != nil {
		return nil, errors.Err(err)
	}
	return modifiedReactions, nil
}
