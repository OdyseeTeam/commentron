package comments

import (
	"crypto/sha256"
	"encoding/hex"
	"math"
	"time"

	"github.com/lbryio/commentron/helper"

	"github.com/spf13/cast"
)

type check int

const (
	checkFrequency check = iota
	ignoreFrequency
)

func createCommentID(comment, channelID string, frequency check) (string, int64, error) {
	timestamp := time.Now().Unix()
	compositeTimestamp := timestamp
	if frequency == checkFrequency {
		// We convert the timestamp from seconds into minutes
		// to prevent spammers from commenting the same BS everywhere.
		compositeTimestamp = int64(math.Floor(float64(timestamp) / 60.0))
	}

	c := sha256.Sum256(helper.CreateDigest(
		[]byte(":"),
		[]byte(comment),
		[]byte(channelID),
		[]byte(cast.ToString(compositeTimestamp))))
	commentID := hex.EncodeToString(c[:])

	err := checkForDuplicate(commentID)
	if err != nil {
		return commentID, timestamp, err
	}

	return commentID, timestamp, nil
}
