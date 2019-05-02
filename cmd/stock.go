package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/BryanSLam/discord-bot/util"
	iex "github.com/jonwho/go-iex"
)

// NewStockCommand TODO: @doc
func NewStockCommand() *Command {
	return &Command{
		Match: func(s string) bool {
			return regexp.MustCompile(`(?i)^!stock [\w.]+$`).MatchString(s)
		},
		Fn: Stock,
	}
}

// Stock TODO: @doc
func Stock(rw io.ReadWriter, logger *util.Logger, m map[string]interface{}) {
	buf, err := ioutil.ReadAll(rw)
	if err != nil {
		rw.Write([]byte(err.Error()))
		return
	}

	slice := strings.Split(string(buf), " ")
	ticker := slice[1]
	iexClient, err := iex.NewClient()
	if err != nil {
		logger.Trace("IEX client initialization failed. Message: " + err.Error())
		rw.Write([]byte(err.Error()))
		return
	}

	logger.Info("Fetching stock info for " + ticker)
	quote, err := iexClient.Quote(ticker, true)
	if err != nil {
		rds, iexErr := iexClient.RefDataSymbols()
		if iexErr != nil {
			logger.Trace("IEX request failed. Message: " + iexErr.Error())
			rw.Write([]byte(iexErr.Error()))
			return
		}

		fuzzySymbols := util.FuzzySearch(ticker, rds.Symbols)

		if len(fuzzySymbols) > 0 {
			fuzzySymbols = fuzzySymbols[:len(fuzzySymbols)%10]
			outputJSON := util.FormatFuzzySymbols(fuzzySymbols)
			rw.Write([]byte(outputJSON))
			return
		}
	}

	if quote == nil {
		logger.Trace(fmt.Sprintf("nil quote from ticker: %s", ticker))
	}

	message := util.FormatQuote(quote)
	rw.Write([]byte(message))
}
