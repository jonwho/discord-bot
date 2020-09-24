package stock

import (
	"bytes"
	"context"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestNew(t *testing.T) {
	// TODO
}

func TestStockExecute(t *testing.T) {
	buf := bytes.NewBuffer([]byte("$tsla"))

	// NIT: copy minimum .env file to project root instead of relative pathing
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Println("Error loading .env -- check if file exists or is valid. Tests will fail unless you have all available ENV vars")
	}
	iexToken := os.Getenv("IEX_SECRET_TOKEN")
	alpacaID := os.Getenv("ALPACA_KEY_ID")
	alpacaKey := os.Getenv("ALPACA_SECRET_KEY")
	stock, _ := New(iexToken, alpacaID, alpacaKey)
	ctx := context.Background()
	stock.Execute(ctx, buf)

	actual := buf.String()
	if actual == "" {
		t.Error("Expect buf to be written to")
	}
}
