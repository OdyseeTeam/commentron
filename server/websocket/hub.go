package websocket

import (
	"encoding/json"
	"sync"

	"github.com/OdyseeTeam/commentron/metrics"

	"github.com/lbryio/lbry.go/v2/extras/errors"

	"github.com/sirupsen/logrus"
)

// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[string][]*Client

	// Broadcast message to all clients
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// Locks the map for read/readwrite access
	clientLock sync.RWMutex
}

func newHub() *Hub {
	hub := &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[string][]*Client),
	}
	go hub.run()
	return hub
}

func (h *Hub) getClients(id string) []*Client {
	return h.clients[id]
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)
		case client := <-h.unregister:
			h.unRegisterClient(client)
		case message := <-h.broadcast:
			h.broadcastToClients(message)
		}
	}
}

func (h *Hub) broadcastToClients(message []byte) {
	h.clientLock.RLock()
	defer h.clientLock.RUnlock()
	for token, clients := range h.clients {
		for _, client := range clients {
			select {
			case client.send <- message:
			default:
				h.clientLock.Lock()
				defer h.clientLock.Unlock()
				close(client.send)
				delete(h.clients, token)
			}
		}
	}
}

func (h *Hub) registerClient(client *Client) {
	h.clientLock.Lock()
	defer h.clientLock.Unlock()
	clients := h.getClients(client.subscriptionID)
	if len(clients) > 0 {
		h.clients[client.subscriptionID] = append(clients, client)
	} else {
		h.clients[client.subscriptionID] = []*Client{client}
	}
	metrics.WSConnections.WithLabelValues(client.subscriptionID).Inc()
}

func (h *Hub) unRegisterClient(client *Client) {
	h.clientLock.Lock()
	defer h.clientLock.Unlock()
	existingClients := h.getClients(client.subscriptionID)
	var newClients []*Client
	for _, c := range existingClients {
		if c == client {
			close(c.send)
			metrics.WSConnections.WithLabelValues(client.subscriptionID).Dec()
			continue
		}
		newClients = append(newClients, c)
	}
	if len(newClients) == 0 {
		delete(h.clients, client.subscriptionID)
	} else {
		h.clients[client.subscriptionID] = newClients
	}
}

var hubNotInitialized = errors.Base("hub is not initialized")

// Broadcast sends a PushNotification to all connected clients
func Broadcast(notification *PushNotification) error {
	message, err := json.Marshal(notification)
	if err != nil {
		return errors.Err(err)
	}
	if hub == nil {
		return hubNotInitialized
	}
	hub.broadcast <- message
	return nil
}

// PushTo sends the PushNotification to the subscribed clients
func PushTo(notification *PushNotification, subscriptionID string) {
	err := pushTo(notification, subscriptionID)
	if err != nil {
		logrus.Error(errors.FullTrace(err))
	}
}

func pushTo(notification *PushNotification, subscriptionID string) error {
	if hub == nil {
		return hubNotInitialized
	}

	hub.clientLock.RLock()
	defer hub.clientLock.RUnlock()
	message, err := json.Marshal(notification)
	if err != nil {
		return errors.Err(err)
	}
	clients := hub.getClients(subscriptionID)
	for _, client := range clients {
		client.send <- message
	}
	return nil
}
