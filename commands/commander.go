package commands

import (
	"fmt"
	"regexp"

	"github.com/BryanSLam/discord-bot/config"
	"github.com/BryanSLam/discord-bot/util"
	dg "github.com/bwmarrin/discordgo"
)

type work func(s *dg.Session, m *dg.MessageCreate)

func Commander() func(s *dg.Session, m *dg.MessageCreate) {
	commandRegex := regexp.MustCompile(`![a-zA-Z]+[ a-zA-Z\"\.]*[ 0-9/]*`)
	pingRegex := regexp.MustCompile(`!ping`)
	stockRegex := regexp.MustCompile(`(?i)^!stock [a-zA-Z\.]+$`)
	erRegex := regexp.MustCompile(`(?i)^!er [a-zA-Z]+$`)
	wizdaddyRegex := regexp.MustCompile(`(?i)^!wizdaddy$`)
	coinRegex := regexp.MustCompile(`(?i)^!coin [a-zA-Z]+$`)

	return func(s *dg.Session, m *dg.MessageCreate) {
		// TODO: refactor logger
		// OPTIONS:
		// 1. global logger
		// 2. check if pkg log supports write streams and if dg has stream to pass
		logger := util.Logger{Session: s, ChannelID: config.GetConfig().BotLogChannelID}

		if commandRegex.MatchString(m.Content) {
			// Ignore all messages created by the bot itself
			// This isn't required in this specific example but it's a good practice.
			if m.Author.ID == s.State.User.ID {
				return
			}

			if pingRegex.MatchString(m.Content) {
				go safelyDo(Ping, s, m, logger)
				return
			}

			if stockRegex.MatchString(m.Content) {
				go safelyDo(Stock, s, m, logger)
				return
			}

			if erRegex.MatchString(m.Content) {
				go safelyDo(Er, s, m, logger)
				return
			}

			if wizdaddyRegex.MatchString(m.Content) {
				go safelyDo(Wizdaddy, s, m, logger)
				return
			}

			if coinRegex.MatchString(m.Content) {
				go safelyDo(Coin, s, m, logger)
				return
			}

			s.ChannelMessageSend(m.ChannelID, config.GetConfig().InvalidCommandMessage)
		}
	}
}

func safelyDo(fn work, s *dg.Session, m *dg.MessageCreate, logger util.Logger) {
	// defer'd funcs will execute before return even if runtime panic
	defer func() {
		// use recover to catch panic so bot doesn't shutdown
		if err := recover(); err != nil {
			logger.Send(util.MentionMaintainers() + " an error has occurred")
			logger.Warn(fmt.Sprintln("function", util.FuncName(fn), "failed:", err))
		}
	}()
	// perform work
	fn(s, m)
}
