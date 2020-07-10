package discordbot

import (
	"context"
	"io"
)

// Command interface receives a Context and ReadWriter and executes its contained logic
type Command interface {
	Execute(ctx context.Context, rw io.ReadWriter) error
}
