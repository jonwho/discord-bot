DC=docker-compose

.PHONY: all
all: build test

.PHONY: build
build:
	$(DC) build

.PHONY: build-discord-bot
build-discord-bot:
	go build -o cmd/discord-bot/discord-bot cmd/discord-bot/main.go

.PHONY: build-discord-bot-linux
build-discord-bot-linux:
	GOOS=linux GOARCH=amd64 go build -o cmd/discord-bot/discord-bot-linux cmd/discord-bot/main.go

.PHONY: run-discord-bot
run-discord-bot:
	./cmd/discord-bot/discord-bot

.PHONY: up
up:
	$(DC) up

.PHONY: down
down:
	$(DC) down

.PHONY: restart
restart:
	$(DC) restart

.PHONY: test-in-docker
test-in-docker:
	$(DC) run bot go test ./...

.PHONY: test
test:
	go test ./...

.PHONY: bash
bash:
	$(DC) run bot bash

.PHONY: music
music:
	$(DC) run musicbot
