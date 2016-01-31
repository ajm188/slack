package slack

import (
	"encoding/json"
)

type messageWrapper struct {
	message *Message
	status  Status
}

func (bot *Bot) OnEvent(event string, handler BotAction) {
	handlers, ok := bot.Handlers[event]
	if !ok {
		handlers make([]BotAction, 0)
	}
	handlers = append(handlers, handler)
	bot.Handlers[event] = handlers
}

func (bot *Bot) OnEventWithSubtype(event, subtype string, handler BotAction) {
	subtypeMap, ok := bot.Subhandlers[event]
	if !ok {
		subtypeMap = make(map[string]([]BotAction))
		bot.Subhandlers[event] = subtypeMap
	}
	handlers, ok := bot.Subhandlers[event][subtype]
	if !ok {
		handlers = make([]BotAction, 0)
	}
	handlers = append(handlers, handler)
	bot.Subhandlers[event][subtype] = handlers
}

func unpackEvent(bytes []byte) (map[string]interface{}, error) {
	var message map[string]interface{}
	err := json.Unmarshal(bytes, &message)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (bot *Bot) handle(event map[string]interface{}) []messageWrapper {
	wrappers := make([]messageWrapper, 0)

	eventType, hasType := event["type"].(string)
	eventSubtype, hasSubtype := event["subtype"].(string)

	if hasSubtype {
		subhandlers, ok := bot.Subhandlers[eventType][eventSubtype]
		if ok {
			for _, subhandler := range subhandlers {
				message, status := subhandler(bot, event)
				wrappers = append(wrappers, messageWrapper{message, status})
			}
		}
	}
	if hasType {
		handlers, ok := bot.Handlers[eventType]
		if ok {
			for _, handler := range handlers {
				message, status := handler(bot, event)
				wrappers = append(wrappers, messageWrapper{message, status})
			}
		}
	}
	return wrappers
}
