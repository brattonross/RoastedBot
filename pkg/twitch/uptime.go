package twitch

import (
	"fmt"
	"time"

	twitch "github.com/gempir/go-twitch-irc"
)

// Uptime command.
type Uptime struct {
	*command
}

// Execute the command.
func (u Uptime) Execute(b *Bot, args []string, channel string, user twitch.User, message twitch.Message) {
	uptime := time.Since(b.start)
	resp := fmt.Sprintf(
		"%s has been running for %d hours, %d minutes, and %d seconds",
		b.config.Username,
		uptime/time.Hour,
		uptime/time.Minute%60,
		uptime/time.Second%60,
	)
	days := int(uptime.Hours() / 24)
	if days > 0 {
		resp = fmt.Sprintf(
			"%s has been running for %d days, %d hours, %d minutes, and %d seconds",
			b.config.Username,
			days,
			uptime/time.Hour%24,
			uptime/time.Minute%60,
			uptime/time.Second%60,
		)
	}
	b.Say(channel, resp)
}
