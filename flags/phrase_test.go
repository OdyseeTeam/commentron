package flags

import (
	"regexp"
	"testing"
)

func TestPhraseRegex(t *testing.T) {
	for _, phraseRE := range flaggedPhrases {
		_, err := regexp.Compile(phraseRE)
		if err != nil {
			t.Error(err)
		}
	}
}
