package cmd

import (
	"io"
	"io/ioutil"
	"regexp"

	"github.com/BryanSLam/discord-bot/util"
)

// NewPingCommand TODO: @doc
func NewPingCommand() *Command {
	return &Command{
		Match: func(s string) bool {
			return regexp.MustCompile(`(?i)^!ping$`).MatchString(s)
		},
		Fn: Ping,
	}
}

// Ping TODO: @doc
func Ping(rw io.ReadWriter, _ *util.Logger, _ map[string]interface{}) {
	ioutil.ReadAll(rw) // empty buffer before writing
	rw.Write([]byte("pong!"))
}
