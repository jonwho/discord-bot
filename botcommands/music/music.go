// +build ignore

package music

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	dg "github.com/bwmarrin/discordgo"
	"github.com/layeh/gopus"
	"github.com/rylio/ytdl"
)

type Music struct {
	httpClient *http.Client
}

// Option used with ctor to modify Stock struct
type Option func(s *Music) error

// WithHTTPClient sets the http.Client
func WithHTTPClient(httpClient *http.Client) Option {
	return func(m *Music) error {
		m.httpClient = httpClient
		return nil
	}
}

func New(options ...Option) (*Music, error) {
	music := &Music{}

	// apply options
	for _, option := range options {
		if err := option(music); err != nil {
			return nil, err
		}
	}

	return music, nil
}

const (
	channels      int    = 2                   // 1 for mono, 2 for stereo
	frameRate     int    = 48000               // audio sampling rate
	frameSize     int    = 960                 // uint16 size of each audio frame
	maxBytes      int    = (frameSize * 2) * 2 // max size of opus data
	testGuildID   string = "422149443266281492"
	testChannelID string = "422149443702358037"
)

func (*Music) PlayTest(s *dg.Session, m *dg.MessageCreate) error {
	dgv, err := s.ChannelVoiceJoin(testGuildID, testChannelID, false, true)
	if err != nil {
		return err
	}

	opusTitle := fmt.Sprintf("ytdl_data/%s.opus", "rzc3_b_KnHc")
	playAudioFile(dgv, opusTitle, make(chan bool))

	return nil
}

// N.B. Discord Voice API requires audio to be encoded with Opus
func (m *Music) Execute(ctx context.Context, rw io.ReadWriter) error {
	// steps
	// 1. download the youtube video
	buf, err := ioutil.ReadAll(rw)
	if err != nil {
		return err
	}
	argSplit := strings.Split(string(buf), " ")
	ytURL := argSplit[1]
	ytdlClient := ytdl.DefaultClient
	videoInfo, err := ytdlClient.GetVideoInfo(ctx, ytURL)
	if err != nil {
		return err
	}
	mp4Title := fmt.Sprintf("ytdl_data/%s.mp4", videoInfo.ID)
	file, err := os.Create(mp4Title)
	if err != nil {
		return err
	}
	defer file.Close()
	err = ytdlClient.Download(ctx, videoInfo, videoInfo.Formats[0], file)
	if err != nil {
		return err
	}

	// 2. strip audio from video with ffmpeg
	opusTitle := fmt.Sprintf("ytdl_data/%s.opus", videoInfo.ID)
	if !fileExists(opusTitle) {
		ffmpegArgStr := fmt.Sprintf("-i %s %s", mp4Title, opusTitle)
		args := strings.Split(ffmpegArgStr, " ")
		var stderr bytes.Buffer
		cmd := exec.Command("ffmpeg", args...)
		cmd.Stderr = &stderr
		err = cmd.Run()
		if err != nil {
			log.Println(stderr.String())
			return errors.New(err.Error() + " " + stderr.String())
		}
	}

	// 3. stream the audio into the discord socket
	run := exec.Command("ffmpeg", "-i", opusTitle, "-f", "s16le", "-ar", strconv.Itoa(frameRate), "-ac", strconv.Itoa(channels), "pipe:1")
	ffmpegout, err := run.StdoutPipe()
	if err != nil {
		return err
	}
	ffmpegbuf := bufio.NewReaderSize(ffmpegout, 16384)

	// Starts the ffmpeg command
	err = run.Start()
	if err != nil {
		return err
	}

	audiobuf, err := ioutil.ReadAll(ffmpegbuf)
	if err != nil {
		return err
	}

	rw.Write(audiobuf)

	return nil
}

func fileExists(filename string) bool {
	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
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

func playAudioFile(v *dg.VoiceConnection, filename string, stop <-chan bool) {
	// Create a shell command "object" to run.
	run := exec.Command("ffmpeg", "-i", filename, "-f", "s16le", "-ar", strconv.Itoa(frameRate), "-ac", strconv.Itoa(channels), "pipe:1")
	ffmpegout, err := run.StdoutPipe()
	if err != nil {
		return
	}

	ffmpegbuf := bufio.NewReaderSize(ffmpegout, 16384)

	// Starts the ffmpeg command
	err = run.Start()
	if err != nil {
		return
	}

	go func() {
		<-stop
		err = run.Process.Kill()
	}()

	// Send "speaking" packet over the voice websocket
	err = v.Speaking(true)
	if err != nil {
	}

	// Send not "speaking" packet over the websocket when we finish
	defer func() {
		err := v.Speaking(false)
		if err != nil {
		}
	}()

	send := make(chan []int16, 2)
	defer close(send)

	close := make(chan bool)
	go func() {
		sendPCM(v, send)
		close <- true
	}()

	for {
		// read data from ffmpeg stdout
		audiobuf := make([]int16, frameSize*channels)
		err = binary.Read(ffmpegbuf, binary.LittleEndian, &audiobuf)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return
		}
		if err != nil {
			return
		}

		// Send received PCM to the sendPCM channel
		select {
		case send <- audiobuf:
		case <-close:
			return
		}
	}
}
