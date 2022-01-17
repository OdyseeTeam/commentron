package commentapi

import (
	"net/http"
	"time"

	"github.com/lbryio/commentron/validator"
	"github.com/lbryio/lbry.go/v2/extras/api"
	"github.com/lbryio/lbry.go/v2/extras/errors"
	v "github.com/lbryio/ozzo-validation"
)

// BlockArgs Arguments to block identities from commenting for both publisher and moderators
type BlockArgs struct {
	ModAuthorization

	//Offender being blocked
	BlockedChannelID   string `json:"blocked_channel_id"`
	BlockedChannelName string `json:"blocked_channel_name"`
	// ID of comment to remove as part of this block
	OffendingCommentID string `json:"offending_comment_id"`
	// Blocks identity from comment universally, requires Admin rights on commentron instance
	BlockAll bool `json:"block_all"`
	// Measured in seconds for the amount of time a channel is blocked for.
	TimeOut uint64 `json:"time_out"`
	// If true will delete all comments of the offender, requires Admin rights on commentron for universal delete
	DeleteAll bool `json:"delete_all"`
}

// Validate validates the data in the list args
func (b BlockArgs) Validate() api.StatusError {
	err := v.ValidateStruct(&b,
		v.Field(&b.BlockedChannelID, validator.ClaimID, v.Required),
		v.Field(&b.BlockedChannelName, v.Required),
		v.Field(&b.ModChannelID, validator.ClaimID, v.Required),
		v.Field(&b.ModChannelName, v.Required),
	)
	if err != nil {
		return api.StatusError{Err: errors.Err(err), Status: http.StatusBadRequest}
	}
	return api.StatusError{}
}

// BlockResponse for the moderation.Block rpc call
type BlockResponse struct {
	DeletedCommentIDs []string `json:"deleted_comment_ids"`
	BannedChannelID   string   `json:"banned_channel_id"`
	AllBlocked        bool     `json:"all_blocked"`
	//Publisher banned from if not universally banned
	BannedFrom *string `json:"banned_from"`
}

// AmIArgs Arguments to check whether a user is a moderator or not
type AmIArgs struct {
	Authorization
}

// AmIResponse for the moderation.AmI rpc call
type AmIResponse struct {
	ChannelName        string            `json:"channel_name"`
	ChannelID          string            `json:"channel_id"`
	Type               string            `json:"type"`
	AuthorizedChannels map[string]string `json:"authorized_channels"`
}

// UnBlockArgs Arguments to un-block identities from commenting for both publisher and moderators
type UnBlockArgs struct {
	ModAuthorization

	//Offender being unblocked
	UnBlockedChannelID   string `json:"un_blocked_channel_id"`
	UnBlockedChannelName string `json:"un_blocked_channel_name"`
	// Unblocks identity from commenting universally, requires Admin rights on commentron instance
	GlobalUnBlock bool `json:"global_un_block"`
}

// UnBlockResponse for the moderation.UnBlock rpc call
type UnBlockResponse struct {
	UnBlockedChannelID string `json:"un_blocked_channel_id"`
	GlobalUnBlock      bool   `json:"global_un_block"`
	//Publisher ban removed from if not universally unblocked
	UnBlockedFrom *string `json:"un_blocked_from"`
}

// BlockedListArgs Arguments to block identities from commenting for both publisher and moderators
type BlockedListArgs struct {
	ModAuthorization
}

// BlockedListResponse for the moderation.Block rpc call
type BlockedListResponse struct {
	BlockedChannels          []BlockedChannel `json:"blocked_channels"`
	DelegatedBlockedChannels []BlockedChannel `json:"delegated_blocked_channels"`
	GloballyBlockedChannels  []BlockedChannel `json:"globally_blocked_channels"`
}

