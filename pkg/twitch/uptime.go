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
		int(uptime.Hours()),
		int(uptime.Minutes())%60,
		int(uptime.Seconds())%60,
	)
	days := int(uptime.Hours() / 24)
	if days > 0 {
		resp = fmt.Sprintf(
			"%s has been running for %d days, %d hours, %d minutes, and %d seconds",
			b.config.Username,
			days,
			int(uptime.Hours())%24,
			int(uptime.Minutes())%60,
			int(uptime.Seconds())%60,
		)
	}
	b.Say(channel, resp)
}
