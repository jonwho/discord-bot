package commands

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/BryanSLam/discord-bot/util"
	dg "github.com/bwmarrin/discordgo"
	"github.com/gocolly/colly"
)

type nexterCommand struct {
	regex *regexp.Regexp
}

func newNextErCommand() nexterCommand {
	return nexterCommand{regexp.MustCompile(`(?i)^!nexter(\s[1-9]\d*)?$`)}
}

func (cmd nexterCommand) match(s string) bool {
	return cmd.regex.MatchString(s)
}

func (cmd nexterCommand) fn(s *dg.Session, m *dg.MessageCreate) {
	logger := util.Logger{Session: s, ChannelID: botLogChannelID}
	logger.Info("nexter test")
	slice := strings.Split(m.Content, " ")

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
		s.ChannelMessageSend(m.ChannelID, message)
		return
	}
	message := fmt.Sprintf("No earnings found %d days from now", days)
	s.ChannelMessageSend(m.ChannelID, message)
}

func visit(url string) []struct {
	Ticker  string
	Company string
	EPS     string
	REV     string
} {
	log.Println("Visiting URL", url)

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
