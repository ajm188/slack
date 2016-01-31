package slack

import (
	"time"
)

func formatTime(t time.Time) string {
	// returns MMDDYYhhmmss
	return t.Format("010206150405")
}
