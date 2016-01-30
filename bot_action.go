package slack

type BotAction func(self *Bot, event map[string]interface{}) (*Message, Status)
