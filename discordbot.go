package discordbot

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	// "github.com/BryanSLam/discord-bot/botcommands"
	"github.com/BryanSLam/discord-bot/commands"
	"github.com/BryanSLam/discord-bot/util"
	dg "github.com/bwmarrin/discordgo"
)

// Bot data container for bot
type Bot struct {
	*dg.Session
	logger *log.Logger
}

// Command interface
type Command interface {
	Execute(ctx context.Context, rw io.ReadWriter) error
}

// MessageCreateHandlerFunc handler for event MESSAGE_CREATE
type MessageCreateHandlerFunc func(s *dg.Session, m *dg.MessageCreate)

// Option modifiers on bot initialization
type Option func(b *Bot) error

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
func New(token string, options ...Option) (*Bot, error) {
	session, err := dg.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	bot := &Bot{
		Session: session,
	}

	for _, option := range options {
		if err := option(bot); err != nil {
			return nil, err
		}
	}

	// Register handlers to the session
	bot.Session.AddHandler(commander)
	bot.Session.AddHandler(bot.userOnly(bot.HandleStock()))
	return bot, err
}

// WithLoggers set the writers for logging
func WithLoggers(writers ...io.Writer) Option {
	return func(b *Bot) error {
		w := io.MultiWriter(writers...)
		logger := log.New(w, "BOT LOG", log.LstdFlags)
		b.logger = logger
		return nil
	}
}

// Open starts a discord session
func (b *Bot) Open() error {
	err := b.Session.Open()
	return err
}

// Close closes a discord session
func (b *Bot) Close() {
	b.Session.Close()
}

func (b *Bot) userOnly(h func(_ *dg.Session, _ *dg.MessageCreate)) func(_ *dg.Session, _ *dg.MessageCreate) {
	log.Println("Called userOnly")
	return func(s *dg.Session, m *dg.MessageCreate) {
		log.Println("Enter userOnly func return")
		if m.Author.ID == s.State.User.ID {
			log.Println("Bot and Author matched")
			return
		}
		h(s, m)
	}
}

// HandleStock is the bot command to call the stock command
// design:
// 1. use a handler to decide if it can call the appropriate command
// 2. if the message matches the pattern then go ahead and make the call
//    to the command
func (b *Bot) HandleStock() func(s *dg.Session, m *dg.MessageCreate) {
	log.Println("Enter HandleStock")

	// if you need a closure set it up here before the func below
	return func(s *dg.Session, m *dg.MessageCreate) {
		log.Println("In HandleStock returned func")
		dr := NewDiscordReader(s, m, "")
		dw := NewDiscordWriter(s, m, "")
		drw := NewDiscordReadWriter(dr, dw)
		drw.Write([]byte("jontest"))
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, time.Second*3)
		defer cancel()
		// read from m and decide if we can call the command
		// symbol := "todo"
		// b.logger.Println("Requesting stock info for ticker", symbol)
		// stock := &botcommands.Stock{}
		// stock.Execute(ctx, drw)
	}
}

// commander return pattern matching handler
func commander(s *dg.Session, m *dg.MessageCreate) {
	if commandRegex.MatchString(m.Content) {
		// Ignore all messages created by the bot itself
		// This isn't required in this specific example but it's a good practice.
		if m.Author.ID == s.State.User.ID {
			return
		}

		dr := NewDiscordReader(s, m, "")
		dw := NewDiscordWriter(s, m, "")
		drw := NewDiscordReadWriter(dr, dw)

		logWriter := NewDiscordWriter(s, nil, botLogChannelID)
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
