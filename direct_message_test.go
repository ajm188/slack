package slack

import (
	"testing"

	log "github.com/Sirupsen/logrus"
)

func TestDirectMessage_failsWithNoToken(t *testing.T) {
	bot := NewBot("")
	message := bot.DirectMessage("andrew", "hello")
	if message != nil {
		t.Errorf("Error. Expecting nil. Got %v.", message)
	}
}

func TestPrivate_logOpenDMError(t *testing.T) {
	log.SetLevel(log.PanicLevel)
	logOpenDMError(nil, "", "") // smoke test that this doesn't panic
}
