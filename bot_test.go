package slack_test

import (
	. "."

	"testing"
)

func TestABotHasAToken(t *testing.T) {
	bot := NewBot("token")
	if bot.Token != "token" {
		t.Errorf("Failure. Expected %s; Got %s", "token", bot.Token)
	}
}
