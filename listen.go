package slack

import (
	"regexp"

	log "github.com/Sirupsen/logrus"
)

// ListenRegexp functions exactly as Listen, but instead takes a compiled
// regexp instead of a string.
func (bot *Bot) ListenRegexp(re *regexp.Regexp, handler BotAction) {
	closure := func(self *Bot, event map[string]interface{}) (*Message, Status) {
		text, ok := event["text"].(string)
		if !ok {
			return nil, Continue
		}
		logger := log.WithFields(log.Fields{
			"text":  text,
			"regex": re,
		})
		if re.MatchString(text) {
			logger.Info("MATCH. Invoking handler.")
			return handler(self, event)
		}
		logger.Info("NO MATCH. Not invoking handler.")
		return nil, Continue
	}
	messageHandlers, ok := bot.Handlers["message"]
	if !ok {
		messageHandlers = make([]BotAction, 0)
	}
	messageHandlers = append(messageHandlers, closure)
	bot.Handlers["message"] = messageHandlers
}

// Listen registers the given handler to fire on "message" events with no
// subtype which match the regexp specified in pattern.
func (bot *Bot) Listen(pattern string, handler BotAction) {
	re := regexp.MustCompile(pattern)
	bot.ListenRegexp(re, handler)
}
