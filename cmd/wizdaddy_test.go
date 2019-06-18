package cmd

import (
	"bytes"
	"testing"

	"github.com/BryanSLam/discord-bot/util"
)

func TestNewWizdaddyCommand(t *testing.T) {
	wizdaddyCmd := NewWizdaddyCommand()
	if wizdaddyCmd.Match("wizdaddy") == true {
		t.Errorf("\nExpected: %v\nActual: %v\n", false, true)
	}
	if wizdaddyCmd.Match(" !wizdaddy") == true {
		t.Errorf("\nExpected: %v\nActual: %v\n", false, true)
	}
	if wizdaddyCmd.Match("!wizdaddy ") == true {
		t.Errorf("\nExpected: %v\nActual: %v\n", false, true)
	}
	if wizdaddyCmd.Match("!wizdaddy") == false {
		t.Errorf("\nExpected: %v\nActual: %v\n", true, false)
	}
}

func TestWizdaddy(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	logWriter := bytes.NewBuffer([]byte{})
	logger := util.NewLogger(logWriter)
	Wizdaddy(buf, logger, nil)
	actual := buf.String()
	if actual == "" {
		t.Errorf("Expected buf to be written to with wizdaddy result but got\n%v\n", actual)
	}
}
