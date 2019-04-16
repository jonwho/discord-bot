package commands

import (
	"regexp"

	dg "github.com/bwmarrin/discordgo"
)

func newClearWatchlistCommand() command {
	return command{
		match: func(s string) bool {
			return regexp.MustCompile(`(?i)^!clearwatchlist$`).MatchString(s)
		},
		fn: clearwatchlist,
	}
}

// remove entire watchlist
func clearwatchlist(s *dg.Session, m *dg.MessageCreate) {
	redisClient.Del(watchlistRedisKey)
	s.ChannelMessageSend(m.ChannelID, "watchlist cleared")
}
