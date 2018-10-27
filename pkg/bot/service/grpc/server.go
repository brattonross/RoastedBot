package grpc

import (
	"github.com/brattonross/roastedbot/pkg/bot"
	context "golang.org/x/net/context"
)

// Service is the bot service.
type Service struct {
	modules map[string]*bot.Module
}

// New creates a new service.
func New(modules map[string]*bot.Module) *Service {
	return &Service{modules}
}

// Modules should return a list of all modules.
func (s *Service) Modules(ctx context.Context, req *ModulesRequest) (*ModulesResponse, error) {
	mods := map[string]*Module{}
	for _, m := range s.modules {
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
