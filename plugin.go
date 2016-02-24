package slack

import (
	log "github.com/Sirupsen/logrus"
)

// PluginFunc is a function which returns a Plugin. This can be used if your
// plugin requires any special setup.
type PluginFunc func() Plugin

// The Plugin interface is used to implement a plugin for a Slack bot.
//
// Name should return the name of the plugin. This is useful for logging.
//
// CanLoad should return true if the plugin can be loaded without error. If
// CanLoad returns false, the plugin will not be loaded.
//
// Load should contain the code necessary to actually load the plugin. It
// receives a reference to the bot, in order to register handlers, or anything
// else the plugin may need to do to load itself into the bot. It also receives
// a variable-length list of arbitrary arguments that the plugin may need to
// load correctly, such as configuration options. Load should return an error
// if there were any issues loading the plugin.
type Plugin interface {
	Name() string
	CanLoad() bool
	Load(*Bot, ...interface{}) error
}

// UsePlugin will load a plugin if the plugin can load. `args` is a
// variable-length list of arguments that will be passed to the plugin's Load
// function.
func (bot *Bot) UsePlugin(pluginF PluginFunc, args ...interface{}) error {
	plugin := pluginF()
	if plugin.CanLoad() {
		if err := plugin.Load(bot, args...); err != nil {
			return err
		}
	} else {
		log.WithFields(log.Fields{
			"plugin": plugin.Name(),
		}).Error("Failed to load plugin.")
	}
	return nil
}
