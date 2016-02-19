package slack

import (
	"net/http"
	"net/url"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
)

const (
	// Version is the semantic version of this library.
	Version = "0.2.0"
)

// Bot encapsulates all the data needed to interact with Slack.
type Bot struct {
	Token        string
	Name         string
	ID           string
	Handlers     map[string]([]BotAction)
	Subhandlers  map[string](map[string]([]BotAction))
	Users        map[string]string
	Channels     map[string]string
	reconnectURL string
}

// NewBot constructs a new bot with the passed-in Slack API token.
func NewBot(token string) *Bot {
	return &Bot{
		Token:        token,
		Name:         "",
		ID:           "",
		Handlers:     make(map[string]([]BotAction)),
		Subhandlers:  make(map[string](map[string]([]BotAction))),
		Users:        make(map[string]string),
		Channels:     make(map[string]string),
		reconnectURL: "",
	}
}

// StoreReconnectURL takes a "url" from an event and stores it. This is done so
// that when Slack migrates a team to a new host, the bot can use the reconnect
// URL to reattach to the team.
func StoreReconnectURL(bot *Bot, event map[string]interface{}) (*Message, Status) {
	bot.reconnectURL = event["url"].(string)
	return nil, Continue
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
	log.WithFields(log.Fields{
		"id":   bot.ID,
		"name": bot.Name,
	}).Info("bot authenticated")
	bot.OnEvent("reconnect_url", StoreReconnectURL)
	for {
		reconnect, err := bot.connect(websocketURL)
		if reconnect && bot.reconnectURL != "" {
			websocketURL = bot.reconnectURL
		} else {
			return err
		}
	}
	return nil
}

func (bot *Bot) connect(websocketURL string) (bool, error) {
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(websocketURL, http.Header{})
	if err != nil {
		return false, err
	}
	return bot.loop(conn), nil
}

func (bot *Bot) loop(conn *websocket.Conn) bool {
	defer conn.Close()
	for {
		messageType, bytes, err := conn.ReadMessage()
		if err != nil {
			// ReadMessage returns an error if the connection is closed
			return false
		}
		if messageType == websocket.BinaryMessage {
			continue // ignore binary messages
		}
		event, err := unpackJSON(bytes)
		if err != nil {
			log.WithFields(log.Fields{
				"raw bytes": bytes,
				"error":     err,
			}).Warn("message could not be unpacked")
			continue
		}
		log.WithFields(log.Fields{
			"event": event,
		}).Info("received event")
		eventType, ok := event["type"]
		if ok && eventType.(string) == "team_migration_started" {
			return true
		}
		wrappers := bot.handle(event)
		closeConnection := sendResponses(wrappers, conn)
		if closeConnection {
			return false
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
