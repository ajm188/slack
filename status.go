package slack

type Status int

const (
	CONTINUE     Status = iota
	SHUTDOWN_NOW        // terminates immediately
	SHUTDOWN            // finishes sending any messages, then terminates
)
