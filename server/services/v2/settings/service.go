package settings

import (
	"net/http"

	"github.com/lbryio/commentron/commentapi"
	"github.com/lbryio/commentron/db"
	"github.com/lbryio/commentron/helper"
	"github.com/lbryio/commentron/model"
	"github.com/lbryio/commentron/server/lbry"

	"github.com/lbryio/errors.go"

	"github.com/btcsuite/btcutil"
	"github.com/volatiletech/sqlboiler/boil"
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
		lbc, err := btcutil.NewAmount(*args.MinTipAmountSuperChat)
		if err != nil {
			return errors.Err(err)
		}
		settings.MinTipAmountSuperChat.SetValid(uint64(lbc.ToBTC()))
	}

	if args.MinTipAmountComment != nil {
		lbc, err := btcutil.NewAmount(*args.MinTipAmountComment)
		if err != nil {
			return errors.Err(err)
		}
		settings.MinTipAmountComment.SetValid(uint64(lbc.ToBTC()))
	}

	if args.CurseJarAmount != nil { // Coming with Appeal process
		settings.CurseJarAmount.SetValid(*args.CurseJarAmount)
	}

	if args.FiltersEnabled != nil { // Future feature to be developed
		settings.IsFiltersEnabled.SetValid(*args.FiltersEnabled)
	}

	err = settings.Update(db.RW, boil.Infer())
	if err != nil {
		return errors.Err(err)
	}

	applySettingsToReply(settings, reply)

	return nil
}

func applySettingsToReply(settings *model.CreatorSetting, reply *commentapi.ListSettingsResponse) {
	if settings.MutedWords.Valid {
		reply.Words = &settings.MutedWords.String
	}
	if settings.CommentsEnabled.Valid {
		reply.CommentsEnabled = &settings.CommentsEnabled.Bool
	}
	if settings.MinTipAmountComment.Valid {
		minTipAmount := btcutil.Amount(settings.MinTipAmountComment.Uint64).ToBTC()
		reply.MinTipAmountComment = &minTipAmount
	}
	if settings.MinTipAmountSuperChat.Valid {
		minTipAmountSuperChat := btcutil.Amount(settings.MinTipAmountSuperChat.Uint64).ToBTC()
		reply.MinTipAmountSuperChat = &minTipAmountSuperChat
	}
	if settings.SlowModeMinGap.Valid {
		reply.SlowModeMinGap = &settings.SlowModeMinGap.Uint64
	}
}
