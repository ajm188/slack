package slack

import (
	"testing"
)

func TestOnEvent(t *testing.T) {
	bot := NewBot("token")
	bot.OnEvent("message", shutdownHandler)
	if len(bot.Handlers) != 1 {
		t.Errorf("Error. Expecting 1 handler. Got %i", len(bot.Handlers))
	}
	if len(bot.Handlers["message"]) != 1 {
		t.Errorf("Error. Expecting 1 handler. Got %i", len(bot.Handlers["message"]))
	}

	handler := bot.Handlers["message"][0]
	_, status := handler(bot, nil)
	if status != Shutdown {
		t.Error("Error. Found a different handler")
	}

	h := func(_ *Bot, _ map[string]interface{}) (*Message, Status) {
		return nil, ShutdownNow
	}
	bot.OnEvent("user_typing", h)
	if len(bot.Handlers) != 2 {
		t.Errorf(
			"Error. Expecting 2 event types to have handlers. Got %i",
			len(bot.Handlers),
		)
	}
	actualH := bot.Handlers["user_typing"][0]
	_, status = actualH(nil, nil)
	if status != ShutdownNow {
		t.Error("Error. Found a different handler")
	}
	bot.OnEvent("user_typing", h)
	userTypingHandlers := bot.Handlers["user_typing"]
	if len(userTypingHandlers) != 2 {
		t.Errorf(
			"Error. Expecting 2 handlers for \"user_typing\" events. Found %i.",
			len(userTypingHandlers),
		)
	}
}

func TestOnEventWithSubtype(t *testing.T) {
	bot := NewBot("token")
	bot.OnEventWithSubtype("message", "channel_join", shutdownHandler)
	if len(bot.Subhandlers) != 1 {
		t.Errorf("Error. Expecting 1 handler. Got %i", len(bot.Subhandlers))
	}
	if len(bot.Subhandlers["message"]["channel_join"]) != 1 {
		t.Errorf("Error. Expecting 1 handler. Got %i",
			len(bot.Subhandlers["message"]["channel_join"]))
	}

	handler := bot.Subhandlers["message"]["channel_join"][0]
	_, status := handler(nil, nil)
	if status != Shutdown {
		t.Error("Error. Found a different handler.")
	}

	bot.OnEventWithSubtype("message", "channel_join", shutdownHandler)
	channelJoinHandlers := bot.Subhandlers["message"]["channel_join"]
	if len(channelJoinHandlers) != 2 {
		t.Errorf("Error. Expecting 2 handlers. Found %i.",
			len(channelJoinHandlers))
	}
}

func TestPrivate_handle_noSubtype(t *testing.T) {
	event := map[string]interface{}{"type": "message"}
	// No handlers
	bot := NewBot("token")
	wrappers := bot.handle(event)
	if len(wrappers) != 0 {
		t.Errorf("Error. Expecting 0 wrappers. Found %i.", len(wrappers))
	}

	// one handler
	h1 := func(_ *Bot, _ map[string]interface{}) (*Message, Status) {
		return nil, Shutdown
	}

	bot.OnEvent("message", h1)
	wrappers = bot.handle(event)
	if len(wrappers) != 1 {
		t.Errorf("Error. Expecting 1 wrapper. Found %i.", len(wrappers))
	}
	if wrappers[0].status != Shutdown {
		t.Errorf("Error. Expecting status %i. Found %i.", Shutdown,
			wrappers[0].status)
	}

	// two handlers
	h2 := func(_ *Bot, _ map[string]interface{}) (*Message, Status) {
		return nil, ShutdownNow
	}
	bot.OnEvent("message", h2)
	wrappers = bot.handle(event)
	if len(wrappers) != 2 {
		t.Errorf("Error. Expecting 2 wrappers. Found %i.", len(wrappers))
	}
}

func TestPrivate_handle_subtype(t *testing.T) {
	event := map[string]interface{}{
		"type":    "message",
		"subtype": "channel_join",
	}
	bot := NewBot("token")

	assert := func(expected, actual int, t *testing.T) {
		if expected != actual {
			t.Errorf("Error. Expecting %i. Got %i.", expected, actual)
		}
	}

	// no handlers
	wrappers := bot.handle(event)
	assert(0, len(wrappers), t)

	// one subevent handler
	bot.OnEventWithSubtype("message", "channel_join", shutdownHandler)
	wrappers = bot.handle(event)
	assert(1, len(wrappers), t)

	// one relevant, one irrelevant
	bot.OnEventWithSubtype("message", "not_relevant", shutdownHandler)
	wrappers = bot.handle(event)
	assert(1, len(wrappers), t)

	// adding regular event handler
	bot.OnEvent("message", shutdownHandler)
	wrappers = bot.handle(event)
	assert(2, len(wrappers), t)

	// second subevent handler
	bot.OnEventWithSubtype("message", "channel_join", shutdownHandler)
	wrappers = bot.handle(event)
	assert(3, len(wrappers), t)
}

func TestPrivate_handle_noType(t *testing.T) {
	event := map[string]interface{}{"foo": "bar"}
	bot := NewBot("token")
	bot.OnEvent("bar", shutdownHandler)
	wrappers := bot.handle(event)
	if len(wrappers) != 0 {
		t.Errorf("Error. Expecting 0 wrappers. Found %i.", len(wrappers))
	}
}
