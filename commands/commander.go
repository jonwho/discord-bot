package commands

import (
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/BryanSLam/discord-bot/config"
	"github.com/BryanSLam/discord-bot/util"
	dg "github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis"
	"github.com/robfig/cron"
)

const (
	dateFormat        string = "1/_2/06"
	redisDateFormat   string = "01/02/06"
	watchlistRedisKey string = "watchlist"
)

type work func(s *dg.Session, m *dg.MessageCreate)

var (
	token       string
	redisClient *redis.Client
	cronner     *cron.Cron
	pst, _      = time.LoadLocation("America/Los_Angeles")
)

var (
	commandRegex        = regexp.MustCompile(`(?i)^![\w]+[\w ".]*[ 0-9/]*$`)
	pingRegex           = regexp.MustCompile(`(?i)^!ping$`)
	stockRegex          = regexp.MustCompile(`(?i)^!stock [\w.]+$`)
	erRegex             = regexp.MustCompile(`(?i)^!er [\w.]+$`)
	wizdaddyRegex       = regexp.MustCompile(`(?i)^!wizdaddy$`)
	coinRegex           = regexp.MustCompile(`(?i)^!coin [\w]+$`)
	remindmeRegex       = regexp.MustCompile(`(?i)^!remindme [\w ]+ (0?[1-9]|1[012])/(0?[1-9]|[12][0-9]|3[01])/(\d\d)$`)
	watchlistRegex      = regexp.MustCompile(`(?i)^!watchlist [\w ]+$`)
	clearwatchlistRegex = regexp.MustCompile(`(?i)^!clearwatchlist$`)
	newsRegex           = regexp.MustCompile(`(?i)^!news [\w.]+$`)
	nexterRegex         = regexp.MustCompile(`(?i)^!nexter(\s[1-9]\d*)?$`)
)

var (
	coinAPIURL            = config.GetConfig().CoinAPIURL
	wizdaddyURL           = config.GetConfig().WizdaddyURL
	invalidCommandMessage = config.GetConfig().InvalidCommandMessage
	botLogChannelID       = config.GetConfig().BotLogChannelID
	earningsWhisperURL    = config.GetConfig().EarningsWhisperURL
)

func init() {
	token = os.Getenv("BOT_TOKEN")

	redisClient = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	cronner = cron.NewWithLocation(pst)
	cronner.Start()
}

func Commander() func(s *dg.Session, m *dg.MessageCreate) {
	return func(s *dg.Session, m *dg.MessageCreate) {

		if commandRegex.MatchString(m.Content) {
			// Ignore all messages created by the bot itself
			// This isn't required in this specific example but it's a good practice.
			if m.Author.ID == s.State.User.ID {
				return
			}

			if pingRegex.MatchString(m.Content) {
				go safelyDo(Ping, s, m)
				return
			}

			if stockRegex.MatchString(m.Content) {
				go safelyDo(Stock, s, m)
				return
			}

			if erRegex.MatchString(m.Content) {
				go safelyDo(Er, s, m)
				return
			}

			if wizdaddyRegex.MatchString(m.Content) {
				go safelyDo(Wizdaddy, s, m)
				return
			}

			if coinRegex.MatchString(m.Content) {
				go safelyDo(Coin, s, m)
				return
			}

			if remindmeRegex.MatchString(m.Content) {
				go safelyDo(Remindme, s, m)
				return
			}

			if watchlistRegex.MatchString(m.Content) {
				go safelyDo(Watchlist, s, m)
				return
			}

			if clearwatchlistRegex.MatchString(m.Content) {
				go safelyDo(ClearWatchlist, s, m)
				return
			}

			if newsRegex.MatchString(m.Content) {
				go safelyDo(News, s, m)
				return
			}

			if nexterRegex.MatchString(m.Content) {
				go safelyDo(NextEr, s, m)
				return
			}

			s.ChannelMessageSend(m.ChannelID, invalidCommandMessage)
		}
	}
}

func safelyDo(fn work, s *dg.Session, m *dg.MessageCreate) {
	logger := util.Logger{Session: s, ChannelID: botLogChannelID}

	// defer'd funcs will execute before return even if runtime panic
	defer func() {
		// use recover to catch panic so bot doesn't shutdown
		if err := recover(); err != nil {
			logger.Send(util.MentionMaintainers() + " an error has occurred")
			logger.Warn(fmt.Sprintln("function", util.FuncName(fn), "failed:", err))
		}
	}()
	// perform work
	fn(s, m)
}
