package comments

import (
	"crypto/sha256"
	"encoding/hex"
	"math"
	"time"

	"github.com/lbryio/commentron/util"

	"github.com/sirupsen/logrus"

	"github.com/spf13/cast"
)

func createCommentId(comment, channelID string) (string, int64, error) {
	// We convert the timestamp from seconds into minutes
	// to prevent spammers from commenting the same BS everywhere.
	timestamp := time.Now().Unix()
	nearestMinute := math.Floor(float64(timestamp) / 60.0)

	logrus.Info("Inputs:", comment, " ", channelID, " ", timestamp, " ", int(nearestMinute))
	c := sha256.Sum256(util.CreateDigest(
		[]byte(":"),
		[]byte(comment),
		[]byte(channelID),
		[]byte(cast.ToString(nearestMinute))))
	return hex.EncodeToString(c[:]), timestamp, nil
}
