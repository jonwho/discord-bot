package commands

import (
	"bytes"
	"testing"

	"github.com/BryanSLam/discord-bot/util"
)

func TestNewNewsCommand(t *testing.T) {
	newsCmd := NewNewsCommand()
	if actual := newsCmd.Match("news"); actual {
		t.Errorf("\nExpected: %v\nActual: %v\n", false, actual)
	}
	if actual := newsCmd.Match("news jd"); actual {
		t.Errorf("\nExpected: %v\nActual: %v\n", false, actual)
	}
	if actual := newsCmd.Match("!news jd"); !actual {
		t.Errorf("\nExpected: %v\nActual: %v\n", true, actual)
	}
}

func TestNews(t *testing.T) {
	buf := bytes.NewBuffer([]byte("!news jd"))
	logWriter := bytes.NewBuffer([]byte{})
	logger := util.NewLogger(logWriter)
	News(buf, logger, nil)
	actual := buf.String()
	if actual == "" {
		t.Errorf("Expected buf to be written to but got\n%v\n", actual)
	}
}
