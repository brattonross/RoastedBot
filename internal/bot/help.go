package bot

import (
	"fmt"

	"github.com/gempir/go-twitch-irc"
)

// Help command.
type Help struct {
	*command
}

// Execute the command.
func (h Help) Execute(b *Bot, args []string, channel string, user twitch.User, message twitch.Message) {
	b.Say(
		channel,
		fmt.Sprintf("%s, to use my commands, mention me at the start or end of your message.", user.DisplayName),
	)
}
