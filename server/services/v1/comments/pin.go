package comments

import (
	"net/http"

	"github.com/OdyseeTeam/commentron/commentapi"
	"github.com/OdyseeTeam/commentron/db"
	"github.com/OdyseeTeam/commentron/helper"
	"github.com/OdyseeTeam/commentron/model"
	"github.com/OdyseeTeam/commentron/server/lbry"
	"github.com/OdyseeTeam/commentron/sockety"

	"github.com/OdyseeTeam/sockety/socketyapi"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/lbryio/lbry.go/v2/extras/errors"
)

func pin(_ *http.Request, args *commentapi.PinArgs) (commentapi.CommentItem, error) {
	var item commentapi.CommentItem
	comment, err := model.Comments(model.CommentWhere.CommentID.EQ(args.CommentID)).One(db.RO)
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

	channel, err := helper.FindOrCreateChannel(args.ChannelID, args.ChannelName)
	if err != nil {
		return item, errors.Err(err)
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
	err = comment.Update(db.RW, boil.Infer())
	if err != nil {
		return item, errors.Err(err)
	}

	item = populateItem(comment, channel, 0)

	pushClaimID := item.ClaimID
	if item.IsProtected {
		pushClaimID = helper.ReverseString(item.ClaimID)
	}
	go sockety.SendNotification(socketyapi.SendNotificationArgs{
		Service: socketyapi.Commentron,
		Type:    "pinned",
		IDs:     []string{pushClaimID, "pins"},
		Data:    map[string]interface{}{"comment": item},
	})
	return item, nil
}
