package earningsreporter

import (
	"bytes"
	"context"
	"log"

	"testing"
)

func TestExecute(t *testing.T) {
	ctx := context.Background()
	buf := bytes.NewBuffer([]byte(""))

	reporter, err := New()
	if err != nil {
		t.Fatal(err)
	}

	err = reporter.Execute(ctx, buf)
	if err != nil {
		t.Fatal(err)
	}

	log.Println(buf)
}
