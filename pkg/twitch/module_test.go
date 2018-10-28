package twitch

import (
	"testing"
	"time"
)

type testCommand struct {
	*command
}

func TestCommand_IsOnCooldown(t *testing.T) {
	c := testCommand{
		&command{
			cooldown: time.Second * 5,
			lastUsed: time.Now(),
		},
	}

	if !c.IsOnCooldown() {
		t.Error("expected command to be on cooldown")
	}

	c.lastUsed = time.Now().Add(-time.Second * 10)
	if c.IsOnCooldown() {
		t.Error("expected command to not be on cooldown")
	}
}
