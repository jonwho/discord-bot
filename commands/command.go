package commands

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/BryanSLam/discord-bot/util"
	"github.com/alpacahq/alpaca-trade-api-go/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/common"
	"github.com/go-redis/redis"
	iex "github.com/jonwho/go-iex/v4"
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
	alpacaPaperURL     string = "https://paper-api.alpaca.markets"
)

var (
	token, iexSecretToken string
	alpacaID, alpacaKey   string
	alpacaClient          *alpaca.Client
	redisClient           *redis.Client
	iexClient             *iex.Client
	cronner               *cron.Cron
	pst, _                = time.LoadLocation("America/Los_Angeles")
)

// TODO: don't use init
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

	alpacaID = os.Getenv("ALPACA_KEY_ID")
	if alpacaID == "" {
		log.Fatalln("Alpaca key id cannot be blank")
		return
	}
	alpacaKey = os.Getenv("ALPACA_SECRET_KEY")
	if alpacaKey == "" {
		log.Fatalln("Alpaca secret key cannot be blank")
		return
	}
	os.Setenv(common.EnvApiKeyID, alpacaID)
	os.Setenv(common.EnvApiSecretKey, alpacaKey)
	alpaca.SetBaseUrl(alpacaPaperURL)
	// N.B. not actually used
	// TODO: figure out how to use it or remove it
	alpacaClient = alpaca.NewClient(common.Credentials())

	redisClient = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	cronner = cron.NewWithLocation(pst)
	cronner.Start()
}
