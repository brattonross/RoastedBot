package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"net"
	
	"github.com/brattonross/roastedbot/twitch"
	tgrpc "github.com/brattonross/roastedbot/twitch/service/grpc"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	configPath := flag.String("config", "config.json", "Location of the configuration json file")

	flag.Parse()

	log.Infof("using config path %s", *configPath)

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
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("failed to unmarshal configuration file")
	}
	log.Info("successfully read config")

	// Initialise bot
	bot := twitch.NewBot(config)
	bot.Init()

	// gRPC server setup
	server := grpc.NewServer()
	service := tgrpc.NewService(bot)
	tgrpc.RegisterBotServiceServer(server, service)
	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatalf("failed to listen on port 1234: %v", err)
	}
	go func() {
		log.Info("roastedbot gRPC server is listening on port 1234")
		log.Error(server.Serve(l))
	}()

	err = bot.Connect()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("connection to twitch failed")
	}
}
