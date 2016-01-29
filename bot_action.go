package slack

type BotAction func(self *Bot, incoming map[string]interface{}) (*Message, Status)
