package stock

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
	iex "github.com/jonwho/go-iex/v4"
)

// Stock dependencies go here
type Stock struct {
	iexToken   string
	iexClient  *iex.Client
	httpClient *http.Client
	alpacaID   string
	alpacaKey  string
}

// Option used with ctor to modify Stock struct
type Option func(s *Stock) error

// WithHTTPClient sets the http.Client
func WithHTTPClient(httpClient *http.Client) Option {
	return func(s *Stock) error {
		s.httpClient = httpClient
		return nil
	}
}

// WithIEXClient sets the iex.Client
func WithIEXClient(iexClient *iex.Client) Option {
	return func(s *Stock) error {
		s.iexClient = iexClient
		return nil
	}
}

// New returns a struct that implements `discordbot.Command`
func New(iexToken, alpacaID, alpacaKey string, options ...Option) (*Stock, error) {
	iexClient, _ := iex.NewClient(iexToken)
	stock := &Stock{
		iexToken:   iexToken,
		iexClient:  iexClient,
		httpClient: http.DefaultClient,
		alpacaID:   alpacaID,
		alpacaKey:  alpacaKey,
	}

	// apply options
	for _, option := range options {
		if err := option(stock); err != nil {
			return nil, err
		}
	}

	return stock, nil
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
		// get quote from IEX
		quoteCh := s.getQuote(ticker)
		// get bar from Alpaca
		barCh := s.getBar(ticker)

		quoteOrErr := <-quoteCh
		barOrErr := <-barCh

		// err check for quote
		if quoteOrErr.err != nil {
			ch <- quoteOrErr.err.Error()
			return
		}
		if quoteOrErr.quote == nil {
			msg := fmt.Sprintf("nil quote from ticker: %s", ticker)
			ch <- msg
			return
		}

		// err check for bar
		if barOrErr.err != nil {
			ch <- barOrErr.err.Error()
			return
		}

		quoteStr := util.FormatStock(quoteOrErr.quote, *barOrErr.bar)
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

// Generator pattern - return a channel that another goroutine will OBSERVE
// use channels to communicate between 2 goroutines
func (s *Stock) getQuote(ticker string) <-chan struct {
	quote *iex.Quote
	err   error
} {
	ch := make(chan struct {
		quote *iex.Quote
		err   error
	}, 1)

	go func() {
		quote, err := s.iexClient.Quote(ticker, &iex.QuoteQueryParams{DisplayPercent: true})
		ch <- struct {
			quote *iex.Quote
			err   error
		}{quote: quote, err: err}
	}()

	return ch
}

type alpacaBar struct {
	Time   int64   `json:"t"`
	Open   float32 `json:"o"`
	High   float32 `json:"h"`
	Low    float32 `json:"l"`
	Close  float32 `json:"c"`
	Volume int32   `json:"v"`
}

func (s *Stock) getBar(ticker string) <-chan struct {
	bar *alpacaBar
	err error
} {
	ch := make(chan struct {
		bar *alpacaBar
		err error
	}, 1)
	dataURL := "https://data.alpaca.markets/v1/bars/day"

	go func() {
		req, err := http.NewRequest(http.MethodGet, dataURL, nil)
		if err != nil {
			ch <- struct {
				bar *alpacaBar
				err error
			}{nil, err}
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
			ch <- struct {
				bar *alpacaBar
				err error
			}{nil, err}
			return
		}
		defer resp.Body.Close()

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			ch <- struct {
				bar *alpacaBar
				err error
			}{nil, err}
			return
		}
		data := map[string][]alpacaBar{}
		err = json.Unmarshal(bodyBytes, &data)
		if err != nil {
			ch <- struct {
				bar *alpacaBar
				err error
			}{nil, err}
			return
		}
		if len(data[ticker]) == 0 {
			err = errors.New("Alpaca API no data found for " + ticker)
			ch <- struct {
				bar *alpacaBar
				err error
			}{nil, err}
			return
		}

		bar := data[ticker][len(data[ticker])-1]
		ch <- struct {
			bar *alpacaBar
			err error
		}{&bar, err}
	}()

	return ch
}
