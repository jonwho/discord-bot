package commands

import (
	"regexp"
	"strings"

	"github.com/BryanSLam/discord-bot/util"
	dg "github.com/bwmarrin/discordgo"
	iex "github.com/jonwho/go-iex"
)

func newErCommand() command {
	return command{
		match: func(s string) bool {
			return regexp.MustCompile(`(?i)^!er [\w.]+$`).MatchString(s)
		},
		fn: er,
	}
}

func er(s *dg.Session, m *dg.MessageCreate) {
	logger := util.Logger{Session: s, ChannelID: botLogChannelID}
	slice := strings.Split(m.Content, " ")
	ticker := slice[1]
	iexClient, err := iex.NewClient()
	if err != nil {
		logger.Trace("IEX client initialization failed. Message: " + err.Error())
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	logger.Info("Fetching earnings report info for " + ticker)
	earnings, err := iexClient.Earnings(ticker)

	if err != nil {
		logger.Trace("IEX request failed. Message: " + err.Error())
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	message := util.FormatEarnings(earnings)

	s.ChannelMessageSend(m.ChannelID, message)
}
