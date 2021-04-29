package flags

import (
	"testing"
)

func TestPhraseRegex(t *testing.T) {
	for _, re := range flaggedPhrases {
		re.MatchString("yolo")
	}
}
