package grpc

import (
	"github.com/brattonross/roastedbot/twitch"
	context "golang.org/x/net/context"
)

// Service is the twitch service.
type Service struct {
	bot *twitch.Bot
}

// NewService creates a new service.
func NewService(b *twitch.Bot) *Service {
	return &Service{b}
}

// Channels should return a list of all channels.
func (s *Service) Channels(ctx context.Context, req *ChannelsRequest) (*ChannelsResponse, error) {
	chans := map[string]*Channel{}
	for _, c := range s.bot.Channels {
		chans[c.Name] = &Channel{
			Name:           c.Name,
			EnabledModules: c.EnabledModules,
			//Modules:        map[string]*Module{},
		}
	}
	return &ChannelsResponse{chans}, nil
}

// Modules should return a list of all modules.
func (s *Service) Modules(ctx context.Context, req *ModulesRequest) (*ModulesResponse, error) {
	return &ModulesResponse{}, nil
}
