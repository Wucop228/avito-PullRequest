-include .env
.PHONY: up down logs run migrate-up migrate-down migrate-steps migrate-goto migrate-version migrate-force

DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)
MIGRATIONS_PATH=./migrations

migrate-up:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" up

migrate-down:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down 1

migrate-steps:
	@if [ -z "$(n)" ]; then echo "Usage: make migrate-steps n=-2 (negative = rollback)"; exit 1; fi
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" steps $(n)

migrate-goto:
	@if [ -z "$(v)" ]; then echo "Usage: make migrate-goto v=3"; exit 1; fi
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" goto $(v)

migrate-version:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" version

migrate-force:
	@if [ -z "$(v)" ]; then echo "Usage: make migrate-force v=3"; exit 1; fi
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" force $(v)

run:
	go run ./cmd/server

up:
	docker compose up -d

down:
	docker compose down

logs:
	docker compose logs -f app