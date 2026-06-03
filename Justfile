set dotenv-load

migrations_dir := "./migrations"

export GOOSE_DRIVER := "postgres"
export GOOSE_DBSTRING := env("POSTGRES_CONNECTION")
export GOOSE_MIGRATION_DIR := migrations_dir

goose := "goose"

dev:
    @go run ./cmd/*.go

build:
    @go build -o bin/app ./cmd

live:
    @watchexec --restart -- go run ./cmd/*.go

db-up:
    @{{ goose }} up

[confirm("Rollback migrations down by one version? (y/N): ")]
db-down:
	@{{ goose }} down

[confirm("Rollback migrations down to the specified version? (y/N): ")]
db-down-to version:
	{{ goose }} down-to {{version }}

db-status:
	@{{ goose }} status

db-create name:
	@{{ goose }} -s create {{ name }} sql

swag-gen:
	@swag init -g main.go -d cmd,internal
