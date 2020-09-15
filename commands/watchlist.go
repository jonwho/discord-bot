package commands

import (
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/BryanSLam/discord-bot/util"
	dg "github.com/bwmarrin/discordgo"
)

func init() {
	// Run on 15 minute interval between hours 6-13 from Monday-Friday
	cronner.AddFunc("*/15 6-13 * * MON-FRI", watchlistCron)
}

// NewWatchlistCommand TODO: @doc
func NewWatchlistCommand() *Command {
	return &Command{
		Match: func(s string) bool {
			return regexp.MustCompile(`(?i)^!watchlist [\w ]+$`).MatchString(s)
		},
		Fn: Watchlist,
	}
}

// Watchlist TODO: @doc
func Watchlist(rw io.ReadWriter, logger *util.Logger, m map[string]interface{}) {
	buf, err := ioutil.ReadAll(rw)
	if err != nil {
		rw.Write([]byte(err.Error()))
		return
	}

	mc := m["messageCreate"].(*dg.MessageCreate)
	channelID := mc.ChannelID

	trimmed := strings.TrimSpace(string(buf))
	slice := strings.Split(trimmed, " ")
	tickers := slice[1:]

	for _, ticker := range tickers {
		logger.Info("Fetching stock info for " + ticker)
		quote, err := iexClient.Quote(ticker, nil)
		if err != nil {
			errStr := fmt.Sprintf("IEX request failed for ticker %s. Message: %s", ticker, err.Error())
			logger.Trace(errStr)
			rw.Write([]byte(errStr))
		} else if quote == nil {
			logger.Trace(fmt.Sprintf("nil quote from ticker: %s", ticker))
		} else {
			redisClient.SAdd(watchlistRedisKey, fmt.Sprintf("%s~*%s", channelID, ticker))
			rw.Write([]byte(fmt.Sprintf("Adding ticker %s to watchlist", ticker)))
		}
	}
}

func watchlistCron() {
	dgSession, _ := dg.New("Bot " + token)
	dgSession.Open()
	defer dgSession.Close()

	tickers := redisClient.SMembers(watchlistRedisKey).Val()

	if len(tickers) > 0 {
		for _, ticker := range tickers {
			split := strings.Split(ticker, "~*")
			channel, symbol := split[0], split[1]
			quote, _ := iexClient.Quote(symbol, nil)
			message := util.FormatQuote(quote)
			dgSession.ChannelMessageSend(channel, message)
		}
	}
}
