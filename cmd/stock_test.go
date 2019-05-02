package cmd

import (
	"bytes"
	"testing"

	"github.com/BryanSLam/discord-bot/util"
)

func TestNewStockCommand(t *testing.T) {
	stockCmd := NewStockCommand()
	if stockCmd.Match("stock") == true {
		t.Errorf("\nExpected: %v\nActual: %v\n", false, true)
	}
	if stockCmd.Match("stock jd") == true {
		t.Errorf("\nExpected: %v\nActual: %v\n", false, true)
	}
	if stockCmd.Match("!stock jd") == false {
		t.Errorf("\nExpected: %v\nActual: %v\n", true, false)
	}
}

func TestStock(t *testing.T) {
	buf := bytes.NewBuffer([]byte("!stock jd"))
	logWriter := bytes.NewBuffer([]byte{})
	logger := util.NewLogger(logWriter)
	Stock(buf, logger, nil)
	actual := buf.String()
	if actual == "" {
		t.Errorf("Expected buf to be written to with stock quote but got\n%v\n", actual)
	}
}
