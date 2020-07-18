package discordbot

import (
	"log"
	// "time"

	dg "github.com/bwmarrin/discordgo"
)

// DiscordWriter implements io.Writer
type DiscordWriter struct {
	session        *dg.Session
	messageCreate  *dg.MessageCreate
	channelID      string
	guildID        string
	isVoiceChannel bool
}

// DiscordWriterOption configures additional properties on DiscordWriter
type DiscordWriterOption func(dw *DiscordWriter) error

func WithIsVoiceChannel(isVoiceChannel bool) DiscordWriterOption {
	return func(dw *DiscordWriter) error {
		dw.isVoiceChannel = isVoiceChannel
		return nil
	}
}

func WithGuildID(guildID string) DiscordWriterOption {
	return func(dw *DiscordWriter) error {
		dw.guildID = guildID
		return nil
	}
}

// NewDiscordWriter returns struct that implements io.Writer
func NewDiscordWriter(s *dg.Session, m *dg.MessageCreate, ch string, options ...DiscordWriterOption) *DiscordWriter {
	dw := &DiscordWriter{
		session:       s,
		messageCreate: m,
		channelID:     ch,
	}

	for _, option := range options {
		if err := option(dw); err != nil {
			return nil
		}
	}

	return dw
}

// Write sends the bytes to the Discord channel which can be text or voice channel
func (w *DiscordWriter) Write(b []byte) (n int, err error) {
	// 1. try voice channel
	// 2. try text channel
	// 3. default to respond to text channel
	if w.isVoiceChannel {
		mute := false
		deaf := false
		voiceChannel, err := w.session.ChannelVoiceJoin(w.guildID, w.channelID, mute, deaf)
		if err != nil {
			log.Println("voice channel error")
			log.Println(err)
		}
		// time.Sleep(time.Millisecond * 250)
		voiceChannel.Speaking(true)
		voiceChannel.OpusSend <- b
		// time.Sleep(time.Millisecond * 250)
		// voiceChannel.Speaking(false)
	} else if w.channelID != "" {
		// if writer has channel then send bytes to it
		w.session.ChannelMessageSend(w.channelID, string(b))
	} else {
		// else send bytes to respond to the channel where message came from
		w.session.ChannelMessageSend(w.messageCreate.ChannelID, string(b))
	}

	return len(b), nil
}
