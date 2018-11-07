package twitch

import (
	"testing"
	"time"

	twitch "github.com/gempir/go-twitch-irc"
)

func TestExecuteNilCommand(t *testing.T) {
	var c *Command
	if err := c.execute(nil, nil, "", twitch.User{}, twitch.Message{}); err == nil {
		t.Error("expected calling execute on a nil Command to throw an error")
	}
}

func TestExecuteRuns(t *testing.T) {
	called := false
	c := &Command{
		Run: func(b *Bot, args []string, channel string, user twitch.User, message twitch.Message) {
			called = true
		},
		Use: "test",
	}
	if err := c.execute(nil, nil, "", twitch.User{}, twitch.Message{}); err != nil {
		t.Errorf("executing Command returned unexpected error: %v", err)
	}
	if !called {
		t.Error("Command's Run method was never called")
	}
}

func TestAssertExecuteNoRun(t *testing.T) {
	c := &Command{Use: "test"}
	if err := c.execute(nil, nil, "", twitch.User{}, twitch.Message{}); err == nil {
		t.Error("expected Command with no Run function to throw an error")
	}
}

func TestAssertExecuteNoUse(t *testing.T) {
	c := &Command{Run: func(b *Bot, args []string, channel string, user twitch.User, message twitch.Message) {}}
	if err := c.execute(nil, nil, "", twitch.User{}, twitch.Message{}); err == nil {
		t.Error("expected Command with no Use to throw an error")
	}
}

func TestIsOnCooldown(t *testing.T) {
	c := &Command{
		Cooldown: time.Second * 5,
		LastUsed: time.Now(),
	}

	if !c.isOnCooldown() {
		t.Error("expected command to be on cooldown")
	}

	c.LastUsed = time.Now().Add(-time.Second * 10)
	if c.isOnCooldown() {
		t.Error("expected command to not be on cooldown")
	}
}

func TestIsOnCooldown_NoCooldown(t *testing.T) {
	c := &Command{LastUsed: time.Now()}
	if c.isOnCooldown() {
		t.Error("expected Command to not be on cooldown")
	}
}

func TestIsOnCooldown_NeverUsed(t *testing.T) {
	c := &Command{Cooldown: time.Second}
	if c.isOnCooldown() {
		t.Error("expected Command to not be on cooldown")
	}
}

func TestMatch(t *testing.T) {
	good := "test"
	bad := "abc"
	c := &Command{Use: good}
	if c.match(bad) {
		t.Errorf("Command unexpectedly matched on string %s", bad)
	}
	if !c.match(good) {
		t.Errorf("Command unexpectedly did not match on string %s", good)
	}
}

func TestMatch_EmptyString(t *testing.T) {
	c := &Command{Use: "test"}
	if c.match("") {
		t.Error("Command unexpectedly matched on empty string")
	}
}
