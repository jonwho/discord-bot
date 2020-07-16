package discordbot

import (
	"context"
	"io/ioutil"
	"regexp"
	"time"

	bmusic "github.com/BryanSLam/discord-bot/botcommands/music"
	"github.com/BryanSLam/discord-bot/util"
	dg "github.com/bwmarrin/discordgo"
)

// NewMusicBot returns bot that only handles music requests
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

	bot.Session.AddHandler(bot.UserOnly(bot.HandleMusic()))
	return bot, nil
}

// TODO: going to need multiple command handlers
func (b *Bot) HandleMusic() func(s *dg.Session, m *dg.MessageCreate) {
	music, _ := bmusic.New()
	// TODO: fix regex later
	musicRegex := regexp.MustCompile(`(?i)^\!music$`)

	return func(s *dg.Session, m *dg.MessageCreate) {
		// TODO: this panic recovery should be captured in middleware
		defer func() {
			if err := recover(); err != nil {
				b.logger.Println([]byte(util.MentionMaintainers(b.maintainers) + " an error has occurred"))
				b.logger.Println(err)
			}

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

			music.Execute(ctx, drw)
		}()
	}
}
