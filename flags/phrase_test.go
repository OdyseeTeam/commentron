package flags

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPhraseRegex(t *testing.T) {
	phrasesThatShouldMatch := []string{
		`follow me as I follow you`,
	}

	for _, phrase := range phrasesThatShouldMatch {
		found := false
		for _, re := range flaggedPhrases {
			if re.MatchString(phrase) {
				found = true
				break
			}
		}
		assert.True(t, found, "Phrase '%s' should match", phrase)
	}
}
