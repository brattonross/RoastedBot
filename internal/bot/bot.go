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
	Execute(b *Bot, args []string, channel string, user twitch.User, message twitch.Message) (string, error)
	Match(s []string) bool
	Name() string
}

// Bot is the bot xD
type Bot struct {
	config Config
	start  time.Time
	client *twitch.Client

	Commands []Command
}

// New creates a new Bot using the given Config.
func New(config Config) *Bot {
	return &Bot{
		config: config,
		client: twitch.NewClient(config.Username, config.OAuth),
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

	for _, channel := range b.config.Channels {
		b.JoinChannel(channel)
		log.WithFields(log.Fields{"channel": channel}).Info("joined channel")
	}
}

// JoinChannel joins the given channel.
func (b *Bot) JoinChannel(channel string) {
	b.client.Join(channel)
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
			log.WithFields(log.Fields{
				"channel": channel,
			}).Info("sending xD message")
			b.client.Say(channel, "xD")
			return
		} else if first == "!bot" {
			log.WithFields(log.Fields{
				"channel": channel,
			}).Info("sending bot message")
			b.client.Say(channel, "I'm roastedb's bot, written in Go pajaH")
			return
		} else if first == "!php" {
			log.WithFields(log.Fields{
				"channel": channel,
			}).Info("sending php message")
			b.client.Say(channel, "PHPDETECTED")
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

		log.WithFields(log.Fields{
			"text": message.Text,
		}).Info("handling message")
		for _, c := range b.Commands {
			if !c.Match(args) {
				continue
			}

			resp, err := c.Execute(b, args, channel, user, message)
			if err != nil {
				log.WithFields(log.Fields{
					"command": c.Name(),
					"error":   err,
				}).Error("error occurred while executing command")
				break
			}

			b.client.Say(channel, resp)
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
