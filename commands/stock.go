package commands

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/BryanSLam/discord-bot/util"
	dg "github.com/bwmarrin/discordgo"
	iex "github.com/jonwho/go-iex"
)

type stockCommand struct {
	regex *regexp.Regexp
}

func newStockCommand() stockCommand {
	return stockCommand{regexp.MustCompile(`(?i)^!stock [\w.]+$`)}
}

func (cmd stockCommand) match(s string) bool {
	return cmd.regex.MatchString(s)
}

func (cmd stockCommand) fn(s *dg.Session, m *dg.MessageCreate) {
	logger := util.Logger{Session: s, ChannelID: botLogChannelID}
	slice := strings.Split(m.Content, " ")
	ticker := slice[1]
	iexClient, err := iex.NewClient()
	if err != nil {
		logger.Trace("IEX client initialization failed. Message: " + err.Error())
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	logger.Info("Fetching stock info for " + ticker)
	quote, err := iexClient.Quote(ticker, true)
	if err != nil {
		rds, iexErr := iexClient.RefDataSymbols()
		if iexErr != nil {
			logger.Trace("IEX request failed. Message: " + iexErr.Error())
			s.ChannelMessageSend(m.ChannelID, iexErr.Error())
			return
		}

		fuzzySymbols := util.FuzzySearch(ticker, rds.Symbols)

		if len(fuzzySymbols) > 0 {
			fuzzySymbols = fuzzySymbols[:len(fuzzySymbols)%10]
			outputJSON := util.FormatFuzzySymbols(fuzzySymbols)
			s.ChannelMessageSend(m.ChannelID, outputJSON)
			return
		}
	}

	if quote == nil {
		logger.Trace(fmt.Sprintf("nil quote from ticker: %s", ticker))
	}

	message := util.FormatQuote(quote)
	s.ChannelMessageSend(m.ChannelID, message)
}
