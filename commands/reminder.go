package commands

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/BryanSLam/discord-bot/config"
	"github.com/BryanSLam/discord-bot/util"
	dg "github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis"
	"github.com/robfig/cron"
)

/*
 * Set up redis table in a way where the keys are the dates of when to pull a reminder
 * and the value is an array of events that need to be reminded on that date
 */

const (
	dateFormat      string = "1/_2/06"
	redisDateFormat string = "01/02/06"
)

var (
	token       string
	redisClient *redis.Client
	cronner     *cron.Cron
	pst, _      = time.LoadLocation("America/Los_Angeles")
)

func init() {
	token = os.Getenv("BOT_TOKEN")

	redisClient = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	cronner = cron.New()
	cronner.Start()
	// Run once at 6:00 AM from Monday-Friday
	cronner.AddFunc("0 6 * * 1-5", todaysReminders)
}

// Remindme creates a reminder entry into datastore (Redis)
func Remindme(s *dg.Session, m *dg.MessageCreate) {
	logger := util.Logger{Session: s, ChannelID: config.GetConfig().BotLogChannelID}

	slice := strings.Split(m.Content, " ")
	date := slice[len(slice)-1]
	// grab message string in between command and date
	msgSlice := slice[1 : len(slice)-1]
	msg := strings.Join(msgSlice, " ")
	reminder := fmt.Sprintf("%s~*REMINDER <@%s>: %s", m.ChannelID, m.Author.Mention(), msg)

	err := addReminder(reminder, date)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		logger.Trace("Remindme request failed: " + err.Error())
		return
	}

	s.ChannelMessageSend(m.ChannelID, "Reminder set for "+date+" at 6:00 AM")
}

// addReminder will function as a stack; pushing new reminder messages for the date key.
// The date should be in the form mm/dd/yy.
func addReminder(message string, date string) error {
	dateCheck, _ := time.Parse(dateFormat, date)
	future := time.Date(dateCheck.Year(), dateCheck.Month(), dateCheck.Day(), 0, 0, 0, 0, pst)
	if future.Before(time.Now().In(pst)) {
		return errors.New("Date has already passed")
	}

	redisClient.RPush(dateCheck.Format(redisDateFormat), message)
	return nil
}

// getReminders will fetch reminder messages datastore for the given date
func getReminders(date string) ([]string, error) {
	_, err := redisClient.LRange(date, 0, -1).Result()
	if err == redis.Nil {
		return nil, errors.New("no reminders for date: " + date)
	}
	reminders := redisClient.LRange(date, 0, -1).Val()

	redisClient.Del(date)
	return reminders, nil
}

// Function run during the daily reminder check
func todaysReminders() {
	dgSession, _ := dg.New("Bot " + token)
	defer dgSession.Close()

	todaysDate := time.Now().In(pst).Format(redisDateFormat)
	reminders, _ := getReminders(todaysDate)
	dgSession.Open()

	for _, reminder := range reminders {
		split := strings.Split(reminder, "~*")
		channel := split[0]
		dgSession.ChannelMessageSend(channel, split[1])
	}
}
