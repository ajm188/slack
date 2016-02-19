package slack

type messageWrapper struct {
	message *Message
	status  Status
}

// OnEvent registers handler to fire on the given type of event.
func (bot *Bot) OnEvent(event string, handler BotAction) {
	handlers, ok := bot.Handlers[event]
	if !ok {
		handlers = make([]BotAction, 0)
	}
	handlers = append(handlers, handler)
	bot.Handlers[event] = handlers
}

// OnEventWithSubtype registers handler to fire on the given type and subtype
// of event.
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

func (bot *Bot) handle(event map[string]interface{}) (wrappers []messageWrapper) {
	eventType, hasType := event["type"].(string)
	eventSubtype, hasSubtype := event["subtype"].(string)

	if hasSubtype {
		subhandlerMap, ok := bot.Subhandlers[eventType]
		if ok {
			subhandlers, ok := subhandlerMap[eventSubtype]
			if ok {
				for _, subhandler := range subhandlers {
					message, status := subhandler(bot, event)
					wrappers = append(wrappers,
						messageWrapper{message, status})
				}
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
	return
}
