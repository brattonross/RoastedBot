package twitch

import (
	"fmt"
	"sync"
	"time"

	"github.com/gempir/go-twitch-irc"
)

// Bot is the bot xD
type Bot struct {
	*twitch.Client

	start         time.Time
	channelsMutex *sync.Mutex
	channels      map[string]*Channel

	Username string
}

// NewBot creates a new Bot using the given Config.
func NewBot(username string, client *twitch.Client) *Bot {
	return &Bot{
		channelsMutex: &sync.Mutex{},
		channels:      make(map[string]*Channel),
		Client:        client,
		start:         time.Now(),
		Username:      username,
	}
}

// AddChannel adds a channel to the bot, but does not join it.
func (b *Bot) AddChannel(name string) error {
	return b.addChannel(newChannel(name))
}

func (b *Bot) addChannel(c *Channel) error {
	if _, ok := b.channels[c.Name]; ok {
		return fmt.Errorf("bot already contains channel with name '%s'", c.Name)
	}
	b.channelsMutex.Lock()
	defer b.channelsMutex.Unlock()

	b.channels[c.Name] = c

	return nil
}

// AddCommand adds a command to the module in the channel.
// If the bot is not currently connected to the channel it will return an error.
// If the module does not already exist, it will be created.
func (b *Bot) AddCommand(channel, module string, c *Command) error {
	ch, ok := b.channels[channel]
	if !ok {
		return fmt.Errorf("bot is not connected to channel '%s'", channel)
	}
	ch.AddCommand(module, c)
	return nil
}

// AddModule adds a new module with the given name to the given channel.
func (b *Bot) AddModule(channel, module string) (*Module, error) {
	ch, ok := b.channels[channel]
	if !ok {
		return nil, fmt.Errorf("channel '%s' is not configured", channel)
	}
	m, err := ch.AddModule(module)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// Channel gets the channel with the given name if the bot
// has it configured, otherwise it returns an error.
func (b *Bot) Channel(name string) (*Channel, error) {
	ch, ok := b.channels[name]
	if !ok {
		return nil, fmt.Errorf("channel %s is not configured", name)
	}
	return ch, nil
}

// Channels returns the channels that the bot is currently connected to.
func (b *Bot) Channels() []Channel {
	chans := []Channel{}
	for _, c := range b.channels {
		chans = append(chans, *c)
	}
	return chans
}

// EnableCommand enables a command in the given channel and module.
// The bot must be connected to the given channel, and the command must exist within the module.
func (b *Bot) EnableCommand(channel, module, command string) error {
	ch, ok := b.channels[channel]
	if !ok {
		return fmt.Errorf("bot is not connected to channel '%s'", channel)
	}
	return ch.EnableCommand(module, command)
}

// EnableModule enables a module in the given channel.
func (b *Bot) EnableModule(channel, module string) error {
	ch, ok := b.channels[channel]
	if !ok {
		return fmt.Errorf("bot is not connected to channel '%s'", channel)
	}
	return ch.EnableModule(module)
}

// DisableCommand disables a command in the given channel and module.
func (b *Bot) DisableCommand(channel, module, command string) error {
	ch, ok := b.channels[channel]
	if !ok {
		return fmt.Errorf("bot is not connected to channel '%s'", channel)
	}
	return ch.DisableCommand(module, command)
}

// DisableModule disables a module in a channel.
func (b *Bot) DisableModule(channel, module string) error {
	ch, ok := b.channels[channel]
	if !ok {
		return fmt.Errorf("bot is not connected to channel '%s'", channel)
	}
	return ch.DisableModule(module)
}

// JoinChannels joins all of the channels in the bot's channel list.
func (b *Bot) JoinChannels() {
	for _, c := range b.channels {
		b.Join(c.Name)
	}
}
