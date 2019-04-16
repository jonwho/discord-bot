package commands

import (
	"regexp"

	dg "github.com/bwmarrin/discordgo"
)

type pingCommand struct {
	regex *regexp.Regexp
}

func newPingCommand() pingCommand {
	return pingCommand{regexp.MustCompile(`(?i)^!ping$`)}
}

func (cmd pingCommand) match(s string) bool {
	return cmd.regex.MatchString(s)
}

func (cmd pingCommand) fn(s *dg.Session, m *dg.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "pong!")
}
