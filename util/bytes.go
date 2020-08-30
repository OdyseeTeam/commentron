package util

func CreateDigest(pieces ...[]byte) []byte {
	var digest []byte
	for _, p := range pieces {
		digest = append(digest, p...)
	}
	return digest
}
