-include .env

APP_NAME=url-shortener
MIGRATE=migrate -path migrations -database "$(DATABASE_URL)"

.PHONY: run build tidy test migrate-up migrate-down docker-up docker-down

run:
	go run ./cmd/api

build:
	go build -o bin/$(APP_NAME) ./cmd/api

tidy:
	go mod tidy

test:
	go test ./...

migrate-up:
	$(MIGRATE) up

migrate-down:
	$(MIGRATE) down 1

docker-up:
	docker compose up -d --build

docker-down:
	docker compose down
