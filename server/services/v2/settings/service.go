package settings

import (
	"net/http"

	"github.com/OdyseeTeam/commentron/commentapi"
	"github.com/OdyseeTeam/commentron/db"
	"github.com/OdyseeTeam/commentron/helper"
	"github.com/OdyseeTeam/commentron/model"
	"github.com/OdyseeTeam/commentron/server/lbry"
	"github.com/OdyseeTeam/commentron/sockety"

	"github.com/btcsuite/btcutil"
	"github.com/lbryio/lbry.go/v2/extras/errors"
	"github.com/lbryio/lbry.go/v2/extras/util"
	"github.com/lbryio/sockety/socketyapi"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

// Service is the service struct defined for the comment package for rpc service "moderation.*"
type Service struct{}

// List returns the list of user settings applicable to them for the creator to manage
func (s *Service) List(r *http.Request, args *commentapi.ListSettingsArgs, reply *commentapi.ListSettingsResponse) error {
	creatorChannel, err := helper.FindOrCreateChannel(args.ChannelID, args.ChannelName)
	if err != nil {
		return errors.Err(err)
	}
	err = lbry.ValidateSignatureAndTS(creatorChannel.ClaimID, args.Signature, args.SigningTS, args.ChannelName)
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
	err = lbry.ValidateSignatureAndTS(creatorChannel.ClaimID, args.Signature, args.SigningTS, args.ChannelName)
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

	if args.ChatOverlay != nil {
		settings.ChatOverlay = *args.ChatOverlay
	}

	if args.ChatOverlayPosition != nil {
		settings.ChatOverlayPosition = *args.ChatOverlayPosition
	}

	if args.ChatRemoveComment != nil {
		chatRemoveComment := int64(*args.ChatRemoveComment)
		settings.ChatRemoveComment = chatRemoveComment
	}

	if args.StickerOverlay != nil {
		settings.StickerOverlay = *args.StickerOverlay
	}

	if args.StickerOverlayKeep != nil {
		settings.StickerOverlayKeep = *args.StickerOverlayKeep
	}

	if args.StickerOverlayRemove != nil {
		stickerOverlayRemove := int64(*args.StickerOverlayRemove)
		settings.StickerOverlayRemove = stickerOverlayRemove
	}

	if args.ViewercountOverlay != nil {
		settings.ViewercountOverlay = *args.ViewercountOverlay
	}

	if args.ViewercountChatBot != nil {
		settings.ViewercountChatBot = *args.ViewercountChatBot
	}

	if args.ViewercountOverlayPosition != nil {
		settings.ViewercountOverlayPosition = *args.ViewercountOverlayPosition
	}

	if args.TipgoalOverlay != nil {
		settings.TipgoalOverlay = *args.TipgoalOverlay
	}

	if args.TipgoalOverlayPosition != nil {
		settings.TipgoalOverlayPosition = *args.TipgoalOverlayPosition
	}

	if args.TipgoalPreviousDonations != nil {
		settings.TipgoalPreviousDonations = *args.TipgoalPreviousDonations
	}

	if args.TipgoalCurrency != nil {
		settings.TipgoalCurrency = *args.TipgoalCurrency
	}

	if args.TipgoalAmount != nil {
		tipGoalAmt := int64(*args.TipgoalAmount)
		settings.TipgoalAmount = tipGoalAmt
	}

	if args.TimeSinceFirstComment != nil {
		timeSinceFirstCmt := int64(*args.TimeSinceFirstComment)
		settings.TimeSinceFirstComment.SetValid(timeSinceFirstCmt)
		if *args.TimeSinceFirstComment == 0 {
			settings.TimeSinceFirstComment.Valid = false
		}
	}

	if args.PrivateShowProtected != nil {
		settings.PrivateShowProtected = *args.PrivateShowProtected
	}

	if args.PublicShowProtected != nil {
		settings.PublicShowProtected = *args.PublicShowProtected
	}

	if args.LivestreamChatMembersOnly != nil {
		settings.LivestreamChatMembersOnly = *args.LivestreamChatMembersOnly
	}
	if args.CommentsMembersOnly != nil {
		settings.CommentsMembersOnly = *args.CommentsMembersOnly
	}

	err = settings.Update(db.RW, boil.Infer())
	if err != nil {
		return errors.Err(err)
	}

	membersOnlyChatToggled := args.ActiveClaimID != nil && (args.LivestreamChatMembersOnly != nil || args.CommentsMembersOnly != nil)

	if membersOnlyChatToggled {
		var data map[string]interface{}
		if args.LivestreamChatMembersOnly != nil {
			data = map[string]interface{}{"LivestreamChatMembersOnly": *args.LivestreamChatMembersOnly}
		} else {
			data = map[string]interface{}{"CommentsMembersOnly": *args.CommentsMembersOnly}
		}

		go sockety.SendNotification(socketyapi.SendNotificationArgs{
			Service: socketyapi.Commentron,
			Type:    "setting",
			IDs:     []string{*args.ActiveClaimID, "settings"},
			Data:    data,
		})
	}

	if !args.FeaturedChannels.IsZero() {
		settings.FeaturedChannels = args.FeaturedChannels
	}
	if !args.HomepageSettings.IsZero() {
		settings.HomepageSettings = args.HomepageSettings
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
	if settings.TimeSinceFirstComment.Valid && authorized {
		tsfc := uint64(settings.TimeSinceFirstComment.Int64)
		reply.TimeSinceFirstComment = &tsfc
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
	reply.ChatOverlay = &settings.ChatOverlay
	reply.ChatOverlayPosition = &settings.ChatOverlayPosition
	chatRemoveComment := uint64(settings.ChatRemoveComment)
	reply.ChatRemoveComment = &chatRemoveComment
	reply.StickerOverlay = &settings.StickerOverlay
	reply.StickerOverlayKeep = &settings.StickerOverlayKeep
	stickerOverlayRemove := uint64(settings.StickerOverlayRemove)
	reply.StickerOverlayRemove = &stickerOverlayRemove
	reply.ViewercountOverlay = &settings.ViewercountOverlay
	reply.ViewercountOverlayPosition = &settings.ViewercountOverlayPosition
	reply.ViewercountChatBot = &settings.ViewercountChatBot
	reply.TipgoalOverlay = &settings.TipgoalOverlay
	reply.TipgoalOverlayPosition = &settings.TipgoalOverlayPosition
	reply.TipgoalPreviousDonations = &settings.TipgoalPreviousDonations
	tipGoalAmount := uint64(settings.TipgoalAmount)
	reply.TipgoalAmount = &tipGoalAmount
	reply.TipgoalCurrency = &settings.TipgoalCurrency
	reply.PublicShowProtected = &settings.PublicShowProtected
	reply.PrivateShowProtected = &settings.PrivateShowProtected
	reply.LivestreamChatMembersOnly = &settings.LivestreamChatMembersOnly
	reply.CommentsMembersOnly = &settings.CommentsMembersOnly
	reply.FeaturedChannels = settings.FeaturedChannels
	reply.HomepageSettings = settings.HomepageSettings
}
