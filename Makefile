include .env

MIGRATION_PATH = "./migrations"
MIGRATION_NAME = "migration"
GOPATH = $(HOME)/go
GOBIN = $(GOPATH)/bin

all: start

start: compose-up migrate	
	go install
	$(GOBIN)/addressBookServer

compose-up:
	docker compose up -d
	@sleep 1

migrate: migrate-up
	
migrate-up: migrate-install 
	$(GOBIN)/migrate -database $(DB_URL) -path $(MIGRATION_PATH) up

migrate-down: migrate-install 
	$(GOBIN)/migrate -database $(DB_URL) -path $(MIGRATION_PATH) down

migrate-create: migrate-install 
	$(GOBIN)/migrate create -ext sql -dir $(MIGRATION_PATH) -seq $(MIGRATION_NAME) 

migrate-install: 
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
