package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/BryanSLam/discord-bot/cmd"
	"github.com/BryanSLam/discord-bot/util"
	dg "github.com/bwmarrin/discordgo"
)

var (
	commandRegex    = regexp.MustCompile(`(?i)^![\w]+[\w ".]*[ 0-9/]*$`)
	commands        []*cmd.Command
	maintainers     []string
	botLogChannelID string
)

func init() {
	maintainers = strings.Split(os.Getenv("MAINTAINERS"), ",")
	botLogChannelID = os.Getenv("BOT_LOG_CHANNEL_ID")

	commands = append(commands,
		cmd.NewPingCommand(),
		cmd.NewStockCommand(),
		cmd.NewErCommand(),
		cmd.NewWizdaddyCommand(),
		cmd.NewCoinCommand(),
		cmd.NewRemindmeCommand(),
		cmd.NewWatchlistCommand(),
		cmd.NewClearWatchlistCommand(),
		cmd.NewNewsCommand(),
		cmd.NewNextErCommand(),
	)
}

func commander(s *dg.Session, m *dg.MessageCreate) {
	if commandRegex.MatchString(m.Content) {
		// Ignore all messages created by the bot itself
		// This isn't required in this specific example but it's a good practice.
		if m.Author.ID == s.State.User.ID {
			return
		}

		dr := util.NewDiscordReader(s, m, "")
		dw := util.NewDiscordWriter(s, m, "")
		drw := util.NewDiscordReadWriter(dr, dw)

		logWriter := util.NewDiscordWriter(s, nil, botLogChannelID)
		logger := util.NewLogger(logWriter)

		for _, c := range commands {
			if c.Match(m.Content) {
				go func() {
					defer func() {
						if err := recover(); err != nil {
							logger.Write([]byte(util.MentionMaintainers(maintainers) + " an error has occurred"))
							logger.Warn(fmt.Sprintln("function", util.FuncName(c.Fn), "failed:", err))
						}
					}()

					// TODO: use context and change command to an interface instead of struct
					mm := map[string]interface{}{}
					mm["messageCreate"] = m

					c.Fn(drw, logger, mm)
				}()
				return
			}
		}

		drw.Write([]byte("Invalid Command"))
	}
}
