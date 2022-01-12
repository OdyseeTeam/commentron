package main

import (
	"github.com/lbryio/commentron/commentapi"
	"github.com/lbryio/commentron/env"
	"github.com/sirupsen/logrus"
)

func main() {
	conf, err := env.NewWithEnvVars() //
	if err != nil {
		logrus.Panic(err)
	}
	testCommenting(conf)
	testReacting(conf)
}

func testReacting(conf *env.Config) {
	c := commentapi.NewClient(conf.TestURL).WithSigning(conf.TestChannel)
	_, err := c.ReactionReact(commentapi.ReactArgs{
		CommentIDs: "0000125eb5349851d3aecfef344092e801a3b84ea5ab055158d2911b3b96dd57",
		Type:       "like",
	})
	if err != nil {
		logrus.Error(err)
	}
}

func testCommenting(conf *env.Config) {
	c := commentapi.NewClient(conf.TestURL).WithSigning(conf.TestChannel)
	_, err := c.CommentCreate(commentapi.CreateArgs{
		CommentText: "Can I create comments?",
		ClaimID:     "ca1ab0d129060ad12774b239976fd3379d571846",
		ChannelID:   c.Channel.ChannelID,
		ChannelName: c.Channel.Name,
	})
	if err != nil {
		logrus.Error(err)
	}
}
