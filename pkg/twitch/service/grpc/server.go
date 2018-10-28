package grpc

import (
	"github.com/brattonross/roastedbot/pkg/twitch"
	context "golang.org/x/net/context"
)

// Service is the twitch service.
type Service struct {
	bot *twitch.Bot
}

// New creates a new service.
func newGrpcService(b *twitch.Bot) *Service {
	return &Service{b}
}

// Modules should return a list of all modules.
func (s *Service) Modules(ctx context.Context, req *ModulesRequest) (*ModulesResponse, error) {
	mods := map[string]*Module{}
	for _, m := range s.bot.Modules {
		mods[m.Name] = &Module{
			Name:     m.Name,
			Enabled:  m.Enabled,
			Commands: []*Command{},
		}
		for _, c := range m.Commands {
			mods[m.Name].Commands = append(mods[m.Name].Commands, &Command{
				Cooldown: c.Cooldown().String(),
				Enabled:  c.Enabled(),
				Name:     c.Name(),
			})
		}
	}

	return &ModulesResponse{
		Modules: mods,
	}, nil
}
