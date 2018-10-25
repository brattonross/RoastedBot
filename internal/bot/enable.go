package bot

import (
	"fmt"
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

	invalidSyntax := "Invalid command syntax. Usage: enable|disable -m module_name [-c command_name]"

	if len(args) < 3 {
		b.Say(channel, invalidSyntax)
		return
	}
	if args[1] != "-m" {
		b.Say(channel, invalidSyntax)
		return
	}

	enable := false
	first := strings.ToLower(args[0])
	if first != "enable" && first != "disable" {
		return
	}
	if first == "enable" {
		enable = true
	}

	module := strings.ToLower(args[2])
	if _, ok := b.Modules[module]; !ok {
		b.Say(channel, fmt.Sprintf("Couldn't find module %s", module))
		return
	}

	// if less than 5 args, no command specified
	if len(args) < 5 {
		b.Modules[module].Enabled = enable

		if enable {
			b.Say(channel, fmt.Sprintf("%s enabled", module))
		} else {
			b.Say(channel, fmt.Sprintf("%s disabled", module))
		}
		return
	}

	if args[3] != "-c" {
		b.Say(channel, invalidSyntax)
		return
	}

	command := strings.ToLower(args[4])
	m := b.Modules[module]
	for _, c := range m.Commands {
		if strings.ToLower(c.Name()) == command {
			if enable {
				c.Enable()
				b.Say(channel, fmt.Sprintf("Enabled command %s in module %s", command, module))
			} else {
				c.Disable()
				b.Say(channel, fmt.Sprintf("Disabled command %s in module %s", command, module))
			}
			return
		}
	}

	b.Say(channel, fmt.Sprintf("Couldn't find command %s in module %s", command, module))
}
