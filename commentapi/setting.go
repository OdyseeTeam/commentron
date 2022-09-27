package commentapi

import (
	"net/http"

	"github.com/OdyseeTeam/commentron/validator"

	"github.com/lbryio/lbry.go/v2/extras/api"
	"github.com/lbryio/lbry.go/v2/extras/errors"
	v "github.com/lbryio/ozzo-validation"
)

// ListSettingsArgs arguments passed to settings.List api
type ListSettingsArgs struct {
	Authorization
}

// ListSettingsResponse returns all the settings for creator/user
type ListSettingsResponse struct {
	// CSV list of containing words to block comment on content
	Words                      *string  `json:"words,omitempty"`
	CommentsEnabled            *bool    `json:"comments_enabled"`
	MinTipAmountComment        *float64 `json:"min_tip_amount_comment"`
	MinTipAmountSuperChat      *float64 `json:"min_tip_amount_super_chat"`
	SlowModeMinGap             *uint64  `json:"slow_mode_min_gap"`
	CurseJarAmount             *uint64  `json:"curse_jar_amount"`
	FiltersEnabled             *bool    `json:"filters_enabled,omitempty"`
	ChatOverlay                *bool    `json:"chat_overlay"`
	ChatOverlayPosition        *string  `json:"chat_overlay_position"`
	ChatRemoveComment          *uint64  `json:"chat_remove_comment"`
	StickerOverlay             *bool    `json:"sticker_overlay"`
	StickerOverlayKeep         *bool    `json:"sticker_overlay_keep"`
	StickerOverlayRemove       *uint64  `json:"sticker_overlay_remove"`
	ViewercountOverlay         *bool    `json:"viewercount_overlay"`
	ViewercountOverlayPosition *string  `json:"viewercount_overlay_position"`
	ViewercountChatBot         *bool    `json:"viewercount_chat_bot"`
	TipgoalOverlay             *bool    `json:"tipgoal_overlay"`
	TipgoalAmount              *uint64  `json:"tipgoal_amount"`
	TipgoalOverlayPosition     *string  `json:"tipgoal_overlay_position"`
	TipgoalPreviousDonations   *bool    `json:"tipgoal_previous_donations"`
	TipgoalCurrency            *string  `json:"tipgoal_currency"`
	TimeSinceFirstComment      *uint64  `json:"time_since_first_comment"`
	PublicShowProtected        *bool    `json:"public_show_protected"`
	PrivateShowProtected       *bool    `json:"private_show_protected"`
	LivestreamChatMembersOnly  *bool    `json:"livestream_chat_members_only"`
	CommentsMembersOnly        *bool    `json:"comments_members_only"`
}

// UpdateSettingsArgs arguments for different settings that could be set
type UpdateSettingsArgs struct {
	Authorization
	CommentsEnabled            *bool    `json:"comments_enabled"`
	MinTipAmountComment        *float64 `json:"min_tip_amount_comment"`
	MinTipAmountSuperChat      *float64 `json:"min_tip_amount_super_chat"`
	SlowModeMinGap             *uint64  `json:"slow_mode_min_gap"`
	CurseJarAmount             *uint64  `json:"curse_jar_amount"`
	FiltersEnabled             *bool    `json:"filters_enabled"`
	ChatOverlay                *bool    `json:"chat_overlay"`
	ChatOverlayPosition        *string  `json:"chat_overlay_position"`
	ChatRemoveComment          *uint64  `json:"chat_remove_comment"`
	StickerOverlay             *bool    `json:"sticker_overlay"`
	StickerOverlayKeep         *bool    `json:"sticker_overlay_keep"`
	StickerOverlayRemove       *uint64  `json:"sticker_overlay_remove"`
	ViewercountOverlay         *bool    `json:"viewercount_overlay"`
	ViewercountOverlayPosition *string  `json:"viewercount_overlay_position"`
	ViewercountChatBot         *bool    `json:"viewercount_chat_bot"`
	TipgoalOverlay             *bool    `json:"tipgoal_overlay"`
	TipgoalAmount              *uint64  `json:"tipgoal_amount"`
	TipgoalOverlayPosition     *string  `json:"tipgoal_overlay_position"`
	TipgoalPreviousDonations   *bool    `json:"tipgoal_previous_donations"`
	TipgoalCurrency            *string  `json:"tipgoal_currency"`
	// Minutes since first comment when users are allowed to comment on your content/livestream
	TimeSinceFirstComment     *uint64 `json:"time_since_first_comment"`
	PrivateShowProtected      *bool   `json:"private_show_protected"`
	PublicShowProtected       *bool   `json:"public_show_protected"`
	LivestreamChatMembersOnly *bool   `json:"livestream_chat_members_only"`
	CommentsMembersOnly       *bool   `json:"comments_members_only"`
	ActiveClaimID             *string `json:"active_claim_id"`
}

