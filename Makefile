DC=docker-compose

.PHONY: all
all: build test

.PHONY: build
build:
	$(DC) build

.PHONY: up
up:
	$(DC) up

.PHONY: down
down:
	$(DC) down

.PHONY: restart
restart:
	$(DC) restart

.PHONY: test
test:
	$(DC) run bot go test ./...

.PHONY: bash
bash:
	$(DC) run bot bash
