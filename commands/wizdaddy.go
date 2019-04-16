package commands

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/BryanSLam/discord-bot/datasource"
	"github.com/BryanSLam/discord-bot/util"
	dg "github.com/bwmarrin/discordgo"
)

func newWizdaddyCommand() command {
	return command{
		match: func(s string) bool {
			return regexp.MustCompile(`(?i)^!wizdaddy$`).MatchString(s)
		},
		fn: wizdaddy,
	}
}

func wizdaddy(s *dg.Session, m *dg.MessageCreate) {
	logger := util.Logger{Session: s, ChannelID: botLogChannelID}
	resp, err := http.Get(wizdaddyURL)

	if err != nil {
		logger.Trace("Wizdaddy request failed. Message: " + err.Error())
		s.ChannelMessageSend(m.ChannelID, "Daddy is down")
		return
	}

	var daddyResponse datasource.WizdaddyResponse
	if err = json.NewDecoder(resp.Body).Decode(&daddyResponse); err != nil {
		logger.Trace("JSON decoding failed. Message: " + err.Error())
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	s.ChannelMessageSend(m.ChannelID,
		fmt.Sprintf("%s %s %s %s", daddyResponse.Symbol,
			daddyResponse.StrikePrice, daddyResponse.ExpirationDate, daddyResponse.Type))
}
