package commands

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"github.com/BryanSLam/discord-bot/datasource"
	"github.com/BryanSLam/discord-bot/util"
	dg "github.com/bwmarrin/discordgo"
)

func newCoinCommand() command {
	return command{
		match: func(s string) bool {
			return regexp.MustCompile(`(?i)^!coin [\w]+$`).MatchString(s)
		},
		fn: coin,
	}
}

func coin(s *dg.Session, m *dg.MessageCreate) {
	slice := strings.Split(m.Content, " ")
	ticker := strings.ToUpper(slice[1])
	coinURL := coinAPIURL + ticker + "&tsyms=USD"
	logger := util.Logger{Session: s, ChannelID: botLogChannelID}

	logger.Info("Fetching coin info for: " + ticker)
	resp, err := http.Get(coinURL)

	if err != nil {
		logger.Trace("Coin request failed. Message: " + err.Error())
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	coin := datasource.Coin{Symbol: ticker}

	if err = json.NewDecoder(resp.Body).Decode(&coin); err != nil || coin.Response == "Error" {
		logger.Trace("JSON decoding failed. Message: " + err.Error())
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	s.ChannelMessageSend(m.ChannelID, coin.OutputJSON())
	defer resp.Body.Close()
}
