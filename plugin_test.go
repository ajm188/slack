package slack

import (
	"testing"
)

type cannotLoad struct{}

func cannotLoadF() Plugin           { return &cannotLoad{} }
func (_ *cannotLoad) Name() string  { return "cannot load" }
func (_ *cannotLoad) CanLoad() bool { return false }
func (_ *cannotLoad) Load(_ *Bot, args ...interface{}) error {
	t := args[0].(testing.T)
	t.Error("Error. Load was called but should not have been")
	return nil
}

type loadsWithError struct{}

func loadsWithErrorF() Plugin           { return &loadsWithError{} }
func (_ *loadsWithError) Error() string { return "implementing error interface" }
func (_ *loadsWithError) Name() string  { return "loads with error" }
func (_ *loadsWithError) CanLoad() bool { return true }
func (plugin *loadsWithError) Load(_ *Bot, args ...interface{}) error {
	return plugin
}

type loadsCleanly struct{}

func loadsCleanlyF() Plugin                                    { return &loadsCleanly{} }
func (_ *loadsCleanly) Name() string                           { return "loads cleanly" }
func (_ *loadsCleanly) CanLoad() bool                          { return true }
func (_ *loadsCleanly) Load(_ *Bot, args ...interface{}) error { return nil }

func TestUsePlugin(t *testing.T) {
	b := NewBot("token")

	err := b.UsePlugin(cannotLoadF, t)
	assert(err == nil, t)

	err = b.UsePlugin(loadsWithErrorF)
	assert(err != nil, t)

	err = b.UsePlugin(loadsCleanlyF)
	assert(err == nil, t)
}
