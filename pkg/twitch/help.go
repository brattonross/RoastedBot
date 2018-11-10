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
func executeHelp(cl *Client, args []string, channel string, user twitch.User, message twitch.Message) {
	cl.Say(
		channel,
		fmt.Sprintf("%s, to use my commands, mention me at the start or end of your message.", user.DisplayName),
	)
}
