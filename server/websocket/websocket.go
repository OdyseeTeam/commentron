package websocket

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/lbryio/lbry.go/v2/extras/api"
	"github.com/lbryio/lbry.go/v2/extras/errors"
	v "github.com/lbryio/ozzo-validation"
	"github.com/lbryio/ozzo-validation/is"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var hub = newHub()

// SubscribeLiveChat is a handler to accept web socket subscriptions for a live chat
func SubscribeLiveChat() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
}

// serveWs handles websocket requests from the peer.
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	params := struct {
		SubscriptionID string `json:"subscription_id"`
	}{}

	err := api.FormValues(r, &params, []*v.FieldRules{
		v.Field(&params.SubscriptionID, v.Required, is.ASCII),
	})

	if err != nil {
		_, writeErr := fmt.Fprintf(w, `%v`, err.Error())
		if writeErr != nil {
			logrus.Error(errors.FullTrace(err), writeErr.Error())
		}
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.Error(errors.FullTrace(err))
		_, err := fmt.Fprintf(w, `%v`, err.Error())
		if err != nil {
			logrus.Error(errors.FullTrace(err))
		}
		return
	}
	client := &Client{subscriptionID: params.SubscriptionID, hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.write()
	go client.read()
}
