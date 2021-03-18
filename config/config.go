package config

import (
	"github.com/lbryio/commentron/db"
	"github.com/lbryio/commentron/env"
	"github.com/lbryio/commentron/helper"
	"github.com/lbryio/commentron/server/lbry"

	"github.com/johntdyer/slackrus"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// InitializeConfiguration inits the base configuration of commentron
func InitializeConfiguration() {
	conf, err := env.NewWithEnvVars()
	if err != nil {
		logrus.Panic(err)
	}
	if viper.GetBool("debugmode") {
		helper.Debugging = true
		logrus.SetLevel(logrus.DebugLevel)
	}
	if viper.GetBool("tracemode") {
		helper.Debugging = true
		logrus.SetLevel(logrus.TraceLevel)
	}
	lbry.Init(conf)
	err = db.Init(conf.MySQLDsn, helper.Debugging)
	if err != nil {
		logrus.Panic(err)
	}
	initSlack(conf)
}

// initSlack initializes the slack connection and posts info level or greater to the set channel.
func initSlack(config *env.Config) {
	slackURL := config.SlackHookURL
	slackChannel := config.SlackChannel
	if slackURL != "" && slackChannel != "" {
		logrus.AddHook(&slackrus.SlackrusHook{
			HookURL:        slackURL,
			AcceptedLevels: slackrus.LevelThreshold(logrus.InfoLevel),
			Channel:        slackChannel,
			IconEmoji:      ":speech_balloon:",
			Username:       "Commentron",
		})
	}
}
