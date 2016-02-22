package slack

import (
	"testing"
)

func assert(condition bool, t *testing.T) {
	if !condition {
		t.Error("Error.")
	}
}

func TestFullName(t *testing.T) {
	var tests = []struct {
		firstName string
		lastName  string
		fullName  string
	}{
		{
			firstName: "Foo",
			lastName:  "Bar",
			fullName:  "Foo Bar",
		},
		{
			firstName: "",
			lastName:  "Bar",
			fullName:  "Bar",
		},
		{
			firstName: "Foo",
			lastName:  "",
			fullName:  "Foo",
		},
		{
			firstName: "",
			lastName:  "",
			fullName:  "",
		},
	}

	for _, test := range tests {
		user := &User{
			ID:        "",
			Nick:      "",
			FirstName: test.firstName,
			LastName:  test.lastName,
		}
		assert(user.FullName() == test.fullName, t)
	}
}

func TestUserFromJSON(t *testing.T) {
	data := map[string]interface{}{
		"id":   "12345",
		"name": "mynick",
		"profile": map[string]interface{}{
			"first_name": "Foo",
			"last_name":  "Bar",
		},
	}

	user := UserFromJSON(data)
	assert(user.ID == "12345", t)
	assert(user.Nick == "mynick", t)
	assert(user.FirstName == "Foo", t)
	assert(user.LastName == "Bar", t)
}
