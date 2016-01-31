package slack

// Error is the struct used to create custom errors that occur within the slack
// package.
type Error struct {
	Message string
}

func (err *Error) Error() string {
	return err.Message
}
