package commands

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/BryanSLam/discord-bot/util"
	dg "github.com/bwmarrin/discordgo"
	iex "github.com/jonwho/go-iex"
)

type newsCommand struct {
	regex *regexp.Regexp
}

func newNewsCommand() newsCommand {
	return newsCommand{regexp.MustCompile(`(?i)^!news [\w.]+$`)}
}

func (cmd newsCommand) match(s string) bool {
	return cmd.regex.MatchString(s)
}

func (cmd newsCommand) fn(s *dg.Session, m *dg.MessageCreate) {
	logger := util.Logger{Session: s, ChannelID: botLogChannelID}
	slice := strings.Split(m.Content, " ")
	ticker := slice[1]
	iexClient, err := iex.NewClient()
	if err != nil {
		logger.Trace("IEX client initialization failed. Message: " + err.Error())
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	logger.Info("Fetching news for " + ticker)
	news, err := iexClient.News(ticker, 3)
	if err != nil {
		logger.Trace("IEX request failed. Message: " + err.Error())
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	if news == nil {
		logger.Trace(fmt.Sprintf("nil news from ticker: %s", ticker))
	}

	for _, e := range news.News {
		s.ChannelMessageSend(m.ChannelID, util.FormatNews(&e))
	}
}
