package github

import (
	"regexp"
	"testing"

	"github.com/google/go-github/github"
)

func TestOpenIssue(t *testing.T) {
	plugin := OpenIssue()
	assert(plugin != nil, t)
}

func TestOpenIssueHasNoExtraEnvVars(t *testing.T) {
	plugin := OpenIssue().(*OpenIssuePlugin)
	assert(len(plugin.extraEnvVars()) == 0, t)
}

func TestOpenIssue_setEnvVar(t *testing.T) {
	plugin := OpenIssue().(*OpenIssuePlugin)
	_, ok := plugin.envVars["foo"]
	assert(!ok, t)
	plugin.setEnvVar("foo", "bar")
	val, ok := plugin.envVars["foo"]
	assert(ok, t)
	assert(val == "bar", t)
}

func TestOpenIssue_Name(t *testing.T) {
	assert(OpenIssue().Name() != "", t)
}

func TestOpenIssue_CanLoad(t *testing.T) {
	env := stubEnv(map[string]string{
		"GH_CLIENT_ID":     "a",
		"GH_CLIENT_SECRET": "b",
		"GH_ACCESS_TOKEN":  "c",
	})
	defer restoreEnv(env)
	assert(env != nil, t)

	assert(OpenIssue().CanLoad(), t)
}

func TestOpenIssue_Load(t *testing.T) {
	dummyIssuesService := github.NewClient(nil).Issues
	plugin := OpenIssue().(*OpenIssuePlugin)
	plugin.issues = dummyIssuesService

	plugin.Load(nil)
	assert(plugin.issues == dummyIssuesService, t)

	env := stubEnv(map[string]string{
		"GH_CLIENT_ID":     "a",
		"GH_CLIENT_SECRET": "b",
		"GH_ACCESS_TOKEN":  "c",
	})
	defer restoreEnv(env)
	assert(env != nil, t)

	plugin.issues = nil
	plugin.Load(nil)
	assert(plugin.issues != nil, t)
}

func TestPrivate_interfaceToRegexp(t *testing.T) {
	var tests = []struct {
		arg      interface{}
		expected *regexp.Regexp
	}{
		{
			arg:      "hello",
			expected: regexp.MustCompile("hello"),
		},
		{
			arg:      regexp.MustCompile("world"),
			expected: regexp.MustCompile("world"),
		},
		{
			arg:      5,
			expected: nil,
		},
	}

	for _, test := range tests {
		actual := interfaceToRegexp(test.arg)
		if test.expected == nil || actual == nil {
			assert(test.expected == actual, t)
		} else {
			actualPrefix, actualFull := actual.LiteralPrefix()
			expectedPrefix, expectedFull := test.expected.LiteralPrefix()
			assert(actualPrefix == expectedPrefix, t)
			assert(actualFull == expectedFull, t)
		}
	}
}

func TestPrivate_interfaceToRegexpWithSubexps(t *testing.T) {
	var tests = []struct {
		arg      interface{}
		n        int
		expected *regexp.Regexp
	}{
		{
			arg:      "(hello)",
			n:        1,
			expected: regexp.MustCompile("(hello)"),
		},
		{
			arg:      "(foo)(bar)",
			n:        3,
			expected: nil,
		},
		{
			arg:      75,
			n:        -1,
			expected: nil,
		},
	}

	for _, test := range tests {
		actual := interfaceToRegexpWithSubexps(test.arg, test.n)
		if test.expected == nil || actual == nil {
			assert(test.expected == actual, t)
		} else {
			actualPrefix, actualFull := actual.LiteralPrefix()
			expectedPrefix, expectedFull := test.expected.LiteralPrefix()
			assert(actualPrefix == expectedPrefix, t)
			assert(actualFull == expectedFull, t)
		}
	}
}
