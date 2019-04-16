package commands

import (
	dg "github.com/bwmarrin/discordgo"
)

type command struct {
	match func(s string) bool
	fn    func(s *dg.Session, m *dg.MessageCreate)
}
