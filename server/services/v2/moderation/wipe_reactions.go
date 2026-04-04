package moderation

import (
	"net/http"

	"github.com/OdyseeTeam/commentron/commentapi"
	"github.com/OdyseeTeam/commentron/db"
	"github.com/OdyseeTeam/commentron/model"
	"github.com/OdyseeTeam/commentron/server/auth"

	"github.com/aarondl/null/v8"
	"github.com/lbryio/lbry.go/v2/extras/api"
	"github.com/lbryio/lbry.go/v2/extras/errors"
)

func wipeReactions(r *http.Request, args *commentapi.WipeReactionsArgs, reply *commentapi.WipeReactionsResponse) error {
	modChannel, _, _, err := auth.ModAuthenticate(r, &args.ModAuthorization)
	if err != nil {
		return err
	}

	isMod, err := modChannel.ModChannelModerators().Exists(db.RO)
	if err != nil {
		return errors.Err(err)
	}
	if !isMod {
		return api.StatusError{Err: errors.Err("cannot wipe reactions without admin privileges"), Status: http.StatusForbidden}
	}

	reactions, err := model.Reactions(
		model.ReactionWhere.ChannelID.EQ(null.StringFrom(args.TargetChannelID)),
	).All(db.RO)
	if err != nil {
		return errors.Err(err)
	}

	reply.TargetChannelID = args.TargetChannelID
	reply.DeletedReactionCount = uint64(len(reactions))
	if len(reactions) == 0 {
		return nil
	}

	err = reactions.DeleteAll(db.RW)
	if err != nil {
		return errors.Err(err)
	}

	return nil
}
