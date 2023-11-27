include .env

MIGRATION_PATH = "./migrations"
MIGRATION_NAME = "migration"

all: start

start: compose-up migrate	
	go install "github.com/serz999/addressBookServer"
	addressBookServer

compose-up:
	docker compose up -d || true
	@sleep 1

migrate: migrate-up
	
migrate-up: migrate-install 
	migrate -database $(DB_URL) -path $(MIGRATION_PATH) up

migrate-down: migrate-install 
	migrate -database $(DB_URL) -path $(MIGRATION_PATH) down

migrate-create: migrate-install 
	migrate create -ext sql -dir $(MIGRATION_PATH) -seq $(MIGRATION_NAME) 

migrate-install: 
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest	
