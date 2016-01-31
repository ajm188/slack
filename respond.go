package slack

import (
	"fmt"
	"regexp"
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
	namePattern := fmt.Sprintf("\\A%s|<@%s>:? ", bot.Name, bot.ID)
	nameRe := regexp.MustCompile(namePattern)
	closure := func(self *Bot, event map[string]interface{}) (*Message, Status) {
		text := event["text"].(string)
		match := nameRe.FindStringIndex(text)
		if match == nil {
			return nil, Continue
		}
		unmatchedText := text[match[1]+1:]
		if re.MatchString(unmatchedText) {
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

// Respond registers the given handler to fire on "message" events with no
// subtype, which address the bot directly and match the given text.
func (bot *Bot) Respond(text string, handler BotAction) {
	re := regexp.MustCompile(text)
	bot.RespondRegexp(re, handler)
}
