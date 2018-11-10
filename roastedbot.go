package roastedbot

import (
	"fmt"
	"strings"
	"time"

	tirc "github.com/gempir/go-twitch-irc"
	log "github.com/sirupsen/logrus"

	"github.com/brattonross/roastedbot/pkg/admin"
	"github.com/brattonross/roastedbot/pkg/twitch"
)

// Config for the bot.
type Config struct {
	Username string   `json:"username"`
	OAuth    string   `json:"oauth"`
	Channels []string `json:"channels"`
}

// Controller is the application controller.
type Controller struct {
	Client *twitch.Client
	Config *Config

	log    *log.Logger
}

// NewController creates a new bot controller.
func NewController(config *Config, log *log.Logger) *Controller {
	irc := tirc.NewClient(config.Username, config.OAuth)
	// TODO: DB driver
	client := twitch.NewClient(config.Username, irc)
	client.OnConnect(func() {
		log.Info("connected to twitch")
	})

	return &Controller{
		client,
		config,
		log,
	}
}

// Connect joins any channels that the bot is configured
// to join, and then connects to twitch.
func (c *Controller) Connect() error {
	defer func() {
		c.log.Debug("disconnecting from twitch...")
		if err := c.Client.Disconnect(); err != nil {
			c.log.Errorf("failed to disconnect elegantly from twitch: %v", err)
		}
	}()

	c.loadChannels()
	for _, channel := range c.Client.Channels() {
		c.addRequiredModules(channel.Name)
	}
	c.Client.JoinChannels()

	c.Client.OnNewMessage(c.onNewMessage)

	return c.Client.Connect()
}

// Disconnect will disconnect the client from twitch.
func (c *Controller) Disconnect() error {
	return c.Client.Disconnect()
}

func (c *Controller) addRequiredModules(channel string) {
	adminModule, err := c.Client.AddModule(channel, "admin")
	if err != nil {
		c.log.Errorf("failed to add admin module: %v", err)
	} else {
		adminModule.AddCommand(admin.EnableCommand)
		adminModule.EnableCommand(admin.EnableCommand.Name)
		c.Client.EnableModule(channel, "admin")
	}

	general, err := c.Client.AddModule(channel, "general")
	if err != nil {
		c.log.Errorf("failed to add general module: %v", err)
	} else {
		general.AddCommand(twitch.HelpCommand)
		general.EnableCommand(twitch.HelpCommand.Name)

		general.AddCommand(twitch.UptimeCommand)
		general.EnableCommand(twitch.UptimeCommand.Name)

		c.Client.EnableModule(channel, "general")
	}
}

// LoadChannels loads the channels that the bot should join on start.
func (c *Controller) loadChannels() {
	for _, ch := range c.Config.Channels {
		if err := c.Client.AddChannel(ch); err != nil {
			c.log.Errorf("failed to load channel: %v", err)
		}
	}
}

func (c *Controller) onNewMessage(channel string, user tirc.User, message tirc.Message) {
	username := strings.ToLower(c.Client.Username)
	if strings.ToLower(user.Username) == username {
		return
	}

	ch, err := c.Client.Channel(channel)
	if err != nil {
		log.WithField("channel", channel).Error("cannot handle message: channel is not configured")
		return
	}

	args := strings.Split(message.Text, " ")
	if len(args) < 1 {
		return
	}

	first := strings.ToLower(args[0])
	last := strings.ToLower(args[len(args)-1])

	if first == "!xd" {
		c.Client.Say(channel, "xD")
		return
	} else if first == "!php" {
		c.Client.Say(channel, "PHPDETECTED")
		return
	}

	// No mentions, don't process
	if !isMention(first, c.Client.Username) && !isMention(last, c.Client.Username) {
		return
	}

	// Only message is a mention of the bot, say hi
	if len(args) == 1 {
		c.Client.Say(channel, fmt.Sprintf("hi %s :)", user.DisplayName))
		return
	}

	if isMention(first, c.Client.Username) {
		args = args[1:]
	} else if isMention(last, c.Client.Username) {
		args = args[:len(args)-1]
	}

	log.WithFields(log.Fields{
		"channel": channel,
		"text":    message.Text,
		"user":    user.DisplayName,
	}).Info("handling message")

	command, module := ch.MatchCommand(args)
	if command == nil {
		return
	}

	if !module.IsCommandEnabled(command.Name) {
		log.WithFields(log.Fields{
			"channel": channel,
			"command": command.Name,
			"module":  module.Name,
			"user":    user.DisplayName,
		}).Info("command is not enabled")
		return
	}
	if command.IsOnCooldown() {
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

		command.Execute(c.Client, args, channel, user, message)
		command.LastUsed = time.Now()
	}()
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
