package slack

import (
	"net/url"
)

// React creates a BotAction which reacts to the passed-in event with emoji.
func React(emoji string) BotAction {
	closure := func(bot *Bot, event map[string]interface{}) (*Message, Status) {
		channel := event["channel"].(string)
		timestamp := event["ts"].(string)
		params := url.Values{}
		params.Set("channel", channel)
		params.Set("timestamp", timestamp)
		params.Set("name", emoji)
		bot.Call("reactions.add", params)
		return nil, Continue
	}
	return closure
}
