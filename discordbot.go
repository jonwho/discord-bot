package discordbot

import (
	"context"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/BryanSLam/discord-bot/commands"
	"github.com/BryanSLam/discord-bot/util"
	dg "github.com/bwmarrin/discordgo"
)

// Bot data container for bot
type Bot struct {
	dg *dg.Session
}

// Command interface
type Command interface {
	Execute(ctx context.Context, rw io.ReadWriter) error
}

var (
	commandRegex    = regexp.MustCompile(`(?i)^![\w]+[\w ".]*[ 0-9/]*$`)
	cmds            []*commands.Command
	maintainers     []string
	botLogChannelID string
)

// TODO: discourage the use of init
// TOOD: use closures or sync.Once to get around one time setup
func init() {
	maintainers = strings.Split(os.Getenv("MAINTAINERS"), ",")
	botLogChannelID = os.Getenv("BOT_LOG_CHANNEL_ID")

	cmds = append(cmds,
		commands.NewPingCommand(),
		commands.NewStockCommand(),
		commands.NewErCommand(),
		commands.NewWizdaddyCommand(),
		commands.NewCoinCommand(),
		commands.NewRemindmeCommand(),
		commands.NewWatchlistCommand(),
		commands.NewClearWatchlistCommand(),
		commands.NewNewsCommand(),
		commands.NewNextErCommand(),
	)
}

// New return new bot service
func New(token string) (*dg.Session, error) {
	session, err := dg.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	// Register handlers to the session
	session.AddHandler(commander)
	return session, err
}

// commander return pattern matching handler
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

		for _, c := range cmds {
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
