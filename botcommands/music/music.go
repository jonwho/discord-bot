package music

import (
	"context"
	"io"
	"net/http"
)

type Music struct {
	httpClient *http.Client
}

// Option used with ctor to modify Stock struct
type Option func(s *Music) error

// WithHTTPClient sets the http.Client
func WithHTTPClient(httpClient *http.Client) Option {
	return func(m *Music) error {
		m.httpClient = httpClient
		return nil
	}
}

func New(options ...Option) (*Music, error) {
	music := &Music{}

	// apply options
	for _, option := range options {
		if err := option(music); err != nil {
			return nil, err
		}
	}

	return music, nil
}

func (m *Music) Execute(ctx context.Context, rw io.ReadWriter) error {
	return nil
}
