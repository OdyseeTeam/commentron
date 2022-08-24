package flags

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPhraseRegex(t *testing.T) {
	phrasesThatShouldMatch := []string{
		`I'M MAKING OVER $420MIL A MONTH WORKING PART TIME. I KEPT HEARING OTHER PEOPLE TELL ME HOW MUCH MONEY THEY CAN MAKE ONLINE SO I DECIDED TO LOOK INTO IT. WELL, IT WAS ALL TRUE AN
Read More Copy here →→ →→ →→
https://www.somerandomcrap.com/`,
		`I am profiting (400$ to 500$/hr )online from my workstation. A month ago I GOT chek of about 30k$, this online work is basic and direct, don't need to go OFFICE, Its home online activity. By then this work opportunity is fbegin your work....★★
13
Copy Here→→→→→ http://Www.somerandomcrap.com/`,
		`I’M MAKING OVER $420MIL A MONTH WORKING PART TIME. I KEPT HEARING OTHER PEOPLE TELL ME HOW MUCH MONEY THEY CAN MAKE ONLINE SO I DECIDED TO LOOK INTO IT. WELL, IT WAS ALL TRUE AN
Read More Copy here →→ →→ →→
https://www.somerandomcrap.com/`,
		`I get paid more than $200 to $400 per hour for working online. I heard about this job 3 months ago and after joining this I have earned easily $30k from this without having online working skills . Simply give it a shot on the accompanying site…
50
Here is I started.…………>> http://Www.somerandomcrap.com/`,
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
