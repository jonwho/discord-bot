package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	dbot "github.com/BryanSLam/discord-bot"
)

func main() {
	maintainers := strings.Split(os.Getenv("MAINTAINERS"), ",")
	botLogChannelID := os.Getenv("BOT_LOG_CHANNEL_ID")

	// Run the program with `go run main.go -t <token>`
	// flag.Parse() will assign to token var
	var token string
	flag.StringVar(&token, "t", "", "Bot Token")
	flag.Parse()

	// If no value was provided from flag look for env var BOT_TOKEN
	if token == "" {
		token = os.Getenv("BOT_TOKEN")
	}

	// If still empty then that's no bueno
	if token == "" {
		log.Fatalln("Bot Token must be set")
	}

	bot, err := dbot.New(
		token,
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
