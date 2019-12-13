package botcommands

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/BryanSLam/discord-bot/util"
	iex "github.com/jonwho/go-iex/v2"
)

// Stock dependencies go here
type Stock struct {
	token     string
	iexClient *iex.Client
}

// NewStock - return a stock command interface to get quotes and stuff
func NewStock(token string) *Stock {
	cli, _ := iex.NewClient(token)
	s := &Stock{token: token, iexClient: cli}
	return s
}

// Execute fetches quote for the ticker symbol
func (s *Stock) Execute(ctx context.Context, rw io.ReadWriter) error {
	buf, err := ioutil.ReadAll(rw)
	if err != nil {
		return err
	}

	ticker := string(buf[1:])
	ticker = strings.ToUpper(ticker)

	ch := make(chan string)

	go func() {
		quote, err := s.iexClient.Quote(ticker, struct {
			DisplayPercent bool `url:"displayPercent,omitempty"`
		}{true})
		if err != nil {
			ch <- err.Error()
			return
		}
		if quote == nil {
			msg := fmt.Sprintf("nil quote from ticker: %s", ticker)
			ch <- msg
			return
		}
		quoteStr := util.FormatQuote(quote)
		ch <- quoteStr
	}()

	select {
	case <-ctx.Done():
		rw.Write([]byte("Stock request timed out"))
		return ctx.Err()
	case quoteOrErr := <-ch:
		rw.Write([]byte(quoteOrErr))
	}

	return nil
}
