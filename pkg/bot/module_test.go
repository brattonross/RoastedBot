package bot

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
			Cooldown: time.Second * 5,
			LastUsed: time.Now(),
		},
	}

	if !c.IsOnCooldown() {
		t.Error("expected command to be on cooldown")
	}

	c.LastUsed = time.Now().Add(-time.Second * 10)
	if c.IsOnCooldown() {
		t.Error("expected command to not be on cooldown")
	}
}
