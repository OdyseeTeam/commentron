package comments

import (
	"database/sql"
	"net/http"

	"github.com/lbryio/commentron/commentapi"
	"github.com/lbryio/commentron/model"
	"github.com/lbryio/commentron/server/lbry"
	"github.com/lbryio/commentron/server/websocket"

	"github.com/lbryio/lbry.go/v2/extras/errors"
	"github.com/volatiletech/sqlboiler/boil"
)

func pin(_ *http.Request, args *commentapi.PinArgs) (commentapi.CommentItem, error) {
	var item commentapi.CommentItem
	comment, err := model.Comments(model.CommentWhere.CommentID.EQ(args.CommentID)).OneG()
	if err != nil {
		return item, errors.Err(err)
	}

	claim, err := lbry.SDK.GetClaim(comment.LbryClaimID)
	if err != nil {
		return item, errors.Err(err)
	}
	if claim == nil {
		return item, errors.Err("could not resolve claim from comment")
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
	claimChannel := claim.SigningChannel
	if claimChannel == nil {
		if claim.ValueType == "channel" {
			claimChannel = claim
		} else {
			return item, errors.Err("claim does not have a signing channel")
		}
	}

	err = lbry.ValidateSignatureFromClaim(claimChannel, args.Signature, args.SigningTS, args.CommentID)
	if err != nil {
		return item, err
	}
	comment.IsPinned = !args.Remove
	err = comment.UpdateG(boil.Infer())
	if err != nil {
		return item, errors.Err(err)
	}

	item = populateItem(comment, channel, 0)
	go pushPinnedItem(item, comment.LbryClaimID)
	return item, nil
}

func pushPinnedItem(item commentapi.CommentItem, claimID string) {
	websocket.PushTo(&websocket.PushNotification{
		Type: "pinned",
		Data: map[string]interface{}{"comment": item},
	}, claimID)

	go sendMessage(item, claimID)

}
