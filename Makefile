.PHONY: up down restart migrate-up migrate-down test test-coverage test-coverage-html swag help

# Load variables from .env if it exists
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

help:
	@echo "Important Commands:"
	@echo "  make up           - Start application (Docker)"
	@echo "  make down         - Stop application"
	@echo "  make restart      - Restart application"
	@echo "  make migrate-up   - Apply database migrations"
	@echo "  make migrate-down - Rollback migrations"
	@echo "  make test           - Run all tests"
	@echo "  make test-coverage  - Run tests and show coverage report"
	@echo "  make test-coverage-html - Run tests and open HTML coverage report"
	@echo "  make swag           - Generate Swagger documentation"

up:
	docker-compose up --build -d

down:
	docker-compose down

restart: down up

migrate-up:
	docker-compose run --rm migrate -path=/migrations/ -database="postgres://$(DB_USER):$(DB_PASSWORD)@db:5432/$(DB_NAME)?sslmode=$(DB_SSLMODE)" up

migrate-down:
	docker-compose run --rm migrate -path=/migrations/ -database="postgres://$(DB_USER):$(DB_PASSWORD)@db:5432/$(DB_NAME)?sslmode=$(DB_SSLMODE)" down 1

test:
	go test -v ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

test-coverage-html: test-coverage
	go tool cover -html=coverage.out

swag:
	$(shell go env GOPATH)/bin/swag init -g cmd/main.go -o docs
