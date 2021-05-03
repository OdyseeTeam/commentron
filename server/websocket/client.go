package websocket

// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/lbryio/lbry.go/v2/extras/errors"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub            *Hub
	subscriptionID string
	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

func (c *Client) handleMessage(message []byte) {
	rm := &PushNotification{}
	err := json.Unmarshal(message, rm)
	if err != nil {
		logrus.Error(errors.FullTrace(err), hex.EncodeToString(message))
	}
	err = rm.process()
	if err != nil {
		logrus.Error(errors.FullTrace(err))
	}
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) read() {
	defer func() {
		c.hub.unregister <- c
		err := c.conn.Close()
		if err != nil {
			logrus.Error(errors.FullTrace(err))
		}
	}()
	c.conn.SetReadLimit(maxMessageSize)
	err := c.conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		panic("could not set read deadline")
	}
	c.conn.SetPongHandler(func(string) error { return c.conn.SetReadDeadline(time.Now().Add(pongWait)) })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseNoStatusReceived, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logrus.Error(errors.FullTrace(errors.Err(fmt.Sprintf("error: %v", err))))
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		go c.handleMessage(message)
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		err := c.conn.Close()
		if err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
			logrus.Error(errors.FullTrace(err))
		}
	}()
	for {
		select {
		case message, ok := <-c.send:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				logrus.Error(errors.FullTrace(err))
			}
			if !ok {
				// The hub closed the channel.
				if err := c.conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil && !errors.Is(err, websocket.ErrCloseSent) {
					if !strings.Contains(err.Error(), "use of closed network connection") {
						logrus.Error(errors.FullTrace(err))
					}
				}
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			if _, err = w.Write(message); err != nil {
				logrus.Error(errors.FullTrace(err))
			}

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				if _, err = w.Write(newline); err != nil {
					logrus.Error(errors.FullTrace(err))
				}
				if _, err = w.Write(<-c.send); err != nil {
					logrus.Error(errors.FullTrace(err))
				}
			}

			if err := w.Close(); err != nil {
				logrus.Error(errors.FullTrace(err))
				return
			}
		case <-ticker.C:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				logrus.Error(errors.FullTrace(err))
			}
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				logrus.Error(errors.FullTrace(err))
				return
			}
		}
	}
}
