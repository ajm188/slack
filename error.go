package slack

type SlackError struct {
	Message string
}

func (err *SlackError) Error() string {
	return err.Message
}
