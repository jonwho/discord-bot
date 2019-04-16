package commands

import (
	dg "github.com/bwmarrin/discordgo"
)

type command interface {
	match(s string) bool
	fn(s *dg.Session, m *dg.MessageCreate)
}
