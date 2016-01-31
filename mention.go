package slack

import (
	"fmt"
)

func (bot *Bot) Mention(nick, text, channel string) *Message {
	fullText := fmt.Sprintf("<@%s>: %s", nick, text)
	return NewMessage(fullText, channel)
}
