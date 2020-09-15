package discordbot

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
	// "time"

	dg "github.com/bwmarrin/discordgo"
	"github.com/layeh/gopus"
)

const (
	channels  int = 2                   // 1 for mono, 2 for stereo
	frameRate int = 48000               // audio sampling rate
	frameSize int = 960                 // uint16 size of each audio frame
	maxBytes  int = (frameSize * 2) * 2 // max size of opus data
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

// WithIsVoiceChannel set flag for voice channel
func WithIsVoiceChannel(isVoiceChannel bool) DiscordWriterOption {
	return func(dw *DiscordWriter) error {
		dw.isVoiceChannel = isVoiceChannel
		return nil
	}
}

// WithGuildID sets the guild ID
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
func (w *DiscordWriter) Write(b []byte) (int, error) {
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
			return 0, err
		}

		// Send "speaking" packet over the voice websocket
		err = voiceChannel.Speaking(true)
		if err != nil {
			return 0, err
		}

		// Send not "speaking" packet over the websocket when we finish
		defer func() {
			voiceChannel.Speaking(false)
		}()

		send := make(chan []int16, 2)
		defer close(send)

		close := make(chan bool)
		go func() {
			sendPCM(voiceChannel, send)
			close <- true
		}()

		buf := bytes.NewReader(b)

		for {
			// read data from ffmpeg stdout
			audiobuf := make([]int16, frameSize*channels)
			err = binary.Read(buf, binary.LittleEndian, &audiobuf)
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				return 0, err
			}
			if err != nil {
				return 0, err
			}

			// Send received PCM to the sendPCM channel
			select {
			case send <- audiobuf:
			case <-close:
				return 0, err
			}
		}
	} else if w.channelID != "" {
		// if writer has channel then send bytes to it
		w.session.ChannelMessageSend(w.channelID, string(b))
	} else {
		// else send bytes to respond to the channel where message came from
		w.session.ChannelMessageSend(w.messageCreate.ChannelID, string(b))
	}

	return len(b), nil
}

// ChannelMessageSend send string to channel
func (w *DiscordWriter) ChannelMessageSend(content string) (*dg.Message, error) {
	return w.session.ChannelMessageSend(w.channelID, content)
}

// Pin takes slice of bytes and creates a PinMessage to send to the channel
func (w *DiscordWriter) Pin(messageID string) error {
	return w.session.ChannelMessagePin(w.channelID, messageID)
}

// SendPCM will receive on the provided channel encode
// received PCM data into Opus then send that to Discordgo
func sendPCM(v *dg.VoiceConnection, pcm <-chan []int16) {
	if pcm == nil {
		return
	}

	opusEncoder, err := gopus.NewEncoder(frameRate, channels, gopus.Audio)

	if err != nil {
		return
	}

	for {
		// read pcm from chan, exit if channel is closed.
		recv, ok := <-pcm
		if !ok {
			return
		}

		// try encoding pcm frame with Opus
		opus, err := opusEncoder.Encode(recv, frameSize, maxBytes)
		if err != nil {
			return
		}

		if v.Ready == false || v.OpusSend == nil {
			return
		}
		// send encoded opus data to the sendOpus channel
		v.OpusSend <- opus
	}
}
