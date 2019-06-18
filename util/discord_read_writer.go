package util

import (
	"io"
	"io/ioutil"
	"strings"

	dg "github.com/bwmarrin/discordgo"
)

// DiscordReadWriter TODO: @doc
type DiscordReadWriter struct {
	*DiscordReader
	*DiscordWriter
}

// DiscordReader TODO: @doc
type DiscordReader struct {
	session       *dg.Session
	channelID     string
	messageCreate *dg.MessageCreate
}

// DiscordWriter TODO: @doc
type DiscordWriter struct {
	session       *dg.Session
	channelID     string
	messageCreate *dg.MessageCreate
}

// NewDiscordReadWriter TODO: @app
func NewDiscordReadWriter(r *DiscordReader, w *DiscordWriter) *DiscordReadWriter {
	return &DiscordReadWriter{r, w}
}

// NewDiscordReader TODO: @doc
func NewDiscordReader(s *dg.Session, m *dg.MessageCreate, ch string) *DiscordReader {
	return &DiscordReader{
		session:       s,
		channelID:     ch,
		messageCreate: m,
	}
}

// NewDiscordWriter TODO: @doc
func NewDiscordWriter(s *dg.Session, m *dg.MessageCreate, ch string) *DiscordWriter {
	return &DiscordWriter{
		session:       s,
		channelID:     ch,
		messageCreate: m,
	}
}

func (r *DiscordReader) Read(b []byte) (int, error) {
	sr := strings.NewReader(r.messageCreate.Content)
	buf, err := ioutil.ReadAll(sr)
	if err != nil {
		return len(buf), err
	}
	copy(b, buf)
	return len(buf), io.EOF
}

func (w *DiscordWriter) Write(b []byte) (n int, err error) {
	if w.channelID != "" {
		w.session.ChannelMessageSend(w.channelID, string(b))
	} else {
		w.session.ChannelMessageSend(w.messageCreate.ChannelID, string(b))
	}

	return len(b), nil
}

// GetSession TODO: @doc
func (drw DiscordReadWriter) GetSession() *dg.Session {
	return drw.DiscordReader.session
}

// GetMessageCreate TODO: @doc
func (drw DiscordReadWriter) GetMessageCreate() *dg.MessageCreate {
	return drw.DiscordReader.messageCreate
}

// GetChannelID TODO: @doc
func (drw DiscordReadWriter) GetChannelID() string {
	return drw.DiscordReader.channelID
}
