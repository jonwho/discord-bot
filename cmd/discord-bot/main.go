package main

import (
	"bytes"
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	dbot "github.com/BryanSLam/discord-bot"
	"github.com/BryanSLam/discord-bot/botcommands/earnings_reporter"
	dg "github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"
)

func init() {
	pst, _ := time.LoadLocation("America/Los_Angeles")
	cronner := cron.New(cron.WithLocation(pst))
	cronner.Start()
	// Run once at 6:00 AM from Monday
	cronner.AddFunc("00 06 * * MON", pinEarningsReports)
}

func main() {
	maintainers := strings.Split(os.Getenv("MAINTAINERS"), ",")
	botLogChannelID := os.Getenv("BOT_LOG_CHANNEL_ID")
	botStockChannelID := os.Getenv("BOT_STOCK_CHANNEL_ID")

	// Run the program with `go run main.go -t <token>`
	// flag.Parse() will assign to token var
	var botToken string
	flag.StringVar(&botToken, "t", "", "Bot Token")
	flag.Parse()

	// If no value was provided from flag look for env var BOT_TOKEN
	if botToken == "" {
		botToken = os.Getenv("BOT_TOKEN")
	}

	// If still empty then that's no bueno
	if botToken == "" {
		log.Fatalln("Bot Token must be set")
	}

	iexToken := os.Getenv("IEX_SECRET_TOKEN")
	if iexToken == "" {
		log.Fatalln("IEX Token cannot be blank")
		return
	}

	alpacaID := os.Getenv("ALPACA_KEY_ID")
	if alpacaID == "" {
		log.Fatalln("Alpaca key id cannot be blank")
		return
	}
	alpacaKey := os.Getenv("ALPACA_SECRET_KEY")
	if alpacaKey == "" {
		log.Fatalln("Alpaca secret key cannot be blank")
		return
	}

	bot, err := dbot.New(
		botToken,
		iexToken,
		alpacaID,
		alpacaKey,
		dbot.WithMaintainers(maintainers),
		dbot.WithBotLogChannelID(botLogChannelID),
		dbot.WithBotStockChannelID(botStockChannelID),
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

// Pins the upcoming earning reports for the week.
func pinEarningsReports() {
	botStockChannelID := os.Getenv("BOT_STOCK_CHANNEL_ID")
	botToken := os.Getenv("BOT_TOKEN")
	dgSession, _ := dg.New("Bot " + botToken)
	dgSession.Open()
	defer dgSession.Close()

	ctx := context.Background()
	dw := dbot.NewDiscordWriter(dgSession, nil, botStockChannelID)
	buf := bytes.NewBuffer([]byte(""))

	reporter, _ := earningsreporter.New()
	err := reporter.Execute(ctx, buf)
	if err != nil {
		log.Println(err)
		dw.Write([]byte("Failed to fetch weekly earnings reports"))
		return
	}

	msg, err := dw.ChannelMessageSend(buf.String())
	if err != nil {
		log.Println(err)
		dw.Write([]byte("Failed to send weekly earnings reports to channel"))
		return
	}

	err = dw.Pin(msg.ID)
	if err != nil {
		log.Println(err)
		dw.Write([]byte("Failed to pin weekly earnings reports to channel"))
		return
	}
}
