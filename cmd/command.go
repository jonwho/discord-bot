package cmd

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/BryanSLam/discord-bot/util"
	"github.com/go-redis/redis"
	iex "github.com/jonwho/go-iex"
	"github.com/robfig/cron"
)

// Command TODO: @doc
type Command struct {
	Match func(s string) bool
	Fn    func(rw io.ReadWriter, l *util.Logger, m map[string]interface{})
}

const (
	dateFormat         string = "1/_2/06"
	redisDateFormat    string = "01/02/06"
	watchlistRedisKey  string = "watchlist"
	coinAPIURL         string = "https://min-api.cryptocompare.com/data/pricemultifull?fsyms="
	wizdaddyURL        string = "http://dev.wizdaddy.io/api/giveItToMeDaddy"
	earningsWhisperURL string = "https://www.earningswhispers.com/calendar?sb=p&t=all"
)

var (
	token, iexSecretToken string
	redisClient           *redis.Client
	iexClient             *iex.Client
	cronner               *cron.Cron
	pst, _                = time.LoadLocation("America/Los_Angeles")
)

func init() {
	var err error
	token = os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatalln("Bot token cannot be blank")
		return
	}

	iexSecretToken = os.Getenv("IEX_SECRET_TOKEN")
	if iexSecretToken == "" {
		log.Fatalln("IEX secret token cannot be blank")
		return
	}
	iexClient, err = iex.NewClient(iexSecretToken)
	if err != nil {
		log.Fatalln("IEX client initialization failed. Message: " + err.Error())
		return
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	cronner = cron.NewWithLocation(pst)
	cronner.Start()
}
