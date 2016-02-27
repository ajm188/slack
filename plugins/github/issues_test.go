package github

import (
	"testing"
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
