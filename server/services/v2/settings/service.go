package settings

import (
	"net/http"

	"github.com/lbryio/lbry.go/extras/util"

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

// List returns the list of user settings applicable to them for the creator to manage
func (s *Service) List(r *http.Request, args *commentapi.ListSettingsArgs, reply *commentapi.ListSettingsResponse) error {
	creatorChannel, err := helper.FindOrCreateChannel(args.ChannelID, args.ChannelName)
	if err != nil {
		return errors.Err(err)
	}
	err = lbry.ValidateSignature(creatorChannel.ClaimID, args.Signature, args.SigningTS, args.ChannelName)
	if err != nil {
		return err
	}
	authorized := true

	settings, err := helper.FindOrCreateSettings(creatorChannel)
	if err != nil {
		return err
	}

	applySettingsToReply(settings, reply, authorized)

	return nil
}

// Get returns the list of creator settings for users
func (s *Service) Get(r *http.Request, args *commentapi.ListSettingsArgs, reply *commentapi.ListSettingsResponse) error {
	creatorChannel, err := helper.FindOrCreateChannel(args.ChannelID, args.ChannelName)
	if err != nil {
		return errors.Err(err)
	}
	authorized := false

	settings, err := helper.FindOrCreateSettings(creatorChannel)
	if err != nil {
		return err
	}

	applySettingsToReply(settings, reply, authorized)

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
	authorized := true

	settings, err := helper.FindOrCreateSettings(creatorChannel)
	if err != nil {
		return err
	}

	if args.CommentsEnabled != nil {
		settings.CommentsEnabled.SetValid(*args.CommentsEnabled)
	}

	if args.SlowModeMinGap != nil {
		settings.SlowModeMinGap.SetValid(*args.SlowModeMinGap)
		if *args.SlowModeMinGap == 0 {
			settings.SlowModeMinGap.Valid = false
		}
	}

	if args.MinTipAmountSuperChat != nil {
		lbc, err := btcutil.NewAmount(*args.MinTipAmountSuperChat)
		if err != nil {
			return errors.Err(err)
		}
		settings.MinTipAmountSuperChat.SetValid(uint64(lbc.ToUnit(btcutil.AmountSatoshi)))
		if lbc == 0.0 {
			settings.MinTipAmountSuperChat.Valid = false
		}
	}

	if args.MinTipAmountComment != nil {
		lbc, err := btcutil.NewAmount(*args.MinTipAmountComment)
		if err != nil {
			return errors.Err(err)
		}
		settings.MinTipAmountComment.SetValid(uint64(lbc.ToUnit(btcutil.AmountSatoshi)))
		if lbc == 0.0 {
			settings.MinTipAmountComment.Valid = false
		}
	}

	if args.CurseJarAmount != nil { // Coming with Appeal process
		settings.CurseJarAmount.SetValid(*args.CurseJarAmount)
		if *args.CurseJarAmount == 0.0 {
			settings.CurseJarAmount.Valid = false
		}
	}

	if args.FiltersEnabled != nil { // Future feature to be developed
		settings.IsFiltersEnabled.SetValid(*args.FiltersEnabled)
	}

	err = settings.Update(db.RW, boil.Infer())
	if err != nil {
		return errors.Err(err)
	}

	applySettingsToReply(settings, reply, authorized)

	return nil
}

func applySettingsToReply(settings *model.CreatorSetting, reply *commentapi.ListSettingsResponse, authorized bool) {
	// RETURN ONLY INF AUTHORIZED TO SEE
	if settings.MutedWords.Valid && authorized {
		reply.Words = &settings.MutedWords.String
	}
	if settings.IsFiltersEnabled.Valid && authorized {
		reply.FiltersEnabled = &settings.IsFiltersEnabled.Bool
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
	if settings.CurseJarAmount.Valid {
		reply.CurseJarAmount = util.PtrToUint64(settings.CurseJarAmount.Uint64)
	}

}
