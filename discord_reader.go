package discordbot

import (
	"io"
	"io/ioutil"
	"strings"

	dg "github.com/bwmarrin/discordgo"
)

// DiscordReader implements io.Reader
type DiscordReader struct {
	session       *dg.Session
	channelID     string
	messageCreate *dg.MessageCreate
}

// NewDiscordReader returns struct that implements io.Reader
func NewDiscordReader(s *dg.Session, m *dg.MessageCreate, ch string) *DiscordReader {
	return &DiscordReader{
		session:       s,
		channelID:     ch,
		messageCreate: m,
	}
}

// Read reads in bytes from a Discord channel
func (r *DiscordReader) Read(b []byte) (int, error) {
	sr := strings.NewReader(r.messageCreate.Content)
	buf, err := ioutil.ReadAll(sr)
	if err != nil {
		return len(buf), err
	}
	copy(b, buf)
	return len(buf), io.EOF
}
