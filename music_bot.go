package discordbot

import (
	"context"
	"io/ioutil"
	"regexp"
	"time"

	bmusic "github.com/BryanSLam/discord-bot/botcommands/music"
	dg "github.com/bwmarrin/discordgo"
)

// NewMusicBot returns bot that only handles music requests
// N.B. looks like you can have multiple sessions with the same bot token
// with this approach can instantiate individual bots to handle different services
// e.g. stock bot, music bot, twitter bot, etc
func NewMusicBot(botToken string, options ...Option) (*Bot, error) {
	session, err := dg.New("Bot " + botToken)
	if err != nil {
		return nil, err
	}

	bot := &Bot{
		Session:  session,
		botToken: botToken,
	}
	for _, option := range options {
		if err := option(bot); err != nil {
			return nil, err
		}
	}

	bot.Session.AddHandler(bot.Unpanic(bot.UserOnly(bot.HandleMusic())))
	return bot, nil
}

// TODO: going to need multiple command handlers
func (b *Bot) HandleMusic() func(s *dg.Session, m *dg.MessageCreate) {
	music, _ := bmusic.New()
	// TODO: fix regex later
	musicRegex := regexp.MustCompile(`(?i)^\!music$`)

	b.logger.Println("Test logger")

	return func(s *dg.Session, m *dg.MessageCreate) {
		dr := NewDiscordReader(s, m, "")
		dw := NewDiscordWriter(s, m, "")
		drw := NewDiscordReadWriter(dr, dw)

		buf, err := ioutil.ReadAll(drw)
		if err != nil {
			drw.Write([]byte(err.Error()))
			return
		}

		if !musicRegex.MatchString(string(buf)) {
			return
		}

		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, time.Second*3)
		defer cancel()

		err = music.Execute(ctx, drw)
		if err != nil {
			drw.Write([]byte(err.Error()))
			return
		}
	}
}
