package commentapi

import "time"

// BlockArgs Arguments to block identities from commenting for both publisher and moderators
type BlockArgs struct {
	//Publisher, Moderator, or Commentron Admin
	ModChannelID   string `json:"mod_channel_id"`
	ModChannelName string `json:"mod_channel_name"`
	//Offender being blocked
	BlockedChannelID   string `json:"blocked_channel_id"`
	BlockedChannelName string `json:"blocked_channel_name"`
	//Creator that Moderator is delegated from. Used for delegated moderation
	CreatorChannelID   string `json:"creator_channel_id"`
	CreatorChannelName string `json:"creator_channel_name"`
	// Blocks identity from comment universally, requires Admin rights on commentron instance
	BlockAll bool `json:"block_all"`
	// If true will delete all comments of the offender, requires Admin rights on commentron for universal delete
	DeleteAll bool   `json:"delete_all"`
	Signature string `json:"signature"`
	SigningTS string `json:"signing_ts"`
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
	ChannelName string `json:"channel_name"`
	ChannelID   string `json:"channel_id"`
	Signature   string `json:"signature"`
	SigningTS   string `json:"signing_ts"`
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
	//Publisher, Moderator, or Commentron Admin
	ModChannelID   string `json:"mod_channel_id"`
	ModChannelName string `json:"mod_channel_name"`
	//Offender being unblocked
	UnBlockedChannelID   string `json:"un_blocked_channel_id"`
	UnBlockedChannelName string `json:"un_blocked_channel_name"`
	//Creator that Moderator is delegated from. Used for delegated moderation
	CreatorChannelID   string `json:"creator_channel_id"`
	CreatorChannelName string `json:"creator_channel_name"`
	// Unblocks identity from commenting universally, requires Admin rights on commentron instance
	GlobalUnBlock bool   `json:"global_un_block"`
	Signature     string `json:"signature"`
	SigningTS     string `json:"signing_ts"`
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
	//Publisher, Moderator or Commentron Admin
	ModChannelID   string `json:"mod_channel_id"`
	ModChannelName string `json:"mod_channel_name"`
	//Creator that Moderator is delegated from. Used for delegated moderation
	CreatorChannelID   string `json:"creator_channel_id"`
	CreatorChannelName string `json:"creator_channel_name"`
	Signature          string `json:"signature"`
	SigningTS          string `json:"signing_ts"`
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
	BlockedByChannelID   string    `json:"blocked_by_channel_id"`
	BlockedByChannelName string    `json:"blocked_by_channel_name"`
	BlockedAt            time.Time `json:"blocked_at"`
}

// AddDelegateArgs Arguments to delagate moderation to another channel for your channel.
type AddDelegateArgs struct {
	ModChannelID       string `json:"mod_channel_id"`
	ModChannelName     string `json:"mod_channel_name"`
	CreatorChannelID   string `json:"creator_channel_id"`
	CreatorChannelName string `json:"creator_channel_name"`
	Signature          string `json:"signature"`
	SigningTS          string `json:"signing_ts"`
}

// RemoveDelegateArgs Arguments to remove a delegated moderator.
type RemoveDelegateArgs struct {
	ModChannelID       string `json:"mod_channel_id"`
	ModChannelName     string `json:"mod_channel_name"`
	CreatorChannelID   string `json:"creator_channel_id"`
	CreatorChannelName string `json:"creator_channel_name"`
	Signature          string `json:"signature"`
	SigningTS          string `json:"signing_ts"`
}

// ListDelegatesArgs Arguments to list delegates
type ListDelegatesArgs struct {
	CreatorChannelID   string `json:"creator_channel_id"`
	CreatorChannelName string `json:"creator_channel_name"`
	Signature          string `json:"signature"`
	SigningTS          string `json:"signing_ts"`
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
