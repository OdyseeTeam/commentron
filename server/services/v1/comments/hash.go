package comments

import (
	"crypto/sha256"
	"encoding/hex"
	"math"
	"time"

	"github.com/lbryio/commentron/helper"
	"github.com/spf13/cast"
)

func createCommentID(comment, channelID string) (string, int64, error) {
	// We convert the timestamp from seconds into minutes
	// to prevent spammers from commenting the same BS everywhere.
	timestamp := time.Now().Unix()
	nearestMinute := math.Floor(float64(timestamp) / 60.0)

	c := sha256.Sum256(helper.CreateDigest(
		[]byte(":"),
		[]byte(comment),
		[]byte(channelID),
		[]byte(cast.ToString(nearestMinute))))
	commentID := hex.EncodeToString(c[:])

	err := checkForDuplicate(commentID)
	if err != nil {
		return commentID, timestamp, err
	}

	return commentID, timestamp, nil
}
