package commentapi

// BlockArgs Arguments to block identities from commenting for both publisher and moderators
type BlockArgs struct {
	//Publisher or Commentron Admin
	ModChannelID   string `json:"mod_channel_id"`
	ModChannelName string `json:"mod_channel_name"`
	//Offender being blocked
	BannedChannelID   string `json:"banned_channel_id"`
	BannedChannelName string `json:"banned_channel_name"`
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
