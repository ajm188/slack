package github

import (
	"regexp"
	"testing"

	"github.com/ajm188/slack"
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
	bot := slack.NewBot("TOKEN")
	dummyIssuesService := github.NewClient(nil).Issues
	plugin := OpenIssue().(*OpenIssuePlugin)
	plugin.issues = dummyIssuesService

	plugin.Load(bot)
	assert(plugin.issues == dummyIssuesService, t)

	env := stubEnv(map[string]string{
		"GH_CLIENT_ID":     "a",
		"GH_CLIENT_SECRET": "b",
		"GH_ACCESS_TOKEN":  "c",
	})
	defer restoreEnv(env)
	assert(env != nil, t)

	plugin.issues = nil
	plugin.Load(bot)
	assert(plugin.issues != nil, t)
}

func TestPrivate_extractOwnerAndRepo(t *testing.T) {
	var tests = []struct {
		text          string
		pattern       string
		err           error
		expectedOwner string
		expectedRepo  string
		expectedErr   bool
	}{
		{
			text:          "happy path",
			pattern:       "(.*) (.*)",
			err:           nil,
			expectedOwner: "happy",
			expectedRepo:  "path",
			expectedErr:   false,
		},
		{
			text:          "issue me default/pattern",
			pattern:       "issue me ([^/ ]+)/([^/ ]+)",
			err:           nil,
			expectedOwner: "default",
			expectedRepo:  "pattern",
			expectedErr:   false,
		},
		{
			text:          "no match",
			pattern:       "(na) (.*)",
			err:           nil,
			expectedOwner: "",
			expectedRepo:  "",
			expectedErr:   true,
		},
		{
			text:          "not enough capture groups",
			pattern:       "(not).*",
			err:           nil,
			expectedOwner: "",
			expectedRepo:  "",
			expectedErr:   true,
		},
		{
			text:          "err is non-nil",
			pattern:       "",
			err:           &repoError{"some error"},
			expectedOwner: "",
			expectedRepo:  "",
			expectedErr:   true,
		},
	}

	for _, test := range tests {
		re := regexp.MustCompile(test.pattern)
		owner, repo, err := extractOwnerAndRepo(test.text, re, test.err)
		if test.expectedErr {
			assert(err != nil, t)
		} else {
			assert(owner == test.expectedOwner, t)
			assert(repo == test.expectedRepo, t)
			assert(err == nil, t)
		}
	}
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
