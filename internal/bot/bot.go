package bot

import (
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

// Command is a command.
type Command interface {
	Execute(b *Bot, args []string, channel string, user twitch.User, message twitch.Message)
	Match(s []string) bool
	Name() string
}

// Message represents a message that the bot sends.
type message struct {
	Channel string
	Text    string
}

// Bot is the bot xD
type Bot struct {
	config Config
	start  time.Time
	client *twitch.Client

	out chan message

	Commands []Command
}

// New creates a new Bot using the given Config.
func New(config Config) *Bot {
	return &Bot{
		config: config,
		client: twitch.NewClient(config.Username, config.OAuth),
		out:    make(chan message),
		start:  time.Now(),
	}
}

// AddCommand adds a command to the bot.
func (b *Bot) AddCommand(c Command) {
	b.Commands = append(b.Commands, c)
}

// Connect will connect the bot to Twitch IRC.
func (b *Bot) Connect() error {
	return b.client.Connect()
}

// Init will initialize the bot with sensible defaults.
func (b *Bot) Init() {
	b.AddCommand(Help{})
	b.AddCommand(Uptime{})

	b.client.OnConnect(func() {
		log.Info("successfully connected to twitch!")
	})
	b.client.OnNewMessage(onNewMessage(b))

	go func(b *Bot) {
		for {
			select {
			case msg := <-b.out:
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

// Say sends a message to the Bot's out channel.
func (b *Bot) Say(m message) {
	b.out <- m
}

func onNewMessage(b *Bot) func(channel string, user twitch.User, message twitch.Message) {
	return func(channel string, user twitch.User, m twitch.Message) {
		username := strings.ToLower(b.config.Username)
		if strings.ToLower(user.Username) == username {
			return
		}

		args := strings.Split(m.Text, " ")
		if len(args) < 1 {
			return
		}

		first := strings.ToLower(args[0])
		last := strings.ToLower(args[len(args)-1])

		if first == "!xd" {
			log.WithFields(log.Fields{"channel": channel}).Info("sending xD message")
			b.Say(message{channel, "xD"})
			return
		} else if first == "!bot" {
			log.WithFields(log.Fields{"channel": channel}).Info("sending bot message")
			b.Say(message{channel, "I'm roastedb's bot, written in Go pajaH"})
			return
		} else if first == "!php" {
			log.WithFields(log.Fields{"channel": channel}).Info("sending php message")
			b.Say(message{channel, "PHPDETECTED"})
			return
		}

		// No mentions, or the only thing in the message is a mention, don't process
		if (!isMention(first, b.config.Username) && !isMention(last, b.config.Username)) || len(args) < 2 {
			return
		}

		if isMention(first, b.config.Username) {
			args = args[1:]
		} else if isMention(last, b.config.Username) {
			args = args[:len(args)-1]
		}

		log.WithFields(log.Fields{"text": m.Text}).Info("handling message")
		for _, c := range b.Commands {
			if !c.Match(args) {
				continue
			}

			go c.Execute(b, args, channel, user, m)
			break
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
