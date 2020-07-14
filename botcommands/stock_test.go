package botcommands

import (
	"bytes"
	"testing"
	// "github.com/dnaeon/go-vcr/cassette"
	// "github.com/dnaeon/go-vcr/recorder"
)

func TestNewStock(t *testing.T) {
	// TODO
}

func TestStockExecute(t *testing.T) {
	buf := bytes.NewBuffer([]byte("$tsla"))
	actual := buf.String()

	// TODO: always error for now
	t.Errorf("Expected buf to be written to got\n%v\n", actual)
}
