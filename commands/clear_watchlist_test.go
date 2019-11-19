package commands

import (
	"bytes"
	"testing"
)

func TestNewClearWatchlistCommand(t *testing.T) {
	clearCmd := NewClearWatchlistCommand()
	if actual := clearCmd.Match("clearwatchlist"); actual {
		t.Errorf("\nExpected: %v\nActual: %v\n", false, actual)
	}
	if actual := clearCmd.Match("clearwatchlist "); actual {
		t.Errorf("\nExpected: %v\nActual: %v\n", false, actual)
	}
	if actual := clearCmd.Match("!clearwatchlist"); !actual {
		t.Errorf("\nExpected: %v\nActual: %v\n", true, actual)
	}
}

func TestClearwatchlist(t *testing.T) {
	buf := bytes.NewBuffer([]byte("!clearwatchlist"))
	Clearwatchlist(buf, nil, nil)
	expected := `watchlist cleared`
	actual := buf.String()
	if expected != actual {
		t.Errorf("\nExpected: %v\nActual: %v\n", expected, actual)
	}
}
