package github

import (
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/ajm188/slack"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	// have to rename so we don't have two github packages
	ghAuth "golang.org/x/oauth2/github"
)

var (
	requiredEnvVars = []string{
		"GH_CLIENT_ID",
		"GH_CLIENT_SECRET",
		"GH_ACCESS_TOKEN",
	}
)

// The Plugin interface is used to define the interface for all GitHub
// plugins. Implementing this interface will allow the plugin to use the common
// `CanLoad` logic.
type Plugin interface {
	slack.Plugin
	// extraEnvVars returns any a list of any additional environment variables
	// this plugin needs.
	extraEnvVars() []string
	// setEnvVar should store the value of an environment variable in this
	// plugin.
	setEnvVar(string, string)
	// getEnvVar should return the value of the given environment variable,
	// according to this plugin.
	getEnvVar(string) string
}

// CanLoad can be used to load any generic GitHub plugin. It will ensure that
// the basic `requiredEnvVars` (GH_CLIENT_ID, GH_CLIENT_SECRET, and
// GH_ACCESS_TOKEN) are set, as well as any extra environment variables the
// plugin may depend on. All of these variables' values will be handed to the
// plugin so it can store them.
//
// If all necessary environment variables are set, then CanLoad will return
// true. Otherwise, it will log at the Error level for each offending
// environment variable and return false.
func CanLoad(plugin Plugin) (ok bool) {
	ok = true
	for _, evar := range append(requiredEnvVars, plugin.extraEnvVars()...) {
		if val := os.Getenv(evar); val == "" {
			ok = false
			log.WithFields(log.Fields{
				"var":    evar,
				"plugin": plugin.Name(),
			}).Error("GitHub plugin missing environment variable. Not loading.")
		} else {
			plugin.setEnvVar(evar, val)
		}
	}
	return
}

// Client returns a github.Client object that can be used to communicate with
// GitHub. It uses the NoContext context from the oauth2 package. It uses the
// plugin to generate a valid token for oauth.
func Client(plugin Plugin) *github.Client {
	oauthConf := OAuthConfig(plugin, "", []string{})
	oauthToken := OAuthToken(plugin, "")
	return github.NewClient(oauthConf.Client(oauth2.NoContext, oauthToken))
}

// OAuthConfig returns an oauth2.Config object that can be used to generate a
// client for communicating GitHub.
func OAuthConfig(plugin Plugin, redirectURL string, scopes []string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     plugin.getEnvVar("GH_CLIENT_ID"),
		ClientSecret: plugin.getEnvVar("GH_CLIENT_SECRET"),
		Endpoint:     ghAuth.Endpoint,
		RedirectURL:  redirectURL,
		Scopes:       scopes,
	}
}

// OAuthToken constructs a basic oauth2.Token with the bare minimum amount of
// information necessary to authenticate with GitHub. It sets the token to
// never expire.
func OAuthToken(plugin Plugin, tokenType string) *oauth2.Token {
	var noExpiry time.Time // sets to zeroed value
	return &oauth2.Token{
		AccessToken:  plugin.getEnvVar("GH_ACCESS_TOKEN"),
		TokenType:    tokenType,
		RefreshToken: "",
		Expiry:       noExpiry,
	}
}
