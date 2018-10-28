package twitch

import (
	"fmt"
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc"
	log "github.com/sirupsen/logrus"
)

// Config for the bot.
type Config struct {
	Username string   `json:"username"`
	OAuth    string   `json:"oauth"`
	Channels []string `json:"channels"`
}

// Bot is the bot xD
type Bot struct {
	config Config
	start  time.Time
	client *twitch.Client

	say chan Message

	Channels map[string]*Channel
}

// NewBot creates a new Bot using the given Config.
func NewBot(config Config) *Bot {
	return &Bot{
		Channels: map[string]*Channel{},
		config:   config,
		client:   twitch.NewClient(config.Username, config.OAuth),
		say:      make(chan Message),
		start:    time.Now(),
	}
}

// AddCommand adds a command to the module in the channel.
// If the bot is not currently connected to the channel it will return an error.
// If the module does not already exist, it will be created.
func (b *Bot) AddCommand(channel, module string, c Command) error {
	ch, ok := b.Channels[channel]
	if !ok {
		return fmt.Errorf("bot is not connected to channel '%s'", channel)
	}
	ch.AddCommand(module, c)
	return nil
}

// Connect will connect the bot to Twitch IRC.
func (b *Bot) Connect() error {
	return b.client.Connect()
}

// EnableCommand enables a command in the given channel and module.
// The bot must be connected to the given channel, and the command must exist within the module.
func (b *Bot) EnableCommand(channel, module, command string) error {
	ch, ok := b.Channels[channel]
	if !ok {
		return fmt.Errorf("bot is not connected to channel '%s'", channel)
	}
	return ch.EnableCommand(module, command)
}

// EnableModule enables a module in the given channel.
func (b *Bot) EnableModule(channel, module string) error {
	ch, ok := b.Channels[channel]
	if !ok {
		return fmt.Errorf("bot is not connected to channel '%s'", channel)
	}
	return ch.EnableModule(module)
}

// DisableCommand disables a command in the given channel and module.
func (b *Bot) DisableCommand(channel, module, command string) error {
	ch, ok := b.Channels[channel]
	if !ok {
		return fmt.Errorf("bot is not connected to channel '%s'", channel)
	}
	return ch.DisableCommand(module, command)
}

// DisableModule disables a module in a channel.
func (b *Bot) DisableModule(channel, module string) error {
	ch, ok := b.Channels[channel]
	if !ok {
		return fmt.Errorf("bot is not connected to channel '%s'", channel)
	}
	return ch.DisableModule(module)
}

// Init will initialize the bot with sensible defaults.
func (b *Bot) Init() {
	b.client.OnConnect(func() {
		log.Info("successfully connected to twitch!")
	})
	b.client.OnNewMessage(onNewMessage(b))

	go func(b *Bot) {
		for {
			select {
			case msg := <-b.say:
				b.client.Say(msg.Channel, msg.Text)
			default:
			}
		}
	}(b)

	for _, channel := range b.config.Channels {
		b.JoinChannel(channel)
		log.WithFields(log.Fields{"channel": channel}).Info("joined channel")
	}
}

// JoinChannel joins the given channel, and initialises the default modules for the channel.
func (b *Bot) JoinChannel(channel string) error {
	if _, ok := b.Channels[channel]; ok {
		return fmt.Errorf("bot is already connected to channel '%s'", channel)
	}
	b.Channels[channel] = newChannel(channel)
	if err := b.Channels[channel].initDefaultModules(); err != nil {
		log.WithField("channel", channel).Warnf("failed to initialise default modules for channel: %v", err)
	}
	b.client.Join(channel)
	return nil
}

// Say sends a message to the Bot's say channel.
func (b *Bot) Say(channel, text string) {
	b.say <- Message{channel, text}
}

func onNewMessage(b *Bot) func(channel string, user twitch.User, message twitch.Message) {
	return func(channel string, user twitch.User, message twitch.Message) {
		username := strings.ToLower(b.config.Username)
		if strings.ToLower(user.Username) == username {
			return
		}

		ch, ok := b.Channels[channel]
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
			log.WithFields(log.Fields{"channel": channel}).Info("sending xD Message")
			b.Say(channel, "xD")
			return
		} else if first == "!bot" {
			log.WithFields(log.Fields{"channel": channel}).Info("sending bot Message")
			b.Say(channel, "I'm roastedb's bot, written in Go pajaH")
			return
		} else if first == "!php" {
			log.WithFields(log.Fields{"channel": channel}).Info("sending php Message")
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

		log.WithField("text", message.Text).Info("handling message")
		for _, m := range ch.Modules {
			enabled, ok := ch.EnabledModules[m.Name]
			if !ok || !enabled {
				continue
			}

			for _, c := range m.Commands {
				if !c.Match(args[0]) {
					continue
				}

				if !c.Enabled() {
					log.WithField("command", c.Name()).Info("command is not enabled")
					return
				}
				if c.IsOnCooldown() {
					log.WithField("command", c.Name()).Warn("command is on cooldown")
					return
				}

				go func(command Command) {
					// TODO: logging middleware?
					log.WithField("command", c.Name()).Info("executing command")
					start := time.Now()
					defer log.WithFields(log.Fields{
						"command": c.Name(),
						"time":    fmt.Sprintf("%dms", time.Now().Sub(start)/time.Millisecond),
					}).Info("finished executing command")
					command.Execute(b, args, channel, user, message)
					command.SetLastUsed()
				}(c)
				return
			}
		}
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
