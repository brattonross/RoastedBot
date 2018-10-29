package twitch

import (
	"fmt"
	"time"

	twitch "github.com/gempir/go-twitch-irc"
)

const defaultModule = "default"

// Module is a named collection of Commands.
type Module struct {
	Commands        map[string]Command
	EnabledCommands map[string]bool
	Name            string
}

func newModule(name string) *Module {
	return &Module{
		Commands:        map[string]Command{},
		EnabledCommands: map[string]bool{},
		Name:            name,
	}
}

// AddCommand adds a command to the module.
func (m *Module) AddCommand(c Command) error {
	if _, ok := m.Commands[c.Name()]; ok {
		return fmt.Errorf("command '%s' already exists in module '%s'", c.Name(), m.Name)
	}
	m.Commands[c.Name()] = c
	return nil
}

// EnableCommand enables a command within the module.
func (m *Module) EnableCommand(command string) error {
	for _, c := range m.Commands {
		if c.Name() == command {
			m.EnabledCommands[command] = true
			return nil
		}
	}
	return fmt.Errorf("command with name '%s' does not exist in module '%s'", command, m.Name)
}

// DisableCommand disables a command in the module.
func (m *Module) DisableCommand(command string) error {
	for _, c := range m.Commands {
		if c.Name() == command {
			m.EnabledCommands[command] = false
			return nil
		}
	}
	return fmt.Errorf("command with name '%s' does not exist in module '%s'", command, m.Name)
}

// Command is a command that a bot can use.
type Command interface {
	// Cooldown of the command.
	Cooldown() time.Duration
	// Executes the command.
	Execute(b *Bot, args []string, channel string, user twitch.User, message twitch.Message)
	// Determines if the command is currently on cooldown.
	IsOnCooldown() bool
	// Checks whether the given args will trigger the command.
	Match(s string) bool
	// Name of the command.
	Name() string
	// Sets the last time that the command was used to the current time.
	SetLastUsed()
}

type command struct {
	cooldown time.Duration
	lastUsed time.Time
	name     string
}

func (c *command) Cooldown() time.Duration {
	return c.cooldown
}

func (c command) IsOnCooldown() bool {
	return time.Now().Add(-c.cooldown).Before(c.lastUsed)
}

func (c command) Match(s string) bool {
	if len(s) < 1 {
		return false
	}
	return s == c.name
}

// Name of the command.
func (c command) Name() string {
	return c.name
}

func (c *command) SetLastUsed() {
	c.lastUsed = time.Now()
}
