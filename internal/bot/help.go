package bot

import (
	"fmt"
	"strings"

	"github.com/gempir/go-twitch-irc"
)

// Help command.
type Help struct{}

// Execute the command.
func (h Help) Execute(b *Bot, args []string, channel string, user twitch.User, m twitch.Message) {
	names := []string{}
	for _, c := range b.Commands {
		if c.Name() == "help" {
			continue
		}
		names = append(names, c.Name())
	}
	b.Say(message{
		channel,
		fmt.Sprintf("%s, to use my commands, mention me at the start or end of your message. Available commands are: %s", user.DisplayName, strings.Join(names, ", ")),
	})
}

// Match determines if the message should trigger the command.
func (h Help) Match(s []string) bool {
	if len(s) < 1 {
		return false
	}
	return strings.ToLower(s[0]) == "help"
}

// Name of the command.
func (h Help) Name() string {
	return "help"
}
