package commands

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/BryanSLam/discord-bot/util"
	dg "github.com/bwmarrin/discordgo"
	iex "github.com/jonwho/go-iex"
)

func init() {
	// Run on 15 minute interval between hours 6-13 from Monday-Friday
	cronner.AddFunc("0 0/15 6-13 * * MON-FRI", watchlistCron)
}

func newWatchlistCommand() command {
	return command{
		match: func(s string) bool {
			return regexp.MustCompile(`(?i)^!watchlist [\w ]+$`).MatchString(s)
		},
		fn: watchlist,
	}
}

// Watchlist tickers to report on on an interval
func watchlist(s *dg.Session, m *dg.MessageCreate) {
	logger := util.Logger{Session: s, ChannelID: botLogChannelID}

	trimmed := strings.TrimSpace(m.Content)
	slice := strings.Split(trimmed, " ")
	tickers := slice[1:]
	iexClient, err := iex.NewClient()
	if err != nil {
		logger.Trace("IEX client initialization failed. Message: " + err.Error())
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	for _, ticker := range tickers {
		logger.Info("Fetching stock info for " + ticker)
		quote, err := iexClient.Quote(ticker, true)
		if err != nil {
			errStr := fmt.Sprintf("IEX request failed for ticker %s. Message: %s", ticker, err.Error())
			logger.Trace(errStr)
			s.ChannelMessageSend(m.ChannelID, errStr)
		} else if quote == nil {
			logger.Trace(fmt.Sprintf("nil quote from ticker: %s", ticker))
		} else {
			redisClient.SAdd(watchlistRedisKey, fmt.Sprintf("%s~*%s", m.ChannelID, ticker))
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Adding ticker %s to watchlist", ticker))
		}
	}
}

func watchlistCron() {
	dgSession, _ := dg.New("Bot " + token)
	dgSession.Open()
	defer dgSession.Close()

	tickers := redisClient.SMembers(watchlistRedisKey).Val()

	if len(tickers) > 0 {
		iexClient, _ := iex.NewClient()
		for _, ticker := range tickers {
			split := strings.Split(ticker, "~*")
			channel, symbol := split[0], split[1]
			quote, _ := iexClient.Quote(symbol, true)
			message := util.FormatQuote(quote)
			dgSession.ChannelMessageSend(channel, message)
		}
	}
}
