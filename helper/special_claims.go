package helper

import (
	"crypto/sha256"
	"encoding/hex"
	"regexp"
	"strings"
)

const (
	nonValidClaimPrefix         = "w00"
	nonValidModerationChannelID = "5e39a2174041333fa4944f62e93b848014b3ac7f"

	wooClaimIDLength         = 40
	wooYtIDHexLength         = 22
	wooClaimChecksumLength   = wooClaimIDLength - len(nonValidClaimPrefix) - wooYtIDHexLength
	wooYtIDAllowedCharsRegex = `^[A-Za-z0-9_-]{11}$`
)

var (
	wooYtIDRegex = regexp.MustCompile(wooYtIDAllowedCharsRegex)
	lowerHexRE   = regexp.MustCompile(`^[0-9a-f]+$`)
)

// IsNonValidClaimID identifies strict Woo/non-SDK claim IDs.
func IsNonValidClaimID(claimID string) bool {
	if !strings.HasPrefix(claimID, nonValidClaimPrefix) {
		return false
	}
	_, ok := getWooYtIDFromClaimID(claimID)
	return ok
}

// ResolveCreatorChannelClaimID maps special claim IDs to the fixed moderation channel.
func ResolveCreatorChannelClaimID(claimID string) string {
	if IsNonValidClaimID(claimID) {
		return nonValidModerationChannelID
	}
	return claimID
}

func getWooYtIDFromClaimID(claimID string) (string, bool) {
	if len(claimID) != wooClaimIDLength || !strings.HasPrefix(claimID, nonValidClaimPrefix) {
		return "", false
	}

	ytHexStart := len(nonValidClaimPrefix)
	ytHexEnd := ytHexStart + wooYtIDHexLength
	ytHex := claimID[ytHexStart:ytHexEnd]
	checksum := claimID[ytHexEnd:]

	if !lowerHexRE.MatchString(ytHex) || !lowerHexRE.MatchString(checksum) {
		return "", false
	}

	ytBytes, err := hex.DecodeString(ytHex)
	if err != nil {
		return "", false
	}
	ytID := string(ytBytes)
	if !wooYtIDRegex.MatchString(ytID) {
		return "", false
	}

	hash := sha256.Sum256([]byte(ytID))
	expectedChecksum := hex.EncodeToString(hash[:])[:wooClaimChecksumLength]
	if checksum != expectedChecksum {
		return "", false
	}

	return ytID, true
}
