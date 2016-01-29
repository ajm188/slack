package slack_test

import (
	. "."
	"testing"
)

func TestSlackError(t *testing.T) {
	err := &SlackError{"test"}
	result := err.Error()
	if result != "test" {
		t.Errorf("Failure. Expecting %s; Got %s", "test", result)
	}
}
