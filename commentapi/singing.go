package commentapi

import "github.com/sirupsen/logrus"

func sign(client *Client, args interface{}) interface{} {
	var updatedArgs interface{}
	switch t := args.(type) {
	case CreateArgs:
		t.ChannelID = client.Channel.ChannelID
		t.ChannelName = client.Channel.Name
		var err error
		t.Signature, t.SigningTS, err = client.Channel.Sign([]byte(t.CommentText))
		if err != nil {
			logrus.Error(err)
		}
		updatedArgs = t
	case ReactArgs:
		t.ChannelID = client.Channel.ChannelID
		t.ChannelName = client.Channel.Name
		var err error
		t.Signature, t.SigningTS, err = client.Channel.Sign([]byte(client.Channel.Name))
		if err != nil {
			logrus.Error(err)
		}
	default:
		logrus.Error("unknown type passed to sign()")
	}
	return updatedArgs
}
