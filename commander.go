package discordbot

import (
	"fmt"
	"regexp"

	"github.com/BryanSLam/discord-bot/util"
	dg "github.com/bwmarrin/discordgo"
)

// Commander return pattern matching handler
func (b *Bot) Commander(s *dg.Session, m *dg.MessageCreate) {
	commandRegex := regexp.MustCompile(`(?i)^![\w]+[\w ".]*[ 0-9/]*$`)
	if commandRegex.MatchString(m.Content) {
		// Ignore all messages created by the bot itself
		// This isn't required in this specific example but it's a good practice.
		if m.Author.ID == s.State.User.ID {
			return
		}

		dr := NewDiscordReader(s, m, "")
		dw := NewDiscordWriter(s, m, "")
		drw := NewDiscordReadWriter(dr, dw)

		logWriter := NewDiscordWriter(s, nil, b.botLogChannelID)
		logger := util.NewLogger(logWriter)

		for _, c := range b.cmds {
			if c.Match(m.Content) {
				go func() {
					defer func() {
						if err := recover(); err != nil {
							logger.Write([]byte(util.MentionMaintainers(b.maintainers) + " an error has occurred"))
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
