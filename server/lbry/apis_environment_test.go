package lbry

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckPerkOptionalEnvironmentForm(t *testing.T) {
	empty, staging := "", "staging"
	for _, environment := range []*string{nil, &empty, &staging} {
		t.Run("environment", func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if err := r.ParseForm(); err != nil {
					t.Error(err)
				}
				values, present := r.PostForm["environment"]
				if present != (environment != nil) || (present && (len(values) != 1 || values[0] != *environment)) {
					t.Errorf("incorrect optional environment: %v", values)
				}
				if r.URL.Path != "/membership_perk/check" || r.PostForm.Get("claim_id") != "fixture-claim" || r.PostForm.Get("type") != "Exclusive content" {
					t.Error("membership request changed")
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"success":true,"data":{"has_access":true}}`))
			}))
			defer server.Close()
			oldURL, oldToken := apiURL, apiToken
			apiURL, apiToken = server.URL, "fixture-token"
			defer func() { apiURL, apiToken = oldURL, oldToken }()
			allowed, err := (apiClient{}).CheckPerk(CheckPerkOptions{ClaimID: "fixture-claim", Type: "Exclusive content", Environment: environment})
			if err != nil || !allowed {
				t.Fatalf("membership form failed: %v, %v", allowed, err)
			}
		})
	}
}
