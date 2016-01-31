package slack

import (
	"net/url"

	log "github.com/Sirupsen/logrus"
)

func (bot *Bot) OpenDirectMessage(userID string) (string, error) {
	payload, err := bot.Call("im.open", url.Values{"user": []string{userID}})
	if err != nil {
		return "", err
	}
	success := payload["ok"].(bool)
	if !success {
		log.WithFields(log.Fields{
			"payload": payload,
			"userID":  userID,
			"nick":    bot.Users[userID],
		}).Error("Failed to open direct message.")
		return "", &SlackError{"could not open direct message"}
	}
	channel := payload["channel"].(map[string]interface{})
	return channel["id"].(string), nil
}
