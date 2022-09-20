package sockety

import (
	"github.com/OdyseeTeam/commentron/config"

	"github.com/lbryio/lbry.go/v2/extras/errors"
	"github.com/lbryio/sockety/socketyapi"

	"github.com/sirupsen/logrus"
)

// Token token used to sent notifications to sockety
var Token string

// URL is the url to connect to an instance of sockety.
var URL = "https://sockety.lbry.com"

var socketyClient *socketyapi.Client

// SendNotification sends the notification to socket using client
func SendNotification(args socketyapi.SendNotificationArgs) {
	if config.SocketyToken == "" || URL == "" {
		return
	}
	defer catchPanic()
	if socketyClient == nil {
		logrus.Debug("initializating sockety client")
		socketyClient = socketyapi.NewClient(URL, config.SocketyToken)
	}
	resp, err := socketyClient.SendNotification(args)
	if err != nil {
		logrus.Error(errors.FullTrace(errors.Prefix("Sockety SendTo: ", err)))
	}
	if resp != nil && resp.Error != nil {
		logrus.Error(errors.FullTrace(errors.Prefix("Sockety SendToResp: ", errors.Base(*resp.Error))))
	}
}

func catchPanic() {
	if r := recover(); r != nil {
		logrus.Error("sockety send recovered from: ", r)
	}
}
