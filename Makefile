include .env

BINARY_NAME=absoluteCinema
MAIN_FILE=./cmd/app/main.go
MIGRATIONS_DIR=./migrations

DATABASE_URL=$(DRIVER_NAME)://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

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

#new-migration:
#	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

clean:
	del /F /Q $(BINARY_NAME)

rebuild: clean build

print_env:
	type .env
