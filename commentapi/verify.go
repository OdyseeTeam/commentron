package commentapi

// SignatureArgs Arguments to verify the signature from a LBRY SDK
type SignatureArgs struct {
	//Channel ID claiming to have signed the signature
	ChannelID string `json:"channel_id"`
	//Claim ID of the content being validated. This is optional.
	ClaimID string `json:"claim_id"`
	//The data payload in Hex that was signed
	DataHex string `json:"data_hex"`
	//Signature and timestamp returned from the channel_sign api of LBRY SDK
	Signature string `json:"signature"`
	SigningTS string `json:"signing_ts"`
}

// SignatureResponse for the verify.Signature call
type SignatureResponse struct {
	IsValid bool `json:"is_valid"`
}
