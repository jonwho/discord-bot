package commands

import (
	"bytes"
	"testing"

	"github.com/BryanSLam/discord-bot/util"
	dg "github.com/bwmarrin/discordgo"
)

func TestNewWatchlistCommand(t *testing.T) {
	watchlistCmd := NewWatchlistCommand()
	if actual := watchlistCmd.Match("watchlist"); actual {
		t.Errorf("\nExpected: %v\nActual: %v\n", false, actual)
	}
	if actual := watchlistCmd.Match("watchlist aapl"); actual {
		t.Errorf("\nExpected: %v\nActual: %v\n", false, actual)
	}
	if actual := watchlistCmd.Match("!watchlist aapl"); !actual {
		t.Errorf("\nExpected: %v\nActual: %v\n", true, actual)
	}
}

func TestWatchlist(t *testing.T) {
	buf := bytes.NewBuffer([]byte("!watchlist aapl"))
	logWriter := bytes.NewBuffer([]byte{})
	logger := util.NewLogger(logWriter)
	messageCreate := &dg.MessageCreate{
		Message: &dg.Message{
			ChannelID: "blahblah",
			Author:    &dg.User{ID: "foobar"},
		},
	}
	Watchlist(buf, logger, map[string]interface{}{"messageCreate": messageCreate})
	expected := `Adding ticker aapl to watchlist`
	actual := buf.String()
	if expected != actual {
		t.Errorf("\nExpected: %v\nActual: %v\n", expected, actual)
	}
}
