package lbry

import (
	"time"

	"github.com/lbryio/lbry.go/v2/extras/errors"
	"github.com/lbryio/lbry.go/v2/extras/jsonrpc"

	"github.com/karlseguin/ccache"
)

// SDKURL is the url the client should use to connect the sdk
var SDKURL string

var claimCache = ccache.New(ccache.Configure().GetsPerPromote(1).MaxSize(10000))

// GetClaim retrieves the channel claim information from the sdk.
func GetClaim(claimID string) (*jsonrpc.Claim, error) {
	cachedValue, err := claimCache.Fetch(claimID, 30*time.Minute, getClaim(claimID))
	if err != nil {
		return nil, err
	}
	if cachedValue.Value() == nil {
		claimCache.Delete(claimID)
		return nil, errors.Err("could not get claim from sdk with claim id %s", claimID)
	}
	return cachedValue.Value().(*jsonrpc.Claim), nil
}

func getClaim(claimID string) func() (interface{}, error) {
	return func() (interface{}, error) {
		c := jsonrpc.NewClient(SDKURL)
		claimSearchResp, err := c.ClaimSearch(nil, &claimID, nil, nil, 1, 1)
		if err != nil {
			return nil, errors.Err(err)
		}
		if len(claimSearchResp.Claims) > 0 {
			channel := claimSearchResp.Claims[0]
			return &channel, nil
		}
		return nil, nil
	}
}

// GetSigningChannelForClaim retrieves the claim for the channel that signed the referenced claim by claim id.
func GetSigningChannelForClaim(claimID string) (*jsonrpc.Claim, error) {
	claim, err := GetClaim(claimID)
	if err != nil {
		return nil, errors.Err(err)
	}
	if claim == nil {
		return nil, errors.Err("could not resolve claim_id %s", claimID)
	}
	claimChannel := claim.SigningChannel
	if claimChannel == nil {
		if claim.ValueType == "channel" {
			claimChannel = claim
		}
	}
	return claimChannel, nil
}
