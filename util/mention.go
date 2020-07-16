package util

import (
	"fmt"
	"strings"
)

// MentionMaintainers TODO: @doc
func MentionMaintainers(maintainers []string) string {
	var mentions []string
	for _, maintainer := range maintainers {
		mentions = append(mentions, fmt.Sprintf("<@%s>", maintainer))
	}
	return strings.Join(mentions, " ")
}
