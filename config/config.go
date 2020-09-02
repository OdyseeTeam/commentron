package config

import (
	"github.com/lbryio/commentron/db"
	"github.com/lbryio/commentron/env"
	"github.com/lbryio/commentron/server/lbry"
	"github.com/lbryio/commentron/util"

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
		util.Debugging = true
		logrus.SetLevel(logrus.DebugLevel)
	}
	if viper.GetBool("tracemode") {
		util.Debugging = true
		logrus.SetLevel(logrus.TraceLevel)
	}
	lbry.SDKURL = conf.SDKUrl
	lbry.APIToken = conf.APIToken
	lbry.APIURL = conf.APIURL
	logrus.Info("DSN: ", conf.MySQLDsn)
	err = db.Init(conf.MySQLDsn, util.Debugging)
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
