package github

import (
	"testing"
)

func TestToken(t *testing.T) {
	token := Token()
	if !token.Expiry.IsZero() {
		t.Error("Error. Expected token expiration to be zero.")
	}
}
