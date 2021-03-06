package admin

import (
	"fmt"
	"strings"
	"time"

	"github.com/brattonross/roastedbot/pkg/twitch"
	tirc "github.com/gempir/go-twitch-irc"
	log "github.com/sirupsen/logrus"
)

// EnableCommand allows modules and commands to be enabled.
var EnableCommand = &twitch.Command{
	Cooldown: time.Second * 1,
	Name:     "enable",
	Run:      executeEnable,
	Use:      "enable",
}

func executeEnable(cl *twitch.Client, args []string, channel string, user tirc.User, message tirc.Message) {
	if strings.ToLower(user.Username) != "roastedb" {
		return
	}

	invalidSyntax := "Invalid command syntax. Usage: enable|disable -m module_name [-c command_name]"

	if len(args) < 3 {
		cl.Say(channel, invalidSyntax)
		return
	}

	enable := false
	first := strings.ToLower(args[0])
	if first != "enable" && first != "disable" {
		cl.Say(channel, invalidSyntax)
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
		cl.Say(channel, invalidSyntax)
		return
	}
	if strings.ToLower(module) == "default" {
		return
	}

	// No command specified - enable/disable module.
	if command == "" {
		if enable {
			if err := cl.EnableModule(channel, module); err != nil {
				log.WithField("module", module).Error(err)
				cl.Say(channel, fmt.Sprintf("Module '%s' does not exist", module))
				return
			}
			cl.Say(channel, fmt.Sprintf("Enabled module '%s'", module))
			return
		}
		if err := cl.DisableModule(channel, module); err != nil {
			log.WithField("module", module).Error(err)
			cl.Say(channel, fmt.Sprintf("Module '%s' does not exist", module))
			return
		}
		cl.Say(channel, fmt.Sprintf("Disabled module '%s'", module))
		return
	}

	// Command specified - enable/disable command.
	if enable {
		if err := cl.EnableCommand(channel, module, command); err != nil {
			log.WithFields(log.Fields{
				module:  module,
				command: command,
			}).Error(err)
			cl.Say(channel, fmt.Sprintf("Command '%s' does not exist in module '%s'", command, module))
			return
		}
		cl.Say(channel, fmt.Sprintf("Enabled command '%s' in module '%s'", command, module))
		return
	}
	if err := cl.DisableCommand(channel, module, command); err != nil {
		log.WithFields(log.Fields{
			module:  module,
			command: command,
		}).Error(err)
		cl.Say(channel, fmt.Sprintf("Command '%s' does not exist in module '%s'", command, module))
		return
	}
	cl.Say(channel, fmt.Sprintf("Disabled command '%s' in module '%s'", command, module))
}
