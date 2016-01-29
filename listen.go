package slack

import (
	"regexp"
)

func (bot *Bot) ListenRegexp(re *regexp.Regexp, handler BotAction) {
}

func (bot *Bot) Listen(pattern string, handler BotAction) {
	re := regexp.MustCompile(pattern)
	bot.ListenRegexp(re, handler)
}
