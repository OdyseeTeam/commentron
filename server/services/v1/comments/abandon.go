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
	"github.com/lbryio/lbry.go/v2/extras/api"
	"github.com/lbryio/lbry.go/v2/extras/errors"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func abandon(args *commentapi.AbandonArgs) (*commentapi.CommentItem, error) {
	loadCommenter := qm.Load(model.CommentRels.Channel)
	comment, err := model.Comments(model.CommentWhere.CommentID.EQ(args.CommentID), loadCommenter).One(db.RO)
	if err != nil {
		return nil, errors.Err(err)
	}
	var commenterChannel *model.Channel
	var modChannel *model.Channel
	var creatorChannel *model.Channel
	if comment.R.Channel == nil {
		return nil, errors.Err("channel id '%s' does not have a channel record", comment.ChannelID.String)
	}
	commenterChannel = comment.R.Channel
	// Handle anonymous content where there's no channel associated
	if args.CreatorChannelID == "" && args.CreatorChannelName == "" {
		// If the content is anonymous, set the modChannel to a default or system channel
		if args.ModChannelName != "" && args.ModChannelID != "" {
			modChannel, _, err = helper.GetModerator(args.ModChannelID, args.ModChannelName, args.ModChannelID, args.ModChannelName)
			if err != nil {
				return nil, err
			}
		}
	}

	// Old versions of desktop app will allow for just creator channel info to be sent for creators to
	// delete comments and mod channel info is newer addition and would not be sent, so we cannot assume
	// it will be sent with request.
	if args.CreatorChannelID != "" && args.CreatorChannelName != "" {
		modChannelID := args.CreatorChannelID
		modChannelName := args.CreatorChannelName
		if args.ModChannelName != "" && args.ModChannelID != "" {
			modChannelID = args.ModChannelID
			modChannelName = args.ModChannelName
		}
		modChannel, creatorChannel, err = helper.GetModerator(modChannelID, modChannelName, args.CreatorChannelID, args.CreatorChannelName)
		if err != nil {
			return nil, err
		}
		content, err := lbry.SDK.GetClaim(comment.LbryClaimID)
		if err != nil {
			return nil, errors.Err(err)
		}
		signingChannelClaimID := content.ClaimID
		if content.SigningChannel != nil {
			signingChannelClaimID = content.SigningChannel.ClaimID
		}
		if signingChannelClaimID != creatorChannel.ClaimID {
			return nil, api.StatusError{Err: errors.Err("you do not have creator authorizations to remove this comment on %s", comment.LbryClaimID), Status: http.StatusBadRequest}
		}
	}

	// if there are neither a mod nor a delegated mod, then we verify if the commenter is the creator which means you're trying to delete your own comment
	if modChannel == nil {
		modChannel = commenterChannel
	}

	err = lbry.ValidateSignatureAndTS(modChannel.ClaimID, args.Signature, args.SigningTS, args.CommentID)
	if err != nil {
		return nil, err
	}
	item := populateItem(comment, commenterChannel, 0)
	err = comment.Delete(db.RW, false)
	if err != nil {
		return nil, errors.Err(err)
	}

	pushClaimID := item.ClaimID
	if item.IsProtected {
		pushClaimID = helper.ReverseString(item.ClaimID)
	}

	go sockety.SendNotification(socketyapi.SendNotificationArgs{
		Service: socketyapi.Commentron,
		Type:    "removed",
		IDs:     []string{pushClaimID, "comments", "deleted"},
		Data:    map[string]interface{}{"comment": item},
	})

	return &item, nil

}
