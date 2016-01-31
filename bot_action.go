package slack

// BotAction represents a handler for an RTM event. A valid BotAction takes a
// reference to the bot, as well as the event that caused it to fire. It should
// then return a reference to a Message (which can be nil), and a Status, which
// instructs the bot's main loop in whether to continue listening and sending
// messages or to terminate.
type BotAction func(self *Bot, event map[string]interface{}) (*Message, Status)
