package github

import (
	"regexp"

	log "github.com/Sirupsen/logrus"
	"github.com/ajm188/slack"
	"github.com/google/go-github/github"
)

// OpenIssuePlugin contains all of the data needed to facilitate opening issues
// via GitHub's API.
type OpenIssuePlugin struct {
	envVars map[string]string
	issues  *github.IssuesService
}

// The OpenIssuePlugin does not require any extra environment variables.
func (_ *OpenIssuePlugin) extraEnvVars() []string {
	return []string{}
}

func (plugin *OpenIssuePlugin) setEnvVar(name, val string) {
	plugin.envVars[name] = val
}

func (plugin *OpenIssuePlugin) getEnvVar(name string) string {
	val, ok := plugin.envVars[name]
	if !ok {
		val = ""
	}
	return val
}

// OpenIssue returns a new OpenIssuePlugin. This function can be registered
// with a *slack.Bot.
func OpenIssue() slack.Plugin {
	return &OpenIssuePlugin{
		envVars: make(map[string]string, 3),
		issues:  nil,
	}
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
	if plugin.issues == nil {
		plugin.issues = Client(plugin).Issues
	}
	repoRe, argsRe := extractPluginArgs(args...)
	bot.RespondRegexp(repoRe, plugin.handler(repoRe, argsRe))
	return nil
}

func (plugin *OpenIssuePlugin) handler(repoRe, argsRe *regexp.Regexp) slack.BotAction {
	return func(bot *slack.Bot, event map[string]interface{}) (*slack.Message, slack.Status) {
		return nil, slack.Continue
	}
}

func extractPluginArgs(args ...interface{}) (repoRe, argsRe *regexp.Regexp) {
	if len(args) > 0 {
		repoRe = interfaceToRegexpWithSubexps(args[0], 2)
	}
	if repoRe == nil {
		repoRe = regexp.MustCompile("issue me ([^/ ]+)/([^/ ]+)")
	}

	if len(args) > 1 {
		argsRe = interfaceToRegexpWithSubexps(args[1], 1)
	}
	if argsRe == nil {
		argsRe = regexp.MustCompile("(\".*?[^\\\\]\")")
	}
	return
}

func interfaceToRegexp(arg interface{}) (re *regexp.Regexp) {
	switch arg.(type) {
	default:
		log.WithFields(log.Fields{
			"arg": arg,
		}).Error("argument must be string or *regexp.Regexp")
	case string:
		re = regexp.MustCompile(arg.(string))
	case *regexp.Regexp:
		re = arg.(*regexp.Regexp)
	}
	return
}

func interfaceToRegexpWithSubexps(arg interface{}, n int) (re *regexp.Regexp) {
	re = interfaceToRegexp(arg)
	if re != nil && re.NumSubexp() != n {
		log.WithFields(log.Fields{
			"regular expression": re,
			"required subexps":   n,
			"arg":                arg,
		}).Error("argument had incorrect number of subexps")
		re = nil
	}
	return
}
