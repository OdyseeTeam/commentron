package websocket

// PushNotification is a message format that tells the client the type of message and the content
type PushNotification struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data,omitempty"`
}
