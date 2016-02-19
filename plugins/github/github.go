package github

// TODO:
// - testing

import (
	"fmt"
	"regexp"
	"time"

	"github.com/ajm188/slack"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	ghAuth "golang.org/x/oauth2/github" // have to rename so we don't have 2 "github"s
)

var (
	ClientID     string
	ClientSecret string
	AccessToken  string
	RedirectURL  string
	Scopes       []string
)

func Config() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     ClientID,
		ClientSecret: ClientSecret,
		Endpoint:     ghAuth.Endpoint,
		RedirectURL:  RedirectURL,
		Scopes:       Scopes,
	}
}

func Token() *oauth2.Token {
	var noExpire time.Time // this sets noExpire to the zero Time value
	return &oauth2.Token{
		AccessToken:  AccessToken,
		TokenType:    "", // uhhh
		RefreshToken: "",
		Expiry:       noExpire,
	}
}

func OpenIssue(bot *slack.Bot) {
	repoRe := regexp.MustCompile("issue me ([^/ ]+)/([^/ ]+)")
	argsRe := regexp.MustCompile("(\".*?[^\\\\]\")")
	oauthClient := Config().Client(oauth2.NoContext, Token())
	ghClient := github.NewClient(oauthClient)
	issues := ghClient.Issues

	handler := func(b *slack.Bot, event map[string]interface{}) (*slack.Message, slack.Status) {
		text := event["text"].(string)
		owner, repo, err := extractOwnerAndRepo(text, repoRe)
		if err != nil {
			return nil, slack.Continue
		}
		issueRequest, err := extractIssueArgs(text, argsRe)
		if err != nil {
			return nil, slack.Continue
		}
		issue, _, err := issues.Create(owner, repo, issueRequest)
		if err != nil {
			return nil, slack.Continue
		}

		user := event["user"].(string)
		channel := event["channel"].(string)
		message := fmt.Sprintf(
			"I created that issue for you. You can view it here: %s",
			*issue.HTMLURL,
		)
		return bot.Mention(user, message, channel), slack.Continue
	}

	bot.RespondRegexp(repoRe, handler)
}

func extractOwnerAndRepo(text string, re *regexp.Regexp) (string, string, error) {
	m := re.FindStringSubmatch(text)
	if m == nil || len(m) < 3 {
		return "", "", &repoError{text}
	}
	return m[1], m[2], nil
}

func removeQuotes(s string) string {
	return s[1 : len(s)-1]
}

func extractIssueArgs(text string, re *regexp.Regexp) (*github.IssueRequest, error) {
	match := re.FindAllString(text, -1)
	m := make([]string, len(match))
	for i, v := range match {
		m[i] = removeQuotes(v)
	}
	if m == nil || len(m) == 0 {
		return nil, &issueError{text}
	}
	var title, body, assignee *string
	title = &m[0]
	if len(m) >= 2 {
		body = &m[1]
	}
	if len(m) >= 3 {
		assignee = &m[2]
	}
	issueState := "open"
	request := github.IssueRequest{
		Title:     title,
		Body:      body,
		Labels:    nil,
		Assignee:  assignee,
		State:     &issueState,
		Milestone: nil,
	}
	return &request, nil
}
