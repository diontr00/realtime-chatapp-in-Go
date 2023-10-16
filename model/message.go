package model

import (
	"encoding/json"
)

// Event is the messages send over the websocket
// Used to differ between different actions
type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

const (
	SendMessageEvent  = "send_message"
	ErrorMessageEvent = "error_message"
)

// The payload sent in the SendMessageEvent
type SendMessagePayload struct {
	Room      string `json:"room"`
	User      string `json:"user"`
	Message   string `json:"message"`
	CreatedAt string `json:"createdAt"`
}
