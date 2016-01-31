package slack

import (
	"time"
)

// Message represents a message to be sent to the Slack RTM API. Messages are
// converted to a JSON-like map before they are written to the websocket.
type Message struct {
	id          string
	messageType string
	channel     string
	text        string
}

// NewMessage constructs a new message object which will send text to channel.
// The Slack RTM API uses IDs to identify messages, so NewMessage uses the
// current time as the identifier.
func NewMessage(text, channel string) *Message {
	return &Message{
		id:          formatTime(time.Now()),
		messageType: "message",
		channel:     channel,
		text:        text,
	}
}

func (m *Message) toMap() map[string]string {
	return map[string]string{
		"id":      m.id,
		"type":    m.messageType,
		"channel": m.channel,
		"text":    m.text,
	}
}
