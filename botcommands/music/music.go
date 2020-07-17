package music

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

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

func (m *Music) Execute(ctx context.Context, rw io.ReadWriter) error {
	// add additional 20 seconds to timeout?
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	// steps
	// 1. download the youtube video
	ytdlClient := ytdl.DefaultClient
	videoInfo, err := ytdlClient.GetVideoInfo(ctx, "https://www.youtube.com/watch?v=rzc3_b_KnHc")
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
	mp3Title := fmt.Sprintf("ytdl_data/%s.mp3", videoInfo.ID)
	ffmpegArgStr := fmt.Sprintf("-i %s -q:a 0 -map a %s", mp4Title, mp3Title)
	args := strings.Split(ffmpegArgStr, " ")
	cmd := exec.Command("ffmpeg", args...)
	err = cmd.Run()
	if err != nil {
		rw.Write([]byte(err.Error()))
	}

	// 3. stream the audio into the discord socket

	rw.Write([]byte("respond to !music command"))

	return nil
}
