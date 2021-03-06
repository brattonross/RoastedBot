package twitch

import (
	"fmt"
	"sync"
)

// Module is a named collection of Commands.
type Module struct {
	commands             map[string]*Command
	commandsMutex        *sync.Mutex
	enabledCommands      map[string]bool
	enabledCommandsMutex *sync.Mutex

	Name string
}

// Create a new module with the given name.
func newModule(name string) *Module {
	return &Module{
		commands:             make(map[string]*Command),
		commandsMutex:        &sync.Mutex{},
		enabledCommands:      make(map[string]bool),
		enabledCommandsMutex: &sync.Mutex{},
		Name:                 name,
	}
}

// AddCommand adds a command to the module.
func (m *Module) AddCommand(c *Command) error {
	m.commandsMutex.Lock()
	defer m.commandsMutex.Unlock()
	if c == nil {
		return fmt.Errorf("attempted to add a nil Command to the module %s", m.Name)
	}
	if _, ok := m.commands[c.Name]; ok {
		return fmt.Errorf("command '%s' already exists in module '%s'", c.Name, m.Name)
	}

	m.commands[c.Name] = c
	return nil
}

// Commands returns all of the commands in this module.
func (m *Module) Commands() []Command {
	commands := []Command{}
	for _, c := range m.commands {
		commands = append(commands, *c)
	}
	return commands
}

// EnableCommand enables a command within the module.
func (m *Module) EnableCommand(command string) error {
	m.enabledCommandsMutex.Lock()
	defer m.enabledCommandsMutex.Unlock()
	for _, c := range m.commands {
		if c.Name == command {
			m.enabledCommands[command] = true
			return nil
		}
	}
	return fmt.Errorf("command with name '%s' does not exist in module '%s'", command, m.Name)
}

// DisableCommand disables a command in the module.
func (m *Module) DisableCommand(command string) error {
	m.enabledCommandsMutex.Lock()
	defer m.enabledCommandsMutex.Unlock()
	for _, c := range m.commands {
		if c.Name == command {
			m.enabledCommands[command] = false
			return nil
		}
	}
	return fmt.Errorf("command with name '%s' does not exist in module '%s'", command, m.Name)
}

// IsCommandEnabled determines if a command is enabled.
func (m *Module) IsCommandEnabled(command string) bool {
	m.enabledCommandsMutex.Lock()
	defer m.enabledCommandsMutex.Unlock()
	enabled, ok := m.enabledCommands[command]
	return ok && enabled
}
