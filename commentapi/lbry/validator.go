package lbry

import "encoding/hex"

// UnhelixifyAndReverse takes a claimID and returns its bytes
func UnhelixifyAndReverse(claimID string) []byte {
	b, err := hex.DecodeString(claimID)
	if err != nil {
		return nil
	}
	return reverseBytes(b)
}

// rev reverses a byte slice. useful for switching endian-ness
func reverseBytes(b []byte) []byte {
	r := make([]byte, len(b))
	for left, right := 0, len(b)-1; left < right; left, right = left+1, right-1 {
		r[left], r[right] = b[right], b[left]
	}
	return r
}
