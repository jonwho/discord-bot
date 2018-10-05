package commands

import (
	"regexp"

	"github.com/BryanSLam/discord-bot/config"
	dg "github.com/bwmarrin/discordgo"
)

func Commander() func(s *dg.Session, m *dg.MessageCreate) {
	commandRegex := regexp.MustCompile(`![a-zA-Z]+[ a-zA-Z\"\.]*[ 0-9/]*`)
	pingRegex := regexp.MustCompile(`!ping`)
	stockRegex := regexp.MustCompile(`(?i)^!stock [a-zA-Z\.]+$`)
	erRegex := regexp.MustCompile(`(?i)^!er [a-zA-Z]+$`)
	wizdaddyRegex := regexp.MustCompile(`(?i)^!wizdaddy$`)
	coinRegex := regexp.MustCompile(`(?i)^!coin [a-zA-Z]+$`)

	return func(s *dg.Session, m *dg.MessageCreate) {
		if commandRegex.MatchString(m.Content) {
			// Ignore all messages created by the bot itself
			// This isn't required in this specific example but it's a good practice.
			if m.Author.ID == s.State.User.ID {
				return
			}

			if pingRegex.MatchString(m.Content) {
				Ping(s, m)
				return
			}

			if stockRegex.MatchString(m.Content) {
				Stock(s, m)
				return
			}

			if erRegex.MatchString(m.Content) {
				Er(s, m)
				return
			}

			if wizdaddyRegex.MatchString(m.Content) {
				Wizdaddy(s, m)
				return
			}

			if coinRegex.MatchString(m.Content) {
				Coin(s, m)
				return
			}

			s.ChannelMessageSend(m.ChannelID, config.GetConfig().InvalidCommandMessage)
		}
	}
}
