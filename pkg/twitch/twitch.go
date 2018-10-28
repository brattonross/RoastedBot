package twitch

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

// Message represents a message that the bot sends.
type Message struct {
	Channel string
	Text    string
}

// Channel represents a twitch channel.
type Channel struct {
	EnabledModules map[string]bool
	Modules        map[string]*Module
	Name           string
}

// AddCommand adds a command to the given module.
// If the module does not exist, it will be created.
func (ch *Channel) AddCommand(module string, c Command) error {
	if _, ok := ch.Modules[module]; !ok {
		ch.Modules[module] = &Module{Name: module}
	}
	return ch.Modules[module].AddCommand(c)
}

// AddModule adds a new module to the channel.
func (ch *Channel) AddModule(name string) (*Module, error) {
	if _, ok := ch.Modules[name]; ok {
		return nil, fmt.Errorf("module '%s' already exists in channel '%s'", name, ch.Name)
	}
	ch.Modules[name] = &Module{Name: name, Commands: map[string]Command{}}
	return ch.Modules[name], nil
}

// EnableCommand enables a command in the given module.
func (ch *Channel) EnableCommand(module, command string) error {
	m, ok := ch.Modules[module]
	if !ok {
		return fmt.Errorf("module with name '%s' does not exist in channel '%s'", module, ch.Name)
	}
	return m.EnableCommand(command)
}

// EnableModule enables a module in the channel.
func (ch *Channel) EnableModule(module string) error {
	if _, ok := ch.Modules[module]; !ok {
		return fmt.Errorf("module with name '%s' does not exist in channel '%s'", module, ch.Name)
	}
	ch.EnabledModules[module] = true
	return nil
}

// DisableCommand disables a command in the given module.
func (ch *Channel) DisableCommand(module, command string) error {
	m, ok := ch.Modules[module]
	if !ok {
		return fmt.Errorf("module with name '%s' does not exist in channel '%s'", module, ch.Name)
	}
	return m.DisableCommand(command)
}

// DisableModule disables a module in the channel.
func (ch *Channel) DisableModule(module string) error {
	if _, ok := ch.Modules[module]; !ok {
		return fmt.Errorf("module with name '%s' does not exist in channel '%s'", module, ch.Name)
	}
	ch.EnabledModules[module] = false
	return nil
}

func newChannel(name string) *Channel {
	return &Channel{
		Modules: map[string]*Module{},
		Name:    name,
	}
}

func (ch *Channel) initDefaultModules() error {
	log.WithField("channel", ch.Name).Info("initialising default module")
	m, err := ch.AddModule(defaultModule)
	if err != nil {
		return err
	}
	m.AddCommand(Help{&command{cooldown: time.Second * 5, enabled: true, name: "help"}})
	m.AddCommand(Uptime{&command{cooldown: time.Second * 5, enabled: true, name: "uptime"}})
	ch.EnabledModules[defaultModule] = true
	return nil
}
