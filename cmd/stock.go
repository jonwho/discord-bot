package cmd

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/BryanSLam/discord-bot/util"
	"net/http"
)

const dataURL string = "https://data.alpaca.markets/v1/bars/1Min"

// NewStockCommand TODO: @doc
func NewStockCommand() *Command {
	return &Command{
		Match: func(s string) bool {
			return regexp.MustCompile(`(?i)^!stock [\w.]+$`).MatchString(s)
		},
		Fn: Stock,
	}
}

// Stock TODO: @doc
func Stock(rw io.ReadWriter, logger *util.Logger, m map[string]interface{}) {
	buf, err := ioutil.ReadAll(rw)
	if err != nil {
		rw.Write([]byte(err.Error()))
		return
	}

	slice := strings.Split(string(buf), " ")
	ticker := strings.ToUpper(slice[1])

	logger.Info("Fetching stock info for " + ticker)
	req, err := http.NewRequest(http.MethodGet, dataURL, nil)
	if err != nil {
		rw.Write([]byte(err.Error()))
		return
	}
	req.Header.Set("APCA-API-KEY-ID", alpacaID)
	req.Header.Set("APCA-API-SECRET-KEY", alpacaKey)
	q := req.URL.Query()
	q.Add("limit", "1")
	q.Add("symbols", ticker)
	req.URL.RawQuery = q.Encode()
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		rw.Write([]byte(err.Error()))
		return
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		rw.Write([]byte(err.Error()))
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
		rw.Write([]byte(err.Error()))
		return
	}
	if len(data[ticker]) == 0 {
		rw.Write([]byte("No data found for " + ticker))
		return
	}

	bar := data[ticker][0]
	message := util.FormatBar(bar, ticker)

	rw.Write([]byte(message))
}
