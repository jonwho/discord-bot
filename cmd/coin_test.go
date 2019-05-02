package cmd

import (
	"bytes"
	"testing"

	"github.com/BryanSLam/discord-bot/util"
)

func TestNewCoin(t *testing.T) {
	coinCmd := NewCoinCommand()
	if coinCmd.Match("coin") == true {
		t.Errorf("\nExpected: %v\nActual: %v\n", false, true)
	}
	if coinCmd.Match("coin btc") == true {
		t.Errorf("\nExpected: %v\nActual: %v\n", false, true)
	}
	if coinCmd.Match(" coin btc") == true {
		t.Errorf("\nExpected: %v\nActual: %v\n", false, true)
	}
	if coinCmd.Match("coin btc ") == true {
		t.Errorf("\nExpected: %v\nActual: %v\n", false, true)
	}
	if coinCmd.Match("!coin btc") == false {
		t.Errorf("\nExpected: %v\nActual: %v\n", true, false)
	}
}

func TestCoin(t *testing.T) {
	buf := bytes.NewBuffer([]byte("!coin btc"))
	logWriter := bytes.NewBuffer([]byte{})
	logger := util.NewLogger(logWriter)
	Coin(buf, logger, nil)
	actual := buf.String()
	if actual == "" {
		t.Errorf("Expected buf to be written to with coin quote but got\n%v\n", actual)
	}
}
