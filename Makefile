DC=docker-compose

all: build test

build:
	$(DC) build

up:
	$(DC) up

down:
	$(DC) down

test:
	$(DC) run bot go test ./...

bash:
	$(DC) run bot bash

.PHONY: all up down test
