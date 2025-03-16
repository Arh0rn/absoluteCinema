include .env

BINARY_NAME=absoluteCinema
MAIN_FILE=./cmd/app/main.go
MIGRATIONS_DIR=./migrations

DATABASE_URL=$(POSTGRES_DRIVER_NAME)://$(POSTGRES_DB_USER):$(POSTGRES_DB_PASSWORD)@$(POSTGRES_DB_HOST):$(POSTGRES_DB_PORT)/$(POSTGRES_DB_NAME)?sslmode=$(POSTGRES_DB_SSLMODE)

build:
	go build -o $(BINARY_NAME) $(MAIN_FILE)

run: build
	./$(BINARY_NAME)

swag:
	swag init -g ./cmd/app/main.go -o ./docs

install-migrate:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

migrate-up:
	migrate -database $(DATABASE_URL) -path $(MIGRATIONS_DIR) up

migrate-down:
	migrate -database $(DATABASE_URL) -path $(MIGRATIONS_DIR) down

clean:
	del /F /Q $(BINARY_NAME)

rebuild: clean build

print.env:
	type .env
