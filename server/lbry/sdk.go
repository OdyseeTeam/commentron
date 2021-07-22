package lbry

import (
	"time"

	"github.com/lbryio/commentron/metrics"

	"github.com/lbryio/lbry.go/v2/extras/errors"
	"github.com/lbryio/lbry.go/v2/extras/jsonrpc"

	"github.com/karlseguin/ccache"
)

var sdkURL string

var claimCache = ccache.New(ccache.Configure().GetsPerPromote(1).MaxSize(10000))

type sdkClient struct{}

// GetClaim retrieves the channel claim information from the sdk.
func (sdk *sdkClient) GetClaim(claimID string) (*jsonrpc.Claim, error) {
	metrics.SDKClaimCache.WithLabelValues("hit").Add(1)
	cachedValue, err := claimCache.Fetch(claimID, 30*time.Minute, sdk.getClaim(claimID))
	if err != nil {
		return nil, errors.Err(err)
	}
	if cachedValue.Value() == nil {
		claimCache.Delete(claimID)
		return nil, errors.Err("could not get claim from sdk with claim id %s", claimID)
	}
	return cachedValue.Value().(*jsonrpc.Claim), nil
}

func (sdk *sdkClient) getClaim(claimID string) func() (interface{}, error) {
	return func() (interface{}, error) {
		metrics.SDKClaimCache.WithLabelValues("miss").Add(1)
		metrics.SDKClaimCache.WithLabelValues("hit").Sub(1)
		defer metrics.SDKCall(time.Now(), "claim-search")
		c := jsonrpc.NewClient(sdkURL)
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
func (sdk *sdkClient) GetSigningChannelForClaim(claimID string) (*jsonrpc.Claim, error) {
	claim, err := sdk.GetClaim(claimID)
	if err != nil {
		return nil, err
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

// GetTx retrieves the transaction details
func (sdk *sdkClient) GetTx(txid string) (*jsonrpc.TransactionSummary, error) {
	c := jsonrpc.NewClient(sdkURL)
	summary, err := c.TransactionShow(txid)
	if err != nil {
		return nil, err
	}
	return summary, nil
}
