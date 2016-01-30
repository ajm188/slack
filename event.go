package slack

import (
	"encoding/json"
)

type messageWrapper struct {
	message *Message
	status Status
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
