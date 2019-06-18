package cmd

import (
	"bytes"
	"testing"

	"github.com/BryanSLam/discord-bot/util"
	dg "github.com/bwmarrin/discordgo"
)

func TestNewRemindmeCommand(t *testing.T) {
	remindmeCmd := NewRemindmeCommand()
	if actual := remindmeCmd.Match("remindme"); actual {
		t.Errorf("\nExpected: %v\nActual: %v\n", false, actual)
	}
	if actual := remindmeCmd.Match("remindme to not fail this test 04/20/69"); actual {
		t.Errorf("\nExpected: %v\nActual: %v\n", false, actual)
	}
	if actual := remindmeCmd.Match("remindme to not fail this test 04/20/69"); actual {
		t.Errorf("\nExpected: %v\nActual: %v\n", false, actual)
	}
	if actual := remindmeCmd.Match("!remindme to not fail this test 04/20/69"); !actual {
		t.Errorf("\nExpected: %v\nActual: %v\n", true, actual)
	}
}

func TestRemindme(t *testing.T) {
	buf := bytes.NewBuffer([]byte("!remindme to not fail this test 04/20/69"))
	logWriter := bytes.NewBuffer([]byte{})
	logger := util.NewLogger(logWriter)
	messageCreate := &dg.MessageCreate{
		Message: &dg.Message{
			ChannelID: "blahblah",
			Author:    &dg.User{ID: "foobar"},
		},
	}
	Remindme(buf, logger, map[string]interface{}{"messageCreate": messageCreate})
	expected := `Date has already passed`
	actual := buf.String()
	if expected != actual {
		t.Errorf("\nExpected: %v\nActual: %v\n", expected, actual)
	}
}
