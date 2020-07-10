package discordbot

import (
	"context"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/BryanSLam/discord-bot/botcommands"
	"github.com/BryanSLam/discord-bot/commands"
	"github.com/BryanSLam/discord-bot/util"
	dg "github.com/bwmarrin/discordgo"
)

// Bot data container for bot
type Bot struct {
	*dg.Session
	logger      *log.Logger
	maintainers []string

	// TODO: thinking about deprecating these
	cmds            []*commands.Command
	botLogChannelID string
}

// Option modifiers on bot initialization
type Option func(b *Bot) error

// New return new bot service
func New(token string, options ...Option) (*Bot, error) {
	session, err := dg.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	bot := &Bot{Session: session}
	for _, option := range options {
		if err := option(bot); err != nil {
			return nil, err
		}
	}

	// TODO: thinking about deprecating this style for botcommands style to isolate dependencies
	bot.cmds = append(bot.cmds,
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

	// Register handlers to the session
	bot.Session.AddHandler(bot.Commander)
	bot.Session.AddHandler(bot.UserOnly(bot.HandleStock()))
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

// WithMaintainers set the maintainers for bot to be DM'd when an issue happens
func WithMaintainers(maintainers []string) Option {
	return func(b *Bot) error {
		b.maintainers = maintainers
		return nil
	}
}

// WithBotLogChannelID set the bot log channel ID
func WithBotLogChannelID(botLogChannelID string) Option {
	return func(b *Bot) error {
		b.botLogChannelID = botLogChannelID
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

func (b *Bot) UserOnly(h func(_ *dg.Session, _ *dg.MessageCreate)) func(_ *dg.Session, _ *dg.MessageCreate) {
	return func(s *dg.Session, m *dg.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}
		h(s, m)
	}
}

// HandleStock is the bot command to call the stock command
func (b *Bot) HandleStock() func(s *dg.Session, m *dg.MessageCreate) {
	token := os.Getenv("IEX_SECRET_TOKEN")
	stock := botcommands.NewStock(token)
	stockRegex := regexp.MustCompile(`(?i)^\$[\w.]+$`)

	return func(s *dg.Session, m *dg.MessageCreate) {
		go func() {
			defer func() {
				if err := recover(); err != nil {
					b.logger.Println([]byte(util.MentionMaintainers(b.maintainers) + " an error has occurred"))
					b.logger.Println(err)
				}
			}()

			dr := NewDiscordReader(s, m, "")
			dw := NewDiscordWriter(s, m, "")
			drw := NewDiscordReadWriter(dr, dw)

			buf, err := ioutil.ReadAll(drw)
			if err != nil {
				drw.Write([]byte(err.Error()))
				return
			}

			if !stockRegex.MatchString(string(buf)) {
				return
			}

			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, time.Second*3)
			defer cancel()

			stock.Execute(ctx, drw)
		}()
	}
}
