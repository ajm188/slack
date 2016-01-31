package slack

import (
	"net/http"
	"net/url"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
)

const (
	// Version is the semantic version of this library.
	Version = "0.1.0"
)

// Bot encapsulates all the data needed to interact with Slack.
type Bot struct {
	Token       string
	Name        string
	ID          string
	Handlers    map[string]([]BotAction)
	Subhandlers map[string](map[string]([]BotAction))
	Users       map[string]string
	Channels    map[string]string
}

// NewBot constructs a new bot with the passed-in Slack API token.
func NewBot(token string) *Bot {
	return &Bot{
		Token:       token,
		Name:        "",
		ID:          "",
		Handlers:    make(map[string]([]BotAction)),
		Subhandlers: make(map[string](map[string]([]BotAction))),
		Users:       make(map[string]string),
		Channels:    make(map[string]string),
	}
}

// Start initiates the bot's interaction with Slack. It obtains a websockect
// URL, connects to it, and then starts the main loop.
func (bot *Bot) Start() error {
	payload, err := bot.Call("rtm.start", url.Values{})
	if err != nil {
		return err
	}
	success, ok := payload["ok"].(bool)
	if !(ok && success) {
		return &Error{"could not connect to RTM API"}
	}
	websocketURL, _ := payload["url"].(string)
	self := payload["self"].(map[string]interface{})
	channels := payload["channels"].([]interface{})
	for _, channelMap := range channels {
		channel := channelMap.(map[string]interface{})
		channelID := channel["id"].(string)
		channelName := channel["name"].(string)
		bot.Channels[channelName] = channelID
		bot.Channels[channelID] = channelName
	}
	users := payload["users"].([]interface{})
	for _, userMap := range users {
		user := userMap.(map[string]interface{})
		userID := user["id"].(string)
		userName := user["name"].(string)
		bot.Users[userName] = userID
		bot.Users[userID] = userName
	}
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
		case Continue:
			if message != nil {
				conn.WriteJSON(message.toMap())
			}
		case Shutdown:
			if message != nil {
				conn.WriteJSON(message.toMap())
			}
			abort = true
		case ShutdownNow:
			return true
		}
	}
	return abort
}
