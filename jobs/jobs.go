package jobs

import (
	"time"

	"github.com/go-co-op/gocron"
	"github.com/lbryio/commentron/db"
	"github.com/lbryio/commentron/model"
	"github.com/sirupsen/logrus"
)

func StartJobs() {
	scheduler := gocron.NewScheduler(time.UTC)

	scheduler.Every(1).Hours().Do(removeFlagged)

}

func removeFlagged() {
	err := model.Comments(model.CommentWhere.IsFlagged.EQ(true)).DeleteAll(db.RW)
	if err != nil {
		logrus.Error("Error removing flagged comments: ", err)
	}
}
