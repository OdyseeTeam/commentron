package flags

import (
	"regexp"

	"github.com/lbryio/commentron/model"

	"github.com/lbryio/lbry.go/v2/extras/errors"
)

func CheckComment(proposedComment *model.Comment) error {
	for _, spammerChannelID := range commentSpammers {
		if proposedComment.ChannelID.String == spammerChannelID {
			proposedComment.IsFlagged = true
		}
	}

	for _, phraseRE := range flaggedPhrases {
		cRegex, err := regexp.Compile(phraseRE)
		if err != nil {
			return errors.Err(err)
		}
		if cRegex.MatchString(proposedComment.Body) {
			proposedComment.IsFlagged = true
		}
	}
	return nil
}

func CheckReaction(proposedReaction *model.Reaction) error {
	for _, spammerChannelID := range reactionSpammers {
		if proposedReaction.ChannelID.String == spammerChannelID {
			proposedReaction.IsFlagged = true
		}
	}
	return nil
}
