DATABASE_URL ?= postgres://sumerki:sumerki@localhost:5432/sumerki?sslmode=disable
JWT_SECRET ?= dev-secret
BACKEND_PORT ?= 8080
GOOSE := go run github.com/pressly/goose/v3/cmd/goose@latest

.PHONY: db-up db-down db-ps db-logs backend-run backend-tidy migrate-up migrate-down migrate-status migrate-reset frontend-install frontend-dev

db-up:
	docker compose up -d postgres

db-down:
	docker compose down

db-ps:
	docker compose ps

db-logs:
	docker compose logs -f postgres

backend-run:
	cd backend && DATABASE_URL="$(DATABASE_URL)" JWT_SECRET="$(JWT_SECRET)" BACKEND_PORT="$(BACKEND_PORT)" go run ./cmd/server

backend-tidy:
	cd backend && go mod tidy

migrate-up:
	cd backend && $(GOOSE) -dir migrations postgres "$(DATABASE_URL)" up

migrate-down:
	cd backend && $(GOOSE) -dir migrations postgres "$(DATABASE_URL)" down

migrate-status:
	cd backend && $(GOOSE) -dir migrations postgres "$(DATABASE_URL)" status

migrate-reset:
	cd backend && $(GOOSE) -dir migrations postgres "$(DATABASE_URL)" reset

frontend-install:
	cd frontend && npm install

frontend-dev:
	cd frontend && npm run dev
