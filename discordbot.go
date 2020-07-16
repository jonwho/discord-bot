package discordbot

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"regexp"
	"time"

	bstock "github.com/BryanSLam/discord-bot/botcommands/stock"
	"github.com/BryanSLam/discord-bot/commands"
	"github.com/BryanSLam/discord-bot/util"
	dg "github.com/bwmarrin/discordgo"
)

// Bot data container for bot
type Bot struct {
	*dg.Session

	botToken    string
	iexToken    string
	alpacaID    string
	alpacaKey   string
	logger      *log.Logger
	maintainers []string

	// TODO: thinking about deprecating these
	cmds            []*commands.Command
	altCmds         []*Command
	botLogChannelID string
}

// Option modifiers on bot initialization
type Option func(b *Bot) error

// New return new bot service
func New(botToken, iexToken, alpacaID, alpacaKey string, options ...Option) (*Bot, error) {
	session, err := dg.New("Bot " + botToken)
	if err != nil {
		return nil, err
	}

	bot := &Bot{
		Session:   session,
		botToken:  botToken,
		iexToken:  iexToken,
		alpacaID:  alpacaID,
		alpacaKey: alpacaKey,
	}
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

	// TODO: rethink this approach
	// maybe leave it up to Option funcs to add services to the bot instead
	// bot.altCmds = append(bot.altCmds, ...somethinghere)

	// Register handlers to the session
	bot.Session.AddHandler(bot.Commander)
	bot.Session.AddHandler(bot.UserOnly(bot.HandleStock()))
	return bot, err
}

// WithLoggers set the writers for logging
// TODO: touch this up to include logger on a discord channel
func WithLoggers(writers ...io.Writer) Option {
	return func(b *Bot) error {
		w := io.MultiWriter(writers...)
		logger := log.New(w, "BOT LOG ", log.LstdFlags)
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
		dw := NewDiscordWriter(b.Session, nil, botLogChannelID)
		w := io.MultiWriter(dw)
		logger := log.New(w, "BOT LOG ", log.LstdFlags)
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

// UserOnly ensure bot doesn't respond to itself
func (b *Bot) UserOnly(h func(_ *dg.Session, _ *dg.MessageCreate)) func(_ *dg.Session, _ *dg.MessageCreate) {
	return func(s *dg.Session, m *dg.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}
		h(s, m)
	}
}

// Unpanic ensure panic is captured and logged
func (b *Bot) Unpanic(h func(_ *dg.Session, _ *dg.MessageCreate)) func(_ *dg.Session, _ *dg.MessageCreate) {
	return func(s *dg.Session, m *dg.MessageCreate) {
		defer func() {
			if rec := recover(); rec != nil {
				var err error
				switch t := rec.(type) {
				case string:
					err = errors.New(t)
				case error:
					err = t
				default:
					err = errors.New("Unknown server error")
				}
				b.logger.Println(util.MentionMaintainers(b.maintainers) + " an error has occurred")
				b.logger.Println(err)
			}
		}()

		h(s, m)
	}
}

// HandleStock is the bot command to call the stock command
// TODO: new idea is to explicitly build the handlers which depend on services being
// DI'd onto the bot
//
// if the service is nil then skip the handler
func (b *Bot) HandleStock() func(s *dg.Session, m *dg.MessageCreate) {
	stock, _ := bstock.New(b.iexToken, b.alpacaID, b.alpacaKey)
	stockRegex := regexp.MustCompile(`(?i)^\$[\w.]+$`)

	// function closure so local variables above only happen once
	return func(s *dg.Session, m *dg.MessageCreate) {
		go func() {
			// TODO: this panic recovery should be captured in middleware
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
