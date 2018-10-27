package bot

import (
	"fmt"
	"strings"
	"time"

	"github.com/brattonross/roastedbot/pkg/bot/service"
	"github.com/gempir/go-twitch-irc"
	log "github.com/sirupsen/logrus"
)

// Config for the bot.
type Config struct {
	Username string   `json:"username"`
	OAuth    string   `json:"oauth"`
	Channels []string `json:"channels"`
}

// Message represents a message that the bot sends.
type Message struct {
	Channel string
	Text    string
}

// Bot is the bot xD
type Bot struct {
	config  Config
	start   time.Time
	client  *twitch.Client
	service service.BotServiceServer

	say chan Message

	Modules map[string]*Module
}

// New creates a new Bot using the given Config.
func New(config Config) *Bot {
	return &Bot{
		config: config,
		client: twitch.NewClient(config.Username, config.OAuth),
		say:    make(chan Message),
		start:  time.Now(),
	}
}

// AddCommand adds a command to the module with the given name.
// If the module does not exist, it is created.
func (b *Bot) AddCommand(module string, c Command) {
	m, ok := b.Modules[module]
	if !ok {
		b.Modules[module] = &Module{Name: module}
	}
	m.AddCommand(c)
}

// Connect will connect the bot to Twitch IRC.
func (b *Bot) Connect() error {
	return b.client.Connect()
}

// EnableCommand enables a command.
func (b *Bot) EnableCommand(module, command string) error {
	module = strings.ToLower(module)
	command = strings.ToLower(command)

	m, ok := b.Modules[module]
	if !ok {
		return fmt.Errorf("module with name '%s' does not exist", module)
	}
	for _, c := range m.Commands {
		if c.Name() == command {
			c.Enable()
			return nil
		}
	}
	return fmt.Errorf("command with name '%s' does not exist in module '%s'", command, module)
}

// EnableModule enables a module.
func (b *Bot) EnableModule(module string) error {
	module = strings.ToLower(module)

	m, ok := b.Modules[module]
	if !ok {
		return fmt.Errorf("module with name '%s' does not exist", module)
	}
	m.Enabled = true
	return nil
}

// DisableCommand disables a command.
func (b *Bot) DisableCommand(module, command string) error {
	module = strings.ToLower(module)
	command = strings.ToLower(command)

	m, ok := b.Modules[module]
	if !ok {
		return fmt.Errorf("module with name '%s' does not exist", module)
	}
	for _, c := range m.Commands {
		if c.Name() == command {
			c.Disable()
			return nil
		}
	}
	return fmt.Errorf("command with name '%s' does not exist in module '%s'", command, module)
}

// DisableModule disables a module.
func (b *Bot) DisableModule(module string) error {
	module = strings.ToLower(module)

	m, ok := b.Modules[module]
	if !ok {
		return fmt.Errorf("module with name '%s' does not exist", module)
	}
	m.Enabled = false
	return nil
}

// Init will initialize the bot with sensible defaults.
func (b *Bot) Init() {
	b.AddCommand(defaultModule, Help{&command{Cooldown: time.Second * 5, enabled: true, name: "help"}})
	b.AddCommand(defaultModule, Uptime{&command{Cooldown: time.Second * 5, enabled: true, name: "uptime"}})

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

// JoinChannel joins the given channel.
func (b *Bot) JoinChannel(channel string) {
	b.client.Join(channel)
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
		for _, m := range b.Modules {
			for _, c := range m.Commands {
				if !c.Match(args) {
					continue
				}

				log.WithField("command", c.Name()).Info("matched command")
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
						"time":    time.Now().Sub(start),
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
