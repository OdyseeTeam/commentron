package settings

import (
	"net/http"

	"github.com/volatiletech/sqlboiler/boil"

	"github.com/lbryio/commentron/model"

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

	applySettingsToReply(settings, reply)

	return nil
}

// Update updates the different settings if passed.
func (s *Service) Update(r *http.Request, args *commentapi.UpdateSettingsArgs, reply *commentapi.ListSettingsResponse) error {
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

	if args.CommentsEnabled != nil {
		settings.CommentsEnabled.SetValid(*args.CommentsEnabled)
	}

	if args.SlowModeMinGap != nil {
		settings.SlowModeMinGap.SetValid(*args.SlowModeMinGap)
	}

	if args.MinTipAmountSuperChat != nil {
		settings.MinTipAmountSuperChat.SetValid(*args.MinTipAmountSuperChat)
	}

	if args.MinTipAmountComment != nil {
		settings.MinTipAmountComment.SetValid(*args.MinTipAmountComment)
	}

	err = settings.UpdateG(boil.Infer())
	if err != nil {
		return errors.Err(err)
	}

	applySettingsToReply(settings, reply)

	return nil
}

func applySettingsToReply(settings *model.CreatorSetting, reply *commentapi.ListSettingsResponse) {
	reply.Words = settings.MutedWords.String
	reply.CommentsEnabled = settings.CommentsEnabled.Bool
	reply.MinTipAmountComment = btcutil.Amount(settings.MinTipAmountComment.Uint64).ToBTC()
	reply.MinTipAmountSuperChat = btcutil.Amount(settings.MinTipAmountSuperChat.Uint64).ToBTC()
	reply.SlowModeMinGap = settings.SlowModeMinGap.Uint64
}
