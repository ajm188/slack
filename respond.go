package slack

import (
	"fmt"
	"regexp"

	log "github.com/Sirupsen/logrus"
)

// Respond creates a BotAction which responds to the passed-in event with text.
func Respond(text string) BotAction {
	closure := func(bot *Bot, event map[string]interface{}) (*Message, Status) {
		user := event["user"].(string)
		channel := event["channel"].(string)
		return bot.Mention(user, text, channel), Continue
	}
	return closure
}

// RespondRegexp functions exactly as Respond, but instead takes a compiled
// regexp instead of a string.
func (bot *Bot) RespondRegexp(re *regexp.Regexp, handler BotAction) {
	closure := func(self *Bot, event map[string]interface{}) (*Message, Status) {
		name := regexp.MustCompile(fmt.Sprintf("\\A%s:? ", self.Name))
		id := regexp.MustCompile(fmt.Sprintf("\\A<@%s>:? ", self.ID))
		maybeText := event["text"]
		if maybeText == nil {
			return nil, Continue
		}
		text := maybeText.(string)
		logger := log.WithFields(log.Fields{
			"text":  text,
			"regex": re,
		})
		match := name.FindStringIndex(text)
		if match == nil {
			match = id.FindStringIndex(text)
			if match == nil {
				logger.Info("NO MENTION. Not invoking handler.")
				return nil, Continue
			}
		}
		unmatchedText := text[match[1]:]
		if re.MatchString(unmatchedText) {
			logger.Info("MATCH. Invoking handler.")
			return handler(self, event)
		}
		logger.Info("NO MATCH. Not invoking handler.")
		return nil, Continue
	}
	bot.OnEvent("message", closure)
}

// Respond registers the given handler to fire on "message" events with no
// subtype, which address the bot directly and match the given text.
func (bot *Bot) Respond(text string, handler BotAction) {
	re := regexp.MustCompile(text)
	bot.RespondRegexp(re, handler)
}
