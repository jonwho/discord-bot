package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/BryanSLam/discord-bot/util"
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

	logger.Info("Fetching stock info for " + ticker)
	quote, err := iexClient.Quote(ticker, nil)
	// TODO: cache the results from reference data before using this to preserve messaging limit
	// if err != nil {
	//   symbols, iexErr := iexClient.ExchangeSymbols()
	//   if iexErr != nil {
	//     logger.Trace("IEX request failed. Message: " + iexErr.Error())
	//     rw.Write([]byte(iexErr.Error()))
	//     return
	//   }
	//
	//   fuzzySymbols := util.FuzzySearch(ticker, symbols)
	//
	//   if len(fuzzySymbols) > 0 {
	//     fuzzySymbols = fuzzySymbols[:len(fuzzySymbols)%10]
	//     outputJSON := util.FormatFuzzySymbols(fuzzySymbols)
	//     rw.Write([]byte(outputJSON))
	//     return
	//   }
	// }

	if quote == nil {
		logger.Trace(fmt.Sprintf("nil quote from ticker: %s", ticker))
		return
	}

	message := util.FormatQuote(quote)
	rw.Write([]byte(message))
}
