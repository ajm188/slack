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

func TestPrivate_toMap(t *testing.T) {
	var tests = []struct {
		message  *Message
		expected map[string]string
	}{
		{
			NewMessage("hello", "world"),
			map[string]string{
				"type":    "message",
				"text":    "hello",
				"channel": "world",
			},
		},
	}

	for _, test := range tests {
		actual := test.message.toMap()
		compareMessages(test.expected, actual, t)
	}
}

func compareMessages(expected, actual map[string]string, t *testing.T) {
	expected["id"] = actual["id"]
	compareMapsString(expected, actual, t)
}

func compareMapsString(expected, actual map[string]string, t *testing.T) {
	if len(expected) != len(actual) {
		t.Errorf("Error. Expected map had %d keys. Actual had %d keys",
			len(expected), len(actual))
		return
	}

	for k, v := range expected {
		actualVal, ok := actual[k]
		if !ok {
			t.Errorf("Error. Actual map missing key %s", k)
			continue
		}
		if actualVal != v {
			t.Errorf(
				"Error. Maps have different values for %s. "+
					"Expected: %v. Actual: %v",
				k, v, actualVal,
			)
		}
	}
}
