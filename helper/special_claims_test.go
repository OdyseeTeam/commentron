package helper

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func buildWooClaimIDForTest(ytID string) string {
	ytHex := hex.EncodeToString([]byte(ytID))
	hash := sha256.Sum256([]byte(ytID))
	checksum := hex.EncodeToString(hash[:])[:wooClaimChecksumLength]
	return nonValidClaimPrefix + ytHex + checksum
}

func TestIsNonValidClaimID(t *testing.T) {
	valid := buildWooClaimIDForTest("dQw4w9WgXcQ")

	testCases := []struct {
		name    string
		claimID string
		want    bool
	}{
		{name: "valid woo claim", claimID: valid, want: true},
		{name: "wrong length", claimID: valid + "0", want: false},
		{name: "wrong prefix", claimID: "x00" + valid[len(nonValidClaimPrefix):], want: false},
		{name: "uppercase prefix", claimID: "W00" + valid[len(nonValidClaimPrefix):], want: false},
		{name: "bad checksum", claimID: valid[:len(valid)-1] + "0", want: false},
		{name: "bad ytid chars after decode", claimID: nonValidClaimPrefix + "0000000000000000000000" + valid[len(nonValidClaimPrefix)+wooYtIDHexLength:], want: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := IsNonValidClaimID(tc.claimID)
			if got != tc.want {
				t.Fatalf("IsNonValidClaimID(%q) = %v, want %v", tc.claimID, got, tc.want)
			}
		})
	}
}

func TestResolveCreatorChannelClaimID(t *testing.T) {
	validWooClaim := buildWooClaimIDForTest("dQw4w9WgXcQ")
	normalClaim := "5e39a2174041333fa4944f62e93b848014b3ac7f"

	gotWoo := ResolveCreatorChannelClaimID(validWooClaim)
	if gotWoo != nonValidModerationChannelID {
		t.Fatalf("ResolveCreatorChannelClaimID(valid woo) = %q, want %q", gotWoo, nonValidModerationChannelID)
	}

	gotNormal := ResolveCreatorChannelClaimID(normalClaim)
	if gotNormal != normalClaim {
		t.Fatalf("ResolveCreatorChannelClaimID(normal) = %q, want %q", gotNormal, normalClaim)
	}
}
