package util

import (
	"fmt"
	"time"

	dg "github.com/bwmarrin/discordgo"
)

type Logger struct {
	Session   *dg.Session
	ChannelID string
}

// Use diff as the codeblock highlighter. Ghetto way of getting text colors.
// # BLUE
// + YELLOW-GREEN
// - RED
func (l *Logger) Info(s string) {
	timezone, _ := time.LoadLocation("America/Los_Angeles")
	t := time.Now().In(timezone)
	logDateTime := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())

	fmtStr := fmt.Sprintf("```md\n# [%s INFO] %s```", logDateTime, s)

	l.Session.ChannelMessageSend(l.ChannelID, fmtStr)
}

func (l *Logger) Trace(s string) {
	timezone, _ := time.LoadLocation("America/Los_Angeles")
	t := time.Now().In(timezone)
	logDateTime := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())

	fmtStr := fmt.Sprintf("```diff\n+ [%s TRACE] %s```", logDateTime, s)

	l.Session.ChannelMessageSend(l.ChannelID, fmtStr)
}

func (l *Logger) Warn(s string) {
	timezone, _ := time.LoadLocation("America/Los_Angeles")
	t := time.Now().In(timezone)
	logDateTime := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())

	fmtStr := fmt.Sprintf("```diff\n- [%s WARN] %s```", logDateTime, s)

	l.Session.ChannelMessageSend(l.ChannelID, fmtStr)
}

func (l *Logger) Send(s string) {
	l.Session.ChannelMessageSend(l.ChannelID, s)
}
