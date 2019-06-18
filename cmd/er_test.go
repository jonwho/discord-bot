package cmd

import (
	"bytes"
	"testing"

	"github.com/BryanSLam/discord-bot/util"
)

func TestNewErCommand(t *testing.T) {
	erCmd := NewErCommand()
	if erCmd.Match("er") == true {
		t.Errorf("\nExpected: %v\nActual: %v\n", false, true)
	}
	if erCmd.Match("er jd") == true {
		t.Errorf("\nExpected: %v\nActual: %v\n", false, true)
	}
	if erCmd.Match(" er jd") == true {
		t.Errorf("\nExpected: %v\nActual: %v\n", false, true)
	}
	if erCmd.Match("er jd ") == true {
		t.Errorf("\nExpected: %v\nActual: %v\n", false, true)
	}
	if erCmd.Match("!er jd") == false {
		t.Errorf("\nExpected: %v\nActual: %v\n", true, false)
	}
}

func TestEr(t *testing.T) {
	buf := bytes.NewBuffer([]byte("!er jd"))
	logWriter := bytes.NewBuffer([]byte{})
	logger := util.NewLogger(logWriter)
	Er(buf, logger, nil)
	actual := buf.String()
	if actual == "" {
		t.Errorf("Expected buf to be written to with stock earnings but got\n%v\n", actual)
	}
}
