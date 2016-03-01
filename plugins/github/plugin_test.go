package github

import (
	"os"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/ajm188/slack"
)

func init() {
	log.SetLevel(log.PanicLevel)
}

func stubEnv(env map[string]string) map[string]string {
	oldEnv := make(map[string]string, len(env))
	for name, value := range env {
		oldEnv[name] = os.Getenv(name)
		if err := os.Setenv(name, value); err != nil {
			return nil
		}
	}
	return oldEnv
}

func restoreEnv(env map[string]string) {
	stubEnv(env)
}

func assert(value bool, t *testing.T) {
	if !value {
		t.Error("HUMILIATING DEFEAT")
	}
}

type dummyPlugin struct{}

func (_ *dummyPlugin) Name() string                              { return "Dummy" }
func (_ *dummyPlugin) CanLoad() bool                             { return false }
func (_ *dummyPlugin) Load(_ *slack.Bot, _ ...interface{}) error { return nil }
func (_ *dummyPlugin) extraEnvVars() []string                    { return []string{} }
func (_ *dummyPlugin) setEnvVar(_, _ string)                     {}
func (_ *dummyPlugin) getEnvVar(_ string) string                 { return "" }

type dummyPluginWithEnvVars struct{}

func (_ *dummyPluginWithEnvVars) Name() string                              { return "Dummy With Env" }
func (_ *dummyPluginWithEnvVars) CanLoad() bool                             { return false }
func (_ *dummyPluginWithEnvVars) Load(_ *slack.Bot, _ ...interface{}) error { return nil }
func (_ *dummyPluginWithEnvVars) extraEnvVars() []string                    { return []string{"DUMMY"} }
func (_ *dummyPluginWithEnvVars) setEnvVar(_, _ string)                     {}
func (_ *dummyPluginWithEnvVars) getEnvVar(_ string) string                 { return "" }

func TestPluginCanLoad(t *testing.T) {
	log.SetLevel(log.PanicLevel)
	assert(!CanLoad(&dummyPlugin{}), t)

	env1 := stubEnv(map[string]string{
		"GH_CLIENT_ID":     "a",
		"GH_CLIENT_SECRET": "b",
		"GH_ACCESS_TOKEN":  "c",
	})
	defer restoreEnv(env1)
	assert(env1 != nil, t)
	assert(CanLoad(&dummyPlugin{}), t)

	assert(!CanLoad(&dummyPluginWithEnvVars{}), t)

	env2 := stubEnv(map[string]string{
		"DUMMY": "d",
	})
	defer restoreEnv(env2)
	assert(env2 != nil, t)
	assert(CanLoad(&dummyPluginWithEnvVars{}), t)
}

func TestPluginClient(t *testing.T) {
	env := stubEnv(map[string]string{
		"GH_CLIENT_ID":     "a",
		"GH_CLIENT_SECRET": "b",
	})
	defer restoreEnv(env)
	plugin := &dummyPlugin{}
	CanLoad(plugin)
	client := Client(plugin)
	assert(client != nil, t)
}

func TestPluginOAuthConfig(t *testing.T) {
	env := stubEnv(map[string]string{
		"GH_CLIENT_ID":     "a",
		"GH_CLIENT_SECRET": "b",
	})
	defer restoreEnv(env)
	plugin := &dummyPlugin{}
	CanLoad(plugin)
	conf := OAuthConfig(plugin, "localhost", []string{"scope1", "scope2"})
	assert(conf != nil, t)
	assert(conf.RedirectURL == "localhost", t)
	assert(len(conf.Scopes) == 2, t)
}

func TestPluginOAuthToken(t *testing.T) {
	env := stubEnv(map[string]string{
		"GH_CLIENT_ID":     "a",
		"GH_CLIENT_SECRET": "b",
	})
	defer restoreEnv(env)
	plugin := &dummyPlugin{}
	CanLoad(plugin)
	token := OAuthToken(plugin, "token type")
	assert(token != nil, t)
	assert(token.TokenType == "token type", t)
	assert(token.Expiry.IsZero(), t)
}
