package commands

import (
	"bytes"
	"testing"

	"github.com/BryanSLam/discord-bot/util"
)

func TestNewNexterCommand(t *testing.T) {
	nexterCmd := NewNextErCommand()
	if actual := nexterCmd.Match("nexter"); actual {
		t.Errorf("\nExpected: %v\nActual: %v\n", false, actual)
	}
	if actual := nexterCmd.Match("nexter 2"); actual {
		t.Errorf("\nExpected: %v\nActual: %v\n", false, actual)
	}
	if actual := nexterCmd.Match("!nexter 2"); !actual {
		t.Errorf("\nExpected: %v\nActual: %v\n", true, actual)
	}
}

func TestNexter(t *testing.T) {
	buf := bytes.NewBuffer([]byte("!nexter 2"))
	logWriter := bytes.NewBuffer([]byte{})
	logger := util.NewLogger(logWriter)
	Nexter(buf, logger, nil)
	actual := buf.String()
	if actual == "" {
		t.Errorf("Expected buf to be written to but got\n%v\n", actual)
	}
}
