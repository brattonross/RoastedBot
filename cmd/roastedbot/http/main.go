package main

import (
	"context"
	"flag"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	
	"github.com/gorilla/mux"
	pb "github.com/brattonross/roastedbot/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type config struct {
	GRPC struct {
		Port int `json:"port"`
	} `json:"grpc"`

	HTTP struct {
		Port int `json:"port"`
	} `json:"http"`
}

func main() {
	configPath := flag.String("config", "api.config.json", "Path of the API configuration")

	flag.Parse()

	// Read config
	b, err := ioutil.ReadFile(*configPath)
	if err != nil {
		log.WithFields(log.Fields{
			"configPath": *configPath,
			"error":      err,
		}).Fatal("failed to read configuration file")
	}
	config := &config{}
	err = json.Unmarshal(b, config)
	if err != nil {
		log.WithField("error", err).Fatal("failed to unmarshal configuration file")
	}

	port := config.GRPC.Port
	target := fmt.Sprintf(":%d", port)
	conn, err := grpc.Dial(target, grpc.WithInsecure())
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"target": target,
		}).Fatal("failed to dial target")
	}
	defer conn.Close()

	client := pb.NewBotServiceClient(conn)

	server := &Server{client}
	addr := fmt.Sprintf(":%d", config.HTTP.Port)
	log.WithField("address", addr).Info("http server listening")
	log.Fatal(server.ListenAndServe(addr))
}

// Server is an http server for the bot api.
type Server struct{
	client pb.BotServiceClient
}

// ListenAndServe begins listening and serving at the given address.
func (s *Server) ListenAndServe(addr string) error {
	r := mux.NewRouter()
	r.HandleFunc("/channels", func(w http.ResponseWriter, req *http.Request) {
		resp, err := s.client.Channels(context.Background(), &pb.ChannelsRequest{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(resp.Channels)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	return http.ListenAndServe(addr, r)
}
