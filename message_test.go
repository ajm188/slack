package slack

import (
	"testing"
)

func TestCreatingANewMessageWithNoError(t *testing.T) {
	m := NewMessage("text", "channel")
	if m == nil {
		t.Error("Failure. Got a nil result.")
	}
}
