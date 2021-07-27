package commentapi

import (
	"net/http"

	"github.com/lbryio/commentron/validator"

	"github.com/lbryio/lbry.go/extras/api"
	"github.com/lbryio/lbry.go/v2/extras/errors"
	v "github.com/lbryio/ozzo-validation"
)

// ListSettingsArgs arguments passed to settings.List api
type ListSettingsArgs struct {
	ChannelName string `json:"channel_name"`
	ChannelID   string `json:"channel_id"`
	Signature   string `json:"signature"`
	SigningTS   string `json:"signing_ts"`
}

// ListSettingsResponse returns all the settings for creator/user
type ListSettingsResponse struct {
	// CSV list of containing words to block comment on content
	Words                 *string  `json:"words,omitempty"`
	CommentsEnabled       *bool    `json:"comments_enabled"`
	MinTipAmountComment   *float64 `json:"min_tip_amount_comment"`
	MinTipAmountSuperChat *float64 `json:"min_tip_amount_super_chat"`
	SlowModeMinGap        *uint64  `json:"slow_mode_min_gap"`
	CurseJarAmount        *uint64  `json:"curse_jar_amount"`
	FiltersEnabled        *bool    `json:"filters_enabled,omitempty"`
}

// UpdateSettingsArgs arguments for different settings that could be set
type UpdateSettingsArgs struct {
	Authorization
	CommentsEnabled       *bool    `json:"comments_enabled"`
	MinTipAmountComment   *float64 `json:"min_tip_amount_comment"`
	MinTipAmountSuperChat *float64 `json:"min_tip_amount_super_chat"`
	SlowModeMinGap        *uint64  `json:"slow_mode_min_gap"`
	CurseJarAmount        *uint64  `json:"curse_jar_amount"`
	FiltersEnabled        *bool    `json:"filters_enabled"`
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
