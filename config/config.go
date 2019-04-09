package config

import (
	"fmt"

	"github.com/tkanos/gonfig"
)

type BotConfig struct {
	CoinAPIURL            string
	WizdaddyURL           string
	InvalidCommandMessage string
	BotLogChannelID       string
	Maintainers           string
	EarningsWhisperURL    string
}

var (
	config BotConfig
)

func GetConfig() *BotConfig {
	// Initalize config if it doesn't already exist
	if config == (BotConfig{}) {
		// Use gonfig to fetch the config variables from config.json
		err := gonfig.GetConf("config.json", &config)
		if err != nil {
			fmt.Println("error fetching config values", err)
			return nil
		}
	}

	return &config
}
