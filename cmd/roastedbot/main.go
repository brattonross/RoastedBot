package main

import (
	"context"
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/brattonross/roastedbot"
	service "github.com/brattonross/roastedbot/pkg/twitch/service/http"
	"github.com/sirupsen/logrus"
)

func main() {
	configPath := flag.String("config", "bot.config.json", "Path of the bot configuration")

	flag.Parse()

	log := logrus.New()

	// Read config
	b, err := ioutil.ReadFile(*configPath)
	if err != nil {
		log.WithFields(logrus.Fields{
			"configPath": *configPath,
			"error":      err,
		}).Fatal("failed to read configuration file")
	}
	config := &roastedbot.Config{}
	err = json.Unmarshal(b, config)
	if err != nil {
		log.WithField("error", err).Fatal("failed to unmarshal configuration file")
	}
	
	controller := roastedbot.NewController(config, log)
	go func() {
		if err = controller.Connect(); err != nil {
			log.Fatalf("fatal error occurred while bot was running: %v", err)
		}
	}()

	handler := service.NewHandler(controller.Client)
	server := &http.Server{
		Addr: ":9001",
		Handler: handler,
	}
	go func() {
		log.Infof("starting http server on port 9001")
		if err := server.ListenAndServe(); err != nil {
			log.Errorf("http server encountered an error: %v", err)
		}
	}()

	// Graceful shutdown
	sigquit := make(chan os.Signal, 1)
	signal.Notify(sigquit, os.Interrupt, syscall.SIGTERM)

	sig := <-sigquit
	log.Infof("caught sig: %+v", sig)
	
	log.Infof("gracefully shutting down controller...")
	if err := server.Shutdown(context.Background()); err != nil {
		log.Errorf("unable to shut down server: %v", err)
	} else {
		log.Info("server stopped")
	}

	log.Infof("gracefully shutting down client...")
	if err := controller.Disconnect(); err != nil {
		log.Errorf("unable to shut down client: %v", err)
	} else {
		log.Infof("client stopped")
	}
}
