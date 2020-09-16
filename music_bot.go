// +build ignore

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
//
// Bot not in a usable state but was fun learning about the OPUS codec.
// TODO:
// 1. Channels to communicate when to stop playing, leaving the voice channel, etc.
// 2. Prevent multiple voice streams (queue up additional songs instead).
// 3. Queue for songs.
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
	musicRegex := regexp.MustCompile(`(?i)^\!music [\w:\/\?\.-_]+$`)

	return func(s *dg.Session, m *dg.MessageCreate) {
		discordReader := NewDiscordReader(s, m, "")

		// find the channel where the message came from
		channel, err := s.State.Channel(m.ChannelID)
		if err != nil {
			b.logger.Println(err)
			return
		}

		// find the guild for that channel
		guild, err := s.State.Guild(channel.GuildID)
		if err != nil {
			b.logger.Println(err)
			return
		}

		// look for the message sender in the guild's voice channels
		var discordVoiceWriter *DiscordWriter
		for _, vs := range guild.VoiceStates {
			if vs.UserID == m.Author.ID {
				discordVoiceWriter = NewDiscordWriter(s, m, vs.ChannelID, WithIsVoiceChannel(true), WithGuildID(guild.ID))
			}
		}
		if discordVoiceWriter == nil {
			b.logger.Println("Could not find a voice channel to join. Try again.")
			return
		}

		buf, err := ioutil.ReadAll(discordReader)
		if err != nil {
			b.logger.Println(err)
			return
		}

		if !musicRegex.MatchString(string(buf)) {
			b.logger.Printf("%s did not match regex\n", string(buf))
			return
		}

		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, time.Minute*5)
		defer cancel()

		drw := NewDiscordReadWriter(discordReader, discordVoiceWriter)
		err = music.Execute(ctx, drw)
		if err != nil {
			b.logger.Println(err)
			return
		}
	}
}
