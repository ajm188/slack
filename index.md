# slack

[![GoDoc](https://godoc.org/github.com/ajm188/slack?status.svg)](https://godoc.org/github.com/ajm188/slack)

## Installation

`go get github.com/ajm188/slack`

## Usage

### Starting the Bot

```go
package main

import (
    "fmt"

    "github.com/ajm188/slack"
)

func main() {
    bot := slack.NewBot("my_slack_token")
    err := bot.Start()
    if err != nil {
        fmt.Println(err)
    }
}
```

Note that you need to pass your Slack token to the bot, so your bot can
authenticate with the Slack API.

### Simple Response

The design of `slack` is built around a handler pattern. The bot listens for
incoming events from the Slack RTM API, and then invokes any handlers for that
type of event.

A valid handler takes a reference to the bot and the event and returns a
reference to a Message, and a Status, which tells the main loop whether or not
to continue interacting with the RTM API.

`slack` has some factory methods for constructing handlers for common
operations, such as responding to a message or posting a reaction to a message.
If you need something more fine-grained, feel free to write your own.

There are also convenience methods for registering handlers specifically for
"message" type events. The Bot defines instance methods `Listen` - which finds
text that matches the given pattern - and `Respond` which finds text that first
mentions the bot by name and then matches the given pattern. These helper
methods have `Regexp` variants which can take a compiled regular expression
directly instead of a string.

```go
package main

// ... imports blah, blah

func main() {
    bot := slack.NewBot("")
    bot.Respond("hi!", slack.Respond("hi!"))
    bot.Start()
}
```

### Adding a Reaction

```go
package main

// imports

func main() {
    bot := slack.NewBot("")
    bot.Listen("ship ?it\\?", slack.React("shipit"))
    bot.Start()
}
```

### Logging

`slack` uses [logrus](https://github.com/Sirupsen/logrus) for logging. Feel
free to set your log level appropriately in an `init` function.
