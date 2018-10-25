package bot

import (
	"strings"

	twitch "github.com/gempir/go-twitch-irc"
)

// Enable allows modules and commands to be enabled and disabled by users.
// TODO: Some kind of authentication probably needs to be in place in the future.
type Enable struct {
	*command
}

// Execute the command.
func (e Enable) Execute(b *Bot, args []string, channel string, user twitch.User, message twitch.Message) {
	if strings.ToLower(user.Username) != "roastedb" {
		return
	}

	argLen := len(args)
	switch {
	case argLen < 2 || argLen == 4:
		b.Say(channel, "Invalid command syntax. Usage: enable|disable -m module_name -c command_name")
	case argLen == 2:
		// Entire module or single command without option tag
	}
}