// BlockedChannel contains information about the blockee blocked by the creator
type BlockedChannel struct {
	BlockedChannelID   string `json:"blocked_channel_id"`
	BlockedChannelName string `json:"blocked_channel_name"`
	//In cases of moderation delegation this could be "other than" the creator
	BlockedByChannelID   string        `json:"blocked_by_channel_id"`
	BlockedByChannelName string        `json:"blocked_by_channel_name"`
	BlockedAt            time.Time     `json:"blocked_at"`
	BlockedFor           time.Duration `json:"banned_for,omitempty"`
	BlcokRemaining       time.Duration `json:"ban_remaining"`
}

// AddDelegateArgs Arguments to delagate moderation to another channel for your channel.
type AddDelegateArgs struct {
	Authorization

	//This is for backwards compatibility, Authorization parameters should be used, not these!
	CreatorChannelID   string `json:"creator_channel_id"`
	CreatorChannelName string `json:"creator_channel_name"`

	//Who is being delegated authority?
	ModChannelID   string `json:"mod_channel_id"`
	ModChannelName string `json:"mod_channel_name"`
}

// Validate validates the data in the AddDelegate args
func (ad *AddDelegateArgs) Validate() api.StatusError {
	if ad.CreatorChannelID != "" {
		ad.ChannelID = ad.CreatorChannelID
	}
	if ad.CreatorChannelName != "" {
		ad.ChannelName = ad.CreatorChannelName
	}
	err := v.ValidateStruct(ad,
		v.Field(&ad.ChannelID, validator.ClaimID, v.Required),
		v.Field(&ad.ChannelName, v.Required),
	)
	if err != nil {
		return api.StatusError{Err: errors.Err(err), Status: http.StatusBadRequest}
	}

	return api.StatusError{}
}

// RemoveDelegateArgs Arguments to remove a delegated moderator.
type RemoveDelegateArgs struct {
	Authorization

	//This is for backwards compatibility, Authorization parameters should be used, not these!
	CreatorChannelID   string `json:"creator_channel_id"`
	CreatorChannelName string `json:"creator_channel_name"`

	//Who is being removed from delegated authority?
	ModChannelID   string `json:"mod_channel_id"`
	ModChannelName string `json:"mod_channel_name"`
}

// Validate validates the data in the RemoveDelegate args
func (rd *RemoveDelegateArgs) Validate() api.StatusError {
	if rd.CreatorChannelID != "" {
		rd.ChannelID = rd.CreatorChannelID
	}
	if rd.CreatorChannelName != "" {
		rd.ChannelName = rd.CreatorChannelName
	}
	err := v.ValidateStruct(rd,
		v.Field(&rd.ChannelID, validator.ClaimID, v.Required),
		v.Field(&rd.ChannelName, v.Required),
	)
	if err != nil {
		return api.StatusError{Err: errors.Err(err), Status: http.StatusBadRequest}
	}

	return api.StatusError{}
}

// ListDelegatesArgs Arguments to list delegates
type ListDelegatesArgs struct {
	Authorization

	CreatorChannelID   string `json:"creator_channel_id"`
	CreatorChannelName string `json:"creator_channel_name"`
}

// Validate validates the data in the ListDelegates args
func (ld *ListDelegatesArgs) Validate() api.StatusError {
	if ld.CreatorChannelID != "" {
		ld.ChannelID = ld.CreatorChannelID
	}
	if ld.CreatorChannelName != "" {
		ld.ChannelName = ld.CreatorChannelName
	}
	err := v.ValidateStruct(ld,
		v.Field(&ld.ChannelID, validator.ClaimID, v.Required),
		v.Field(&ld.ChannelName, v.Required),
	)
	if err != nil {
		return api.StatusError{Err: errors.Err(err), Status: http.StatusBadRequest}
	}

	return api.StatusError{}
}

// ListDelegateResponse response for modifying the delegates
type ListDelegateResponse struct {
	Delegates []Delegate
}

// Delegate a particular channel thats delegated moderation capabilities
type Delegate struct {
	ChannelID   string `json:"channel_id"`
	ChannelName string `json:"channel_name"`
}
