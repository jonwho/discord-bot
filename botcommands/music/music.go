package music

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

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

// N.B. Discord Voice API requires audio to be encoded with Opus
func (m *Music) Execute(ctx context.Context, rw io.ReadWriter) error {
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
	opusTitle := fmt.Sprintf("ytdl_data/%s.opus", videoInfo.ID)
	if !fileExists(opusTitle) {
		// ffmpegArgStr := fmt.Sprintf("-i %s -q:a 0 -map a %s", mp4Title, opusTitle)
		ffmpegArgStr := fmt.Sprintf("-i %s %s", mp4Title, opusTitle)
		args := strings.Split(ffmpegArgStr, " ")
		var stderr bytes.Buffer
		cmd := exec.Command("ffmpeg", args...)
		cmd.Stderr = &stderr
		err = cmd.Run()
		if err != nil {
			log.Println(stderr.String())
			return err
		}
	}

	// 3. stream the audio into the discord socket
	opusFile, err := os.Open(opusTitle)
	if err != nil {
		return err
	}
	var opuslen int16
	// var buffer = make([][]byte, 0)
	// for i := 0; i < 100; i++ {
	for {
		log.Println("BEFORE ", opuslen)
		// Read opus frame length from file
		// N.B. input and output are little-endian signed 16-bit PCM (pulse code modulation) files
		err = binary.Read(opusFile, binary.LittleEndian, &opuslen)
		log.Println("AFTER ", opuslen)

		// EOF stop
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			log.Println("EOF OR ErrUnexpectedEOF")
			log.Println(err)
			err := opusFile.Close()
			if err != nil {
				log.Println("FILE CLOSE ERROR")
				log.Println(err)
				return err
			}
			break
		}

		// report this err
		if err != nil {
			log.Println("WEIRD ERROR")
			log.Println(err)
			return err
		}

		// Read bytes
		inBuf := make([]byte, binary.MaxVarintLen64)
		binary.PutVarint(inBuf, int64(opuslen))
		// err = binary.Read(opusFile, binary.LittleEndian, &inBuf)

		// EOF errors should not exist
		// if err != nil {
		//   log.Println("ERR READING OPUS BYTES")
		//   log.Println(err)
		//   return err
		// }

		// buffer = append(buffer, inBuf)
		rw.Write(inBuf)
	}

	// counter := 1
	// for _, buf := range buffer {
	//   log.Println("looping buffer ", counter)
	//   counter++
	//   log.Println(buf)
	//   rw.Write(buf)
	// }

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
