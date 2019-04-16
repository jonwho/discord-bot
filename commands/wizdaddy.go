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

type wizdaddyCommand struct {
	regex *regexp.Regexp
}

func newWizdaddyCommand() wizdaddyCommand {
	return wizdaddyCommand{regexp.MustCompile(`(?i)^!wizdaddy$`)}
}

func (cmd wizdaddyCommand) match(s string) bool {
	return cmd.regex.MatchString(s)
}

func (cmd wizdaddyCommand) fn(s *dg.Session, m *dg.MessageCreate) {
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
