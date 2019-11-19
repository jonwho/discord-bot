package commands

import (
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"

	"github.com/BryanSLam/discord-bot/util"
	"github.com/gocolly/colly"
)

// NewNextErCommand TODO: @doc
func NewNextErCommand() *Command {
	return &Command{
		Match: func(s string) bool {
			return regexp.MustCompile(`(?i)^!nexter(\s[1-9]\d*)?$`).MatchString(s)
		},
		Fn: Nexter,
	}
}

// Nexter TODO: @doc
func Nexter(rw io.ReadWriter, logger *util.Logger, _ map[string]interface{}) {
	buf, err := ioutil.ReadAll(rw)
	if err != nil {
		rw.Write([]byte(err.Error()))
		return
	}

	slice := strings.Split(string(buf), " ")

	days := 1
	if len(slice) > 1 {
		var err error
		days, err = strconv.Atoi(slice[1])
		if err != nil {
			panic(err)
		}
	}

	url := fmt.Sprintf("%s&d=%d", earningsWhisperURL, days)
	upcomingEarnings := visit(url)
	if len(upcomingEarnings) > 0 {
		message := util.FormatUpcomingErs(upcomingEarnings)
		rw.Write([]byte(message))
		return
	}
	message := fmt.Sprintf("No earnings found %d days from now", days)
	rw.Write([]byte(message))
}

func visit(url string) []struct {
	Ticker  string
	Company string
	EPS     string
	REV     string
} {

	ers := []struct {
		Ticker  string
		Company string
		EPS     string
		REV     string
	}{}

	c := colly.NewCollector()

	c.OnHTML("li.cor.bmo", func(e *colly.HTMLElement) {
		er := struct {
			Ticker  string
			Company string
			EPS     string
			REV     string
		}{
			e.ChildText("div.ticker"),
			e.ChildText("div.company"),
			e.ChildText("div.estimate"),
			e.ChildText("div.revest"),
		}
		ers = append(ers, er)
	})

	c.Visit(url)

	return ers
}
