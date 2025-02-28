package config

import (
	"strings"

	"github.com/OdyseeTeam/commentron/db"
	"github.com/OdyseeTeam/commentron/env"
	"github.com/OdyseeTeam/commentron/helper"
	"github.com/OdyseeTeam/commentron/jobs/commentclassification"
	"github.com/johntdyer/slackrus"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// SocketyToken token used to communicate with Sockety
var SocketyToken string

// IsTestMode turns off validations for local testing
var IsTestMode bool

// InitializeConfiguration inits the base configuration of commentron
func InitializeConfiguration(conf *env.Config) {

	IsTestMode = conf.IsTestMode
	if viper.GetBool("debugmode") {
		helper.Debugging = true
		logrus.SetLevel(logrus.DebugLevel)
	}
	if viper.GetBool("tracemode") {
		helper.Debugging = true
		logrus.SetLevel(logrus.TraceLevel)
	}

	err := db.Init(conf.MySQLDsnRO, conf.MySQLDsnRW, helper.Debugging)
	if err != nil {
		logrus.Panic(err)
	}
	commentclassification.Init(conf)
	initSlack(conf)
	initStripe(conf)
	SocketyToken = conf.SocketyToken

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
			Filters: []slackrus.Filter{
				func(entry *logrus.Entry) bool {
					return !strings.Contains(entry.Message, "could not get claim from sdk")
				},
			},
		})
	}
}
