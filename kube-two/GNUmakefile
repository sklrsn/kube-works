.DEFAULT_GOAL:all

.PHONY: all
all: build up

.PHONY: build
build:
	docker compose -f docker-compose.yaml build

.PHONY: up
up:
	docker compose -f docker-compose.yaml up

.PHONY: down
down:
	docker compose -f docker-compose.yaml down
