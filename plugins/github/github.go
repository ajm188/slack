/*
Package github provides a plugin for building Github integrations into
github.com/ajm188/slack.

The package exports a number of package-wide variables, which may be used to
configure the parameters used for authenticating with the Github API. See the
documentation on each variable for what is is used for.

Here is an example of how one might use this library to register a hook that
opens Github issues:

    import (
        "github.com/ajm188/slack"
        "github.com/ajm188/slack/plugins/github"
    )

    func main() {
        bot := slack.NewBot(myToken)
        // configure auth for github plugin
        github.OpenIssue(bot, nil)
    }
*/
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
	// ClientID is issued by Github when you register an application. You
	// should register an application for your bot, and set the ID before
	// creating a client to authenticate with Github.
	ClientID string
	// ClientSecret is the secret key used to verify your registered
	// application. You should set the secret before creating a client to
	// authenticate with Github.
	ClientSecret string
	// AccessToken is an OAuth token for the user you want your Github
	// interactions to be performed as. The bot will be commenting, opening
	// issues, etc as the user who owns this token. You should set an access
	// token before creating a client to authenticate with Github.
	AccessToken string
	// RedirectURL is the URL that Github should redirect to after a successful
	// web authentication. Since the bot does not perform web-based
	// authentication, this is likely a useless field. More information can be
	// found at https://developer.github.com/v3/oauth/#redirect-urls.
	RedirectURL string
	// Scopes is the list of scopes that the OAuth token should be limited
	// to. Since the user that created the token can specify the scopes
	// available to the token when they create it, this field is probably also
	// useless.
	Scopes []string
	// SharedClient is a variable for sharing a single OAuth client among
	// various handlers. For example, when a call to OpenIssue is made, you may
	// pass in a client of your own, if you want the issue hook to be handled
	// by a different Github user. If you do not pass in a client, then the
	// various hook methods will fall back to using this shared client.
	SharedClient *github.Client
)

// DefaultClient constructs a Github client based on the variables set in this
// package (ClientID, ClientSecret, AccessToken). This can be used to quickly
// create a client when you don't need any customization to the underlying
// oauth client. It uses the NoContext context from the oauth2 package. See the
// Token function for the Token it will use.
func DefaultClient() *github.Client {
	return github.NewClient(Config().Client(oauth2.NoContext, Token()))
}

// Config returns an oauth config object that can be used to generate a client
// for communicating with Github.
func Config() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     ClientID,
		ClientSecret: ClientSecret,
		Endpoint:     ghAuth.Endpoint,
		RedirectURL:  RedirectURL,
		Scopes:       Scopes,
	}
}

// Token constructs a basic token, with the bare minimum amount of information
// necessary to authenticate with Github. It uses the package-wide AccessToken,
// and sets the token to never expire. "TokenType" and "RefreshToken" fields
// are left blank.
func Token() *oauth2.Token {
	var noExpire time.Time // this sets noExpire to the zero Time value
	return &oauth2.Token{
		AccessToken:  AccessToken,
		TokenType:    "", // uhhh
		RefreshToken: "",
		Expiry:       noExpire,
	}
}

// OpenIssue registers a handler that will cause the bot to open a github issue
// based on the event text.
//
// The handler is registered as a "Respond", not a "Listen" (see the docs for
// github.com/ajm188/slack for the difference). The pattern which will cause
// the handler to fire has the form 'issue me // <owner>/<repo> "<title>"
// ("<body>" ("<assignee>")?)?'.
//
// The function takes as arguments the bot to which it should register the
// handler, and a reference to a client that can authenticate with Github. If
// no client is provided, then OpenIssue will fall back to using the
// package-wide SharedClient.
//
// Users should note that an attempt to assign an issue to a Github user that
// is not a "contributor" on the repository will result in a 422 returned by
// the Github API. This will prevent the issue from being created.
//
// When an issue has successfully been created, the bot will reply to the user
// which triggered the handler with a link to the issue.
func OpenIssue(bot *slack.Bot, client *github.Client) {
	repoRe := regexp.MustCompile("issue me ([^/ ]+)/([^/ ]+)")
	argsRe := regexp.MustCompile("(\".*?[^\\\\]\")")
	if client == nil {
		client = SharedClient
	}
	issues := client.Issues

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
		user := event["user"].(string)
		channel := event["channel"].(string)
		if err != nil {
			message := fmt.Sprintf(
				"I had some trouble opening an issue. Here was the error I got:\n%v",
				err)
			return bot.Mention(user, message, channel), slack.Continue
		}

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
