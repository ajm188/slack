package github

import (
	"fmt"
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
		var owner, repo string
		var title, body, assignee *string
		var err error
		var request *github.IssueRequest

		text := event["text"].(string)
		userID := event["user"].(string)

		owner, repo, err = extractOwnerAndRepo(text, repoRe, nil)
		title, body, assignee, err = extractIssueArgs(text, argsRe, err)
		err = appendSlackSignature(body, bot.Users[userID], err)
		request, err = buildIssueRequest(title, body, assignee, err)

		if err != nil {
			return nil, slack.Continue
		}

		var response string

		issue, _, err := plugin.issues.Create(owner, repo, request)
		if err != nil {
			response = fmt.Sprintf(
				"I had some trouble opening that issue. Here was the error I got:\n%v",
				err)
		} else {
			response = fmt.Sprintf(
				"I opened that issue for you. You can view it here: %s.",
				*issue.HTMLURL)
		}
		message := bot.Mention(userID, response, event["channel"].(string))
		status := slack.Continue
		return message, status
	}
}

func buildIssueRequest(title, body, assignee *string, err error) (*github.IssueRequest, error) {
	if err != nil {
		return nil, err
	}
	open := "open"
	request := &github.IssueRequest{
		Title:     title,
		Body:      body,
		Labels:    nil,
		Assignee:  assignee,
		State:     &open,
		Milestone: nil,
	}
	return request, nil
}

func appendSlackSignature(body *string, user *slack.User, err error) error {
	if err != nil {
		return err
	}
	if user == nil {
		return &repoError{"could not find user"}
	}
	fullName := user.FullName()
	name := user.Nick
	if fullName != "" {
		name += fmt.Sprintf(" (%s)", fullName)
	}
	*body += fmt.Sprintf("\n\n Issue created via slack on behalf of %s.", name)
	return nil
}

func extractIssueArgs(text string, re *regexp.Regexp, err error) (title, body, assignee *string, e error) {
	e = err
	if e != nil {
		return
	}
	match := re.FindAllString(text, -1)
	if match == nil || len(match) == 0 {
		e = &issueError{text}
		return
	}

	m := make([]string, len(match))
	for i, v := range match {
		m[i] = v[1 : len(v)-1] // strip off the quotes
	}

	title = &m[0]
	if len(m) > 1 {
		body = &m[1]
	}
	if len(m) > 2 {
		assignee = &m[2]
	}
	return
}

func extractOwnerAndRepo(text string, re *regexp.Regexp, err error) (string, string, error) {
	if err != nil {
		return "", "", err
	}
	m := re.FindStringSubmatch(text)
	if m == nil || len(m) != 3 {
		return "", "", &repoError{text}
	}
	return m[1], m[2], nil
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
