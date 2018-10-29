package twitch

import (
	"fmt"
	"strings"

	twitch "github.com/gempir/go-twitch-irc"
	log "github.com/sirupsen/logrus"
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

	enable := false
	first := strings.ToLower(args[0])
	if first != "enable" && first != "disable" {
		b.Say(channel, invalidSyntax)
		return
	}
	if first == "enable" {
		enable = true
	}
	args = args[1:]

	// Iterate over args index, locate m and c
	var module string
	var command string
	for i := 0; i < len(args); i++ {
		cur := strings.ToLower(args[i])
		if (cur == "-m" || cur == "--module") && i+1 < len(args) {
			module = args[i+1]
		} else if (cur == "-c" || cur == "--command") && i+1 < len(args) {
			command = args[i+1]
		}
	}

	if module == "" {
		b.Say(channel, invalidSyntax)
		return
	}

	// No command specified - enable/disable module.
	if command == "" {
		if enable {
			if err := b.EnableModule(channel, module); err != nil {
				log.WithField("module", module).Error(err)
				b.Say(channel, fmt.Sprintf("Module '%s' does not exist", module))
				return
			}
			b.Say(channel, fmt.Sprintf("Enabled module '%s'", module))
			return
		}
		if err := b.DisableModule(channel, module); err != nil {
			log.WithField("module", module).Error(err)
			b.Say(channel, fmt.Sprintf("Module '%s' does not exist", module))
			return
		}
		b.Say(channel, fmt.Sprintf("Disabled module '%s'", module))
		return
	}

	// Command specified - enable/disable command.
	if enable {
		if err := b.EnableCommand(channel, module, command); err != nil {
			log.WithFields(log.Fields{
				module:  module,
				command: command,
			}).Error(err)
			b.Say(channel, fmt.Sprintf("Command '%s' does not exist in module '%s'", command, module))
			return
		}
		b.Say(channel, fmt.Sprintf("Enabled command '%s' in module '%s'", command, module))
		return
	}
	if err := b.DisableCommand(channel, module, command); err != nil {
		log.WithFields(log.Fields{
			module:  module,
			command: command,
		}).Error(err)
		b.Say(channel, fmt.Sprintf("Command '%s' does not exist in module '%s'", command, module))
		return
	}
	b.Say(channel, fmt.Sprintf("Disabled command '%s' in module '%s'", command, module))
}

// Match checks if the command should execute.
func (e Enable) Match(s string) bool {
	return strings.ToLower(s) == "enable" || strings.ToLower(s) == "disable"
}
