package slack

import (
	"net/url"

	log "github.com/Sirupsen/logrus"
)

func (bot *Bot) DirectMessage(userID, text string) *Message {
	dm, err := bot.OpenDirectMessage(userID)
	if err != nil {
		return nil
	}
	return NewMessage(text, dm)
}

func (bot *Bot) OpenDirectMessage(userID string) (string, error) {
	payload, err := bot.Call("im.open", url.Values{"user": []string{userID}})
	if err != nil {
		return "", err
	}
	success := payload["ok"].(bool)
	if !success {
		logOpenDMError(payload, userID, bot.Users[userID])
		return "", &Error{"could not open direct message"}
	}
	channel := payload["channel"].(map[string]interface{})
	return channel["id"].(string), nil
}

func logOpenDMError(payload map[string]interface{}, userID, nick string) {
	log.WithFields(log.Fields{
		"payload": payload,
		"userID":  userID,
		"nick":    nick,
	}).Error("Failed to open direct message.")
}
