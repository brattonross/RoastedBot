package twitch

import (
	"fmt"
	"time"

	twitch "github.com/gempir/go-twitch-irc"
)

// HelpCommand prints help for the bot.
var HelpCommand = &Command{
	Cooldown: time.Second * 5,
	Name:     "help",
	Run:      executeHelp,
	Use:      "help",
}

// Execute the command.
func executeHelp(b *Bot, args []string, channel string, user twitch.User, message twitch.Message) {
	b.Say(
		channel,
		fmt.Sprintf("%s, to use my commands, mention me at the start or end of your message.", user.DisplayName),
	)
}
