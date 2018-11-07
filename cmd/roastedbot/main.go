package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	
	"github.com/brattonross/roastedbot/twitch"
	pb "github.com/brattonross/roastedbot/proto"
	tgrpc "github.com/brattonross/roastedbot/twitch/service/grpc"
	tirc "github.com/gempir/go-twitch-irc"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	configPath := flag.String("config", "bot.config.json", "Path of the bot configuration")

	flag.Parse()

	// Read config
	b, err := ioutil.ReadFile(*configPath)
	if err != nil {
		log.WithFields(log.Fields{
			"configPath": *configPath,
			"error":      err,
		}).Fatal("failed to read configuration file")
	}
	config := twitch.Config{}
	err = json.Unmarshal(b, &config)
	if err != nil {
		log.WithField("error", err).Fatal("failed to unmarshal configuration file")
	}

	client := tirc.NewClient(config.Username, config.OAuth)

	// Initialise bot
	// TODO: DB driver
	bot := twitch.NewBot(config, client, nil)

	bot.OnConnect(func() {
		log.Info("connected to twitch")
	})

	bot.LoadChannels()
	bot.JoinChannels()

	// gRPC server setup
	port := config.GRPC.Port
	server := grpc.NewServer()
	service := tgrpc.NewService(bot)
	pb.RegisterBotServiceServer(server, service)
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen on port %d: %v", port, err)
	}
	go func() {
		log.Infof("roastedbot gRPC server is listening on port %d", port)
		log.Error(server.Serve(l))
	}()

	defer bot.Disconnect()
	err = bot.Connect()
	if err != nil {
		log.WithField("error", err).Fatal("connection to twitch failed")
	}
}
