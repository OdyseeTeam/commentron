package jobs

import (
	"time"

	"github.com/OdyseeTeam/commentron/jobs/commentclassification"

	"github.com/OdyseeTeam/commentron/db"
	"github.com/OdyseeTeam/commentron/metrics"
	"github.com/OdyseeTeam/commentron/model"

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

	// Use gocron for the frequently polling comment classification job
	// because you would expect to find it in periodic jobs.
	_, err = scheduler.Every(1).Minute().Do(commentclassification.PollAndClassifyNewComments)
	if err != nil {
		logrus.Error(err)
	}
	scheduler.StartAsync()
}

func removeFlaggedComments() {
	defer metrics.Job(time.Now(), "remove_flagged_comments")
	err := model.Comments(model.CommentWhere.IsFlagged.EQ(true)).DeleteAll(db.RW, false)
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
