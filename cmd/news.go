package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/BryanSLam/discord-bot/util"
)

// NewNewsCommand TODO: @doc
func NewNewsCommand() *Command {
	return &Command{
		Match: func(s string) bool {
			return regexp.MustCompile(`(?i)^!news [\w.]+$`).MatchString(s)
		},
		Fn: News,
	}
}

// News TODO: @doc
func News(rw io.ReadWriter, logger *util.Logger, _ map[string]interface{}) {
	buf, err := ioutil.ReadAll(rw)
	if err != nil {
		rw.Write([]byte(err.Error()))
		return
	}

	slice := strings.Split(string(buf), " ")
	ticker := slice[1]

	logger.Info("Fetching news for " + ticker)
	news, err := iexClient.News(ticker, 3)
	if err != nil {
		logger.Trace("IEX request failed. Message: " + err.Error())
		rw.Write([]byte(err.Error()))
		return
	}

	if news == nil {
		logger.Trace(fmt.Sprintf("nil news from ticker: %s", ticker))
	}

	rw.Write([]byte(util.FormatNews(news)))
}
