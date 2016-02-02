/*
Package slack provides a library for interacting with the Slack API and
building custom bots.

Basics

The basic workflow for writing a bot goes as follows: first, create a new Bot
object with your Slack API token; next, register any callbacks you want
(outlined in further detail below); last, connect the bot to Slack and let it
run forever. This looks like:

	bot := slack.NewBot(myToken)
	// register callbacks here
	bot.Start()

A bot requires a Slack API token in order to connect to Slack, which you can
find under the Custom Integrations for your Slack team. It's worth noting that
a bot cannot add or remove itself from channels; this has to be done by you
when you configure the bot.

When the bot connects, it will collect information about users and channels,
and store this information in the Users and Channels maps. The reason for this
is that Slack does not deal with channels and users in terms of their names
(this is a good thing - channel names and nicks can change), but by a unique
ID. The Users and Channels maps map in both directions; so given the
human-readable name, they will return the ID, and given the ID they will return
the human-readable name.

Slack RTM Basics

Slack provides a Real Time Messaging (RTM) API, for interacting with a Slack
channel programmatically. The important thing to know is that all data is in a
JSON format. See https://api.slack.com/rtm for more information.

Design

Communication with the RTM API is done via websockets. Package slack uses
https://github.com/gorilla/websocket for websockets. From their documentation:
"Connections support one concurrent reader and one concurrent writer.
Applications are responsible for ensuring that no more than one goroutine calls
the write methods (NextWriter, SetWriteDeadline, WriteMessage, WriteJSON)
concurrently and that no more than one goroutine calls the read methods
(NextReader, SetReadDeadline, ReadMessage, ReadJSON, SetPongHandler,
SetPingHandler) concurrently."

For this reason, BotActions (the type for event handlers) do not take a
reference to the websocket connection. Instead, a BotAction takes a reference
to the bot and the event that caused the handler to fire, and it should return
a tuple of (*Message, Status). If the reference to the message is nil, then
nothing will be written into the connection. The Status indicates to the bot
how it should continue to process. See the documentation on the Status values
for more information.

The main loop listens for incoming events from the RTM websocket, and then
calls any handlers that are registered to handle that kind of event. It then
writes any non-nil responses into the websocket, and - depending on the various
status values - may terminate or continue looping.

Events

The Slack RTM API defines a large number of events, which are listed at
https://api.slack.com/events. Note that some events have subtypes. Thus, the
bot supports two general purpose methods for registering an event handler,
which look like:

	bot := NewBot(myToken)
	// Fires on any event with type `eventType`
	bot.OnEvent(eventType, myHandler)
	// Fires only on events with type `eventType` and subtype `eventSubtype'.
	// Events with type `eventType` and no subtype will never cause this
	// handler to fire.
	bot.OnEventWithSubtype(eventType, eventSubtype, myOtherHandler)

Since messages are the most common kind of event, instances of Bot have two
helper methods for registering handlers for messages: "Listen" and "Respond".

Listen takes a pattern and a BotAction, and only invokes the given handler if
the message text matches the regular expression defined by the pattern. It has
a variant, ListenRegexp, which does the same but takes a compiled regular
expression rather than a string pattern.

Respond also takes a pattern and a BotAction, and only invokes the given
handler if the message text "mentions" the bot, and the rest of the text
matches the regular expression defined by the pattern. For a message to
"mention" the bot, the message must begin with the bot's name. The leading "@"
that is commonly used in Slack is optional, as is the trailing ": ". The text
without the portion that was considered part of the "mention" is then compared
against the pattern. Respond also has a variant, RespondRegexp, which does
exactly what you would expect.

Common BotActions

Package slack provides a few helper functions for generating BotAction handlers
for common tasks.

"Respond" creates a handler which will reply to a "message" event with the
specified text. So, if a user named "@example" triggers the handler, the bot
will say "@example: <text>".

"React" creates a handler which will post a Slack reaction to a "message" event
with the specified emoji name. Note that you do not need to put the colons
around the emoji name, unlike what you would need to manually do in Slack to
produce the emoji.
*/
package slack
