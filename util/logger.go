package util

import (
	"fmt"
	"io"
	"time"
)

// Logger TODO: @doc
type Logger struct {
	io.Writer
}

// NewLogger TODO: @doc
func NewLogger(w io.Writer) *Logger {
	return &Logger{w}
}

// Use diff as the codeblock highlighter. Ghetto way of getting text colors.
// # BLUE
// + YELLOW-GREEN
// - RED

// Info TODO: @doc
func (l *Logger) Info(v ...interface{}) {
	timezone, _ := time.LoadLocation("America/Los_Angeles")
	t := time.Now().In(timezone)
	logDateTime := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())

	fmtStr := fmt.Sprintf("```md\n# [%s INFO] %s```", logDateTime, v)

	l.Write([]byte(fmtStr))
}

// Trace TODO: @doc
func (l *Logger) Trace(s string) {
	timezone, _ := time.LoadLocation("America/Los_Angeles")
	t := time.Now().In(timezone)
	logDateTime := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())

	fmtStr := fmt.Sprintf("```diff\n+ [%s TRACE] %s```", logDateTime, s)

	l.Write([]byte(fmtStr))
}

// Warn TODO: @doc
func (l *Logger) Warn(s string) {
	timezone, _ := time.LoadLocation("America/Los_Angeles")
	t := time.Now().In(timezone)
	logDateTime := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())

	fmtStr := fmt.Sprintf("```diff\n- [%s WARN] %s```", logDateTime, s)

	l.Write([]byte(fmtStr))
}

// Send TODO: @doc
func (l *Logger) Send(s string) {
	l.Write([]byte(s))
}
