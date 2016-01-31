package slack

import (
	"net/url"
)

func (bot *Bot) OpenDirectMessage(userID string) (map[string]interface{}, error) {
	return bot.Call("im.open", url.Values{"user": []string{userID}})
}
