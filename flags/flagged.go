package flags

import (
	"github.com/OdyseeTeam/commentron/model"
)

// CheckComment checks and flags comments for deletion due to spam or key phrases
func CheckComment(proposedComment *model.Comment) {
	if _, found := commentSpammers[proposedComment.ChannelID.String]; found {
		proposedComment.IsFlagged = true
	}

	for _, re := range flaggedPhrases {
		if re.MatchString(proposedComment.Body) {
			proposedComment.IsFlagged = true
			break
		}
	}
}

// CheckReaction checks reactions for spammers and flags reaction for deletion.
func CheckReaction(proposedReaction *model.Reaction) error {
	if _, found := reactionSpammers[proposedReaction.ChannelID.String]; found {
		proposedReaction.IsFlagged = true
	}
	return nil
}
