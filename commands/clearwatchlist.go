package commands

import (
	"regexp"

	dg "github.com/bwmarrin/discordgo"
)

type clearwatchlistCommand struct {
	regex *regexp.Regexp
}

func newClearWatchlistCommand() clearwatchlistCommand {
	return clearwatchlistCommand{regexp.MustCompile(`(?i)^!clearwatchlist$`)}
}

func (cmd clearwatchlistCommand) match(s string) bool {
	return cmd.regex.MatchString(s)
}

// ClearWatchlist remove entire watchlist
func (cmd clearwatchlistCommand) fn(s *dg.Session, m *dg.MessageCreate) {
	redisClient.Del(watchlistRedisKey)
	s.ChannelMessageSend(m.ChannelID, "watchlist cleared")
}
