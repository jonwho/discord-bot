package earningsreport

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/BryanSLam/discord-bot/util"
	iex "github.com/jonwho/go-iex/v4"
)

type EarningsReport struct {
	iexToken   string
	iexClient  *iex.Client
	httpClient *http.Client
}

type Option func(er *EarningsReport) error

// WithHTTPClient sets the http.Client
func WithHTTPClient(httpClient *http.Client) Option {
	return func(er *EarningsReport) error {
		er.httpClient = httpClient
		return nil
	}
}

// WithIEXClient sets the iex.Client
func WithIEXClient(iexClient *iex.Client) Option {
	return func(er *EarningsReport) error {
		er.iexClient = iexClient
		return nil
	}
}

// New returns a struct that implements `discordbot.Command`
func New(iexToken string, options ...Option) (*EarningsReport, error) {
	iexClient, _ := iex.NewClient(iexToken)
	er := &EarningsReport{
		httpClient: http.DefaultClient,
		iexToken:   iexToken,
		iexClient:  iexClient,
	}

	// apply options
	for _, option := range options {
		if err := option(er); err != nil {
			return nil, err
		}
	}

	return er, nil
}

// Execute fetches earnings report(s) for the ticker symbol
func (er *EarningsReport) Execute(ctx context.Context, rw io.ReadWriter) error {
	buf, err := ioutil.ReadAll(rw)
	if err != nil {
		return err
	}

	bufSplit := strings.Split(string(buf[1:]), " ")
	ticker := bufSplit[0]
	ticker = strings.ToUpper(ticker)

	earningsCh := er.getEarnings(ticker)

	select {
	case <-ctx.Done():
		rw.Write([]byte("Earnings request timed out"))
		return ctx.Err()
	case report := <-earningsCh:
		if report.err != nil {
			rw.Write([]byte(err.Error()))
		}
		for _, message := range util.FormatAllEarnings(report.earnings) {
			rw.Write([]byte(message))
		}
	}

	return nil
}

func (er *EarningsReport) getEarnings(ticker string) <-chan struct {
	earnings *iex.Earnings
	err      error
} {
	ch := make(chan struct {
		earnings *iex.Earnings
		err      error
	})

	go func() {
		earnings, err := er.iexClient.Earnings(ticker, nil)
		ch <- struct {
			earnings *iex.Earnings
			err      error
		}{earnings: earnings, err: err}
	}()

	return ch
}
