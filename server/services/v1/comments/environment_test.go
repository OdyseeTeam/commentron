package comments

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/OdyseeTeam/commentron/commentapi"
	"github.com/OdyseeTeam/commentron/server/lbry"
	"github.com/lbryio/lbry.go/v2/extras/jsonrpc"
)

const environmentClaim = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
const environmentChannel = "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"

type environmentAPI struct {
	lbry.APIClient
	calls []lbry.CheckPerkOptions
	allow bool
	err   error
}

func (a *environmentAPI) CheckPerk(options lbry.CheckPerkOptions) (bool, error) {
	a.calls = append(a.calls, options)
	return a.allow, a.err
}

type environmentSDK struct{ lbry.SDKClient }

func (environmentSDK) GetClaim(string) (*jsonrpc.Claim, error) {
	return &jsonrpc.Claim{}, nil
}

func installEnvironmentClients(t *testing.T, api *environmentAPI) {
	t.Helper()
	oldAPI, oldSDK := lbry.API, lbry.SDK
	lbry.API, lbry.SDK = api, environmentSDK{}
	t.Cleanup(func() { lbry.API, lbry.SDK = oldAPI, oldSDK })
}

func TestMembershipEnvironmentRemainsOptional(t *testing.T) {
	empty, staging := "", "staging"
	for _, environment := range []*string{nil, &empty, &staging} {
		for name, check := range map[string]func(string, string, *string) (bool, error){
			"Exclusive content": HasAccessToProtectedContent,
			"Members-only chat": HasAccessToProtectedChat,
		} {
			t.Run(name, func(t *testing.T) {
				api := &environmentAPI{allow: true}
				installEnvironmentClients(t, api)
				allowed, err := check(environmentClaim, environmentChannel, environment)
				if err != nil || !allowed || len(api.calls) != 1 {
					t.Fatalf("membership check = %v, %v, calls %d", allowed, err, len(api.calls))
				}
				got := api.calls[0]
				if got.Environment != environment || got.Type != name || got.ClaimID != environmentClaim || got.ChannelClaimID != environmentChannel {
					t.Fatalf("membership options changed: %#v", got)
				}
			})
		}
	}
}

func TestMembershipEnvironmentPreservesDeniedAndFailedChecks(t *testing.T) {
	failed := errors.New("membership service unavailable")
	for _, failure := range []error{nil, failed} {
		api := &environmentAPI{err: failure}
		installEnvironmentClients(t, api)
		for _, check := range []func(string, string, *string) (bool, error){HasAccessToProtectedContent, HasAccessToProtectedChat} {
			allowed, err := check(environmentClaim, environmentChannel, nil)
			if allowed || err != failure {
				t.Fatalf("failed membership check = %v, %v", allowed, err)
			}
		}
	}
}

func TestLegacyProtectedListsWithoutEnvironmentDoNotPanic(t *testing.T) {
	api := &environmentAPI{}
	installEnvironmentClients(t, api)
	claim, channel := environmentClaim, environmentChannel
	request := httptest.NewRequest("POST", "/api", nil)
	if err := getCachedList(request, &commentapi.ListArgs{IsProtected: true, ClaimID: &claim, RequestorChannelID: &channel}, &commentapi.ListResponse{}); err == nil {
		t.Fatal("denied protected comments became readable")
	}
	if err := getCachedSuperChatList(request, &commentapi.SuperListArgs{IsProtected: true, ClaimID: &claim, RequestorChannelID: &channel}, &commentapi.SuperListResponse{}); err == nil {
		t.Fatal("denied protected superchats became readable")
	}
	if len(api.calls) != 2 || api.calls[0].Environment != nil || api.calls[1].Environment != nil {
		t.Fatal("legacy requests did not preserve the omitted environment")
	}
}

func TestMembershipEnvironmentKeepsNonCanonicalClaimHandling(t *testing.T) {
	api := &environmentAPI{}
	installEnvironmentClients(t, api)
	ytID := "fixture0001"
	checksum := sha256.Sum256([]byte(ytID))
	claim := "w00" + hex.EncodeToString([]byte(ytID)) + hex.EncodeToString(checksum[:])[:15]
	for _, check := range []func(string, string, *string) (bool, error){HasAccessToProtectedContent, HasAccessToProtectedChat} {
		allowed, err := check(claim, environmentChannel, nil)
		if err != nil || !allowed {
			t.Fatalf("noncanonical claim behavior changed: %v, %v", allowed, err)
		}
	}
	if len(api.calls) != 0 {
		t.Fatal("noncanonical claim unexpectedly reached membership service")
	}
}
