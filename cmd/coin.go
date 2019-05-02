package cmd

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/BryanSLam/discord-bot/datasource"
	"github.com/BryanSLam/discord-bot/util"
)

// NewCoinCommand TODO: @doc
func NewCoinCommand() *Command {
	return &Command{
		Match: func(s string) bool {
			return regexp.MustCompile(`(?i)^!coin [\w]+$`).MatchString(s)
		},
		Fn: Coin,
	}
}

// Coin TODO: @doc
func Coin(rw io.ReadWriter, logger *util.Logger, _ map[string]interface{}) {
	buf, err := ioutil.ReadAll(rw)
	if err != nil {
		rw.Write([]byte(err.Error()))
		return
	}

	slice := strings.Split(string(buf), " ")
	ticker := strings.ToUpper(slice[1])
	coinURL := coinAPIURL + ticker + "&tsyms=USD"

	logger.Info("Fetching coin info for: " + ticker)
	resp, err := http.Get(coinURL)

	if err != nil {
		logger.Trace("Coin request failed. Message: " + err.Error())
		rw.Write([]byte(err.Error()))
		return
	}

	coin := datasource.Coin{Symbol: ticker}

	if err = json.NewDecoder(resp.Body).Decode(&coin); err != nil || coin.Response == "Error" {
		logger.Trace("JSON decoding failed. Message: " + err.Error())
		rw.Write([]byte(err.Error()))
		return
	}

	rw.Write([]byte(coin.OutputJSON()))
	defer resp.Body.Close()
}
