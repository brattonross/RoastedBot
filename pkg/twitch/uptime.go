package twitch

import (
	"fmt"
	"time"

	twitch "github.com/gempir/go-twitch-irc"
)

// UptimeCommand prints the bot's uptime.
var UptimeCommand = &Command{
	Cooldown: time.Second * 2,
	Name:     "uptime",
	Run:      executeUptime,
	Use:      "uptime",
}

// Execute the command.
func executeUptime(cl *Client, args []string, channel string, user twitch.User, message twitch.Message) {
	uptime := time.Since(cl.start)
	resp := fmt.Sprintf(
		"%s has been running for %d hours, %d minutes, and %d seconds",
		cl.Username,
		uptime/time.Hour,
		uptime/time.Minute%60,
		uptime/time.Second%60,
	)
	days := int(uptime.Hours() / 24)
	if days > 0 {
		resp = fmt.Sprintf(
			"%s has been running for %d days, %d hours, %d minutes, and %d seconds",
			cl.Username,
			days,
			uptime/time.Hour%24,
			uptime/time.Minute%60,
			uptime/time.Second%60,
		)
	}
	cl.Say(channel, resp)
}
