package twitch

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gempir/go-twitch-irc"
	log "github.com/sirupsen/logrus"
)

// Config for the bot.
type Config struct {
	Username string   `json:"username"`
	OAuth    string   `json:"oauth"`
	Channels []string `json:"channels"`
	GRPC     struct {
		Port int `json:"port"`
	} `json:"grpc"`
}

// Bot is the bot xD
type Bot struct {
	config Config
	start  time.Time

	channelsMutex *sync.Mutex
	channels      map[string]*Channel

	client *twitch.Client
	db     *sql.DB
}

// NewBot creates a new Bot using the given Config.
func NewBot(config Config, client *twitch.Client, db *sql.DB) *Bot {
	b := &Bot{
		channelsMutex: &sync.Mutex{},
		channels:      make(map[string]*Channel),
		config:        config,
		client:        client,
		db:            db,
		start:         time.Now(),
	}

	b.client.OnNewMessage(onNewMessage(b))

	return b
}

// AddChannel adds a channel to the bot, but does not join it.
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

// Channels returns the channels that the bot is currently connected to.
func (b *Bot) Channels() []Channel {
	chans := []Channel{}
	for _, c := range b.channels {
		chans = append(chans, *c)
	}
	return chans
}

// Connect will connect the bot to Twitch IRC.
func (b *Bot) Connect() error {
	return b.client.Connect()
}

// Disconnect disconnects the bot from Twitch IRC.
func (b *Bot) Disconnect() error {
	return b.client.Disconnect()
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

// LoadChannels laods channels that should be joined on start.
func (b *Bot) LoadChannels() {
	// TODO: fetch channels from db here

	for _, c := range b.config.Channels {
		if err := b.addChannel(newChannel(c)); err != nil {
			log.WithField("channel", c).Error(err)
		}
	}
}

// JoinChannel joins the given channel, and initialises the default modules for the channel.
func (b *Bot) JoinChannel(channel string) {
	b.client.Join(channel)
}

// JoinChannels joins all of the channels in the bot's channel list.
func (b *Bot) JoinChannels() {
	for _, c := range b.channels {
		b.JoinChannel(c.Name)
	}
}

// OnConnect sets the callback for when the bot connects to twitch.
func (b *Bot) OnConnect(f func()) {
	b.client.OnConnect(f)
}

// Say sends the given text to the channel.
func (b *Bot) Say(channel, text string) {
	b.client.Say(channel, text)
}

func onNewMessage(b *Bot) func(channel string, user twitch.User, message twitch.Message) {
	return func(channel string, user twitch.User, message twitch.Message) {
		username := strings.ToLower(b.config.Username)
		if strings.ToLower(user.Username) == username {
			return
		}

		ch, ok := b.channels[channel]
		if !ok {
			log.WithField("channel", channel).Error("cannot handle message: bot is not connected to channel")
			return
		}

		args := strings.Split(message.Text, " ")
		if len(args) < 1 {
			return
		}

		first := strings.ToLower(args[0])
		last := strings.ToLower(args[len(args)-1])

		if first == "!xd" {
			b.Say(channel, "xD")
			return
		} else if first == "!php" {
			b.Say(channel, "PHPDETECTED")
			return
		}

		// No mentions, don't process
		if !isMention(first, b.config.Username) && !isMention(last, b.config.Username) {
			return
		}

		// Only message is a mention of the bot, say hi
		if len(args) == 1 {
			b.Say(channel, fmt.Sprintf("hi %s :)", user.DisplayName))
			return
		}

		if isMention(first, b.config.Username) {
			args = args[1:]
		} else if isMention(last, b.config.Username) {
			args = args[:len(args)-1]
		}

		log.WithFields(log.Fields{
			"channel": channel,
			"text":    message.Text,
			"user":    user.DisplayName,
		}).Info("handling message")

		command, module := ch.matchCommand(args)
		if command == nil {
			return
		}

		if !module.isCommandEnabled(command.Name) {
			log.WithFields(log.Fields{
				"channel": channel,
				"command": command.Name,
				"module":  module.Name,
				"user":    user.DisplayName,
			}).Info("command is not enabled")
			return
		}
		if command.isOnCooldown() {
			log.WithFields(log.Fields{
				"channel": channel,
				"command": command.Name,
				"module":  module.Name,
				"user":    user.DisplayName,
			}).Info("command is on cooldown")
			return
		}

		go func() {
			log.WithFields(log.Fields{
				"channel": channel,
				"command": command.Name,
				"module":  module.Name,
				"user":    user.DisplayName,
			}).Info("executing command")
			start := time.Now()

			defer log.WithFields(log.Fields{
				"channel": channel,
				"command": command.Name,
				"delta":   fmt.Sprintf("%dms", time.Now().Sub(start)/time.Millisecond),
				"module":  module.Name,
				"user":    user.DisplayName,
			}).Info("finished executing command")

			command.execute(b, args, channel, user, message)
			command.LastUsed = time.Now()
		}()
	}
}

func isMention(s, username string) bool {
	if strings.HasPrefix(s, "@") {
		s = s[1:]
	}
	if strings.HasSuffix(s, ",") {
		s = s[:len(s)-1]
	}
	return strings.ToLower(username) == s
}
