package slack

import (
	"regexp"
	"testing"

	log "github.com/Sirupsen/logrus"
)

func TestRespondGeneratedClosure(t *testing.T) {
	log.SetLevel(1) // TODO: see listen_test.go
	text := "hi there"
	var tests = []struct {
		user, channel   string
		expectedMessage map[string]string
		expectedStatus  Status
	}{
		{
			"andrew",
			"general",
			map[string]string{
				"text":    "<@andrew>: hi there",
				"type":    "message",
				"channel": "general",
			},
			Continue,
		},
	}

	for _, test := range tests {
		responseHandler := Respond(text)
		bot := NewBot("token")
		event := map[string]interface{}{"user": test.user, "channel": test.channel}
		actualMessage, actualStatus := responseHandler(bot, event)
		if test.expectedMessage == nil {
			if actualMessage != nil {
				t.Errorf("Error. Expected nil. Got %v.", actualMessage)
			}
		} else if actualMessage == nil {
			t.Errorf("Error. Expected %v. Got nil.", test.expectedMessage)
		} else {
			compareMessages(test.expectedMessage, actualMessage.toMap(), t)
		}
		if test.expectedStatus != actualStatus {
			t.Errorf("Error. Expected %i. Got %i", test.expectedStatus, actualStatus)
		}
	}
}

func TestRespond(t *testing.T) {
	log.SetLevel(1) // TODO: see above
	var tests = []struct {
		pattern         string
		eventText       string
		expectedMessage *Message
		expectedStatus  Status
	}{
		{"hello", "testbot: hello", shutdownMessage, Shutdown},
		{"hello", "testbot: goodbye", nil, Continue},
		{"hel*", "testbot: he", shutdownMessage, Shutdown},
		{"hel*", "testbot hellllll", shutdownMessage, Shutdown},
		{"hel*", "testboot: hel", nil, Continue},
	}

	for _, test := range tests {
		bot := NewBot("token")
		bot.Name = "testbot"
		bot.Respond(test.pattern, shutdownHandler)
		handler := bot.Handlers["message"][0]
		event := map[string]interface{}{"text": test.eventText}
		actualMessage, actualStatus := handler(bot, event)
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

func TestRespondNoEventText(t *testing.T) {
	log.SetLevel(1) // TODO: see above

	bot := NewBot("token")
	bot.Name = "testbot"
	bot.Respond("hi", shutdownHandler)
	handler := bot.Handlers["message"][0]
	event := map[string]interface{}{}
	actualMessage, actualStatus := handler(bot, event)
	if actualMessage != nil {
		t.Errorf("Error. Expected nil. Got %v.", actualMessage)
	}
	if Continue != actualStatus {
		t.Errorf("Error. Expected %i. Got %i", Continue, actualStatus)
	}
}

func TestRespondRegexp(t *testing.T) {
	log.SetLevel(1) // TODO: see above
	re := regexp.MustCompile("lo?l")

	var tests = []struct {
		eventText       string
		expectedMessage *Message
		expectedStatus  Status
	}{
		{"testbot lol", shutdownMessage, Shutdown},
		{"goodbye", nil, Continue},
		{"tastbot lol", nil, Continue},
		{"testbot: ll", shutdownMessage, Shutdown},
		{"halp", nil, Continue},
	}

	for _, test := range tests {
		bot := NewBot("token")
		bot.Name = "testbot"
		bot.RespondRegexp(re, shutdownHandler)
		handler := bot.Handlers["message"][0]
		event := map[string]interface{}{"text": test.eventText}
		actualMessage, actualStatus := handler(bot, event)
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
