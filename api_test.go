package slack

import (
	"testing"
)

func TestPrivate_unpackJSON(t *testing.T) {
	var tests = []struct {
		body             string
		expectedResponse map[string]interface{}
		expectedError    bool
	}{
		{
			"{\"hello\": 1}",
			map[string]interface{}{"hello": float64(1)},
			false,
		},
		{
			"{\"hello\": \"world\", \"json\": \"data\"}",
			map[string]interface{}{"hello": "world", "json": "data"},
			false,
		},
		{"not valid json", nil, true},
		{"{'single quotes': 'not allowed'}", nil, true},
		{"{\"not delimited\": 5", nil, true},
	}

	for _, test := range tests {
		body := []byte(test.body)
		json, err := unpackJSON(body)
		compareMaps(test.expectedResponse, json, t)
		if test.expectedError && err == nil {
			t.Error("Error. Was expecting error, but did not get one.")
		} else if !test.expectedError && err != nil {
			t.Errorf("Error. Was not expecting an error, but found %v", err)
		}
	}
}

func compareMaps(expected, actual map[string]interface{}, t *testing.T) {
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
