package cmd

import (
	"io"
	"io/ioutil"
	"regexp"

	"github.com/BryanSLam/discord-bot/util"
)

// NewClearWatchlistCommand TODO: @doc
func NewClearWatchlistCommand() *Command {
	return &Command{
		Match: func(s string) bool {
			return regexp.MustCompile(`(?i)^!clearwatchlist$`).MatchString(s)
		},
		Fn: Clearwatchlist,
	}
}

// Clearwatchlist TODO: @doc
func Clearwatchlist(rw io.ReadWriter, _ *util.Logger, _ map[string]interface{}) {
	ioutil.ReadAll(rw) // empty buffer before writing
	redisClient.Del(watchlistRedisKey)
	rw.Write([]byte("watchlist cleared"))
}
