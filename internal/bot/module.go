package bot

import (
	"strings"
	"time"

	twitch "github.com/gempir/go-twitch-irc"
)

const defaultModule = "default"

// Module is a named collection of Commands.
type Module struct {
	Commands []Command
	Enabled  bool
	Name     string
}

// AddCommand adds a command to the module.
func (m *Module) AddCommand(c Command) {
	m.Commands = append(m.Commands, c)
}

// Command is a command that a bot can use.
type Command interface {
	// Enables the command.
	Enable()
	// Determines if the command is currently enabled.
	Enabled() bool
	// Executes the command.
	Execute(b *Bot, args []string, channel string, user twitch.User, message twitch.Message)
	// Disables the command.
	Disable()
	// Determines if the command is currently on cooldown.
	IsOnCooldown() bool
	// Checks whether the given args will trigger the command.
	Match(s []string) bool
	// Name of the command.
	Name() string
	// Sets the last time that the command was used to the current time.
	SetLastUsed()
}

type command struct {
	enabled  bool
	Cooldown time.Duration
	LastUsed time.Time
	name     string
}

func (c *command) Enable() {
	c.enabled = true
}

func (c command) Enabled() bool {
	return c.enabled
}

func (c *command) Disable() {
	c.enabled = false
}

func (c command) IsOnCooldown() bool {
	return time.Now().Add(-c.Cooldown).Before(c.LastUsed)
}

func (c command) Match(s []string) bool {
	if len(s) < 1 {
		return false
	}
	return strings.ToLower(s[0]) == c.name
}

// Name of the command.
func (c command) Name() string {
	return c.name
}

func (c *command) SetLastUsed() {
	c.LastUsed = time.Now()
}
