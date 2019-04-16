package commands

import (
	"regexp"

	dg "github.com/bwmarrin/discordgo"
)

func newPingCommand() command {
	return command{
		match: func(s string) bool {
			return regexp.MustCompile(`(?i)^!ping$`).MatchString(s)
		},
		fn: ping,
	}
}

func ping(s *dg.Session, m *dg.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "pong!")
}
