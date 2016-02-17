package slack

import (
	"testing"
)

func TestSlackError(t *testing.T) {
	err := &Error{"test"}
	result := err.Error()
	if result != "test" {
		t.Errorf("Failure. Expecting %s; Got %s", "test", result)
	}
}
