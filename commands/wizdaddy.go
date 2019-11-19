package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/BryanSLam/discord-bot/datasource"
	"github.com/BryanSLam/discord-bot/util"
)

// NewWizdaddyCommand TODO: @doc
func NewWizdaddyCommand() *Command {
	return &Command{
		Match: func(s string) bool {
			return regexp.MustCompile(`(?i)^!wizdaddy$`).MatchString(s)
		},
		Fn: Wizdaddy,
	}
}

// Wizdaddy TODO: @doc
func Wizdaddy(rw io.ReadWriter, logger *util.Logger, _ map[string]interface{}) {
	resp, err := http.Get(wizdaddyURL)

	if err != nil {
		logger.Trace("Wizdaddy request failed. Message: " + err.Error())
		rw.Write([]byte("Daddy is down"))
		return
	}

	var daddyResponse datasource.WizdaddyResponse
	if err = json.NewDecoder(resp.Body).Decode(&daddyResponse); err != nil {
		logger.Trace("JSON decoding failed. Message: " + err.Error())
		rw.Write([]byte(err.Error()))
		return
	}

	rw.Write([]byte(fmt.Sprintf("%s %s %s %s", daddyResponse.Symbol,
		daddyResponse.StrikePrice, daddyResponse.ExpirationDate, daddyResponse.Type)))
}
