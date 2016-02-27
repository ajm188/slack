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

type Plugin interface {
	slack.Plugin
	extraEnvVars() []string
	setEnvVar(string, string)
}

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
