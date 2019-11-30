package botcommands

import (
	"context"
	"io"
)

// Stock dependencies go here
type Stock struct {
}

// Execute fetches quote for the ticker symbol
func (s *Stock) Execute(ctx context.Context, rw io.ReadWriter) error {
	ch := make(chan string)

	go func() {
		ch <- "foobar"
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case quote := <-ch:
		rw.Write([]byte(quote))
	}

	return nil
}
