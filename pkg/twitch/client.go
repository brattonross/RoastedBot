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

	channels      map[string]*Channel
	channelsMutex *sync.Mutex
	rateLimit     <-chan time.Time
	start         time.Time

	Username string
}

// NewClient creates a new Client using the given Config.
func NewClient(username string, client *twitch.Client) *Client {
	return &Client{
		channelsMutex: &sync.Mutex{},
		channels:      make(map[string]*Channel),
		Client:        client,
		rateLimit:     time.Tick(time.Millisecond * 1500),
		start:         time.Now(),
		Username:      username,
	}
}

// AddChannel adds a channel to the Client, but does not join it.
func (cl *Client) AddChannel(name string) error {
	return cl.addChannel(newChannel(name))
}

func (cl *Client) addChannel(ch *Channel) error {
	if _, ok := cl.channels[ch.Name]; ok {
		return fmt.Errorf("Client already contains channel with name '%s'", ch.Name)
	}
	cl.channelsMutex.Lock()
	defer cl.channelsMutex.Unlock()

	cl.channels[ch.Name] = ch

	return nil
}

// AddCommand adds a command to the module in the channel.
// If the Client is not currently connected to the channel it will return an error.
// If the module does not already exist, it will be created.
func (cl *Client) AddCommand(channel, module string, c *Command) error {
	ch, ok := cl.channels[channel]
	if !ok {
		return fmt.Errorf("Client is not connected to channel '%s'", channel)
	}
	ch.AddCommand(module, c)
	return nil
}

// AddModule adds a new module with the given name to the given channel.
func (cl *Client) AddModule(channel, module string) (*Module, error) {
	ch, ok := cl.channels[channel]
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
func (cl *Client) Channel(name string) (*Channel, error) {
	ch, ok := cl.channels[name]
	if !ok {
		return nil, fmt.Errorf("channel %s is not configured", name)
	}
	return ch, nil
}

// Channels returns the channels that the Client is currently connected to.
func (cl *Client) Channels() []Channel {
	chans := []Channel{}
	for _, c := range cl.channels {
		chans = append(chans, *c)
	}
	return chans
}

// EnableCommand enables a command in the given channel and module.
// The Client must be connected to the given channel, and the command must exist within the module.
func (cl *Client) EnableCommand(channel, module, command string) error {
	ch, ok := cl.channels[channel]
	if !ok {
		return fmt.Errorf("Client is not connected to channel '%s'", channel)
	}
	return ch.EnableCommand(module, command)
}

// EnableModule enables a module in the given channel.
func (cl *Client) EnableModule(channel, module string) error {
	ch, ok := cl.channels[channel]
	if !ok {
		return fmt.Errorf("Client is not connected to channel '%s'", channel)
	}
	return ch.EnableModule(module)
}

// DisableCommand disables a command in the given channel and module.
func (cl *Client) DisableCommand(channel, module, command string) error {
	ch, ok := cl.channels[channel]
	if !ok {
		return fmt.Errorf("Client is not connected to channel '%s'", channel)
	}
	return ch.DisableCommand(module, command)
}

// DisableModule disables a module in a channel.
func (cl *Client) DisableModule(channel, module string) error {
	ch, ok := cl.channels[channel]
	if !ok {
		return fmt.Errorf("Client is not connected to channel '%s'", channel)
	}
	return ch.DisableModule(module)
}

// JoinChannels joins all of the channels in the Client's channel list.
func (cl *Client) JoinChannels() {
	for _, c := range cl.channels {
		cl.Join(c.Name)
	}
}

// Say will send a message to twitch irc.
// Messages are limited to sending as per the rate limit.
func (cl *Client) Say(channel, text string) {
	<-cl.rateLimit
	cl.Client.Say(channel, text)
}
