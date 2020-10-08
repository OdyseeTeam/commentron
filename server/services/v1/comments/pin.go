package comments

import (
	"database/sql"
	"net/http"

	"github.com/lbryio/commentron/commentapi"
	"github.com/lbryio/commentron/model"
	"github.com/lbryio/commentron/server/lbry"

	"github.com/lbryio/lbry.go/v2/extras/errors"

	"github.com/volatiletech/sqlboiler/boil"
)

func pin(_ *http.Request, args *commentapi.PinArgs) (commentapi.CommentItem, error) {
	var item commentapi.CommentItem
	comment, err := model.Comments(model.CommentWhere.CommentID.EQ(args.CommentID)).OneG()
	if err != nil {
		return item, errors.Err(err)
	}
	channel, err := model.Channels(model.ChannelWhere.ClaimID.EQ(args.ChannelID)).OneG()
	if errors.Is(err, sql.ErrNoRows) {
		channel = &model.Channel{
			ClaimID: args.ChannelID,
			Name:    args.ChannelName,
		}
		err = nil
		err := channel.InsertG(boil.Infer())
		if err != nil {
			return item, errors.Err(err)
		}
	}
	err = lbry.ValidateSignature(args.ChannelID, args.Signature, args.SigningTS, args.CommentID)
	if err != nil {
		return item, err
	}
	comment.IsPinned = !args.RemovePin
	err = comment.UpdateG(boil.Infer())
	if err != nil {
		return item, errors.Err(err)
	}
	return populateItem(comment, channel), nil
}
