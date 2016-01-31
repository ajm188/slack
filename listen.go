package slack

import (
	"regexp"
)

func (bot *Bot) ListenRegexp(re *regexp.Regexp, handler BotAction) {
	closure := func(self *Bot, event map[string]interface{}) (*Message, Status) {
		text, ok := event["text"].(string)
		if !ok {
			return nil, Continue
		}
		if re.MatchString(text) {
			return handler(self, event)
		}
		return nil, Continue
	}
	messageHandlers, ok := bot.Handlers["message"]
	if !ok {
		messageHandlers = make([]BotAction, 0)
	}
	messageHandlers = append(messageHandlers, closure)
	bot.Handlers["message"] = messageHandlers
}

func (bot *Bot) Listen(pattern string, handler BotAction) {
	re := regexp.MustCompile(pattern)
	bot.ListenRegexp(re, handler)
}
