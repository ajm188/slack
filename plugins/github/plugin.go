package github

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/ajm188/slack"
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
