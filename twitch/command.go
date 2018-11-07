package twitch

import (
	"fmt"
	"time"

	twitch "github.com/gempir/go-twitch-irc"
)

// Command is a command that a bot can use.
type Command struct {
	// Cooldown of the command.
	Cooldown time.Duration
	// Function to run when the command is executed.
	Run func(b *Bot, args []string, channel string, user twitch.User, message twitch.Message)
	// The last time that the command was invoked successfully.
	LastUsed time.Time
	// Name of the command.
	Name string
	// Usage of the command.
	Use string
}

// Executes the command.
func (c *Command) execute(b *Bot, args []string, channel string, user twitch.User, message twitch.Message) error {
	if c == nil {
		return fmt.Errorf("attempted to execute a nil Command")
	}
	if c.Run == nil || c.Use == "" {
		return fmt.Errorf("attempted to execute an unconfigured Command")
	}
	c.Run(b, args, channel, user, message)
	return nil
}

func (c Command) isOnCooldown() bool {
	return time.Now().Add(-c.Cooldown).Before(c.LastUsed)
}

func (c Command) match(s string) bool {
	if len(s) < 1 {
		return false
	}
	return s == c.Use
}
