package cmd

import (
	"bytes"
	"testing"

	"github.com/BryanSLam/discord-bot/util"
)

func TestNewPingCommand(t *testing.T) {
	pingCmd := NewPingCommand()
	if pingCmd.Match("ping") == true {
		t.Errorf("Expected: %v\nActual: %v\n", false, true)
	}
	if pingCmd.Match("!ping ") == true {
		t.Errorf("Expected: %v\nActual: %v\n", false, true)
	}
	if pingCmd.Match(" !ping") == true {
		t.Errorf("Expected: %v\nActual: %v\n", false, true)
	}
	if pingCmd.Match("!ping") == false {
		t.Errorf("Expected: %v\nActual: %v\n", true, false)
	}
	if pingCmd.Match("!PING") == false {
		t.Errorf("Expected: %v\nActual: %v\n", true, false)
	}
	if pingCmd.Match("!pInG") == false {
		t.Errorf("Expected: %v\nActual: %v\n", true, false)
	}
}

func TestPing(t *testing.T) {
	buf := bytes.NewBuffer([]byte("!ping"))
	logWriter := bytes.NewBuffer([]byte{})
	logger := util.NewLogger(logWriter)
	Ping(buf, logger, nil)
	expected := `pong!`
	actual := buf.String()
	if expected != actual {
		t.Errorf("Expected: %v\nActual: %v\n", expected, actual)
	}
}
