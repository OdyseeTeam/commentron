package tests

import (
	"fmt"
	"time"

	"github.com/lbryio/commentron/commentapi"

	"github.com/lbryio/lbry.go/extras/util"
	"github.com/lbryio/lbry.go/v2/extras/errors"

	"github.com/sirupsen/logrus"
)

var address = fmt.Sprintf("http://%s:%d/api", "localhost", 5900)

// Launch launches the e2e tests for commentron
func Launch() {
	err := createComment("testing 1")
	if err != nil {
		logrus.Fatal(err)
	}
	err = createComment("testing 2")
	if err != nil {
		logrus.Fatal(err)
	}
	list, err := listComments()
	if err != nil {
		logrus.Fatal(err)
	}
	if len(list.Items) < 2 {
		logrus.Fatal("should have at least 2 comments")
	}
}

func listComments() (*commentapi.ListResponse, error) {
	c := commentapi.NewClient(address)
	r, err := c.CommentList(commentapi.ListArgs{
		ClaimID: util.PtrToString("abe3c90453fd481383acb4e3d243e2f4efd43e02"),
	})
	return r, errors.Err(err)
}

func createComment(comment string) error {
	c := commentapi.NewClient(address)
	_, err := c.CommentCreate(commentapi.CreateArgs{
		CommentText: comment,
		ClaimID:     "abe3c90453fd481383acb4e3d243e2f4efd43e02",
		ChannelID:   "599617e276c2704a3bff888991bd8a018df672a9",
		ChannelName: "@LBRYBeamer",
		Signature:   "f2e71944954c7e1b5b3a2d01d82c8830e3b9e3ec4d9bb152369dca0707a58e8b390fbdc301539383c398287c08e1c4054198a7730f1fa945c272085b47ec7928",
		SigningTS:   fmt.Sprintf("%d", time.Now().Unix()),
	})
	return errors.Err(err)
}
