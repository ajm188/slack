package slack

import (
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

type Bot struct {
	Token string
}

func NewBot(token string) *Bot {
	return &Bot{Token: token}
}

func (bot *Bot) Start() error {
	payload, err := bot.Call("rtm.start", url.Values{})
	if err != nil {
		return err
	}
	success, ok := payload["ok"].(bool)
	if !(ok && success) {
		return &SlackError{"could not connect to RTM API"}
	}
	websocketURL, _ := payload["url"].(string)
	return bot.connect(websocketURL)
}

func (bot *Bot) connect(websocketURL string) error {
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(websocketURL, http.Header{})
	if err != nil {
		return err
	}
	bot.loop(conn)
	return nil
}

func (bot *Bot) loop(conn *websocket.Conn) {
	for {
		messageType, bytes, err := conn.ReadMessage()
		if err != nil {
			// ReadMessage returns an error if the connection is closed
			conn.Close()
			return
		}
		if messageType == websocket.BinaryMessage {
			continue // ignore binary messages
		}
	}
}
