package helper

// CreateDigest utility function for grouping multiple sets of bytes. Largely used for signature verification
func CreateDigest(pieces ...[]byte) []byte {
	var digest []byte
	for _, p := range pieces {
		digest = append(digest, p...)
	}
	return digest
}
