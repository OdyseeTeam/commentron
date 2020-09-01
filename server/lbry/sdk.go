package lbry

import (
	"github.com/lbryio/lbry.go/v2/extras/errors"
	"github.com/lbryio/lbry.go/v2/extras/jsonrpc"
)

// SDKURL is the url the client should use to connect the sdk
var SDKURL string

// GetChannelClaim retrieves the channel claim information from the sdk.
func GetChannelClaim(claimID string) (*jsonrpc.Claim, error) {
	c := jsonrpc.NewClient(SDKURL)
	claimSearchResp, err := c.ClaimSearch(nil, &claimID, nil, nil, 1, 1)
	if err != nil {
		return nil, errors.Err(err)
	}
	if len(claimSearchResp.Claims) > 0 {
		channel := claimSearchResp.Claims[0]
		return &channel, nil
	}
	return nil, errors.Err("could not get channel claim from sdk with claim id %s", claimID)
}
