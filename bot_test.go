package slack

import (
	"testing"
)

func TestABotHasAToken(t *testing.T) {
	bot := NewBot("token")
	if bot.Token != "token" {
		t.Errorf("Failure. Expected %s; Got %s", "token", bot.Token)
	}
}

func TestPrivate_connect(t *testing.T) {
	bot := NewBot("token")
	err := bot.connect("junk")
	if err == nil {
		t.Error("Error. Expecting error. Got nil")
	}
}
