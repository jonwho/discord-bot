package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	dbot "github.com/BryanSLam/discord-bot"
)

func main() {
	var botToken string
	// If no value was provided from flag look for env var BOT_TOKEN
	botToken = os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatalln("Bot Token must be set")
	}

	maintainers := strings.Split(os.Getenv("MAINTAINERS"), ",")
	botLogChannelID := os.Getenv("BOT_LOG_CHANNEL_ID")
	bot, err := dbot.NewMusicBot(
		botToken,
		dbot.WithMaintainers(maintainers),
		dbot.WithBotLogChannelID(botLogChannelID),
	)
	if err != nil {
		log.Fatalln("error creating Discord session,", err)
	}

	// Open a websocket connection to Discord and begin listening.
	err = bot.Open()
	if err != nil {
		log.Fatalln("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	log.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	bot.Close()
}
