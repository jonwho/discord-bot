package botcommands

import (
	"bytes"
	"context"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestNewStock(t *testing.T) {
	// TODO
}

func TestStockExecute(t *testing.T) {
	buf := bytes.NewBuffer([]byte("$tsla"))

	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalln("Error loading .env -- check if file exists or is valid")
	}
	iexToken := os.Getenv("IEX_SECRET_TOKEN")
	alpacaID := os.Getenv("ALPACA_KEY_ID")
	alpacaKey := os.Getenv("ALPACA_SECRET_KEY")
	stock := NewStock(iexToken, alpacaID, alpacaKey)
	ctx := context.Background()
	stock.Execute(ctx, buf)

	actual := buf.String()
	if actual == "" {
		t.Error("Expect buf to be written to")
	}
}
