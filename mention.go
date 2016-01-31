package slack

import (
	"fmt"
)

// Mention constructs a Message which mentions nick with text in channel and
// returns a reference to it.
func (bot *Bot) Mention(nick, text, channel string) *Message {
	fullText := fmt.Sprintf("<@%s>: %s", nick, text)
	return NewMessage(fullText, channel)
}
