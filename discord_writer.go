package discordbot

import (
	dg "github.com/bwmarrin/discordgo"
)

// DiscordWriter implements io.Writer
type DiscordWriter struct {
	session       *dg.Session
	messageCreate *dg.MessageCreate
	channelID     string
}

// NewDiscordWriter returns struct that implements io.Writer
func NewDiscordWriter(s *dg.Session, m *dg.MessageCreate, ch string) *DiscordWriter {
	return &DiscordWriter{
		session:       s,
		messageCreate: m,
		channelID:     ch,
	}
}

// Write sends the bytes to the Discord channel
func (w *DiscordWriter) Write(b []byte) (n int, err error) {
	if w.channelID != "" {
		// if writer has channel then send bytes to it
		w.session.ChannelMessageSend(w.channelID, string(b))
	} else {
		// else send bytes to respond to the channel where message came from
		w.session.ChannelMessageSend(w.messageCreate.ChannelID, string(b))
	}

	return len(b), nil
}
