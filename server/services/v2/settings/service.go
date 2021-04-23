package settings

import (
	"net/http"

	"github.com/btcsuite/btcutil"

	"github.com/lbryio/commentron/commentapi"
	"github.com/lbryio/commentron/helper"
	"github.com/lbryio/commentron/server/lbry"

	"github.com/lbryio/errors.go"
)

// Service is the service struct defined for the comment package for rpc service "moderation.*"
type Service struct{}

// List returns the list of user settings applicable to them.
func (s *Service) List(r *http.Request, args *commentapi.ListSettingsArgs, reply *commentapi.ListSettingsResponse) error {
	creatorChannel, err := helper.FindOrCreateChannel(args.ChannelID, args.ChannelName)
	if err != nil {
		return errors.Err(err)
	}
	err = lbry.ValidateSignature(creatorChannel.ClaimID, args.Signature, args.SigningTS, args.ChannelName)
	if err != nil {
		return err
	}

	settings, err := helper.FindOrCreateSettings(creatorChannel)
	if err != nil {
		return err
	}

	reply.Words = settings.MutedWords.String
	reply.CommentsEnabled = settings.CommentsEnabled.Bool
	reply.MinTipAmountComment = btcutil.Amount(settings.MinTipAmountComment.Uint64).ToBTC()
	reply.MinTipAmountSuperChat = btcutil.Amount(settings.MinTipAmountSuperChat.Uint64).ToBTC()
	reply.SlowModeMinGap = settings.SlowModeMinGap.Uint64

	return nil
}
