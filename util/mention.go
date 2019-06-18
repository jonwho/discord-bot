package util

import (
	"fmt"
	"strings"
)

// MentionMaintainers TODO: @doc
func MentionMaintainers(maintainers []string) string {
	for index, maintainer := range maintainers {
		maintainers[index] = fmt.Sprintf("<@%s>", maintainer)
	}
	return strings.Join(maintainers, " ")
}
