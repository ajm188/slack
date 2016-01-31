package slack

// Status is an enumerated type used to communicate between BotActions and the
// bot's main loop.
type Status int

const (
	// Continue indicates to keep listening and sending messages.
	Continue Status = iota
	// ShutdownNow indicates to terminate immediately and not send any messages
	// from possible downstream BotActions.
	ShutdownNow
	// Shutdown indicates to finish sending any messages from possible
	// downstream BotActions and then terminate.
	Shutdown
)
