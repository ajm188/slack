package slack

import (
	"net/http"
	"net/url"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
)

type Bot struct {
	Token       string
	Name        string
	ID          string
	Handlers    map[string]([]BotAction)
	Subhandlers map[string](map[string]([]BotAction))
}

func NewBot(token string) *Bot {
	return &Bot{
		Token:       token,
		Name:        "",
		ID:          "",
		Handlers:    make(map[string]([]BotAction)),
		Subhandlers: make(map[string](map[string]([]BotAction))),
	}
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
	self := payload["self"].(map[string]interface{})
	bot.Name = self["name"].(string)
	bot.ID = self["id"].(string)
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
		event, err := unpackEvent(bytes)
		if err != nil {
			log.WithFields(log.Fields{
				"raw bytes": bytes,
				"error":     err,
			}).Warn("message could not be unpacked")
		}
		wrappers := bot.handle(event)
		closeConnection := sendResponses(wrappers, conn)
		if closeConnection {
			conn.Close()
			return
		}
	}
}

func sendResponses(wrappers []messageWrapper, conn *websocket.Conn) bool {
	abort := false
	for _, wrapper := range wrappers {
		message := wrapper.message
		switch wrapper.status {
		case CONTINUE:
			if message != nil {
				conn.WriteJSON(message.toMap())
			}
		case SHUTDOWN:
			if message != nil {
				conn.WriteJSON(message.toMap())
			}
			abort = true
		case SHUTDOWN_NOW:
			return true
		}
	}
	return abort
}
