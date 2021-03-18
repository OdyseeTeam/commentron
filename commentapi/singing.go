package commentapi

import "github.com/sirupsen/logrus"

func sign(client *Client, args interface{}) interface{} {
	var updatedArgs interface{}
	var err error
	switch t := args.(type) {
	case CreateArgs:
		t.ChannelID = client.Channel.ChannelID
		t.ChannelName = client.Channel.Name
		t.Signature, t.SigningTS, err = client.Channel.Sign([]byte(t.CommentText))
		updatedArgs = t
	case ReactArgs:
		t.ChannelID = client.Channel.ChannelID
		t.ChannelName = client.Channel.Name
		t.Signature, t.SigningTS, err = client.Channel.Sign([]byte(client.Channel.Name))
	default:
		if err != nil {
			logrus.Panic("unknown type")
		}
	}
	return updatedArgs
}
