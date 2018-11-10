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
// The bot and it's API are managed within the controller.
type Controller struct {
	bot    *twitch.Bot
	config *Config
	log    *log.Logger
}

// NewController creates a new bot controller.
func NewController(config *Config, log *log.Logger) *Controller {
	client := tirc.NewClient(config.Username, config.OAuth)
	// Initialise bot
	// TODO: DB driver
	bot := twitch.NewBot(config.Username, client)
	bot.OnConnect(func() {
		log.Info("connected to twitch")
	})
	return &Controller{
		bot:    bot,
		config: config,
		log:    log,
	}
}

// StartBot joins any channels that the bot is configured
// to connect to, and then connects to twitch.
func (c *Controller) StartBot() error {
	defer func() {
		c.log.Debug("disconnecting from twitch...")
		if err := c.bot.Disconnect(); err != nil {
			c.log.Errorf("failed to disconnect elegantly from twitch: %v", err)
		}
	}()

	c.loadChannels()
	c.bot.JoinChannels()

	for _, channel := range c.bot.Channels() {
		c.addRequiredModules(channel.Name)
	}

	c.bot.OnNewMessage(c.onNewMessage)

	return c.bot.Connect()
}

func (c *Controller) addRequiredModules(channel string) {
	adminModule, err := c.bot.AddModule(channel, "admin")
	if err != nil {
		c.log.Errorf("failed to add admin module: %v", err)
	} else {
		adminModule.AddCommand(admin.EnableCommand)
		adminModule.EnableCommand(admin.EnableCommand.Name)
		c.bot.EnableModule(channel, "admin")
	}

	general, err := c.bot.AddModule(channel, "general")
	if err != nil {
		c.log.Errorf("failed to add general module: %v", err)
	} else {
		general.AddCommand(twitch.HelpCommand)
		general.EnableCommand(twitch.HelpCommand.Name)

		general.AddCommand(twitch.UptimeCommand)
		general.EnableCommand(twitch.UptimeCommand.Name)

		c.bot.EnableModule(channel, "general")
	}
}

// LoadChannels loads the channels that the bot should join on start.
func (c *Controller) loadChannels() {
	for _, ch := range c.config.Channels {
		if err := c.bot.AddChannel(ch); err != nil {
			c.log.Errorf("failed to load channel: %v", err)
		}
	}
}

func (c *Controller) onNewMessage(channel string, user tirc.User, message tirc.Message) {
	username := strings.ToLower(c.bot.Username)
	if strings.ToLower(user.Username) == username {
		return
	}

	ch, err := c.bot.Channel(channel)
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
		c.bot.Say(channel, "xD")
		return
	} else if first == "!php" {
		c.bot.Say(channel, "PHPDETECTED")
		return
	}

	// No mentions, don't process
	if !isMention(first, c.bot.Username) && !isMention(last, c.bot.Username) {
		return
	}

	// Only message is a mention of the bot, say hi
	if len(args) == 1 {
		c.bot.Say(channel, fmt.Sprintf("hi %s :)", user.DisplayName))
		return
	}

	if isMention(first, c.bot.Username) {
		args = args[1:]
	} else if isMention(last, c.bot.Username) {
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

		command.Execute(c.bot, args, channel, user, message)
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
