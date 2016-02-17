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
		actualMessage := bot.Mention(test.nick, test.text, test.channel)
		actual := actualMessage.toMap()
		test.expected["id"] = actual["id"]
		compareMapsString(test.expected, actual, t)
	}
}
