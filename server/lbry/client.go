package lbry

import (
	"github.com/lbryio/commentron/env"
	"github.com/lbryio/lbry.go/v2/extras/jsonrpc"
	"github.com/sirupsen/logrus"
)

// SDK is the API entrypoint for LBRYSDK calls
var SDK SDKClient

// API is the API entrypoint for internal-api calls
var API APIClient

// SDKClient is the interface type for SDK call
type SDKClient interface {
	GetTx(string) (*jsonrpc.TransactionSummary, error)
	GetClaim(string) (*jsonrpc.Claim, error)
	GetSigningChannelForClaim(string) (*jsonrpc.Claim, error)
}

// NotifyOptions Are the options used to construct the comment event api signature.
type NotifyOptions struct {
	ActionType string
	CommentID  string
	ChannelID  *string
	ParentID   *string
	Comment    *string
	ClaimID    string
	Amount     uint64
}

// APIClient is the interface type for internal-api calls
type APIClient interface {
	Notify(NotifyOptions)
}

// Init initializes the configuration of the LBRY clients and allows for mock clients for testing
func Init(conf *env.Config) {
	SDK = &mockSDK{}
	API = &mockAPI{}
	if conf.SDKUrl != "" {
		SDK = &sdkClient{}
		sdkURL = conf.SDKUrl
	}
	if conf.APIURL != "" && conf.APIToken != "" {
		apiToken = conf.APIToken
		apiURL = conf.APIURL
		API = apiClient{}
	}

	if conf.TestChannel != "" {
		err := setSerializedTestChannel(conf.TestChannel)
		if err != nil {
			logrus.Panic(err)
		}
	}
}

type mockSDK struct{}

func (m *mockSDK) GetClaim(claimID string) (*jsonrpc.Claim, error) {
	return nil, nil
}

func (m *mockSDK) GetSigningChannelForClaim(channelClaimID string) (*jsonrpc.Claim, error) {
	return nil, nil
}

func (m *mockSDK) GetTx(txid string) (*jsonrpc.TransactionSummary, error) {
	return nil, nil
}

type mockAPI struct{}

func (m *mockAPI) Notify(options NotifyOptions) {

}
