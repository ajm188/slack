package slack

import (
	"testing"
)

func TestMention(t *testing.T) {
	bot := NewBot("")
	var tests = []struct {
		nick, text, channel string
		expected            map[string]string
	}{
		{
			"andrew",
			"hello",
			"world",
			map[string]string{
				"text":    "<@andrew>: hello",
				"channel": "world",
				"type":    "message",
			},
		},
	}

	for _, test := range tests {
		actual := bot.Mention(test.nick, test.text, test.channel)
		compareMessages(test.expected, actual.toMap(), t)
	}
}
