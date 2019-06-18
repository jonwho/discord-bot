package cmd

import (
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/BryanSLam/discord-bot/util"
	iex "github.com/jonwho/go-iex"
)

// NewErCommand TODO: @doc
func NewErCommand() *Command {
	return &Command{
		Match: func(s string) bool {
			return regexp.MustCompile(`(?i)^!er [\w.]+$`).MatchString(s)
		},
		Fn: Er,
	}
}

// Er TODO: @doc
func Er(rw io.ReadWriter, logger *util.Logger, _ map[string]interface{}) {
	buf, err := ioutil.ReadAll(rw)
	if err != nil {
		rw.Write([]byte(err.Error()))
		return
	}
	slice := strings.Split(string(buf), " ")
	ticker := slice[1]
	iexSecretToken := os.Getenv("IEX_SECRET_TOKEN")
	iexClient, err := iex.NewClient(iexSecretToken)
	if err != nil {
		logger.Trace("IEX client initialization failed. Message: " + err.Error())
		rw.Write([]byte(err.Error()))
		return
	}

	logger.Info("Fetching earnings report info for " + ticker)
	earnings, err := iexClient.Earnings(ticker, nil)

	if err != nil {
		logger.Trace("IEX request failed. Message: " + err.Error())
		rw.Write([]byte(err.Error()))
		return
	}

	message := util.FormatEarnings(earnings)

	rw.Write([]byte(message))
}
