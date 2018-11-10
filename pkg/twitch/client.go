package twitch

import (
	"fmt"
	"sync"
	"time"

	"github.com/gempir/go-twitch-irc"
)

// Client is a wrapper of go-twitch-irc Client.
type Client struct {
	*twitch.Client

	start         time.Time
	channelsMutex *sync.Mutex
	channels      map[string]*Channel

	Username string
}

// NewClient creates a new Client using the given Config.
func NewClient(username string, client *twitch.Client) *Client {
	return &Client{
		channelsMutex: &sync.Mutex{},
		channels:      make(map[string]*Channel),
		Client:        client,
		start:         time.Now(),
		Username:      username,
	}
}

// AddChannel adds a channel to the Client, but does not join it.
func (b *Client) AddChannel(name string) error {
	return b.addChannel(newChannel(name))
}

func (b *Client) addChannel(c *Channel) error {
	if _, ok := b.channels[c.Name]; ok {
		return fmt.Errorf("Client already contains channel with name '%s'", c.Name)
	}
	b.channelsMutex.Lock()
	defer b.channelsMutex.Unlock()

	b.channels[c.Name] = c

	return nil
}

// AddCommand adds a command to the module in the channel.
// If the Client is not currently connected to the channel it will return an error.
// If the module does not already exist, it will be created.
func (b *Client) AddCommand(channel, module string, c *Command) error {
	ch, ok := b.channels[channel]
	if !ok {
		return fmt.Errorf("Client is not connected to channel '%s'", channel)
	}
	ch.AddCommand(module, c)
	return nil
}

// AddModule adds a new module with the given name to the given channel.
func (b *Client) AddModule(channel, module string) (*Module, error) {
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

// Channel gets the channel with the given name if the Client
// has it configured, otherwise it returns an error.
func (b *Client) Channel(name string) (*Channel, error) {
	ch, ok := b.channels[name]
	if !ok {
		return nil, fmt.Errorf("channel %s is not configured", name)
	}
	return ch, nil
}

// Channels returns the channels that the Client is currently connected to.
func (b *Client) Channels() []Channel {
	chans := []Channel{}
	for _, c := range b.channels {
		chans = append(chans, *c)
	}
	return chans
}

// EnableCommand enables a command in the given channel and module.
// The Client must be connected to the given channel, and the command must exist within the module.
func (b *Client) EnableCommand(channel, module, command string) error {
	ch, ok := b.channels[channel]
	if !ok {
		return fmt.Errorf("Client is not connected to channel '%s'", channel)
	}
	return ch.EnableCommand(module, command)
}

// EnableModule enables a module in the given channel.
func (b *Client) EnableModule(channel, module string) error {
	ch, ok := b.channels[channel]
	if !ok {
		return fmt.Errorf("Client is not connected to channel '%s'", channel)
	}
	return ch.EnableModule(module)
}

// DisableCommand disables a command in the given channel and module.
func (b *Client) DisableCommand(channel, module, command string) error {
	ch, ok := b.channels[channel]
	if !ok {
		return fmt.Errorf("Client is not connected to channel '%s'", channel)
	}
	return ch.DisableCommand(module, command)
}

// DisableModule disables a module in a channel.
func (b *Client) DisableModule(channel, module string) error {
	ch, ok := b.channels[channel]
	if !ok {
		return fmt.Errorf("Client is not connected to channel '%s'", channel)
	}
	return ch.DisableModule(module)
}

// JoinChannels joins all of the channels in the Client's channel list.
func (b *Client) JoinChannels() {
	for _, c := range b.channels {
		b.Join(c.Name)
	}
}
