package slack

import (
	"regexp"
	"testing"

	log "github.com/Sirupsen/logrus"
)

var shutdownMessage *Message = NewMessage("shutdown", "general")

func shutdownHandler(_ *Bot, _ map[string]interface{}) (*Message, Status) {
	return shutdownMessage, Shutdown
}

func TestListen(t *testing.T) {
	log.SetLevel(1) // TODO: find the actual name
	var tests = []struct {
		pattern         string
		eventText       string
		expectedMessage *Message
		expectedStatus  Status
	}{
		{"hello", "hello", shutdownMessage, Shutdown},
		{"hello", "goodbye", nil, Continue},
		{"hel*", "he", shutdownMessage, Shutdown},
		{"hel*", "hellllll", shutdownMessage, Shutdown},
		{"hel*", "halp", nil, Continue},
	}

	for _, test := range tests {
		bot := NewBot("token")
		bot.Listen(test.pattern, shutdownHandler)
		handler := bot.Handlers["message"][0]
		event := map[string]interface{}{"text": test.eventText}
		actualMessage, actualStatus := handler(nil, event)
		if test.expectedMessage == nil {
			if actualMessage != nil {
				t.Errorf("Error. Expected nil. Got %v.", actualMessage)
			}
		} else if actualMessage == nil {
			t.Errorf("Error. Expected %v. Got nil.", test.expectedMessage)
		} else {
			compareMessages(test.expectedMessage.toMap(), actualMessage.toMap(), t)
		}
		if test.expectedStatus != actualStatus {
			t.Errorf("Error. Expected %i. Got %i", test.expectedStatus, actualStatus)
		}
	}
}

func TestListenRegexp(t *testing.T) {
	log.SetLevel(1) // TODO: see above
	re := regexp.MustCompile("lo?l")

	var tests = []struct {
		eventText       string
		expectedMessage *Message
		expectedStatus  Status
	}{
		{"ll", shutdownMessage, Shutdown},
		{"lol", shutdownMessage, Shutdown},
		{"ol", nil, Continue},
	}

	for _, test := range tests {
		bot := NewBot("token")
		bot.ListenRegexp(re, shutdownHandler)
		handler := bot.Handlers["message"][0]
		event := map[string]interface{}{"text": test.eventText}
		actualMessage, actualStatus := handler(nil, event)
		if test.expectedMessage == nil {
			if actualMessage != nil {
				t.Errorf("Error. Expected nil. Got %v.", actualMessage)
			}
		} else if actualMessage == nil {
			t.Errorf("Error. Expected %v. Got nil.", test.expectedMessage)
		} else {
			compareMessages(test.expectedMessage.toMap(), actualMessage.toMap(), t)
		}
		if test.expectedStatus != actualStatus {
			t.Errorf("Error. Expected %i. Got %i", test.expectedStatus, actualStatus)
		}
	}
}
