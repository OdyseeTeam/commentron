package jobs

import (
	"time"

	"github.com/lbryio/commentron/db"
	"github.com/lbryio/commentron/metrics"
	"github.com/lbryio/commentron/model"

	"github.com/go-co-op/gocron"
	"github.com/sirupsen/logrus"
)

// StartJobs runs all the background jobs of Commentron
func StartJobs() {
	scheduler := gocron.NewScheduler(time.UTC)

	_, err := scheduler.Every(1).Hours().Do(removeFlaggedComments)
	if err != nil {
		logrus.Error(err)
	}
	_, err = scheduler.Every(1).Hours().Do(removeFlaggedReactions)
	if err != nil {
		logrus.Error(err)
	}
	scheduler.StartAsync()
}

func removeFlaggedComments() {
	defer metrics.Job(time.Now(), "remove_flagged_comments")
	err := model.Comments(model.CommentWhere.IsFlagged.EQ(true)).DeleteAll(db.RW)
	if err != nil {
		logrus.Error("Error removing flagged comments: ", err)
	}
}

func removeFlaggedReactions() {
	defer metrics.Job(time.Now(), "remove_flagged_reactions")
	err := model.Reactions(model.ReactionWhere.IsFlagged.EQ(true)).DeleteAll(db.RW)
	if err != nil {
		logrus.Error("Error removing flagged reactions: ", err)
	}
}
