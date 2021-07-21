package commentapi

// Authorization parameters for calls requiring user authentication
type Authorization struct {
	ChannelName string `json:"channel_name"`
	ChannelID   string `json:"channel_id"`
	Signature   string `json:"signature"`
	SigningTS   string `json:"signing_ts"`
}
