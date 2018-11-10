package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"

	"github.com/brattonross/roastedbot"
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
	if err = controller.StartBot(); err != nil {
		log.Fatalf("fatal error occurred while bot was running: %v", err)
	}
}
