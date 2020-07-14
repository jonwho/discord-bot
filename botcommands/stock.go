package botcommands

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/BryanSLam/discord-bot/util"
	iex "github.com/jonwho/go-iex/v3"
)

// Stock dependencies go here
type Stock struct {
	iexToken  string
	iexClient *iex.Client
	alpacaID  string
	alpacaKey string
}

// NewStock - return a stock command interface to get quotes and stuff
func NewStock(iexToken, alpacaID, alpacaKey string) *Stock {
	iexClient, _ := iex.NewClient(iexToken)
	s := &Stock{
		iexToken:  iexToken,
		iexClient: iexClient,
		alpacaID:  alpacaID,
		alpacaKey: alpacaKey,
	}
	return s
}

// Execute fetches quote for the ticker symbol
func (s *Stock) Execute(ctx context.Context, rw io.ReadWriter) error {
	dataURL := "https://data.alpaca.markets/v1/bars/day"

	buf, err := ioutil.ReadAll(rw)
	if err != nil {
		return err
	}

	ticker := string(buf[1:])
	ticker = strings.ToUpper(ticker)

	ch := make(chan string)

	go func() {
		// get quote from IEX
		quote, err := s.iexClient.Quote(ticker, &iex.QuoteQueryParams{DisplayPercent: true})
		if err != nil {
			ch <- err.Error()
			return
		}
		if quote == nil {
			msg := fmt.Sprintf("nil quote from ticker: %s", ticker)
			ch <- msg
			return
		}

		// get bar from Alpaca
		req, err := http.NewRequest(http.MethodGet, dataURL, nil)
		if err != nil {
			ch <- err.Error()
			return
		}
		req.Header.Set("APCA-API-KEY-ID", s.alpacaID)
		req.Header.Set("APCA-API-SECRET-KEY", s.alpacaKey)
		q := req.URL.Query()
		q.Add("limit", "1")
		q.Add("symbols", ticker)
		req.URL.RawQuery = q.Encode()
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			ch <- err.Error()
			return
		}
		defer resp.Body.Close()

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			ch <- err.Error()
			return
		}
		data := map[string][]struct {
			Time   int64   `json:"t"`
			Open   float32 `json:"o"`
			High   float32 `json:"h"`
			Low    float32 `json:"l"`
			Close  float32 `json:"c"`
			Volume int32   `json:"v"`
		}{}
		err = json.Unmarshal(bodyBytes, &data)
		if err != nil {
			ch <- err.Error()
			return
		}
		if len(data[ticker]) == 0 {
			ch <- errors.New("Alpaca API no data found for " + ticker).Error()
			return
		}

		bar := data[ticker][len(data[ticker])-1]
		if err != nil {
			ch <- err.Error()
			return
		}

		quoteStr := util.FormatStock(quote, bar)
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
