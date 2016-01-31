package slack

import (
	"time"
)

type Message struct {
	id          string
	messageType string
	channel     string
	text        string
}

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
