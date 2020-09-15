package earningsreporter

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"time"
	// "log"

	"github.com/gocolly/colly"
)

const (
	earningsWhisperURL string = "https://www.earningswhispers.com/calendar?sb=c&t=all"
	dateFormat         string = "1/_2/06"
)

var (
	pst, _ = time.LoadLocation("America/Los_Angeles")
)

// EarningsReporter finds companies that will soon report earnings.
type EarningsReporter struct{}

type report struct {
	Ticker  string
	Company string
}

// New returns pointer to EarningsReporter.
func New() (*EarningsReporter, error) {
	er := &EarningsReporter{}
	return er, nil
}

// Execute perform task of finding stock earnings. Assumption is that this
// is ran on a Monday.
func (er *EarningsReporter) Execute(ctx context.Context, rw io.ReadWriter) error {
	buf, err := ioutil.ReadAll(rw)
	if err != nil {
		return err
	}
	if len(buf) < 1 {
		// TODO: maybe do something with input in the future
	}

	day := 24 * time.Hour
	mon := time.Now().In(pst)
	tue := mon.Add(day)
	wed := mon.Add(2 * day)
	thu := mon.Add(3 * day)
	fri := mon.Add(4 * day)

	days := []string{
		mon.Format(dateFormat),
		tue.Format(dateFormat),
		wed.Format(dateFormat),
		thu.Format(dateFormat),
		fri.Format(dateFormat),
	}

	weeklyReports := map[string][]report{
		mon.Format(dateFormat): []report{},
		tue.Format(dateFormat): []report{},
		wed.Format(dateFormat): []report{},
		thu.Format(dateFormat): []report{},
		fri.Format(dateFormat): []report{},
	}

	urls := []string{}
	for i := 0; i < 4; i++ {
		url := fmt.Sprintf("%s&d=%d", earningsWhisperURL, i)
		urls = append(urls, url)
	}

	for idx, url := range urls {
		dayKey := days[idx]
		weeklyReports[dayKey] = visit(url)
	}

	reportCount := 0
	for _, dayReports := range weeklyReports {
		reportCount += len(dayReports)
	}

	weeklyReportsFormatted := []string{}
	for dayKey, dayReports := range weeklyReports {
		weeklyReportsFormatted = append(weeklyReportsFormatted, formatDailyReport(dayKey, dayReports))
	}

	allReportsFormatted := strings.Join(weeklyReportsFormatted, "\n\n")
	outputReport := fmt.Sprintf("```\n%s\n```\n", allReportsFormatted)
	rw.Write([]byte(outputReport))

	return nil
}

func visit(url string) []report {
	ers := []report{}
	c := colly.NewCollector()

	c.OnHTML("[id*='T-']", func(e *colly.HTMLElement) {
		er := report{
			Ticker:  e.ChildText("div.ticker"),
			Company: e.ChildText("div.company"),
		}
		ers = append(ers, er)
	})

	c.OnRequest(func(r *colly.Request) {
		// debug before each request
		// log.Println("Visiting", r.URL.String())
	})

	c.Visit(url)

	return ers
}

func formatReport(rp report) string {
	return fmt.Sprintf("Ticker: %s\nCompany: %s\n", rp.Ticker, rp.Company)
}

func formatDailyReport(dayKey string, reports []report) string {
	formattedReports := []string{}

	var outputStr string
	if len(reports) == 0 {
		outputStr = "No earnings to report"
	} else {
		for _, rp := range reports {
			formattedReports = append(formattedReports, formatReport(rp))
		}
		outputStr = strings.Join(formattedReports, "\n")
	}

	return fmt.Sprintf("%s\n%s\n", dayKey, outputStr)
}
