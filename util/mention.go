package util

import (
	"fmt"
	"strings"

	"github.com/BryanSLam/discord-bot/config"
)

func MentionMaintainers() string {
	maintainers := strings.Split(config.GetConfig().Maintainers, ",")
	for index, maintainer := range maintainers {
		maintainers[index] = fmt.Sprintf("<@%s>", maintainer)
	}
	return strings.Join(maintainers, " ")
}
