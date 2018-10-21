package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"

	bi "github.com/brattonross/roastedbot/internal/bot"
	log "github.com/sirupsen/logrus"
)

func main() {
	configPath := flag.String("config", "config.json", "Location of the configuration json file")

	flag.Parse()

	log.Infof("using config path %s", *configPath)

	b, err := ioutil.ReadFile(*configPath)
	if err != nil {
		log.WithFields(log.Fields{
			"configPath": *configPath,
			"error":      err,
		}).Fatal("failed to read configuration file")
	}

	config := bi.Config{}
	err = json.Unmarshal(b, &config)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("failed to unmarshal configuration file")
	}
	log.Info("successfully read config")

	bot := bi.New(config)
	bot.Init()

	err = bot.Connect()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("connection to twitch failed")
	}
}
