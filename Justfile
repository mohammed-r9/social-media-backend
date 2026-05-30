set dotenv-load
migrations := "./migrations"

dev:
	@go run ./cmd/*.go

build:
	@go build -o bin/app ./cmd

live:
	@watchexec --restart -- go run ./cmd/*.go

migrate action:
	@goose -dir {{migrations}} postgres "$POSTGRES_CONNECTION" {{action}}
