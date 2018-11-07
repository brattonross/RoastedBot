package twitch

import (
	"fmt"
	"sync"
)

// Channel represents a twitch channel.
type Channel struct {
	enabledModules map[string]bool
	modulesMutex   *sync.Mutex
	modules        map[string]*Module

	Name string
}

const defaultModule = "default"

func newChannel(name string) *Channel {
	ch := &Channel{
		enabledModules: make(map[string]bool),
		modulesMutex:   &sync.Mutex{},
		modules:        make(map[string]*Module),
		Name:           name,
	}

	ch.addDefaultModule()

	return ch
}

// AddCommand adds a command to the given module.
// If the module does not exist, it will be created.
func (ch *Channel) AddCommand(module string, c *Command) error {
	if _, ok := ch.modules[module]; !ok {
		ch.modules[module] = newModule(module)
	}
	return ch.modules[module].AddCommand(c)
}

func (ch *Channel) addDefaultModule() {
	m, _ := ch.AddModule(defaultModule)

	m.AddCommand(helpCommand)
	m.EnableCommand(helpCommand.Name)

	m.AddCommand(uptimeCommand)
	m.EnableCommand(uptimeCommand.Name)

	m.AddCommand(enableCommand)
	m.EnableCommand(enableCommand.Name)

	// Default module should always be enabled.
	ch.EnableModule(defaultModule)
}

// AddModule adds a new module to the channel.
func (ch *Channel) AddModule(name string) (*Module, error) {
	if _, ok := ch.modules[name]; ok {
		return nil, fmt.Errorf("module '%s' already exists in channel '%s'", name, ch.Name)
	}
	ch.modules[name] = newModule(name)
	return ch.modules[name], nil
}

// EnableCommand enables a command in the given module.
func (ch *Channel) EnableCommand(module, command string) error {
	m, ok := ch.modules[module]
	if !ok {
		return fmt.Errorf("module with name '%s' does not exist in channel '%s'", module, ch.Name)
	}
	return m.EnableCommand(command)
}

// EnableModule enables a module in the channel.
func (ch *Channel) EnableModule(module string) error {
	if _, ok := ch.modules[module]; !ok {
		return fmt.Errorf("module with name '%s' does not exist in channel '%s'", module, ch.Name)
	}
	ch.enabledModules[module] = true
	return nil
}

// EnabledModules returns a list of the names of all enabled modules.
func (ch *Channel) EnabledModules() []string {
	mods := []string{}
	for _, m := range ch.modules {
		mods = append(mods, m.Name)
	}
	return mods
}

// DisableCommand disables a command in the given module.
func (ch *Channel) DisableCommand(module, command string) error {
	m, ok := ch.modules[module]
	if !ok {
		return fmt.Errorf("module with name '%s' does not exist in channel '%s'", module, ch.Name)
	}
	return m.DisableCommand(command)
}

// DisableModule disables a module in the channel.
func (ch *Channel) DisableModule(module string) error {
	if _, ok := ch.modules[module]; !ok {
		return fmt.Errorf("module with name '%s' does not exist in channel '%s'", module, ch.Name)
	}
	ch.enabledModules[module] = false
	return nil
}

// IsModuleEnabled determines if the module is enabled.
func (ch *Channel) isModuleEnabled(module string) bool {
	enabled, ok := ch.enabledModules[module]
	return ok && enabled
}

// MatchCommand returns the Command that is triggered by the given args,
// as well as the module that it belongs to.
func (ch *Channel) matchCommand(args []string) (command *Command, module *Module) {
	for _, m := range ch.modules {
		if !ch.isModuleEnabled(m.Name) {
			continue
		}

		for _, c := range m.commands {
			if !c.match(args[0]) {
				continue
			}
			return c, m
		}
	}
	return nil, nil
}

// Modules returns a list of all modules in the channel.
func (ch *Channel) Modules() []Module {
	mods := []Module{}
	for _, m := range ch.modules {
		mods = append(mods, *m)
	}
	return mods
}