// Validate validates the data in the args
func (u UpdateSettingsArgs) Validate() api.StatusError {
	err := v.ValidateStruct(&u,
		v.Field(&u.ChatOverlayPosition, v.In("Left", "Right")),
		v.Field(&u.ViewercountOverlayPosition, v.In("Top Left", "Top Center", "Top Right", "Bottom Left", "Bottom Center", "Bottom Right")),
		v.Field(&u.TipgoalOverlayPosition, v.In("Top", "Bottom")),
		v.Field(&u.TipgoalCurrency, v.In("LBC", "FIAT")),
		v.Field(&u.ActiveClaimID, validator.ClaimID),
	)
	if err != nil {
		return api.StatusError{Err: errors.Err(err), Status: http.StatusBadRequest}
	}
	return api.StatusError{}
}

// BlockWordArgs arguments passed to settings.BlockWord. Appends to list
type BlockWordArgs struct {
	Authorization
	// CSV list of containing words to block comment on content
	Words string `json:"words"`
}

// Validate validates the data in the args
func (b BlockWordArgs) Validate() api.StatusError {
	err := v.ValidateStruct(&b,
		v.Field(&b.ChannelID, validator.ClaimID, v.Required),
		v.Field(&b.ChannelName, v.Required),
		v.Field(&b.Words),
		v.Field(&b.Signature, v.Required),
		v.Field(&b.SigningTS, v.Required),
	)
	if err != nil {
		return api.StatusError{Err: errors.Err(err), Status: http.StatusBadRequest}
	}
	return api.StatusError{}
}

// BlockWordRespose result from BlockWord,UnBlockWord, ListBlockedWords. Lists the words added/removed or all.
type BlockWordRespose struct {
	//If added to list, removed from list, or list all
	WordList  []string `json:"word_list"`
	Signature string   `json:"signature"`
	SigningTS string   `json:"signing_ts"`
}

// UnBlockWordArgs arguments passed to settings.UnBlockWord. Removes if exists
type UnBlockWordArgs struct {
	Authorization
	// CSV list of containing words to block comment on content
	Words string `json:"words"`
}

// Validate validates the data in the args
func (b UnBlockWordArgs) Validate() api.StatusError {
	err := v.ValidateStruct(&b,
		v.Field(&b.ChannelID, validator.ClaimID, v.Required),
		v.Field(&b.ChannelName, v.Required),
		v.Field(&b.Words),
		v.Field(&b.Signature, v.Required),
		v.Field(&b.SigningTS, v.Required),
	)
	if err != nil {
		return api.StatusError{Err: errors.Err(err), Status: http.StatusBadRequest}
	}
	return api.StatusError{}
}

// ListBlockedWordsArgs lists all the blocked words for the channel
type ListBlockedWordsArgs struct {
	Authorization
}

// Validate validates the data in the args
func (b ListBlockedWordsArgs) Validate() api.StatusError {
	err := v.ValidateStruct(&b,
		v.Field(&b.ChannelID, validator.ClaimID, v.Required),
		v.Field(&b.ChannelName, v.Required),
		v.Field(&b.Signature, v.Required),
		v.Field(&b.SigningTS, v.Required),
	)
	if err != nil {
		return api.StatusError{Err: errors.Err(err), Status: http.StatusBadRequest}
	}
	return api.StatusError{}
}
