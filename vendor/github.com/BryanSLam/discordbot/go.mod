module github.com/BryanSLam/discord-bot

go 1.13

require (
	github.com/BryanSLam/discordbot v0.0.0-00010101000000-000000000000
	github.com/alpacahq/alpaca-trade-api-go v1.3.8
	github.com/bwmarrin/discordgo v0.20.1
	github.com/go-redis/redis v6.15.6+incompatible
	github.com/gocolly/colly v1.2.0
	github.com/jonwho/go-iex/v2 v2.7.1
	github.com/robfig/cron v1.2.0
	golang.org/x/text v0.3.2
)

replace github.com/BryanSLam/discordbot => ./
