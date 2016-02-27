package github

import (
	"github.com/ajm188/slack"
)

// OpenIssuePlugin contains all of the data needed to facilitate opening issues
// via GitHub's API.
type OpenIssuePlugin struct {
	envVars map[string]string
}

// The OpenIssuePlugin does not require any extra environment variables.
func (_ *OpenIssuePlugin) extraEnvVars() []string {
	return []string{}
}

func (plugin *OpenIssuePlugin) setEnvVar(name, val string) {
	plugin.envVars[name] = val
}

// OpenIssue returns a new OpenIssuePlugin. This function can be registered
// with a *slack.Bot.
func OpenIssue() slack.Plugin {
	return &OpenIssuePlugin{envVars: make(map[string]string, 3)}
}

// Name returns the name of the OpenIssuePlugin.
func (_ *OpenIssuePlugin) Name() string {
	return "Open Issues"
}

// CanLoad uses the package-default loading mechanism, returning true if the
// procedure succeeded and false otherwise.
func (plugin *OpenIssuePlugin) CanLoad() (ok bool) {
	return CanLoad(plugin)
}

// Load loads the OpenIssuePlugin into the bot.
func (plugin *OpenIssuePlugin) Load(bot *slack.Bot, args ...interface{}) error {
	return nil
}
