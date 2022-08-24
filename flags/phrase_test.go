package flags

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPhraseRegex(t *testing.T) {
	phrasesThatShouldMatch := []string{
		`I'M MAKING OVER $420MIL A MONTH WORKING PART TIME. I KEPT HEARING OTHER PEOPLE TELL ME HOW MUCH MONEY THEY CAN MAKE ONLINE SO I DECIDED TO LOOK INTO IT. WELL, IT WAS ALL TRUE AN
Read More Copy here →→ →→ →→
https://www.salaryto.com/`,
	}

	for _, phrase := range phrasesThatShouldMatch {
		found := false
		for _, re := range flaggedPhrases {
			if re.MatchString(phrase) {
				found = true
				break
			}
		}
		assert.True(t, found)
	}
}
