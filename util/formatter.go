package util

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	iex "github.com/jonwho/go-iex"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var (
	pst, _ = time.LoadLocation("America/Los_Angeles")
)

func FormatQuote(quote *iex.Quote) string {
	stringOrder := []string{
		"Symbol",
		"Company Name",
		"Current",
		"High",
		"Low",
		"Open",
		"Close",
		"Change % (1 day)",
		"Delta",
		"Volume",
	}

	var current float32
	var changePercent float32
	var delta float32

	if outsideNormalTradingHours() {
		current = quote.ExtendedPrice
		changePercent = quote.ExtendedChangePercent
		delta = current - quote.Close
	} else {
		current = quote.LatestPrice
		changePercent = quote.ChangePercent
		delta = current - quote.Close
	}

	outputMap := map[string]string{
		"Symbol":           quote.Symbol,
		"Company Name":     quote.CompanyName,
		"Current":          fmt.Sprintf("%#v", current),
		"High":             fmt.Sprintf("%#v", quote.High),
		"Low":              fmt.Sprintf("%#v", quote.Low),
		"Open":             fmt.Sprintf("%#v", quote.Open),
		"Close":            fmt.Sprintf("%#v", quote.Close),
		"Change % (1 day)": fmt.Sprintf("%#v", changePercent) + " %",
		"Delta":            fmt.Sprintf("%#v", Round(float64(delta))),
		"Volume":           fmt.Sprintf("%#v", quote.LatestVolume),
	}

	printer := message.NewPrinter(language.English)
	fmtStr := "```\n"

	for _, k := range stringOrder {
		if k == "Volume" {
			i, _ := strconv.Atoi(outputMap[k])
			fmtStr += printer.Sprintf("%-17s %d\n", k, i)
		} else {
			fmtStr += printer.Sprintf("%-17s %-20s\n", k, outputMap[k])
		}
	}

	fmtStr += "```\n"

	return fmtStr
}

func FormatEarnings(earnings *iex.Earnings) string {
	stringOrder := []string{
		"Symbol",
		"Actual EPS",
		"Estimated EPS",
		"EPS delta",
		"Announce Time",
		"Fiscal Start Date",
		"Fiscal End Date",
		"Fiscal Period",
	}

	if len(earnings.Earnings) == 0 {
		return "No earnings to report for " + earnings.Symbol
	}

	recentEarnings := earnings.Earnings[0]

	outputMap := map[string]string{
		"Symbol":            earnings.Symbol,
		"Actual EPS":        fmt.Sprintf("%#v", recentEarnings.ActualEPS),
		"Estimated EPS":     fmt.Sprintf("%#v", recentEarnings.EstimatedEPS),
		"EPS delta":         fmt.Sprintf("%#v", recentEarnings.EPSSurpriseDollar),
		"Announce Time":     recentEarnings.AnnounceTime,
		"Fiscal Start Date": recentEarnings.FiscalEndDate,
		"Fiscal End Date":   recentEarnings.EPSReportDate,
		"Fiscal Period":     recentEarnings.FiscalPeriod,
		"Year Ago EPS":      fmt.Sprintf("%#v", recentEarnings.YearAgo),
	}

	if strings.ToLower(outputMap["Announce Time"]) == "bto" {
		outputMap["Announce Time"] = "Before Trading Open"
	} else if strings.ToLower(outputMap["Announce Time"]) == "amc" {
		outputMap["Announce Time"] = "After Market Close"
	} else if strings.ToLower(outputMap["Announce Time"]) == "dmt" {
		outputMap["Announce Time"] = "During Market Trading"
	}

	printer := message.NewPrinter(language.English)
	fmtStr := "```\n"

	for _, k := range stringOrder {
		fmtStr += printer.Sprintf("%-17s %-20s\n", k, outputMap[k])
	}

	fmtStr += "```\n"

	return fmtStr
}

func FormatFuzzySymbols(symbols []iex.SymbolDTO) string {
	printer := message.NewPrinter(language.English)
	fmtStr := "```\n"
	fmtStr += "Could not find symbol you requested. Did you mean one of these symbols?\n\n"

	for _, symbol := range symbols {
		fmtStr += printer.Sprintf("%-5s %-20s\n", symbol.Symbol, symbol.Name)
	}
	fmtStr += "```\n"

	return fmtStr
}

func outsideNormalTradingHours() bool {
	now := time.Now().In(pst)

	return now.Hour() >= 13 || (now.Hour() <= 6 && now.Minute() <= 30)
}
