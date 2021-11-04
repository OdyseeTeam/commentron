package commentapi

// Authorization parameters for calls requiring user authentication
type Authorization struct {
	ChannelName string `json:"channel_name"`
	ChannelID   string `json:"channel_id"`
	Signature   string `json:"signature"`
	SigningTS   string `json:"signing_ts"`
}

// ModAuthorization parameters for calls requiring creator/moderator authentication
type ModAuthorization struct {
	//Publisher, Moderator or Commentron Admin
	ModChannelID   string `json:"mod_channel_id"`
	ModChannelName string `json:"mod_channel_name"`
	//Creator that Moderator is delegated from. Used for delegated moderation
	CreatorChannelID   string `json:"creator_channel_id"`
	CreatorChannelName string `json:"creator_channel_name"`
	Signature          string `json:"signature"`
	SigningTS          string `json:"signing_ts"`
}
