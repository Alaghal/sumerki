DATABASE_URL ?= postgres://sumerki:sumerki@localhost:5432/sumerki?sslmode=disable
JWT_SECRET ?= dev-secret
BACKEND_PORT ?= 8080
GOCACHE ?= $(CURDIR)/.cache/go-build
GOOSE := GOCACHE="$(GOCACHE)" go run github.com/pressly/goose/v3/cmd/goose@latest

.PHONY: db-up db-down db-ps db-logs backend-run backend-tidy seed-dev migrate-up migrate-down migrate-status migrate-reset reset-db smoke-api test-backend test-frontend test-all playtest-setup playtest-check playtest-reset frontend-install frontend-dev

db-up:
	docker compose up -d postgres

db-down:
	docker compose down

db-ps:
	docker compose ps

db-logs:
	docker compose logs -f postgres

backend-run:
	cd backend && DATABASE_URL="$(DATABASE_URL)" JWT_SECRET="$(JWT_SECRET)" BACKEND_PORT="$(BACKEND_PORT)" GOCACHE="$(GOCACHE)" go run ./cmd/server

backend-tidy:
	cd backend && go mod tidy

seed-dev:
	cd backend && DATABASE_URL="$(DATABASE_URL)" GOCACHE="$(GOCACHE)" go run ./cmd/seed-dev

migrate-up:
	cd backend && $(GOOSE) -dir migrations postgres "$(DATABASE_URL)" up

migrate-down:
	cd backend && $(GOOSE) -dir migrations postgres "$(DATABASE_URL)" down

migrate-status:
	cd backend && $(GOOSE) -dir migrations postgres "$(DATABASE_URL)" status

migrate-reset:
	cd backend && $(GOOSE) -dir migrations postgres "$(DATABASE_URL)" reset

reset-db:
	@echo "WARNING: reset-db is destructive and intended for local Docker development only."
	docker compose exec -T postgres dropdb -U sumerki --if-exists sumerki
	docker compose exec -T postgres createdb -U sumerki sumerki

smoke-api:
	API_BASE_URL="$(API_BASE_URL)" scripts/smoke-api.sh

test-backend:
	cd backend && GOCACHE="$(GOCACHE)" go test ./...

test-frontend:
	cd frontend && npm run typecheck
	cd frontend && npm run build

test-all: test-backend test-frontend

playtest-setup: migrate-up seed-dev

playtest-check: test-all
	@echo "Start the backend separately, then run: make smoke-api"

playtest-reset:
	@echo "WARNING: playtest-reset is destructive and intended for local Docker development only."
	$(MAKE) reset-db
	$(MAKE) playtest-setup

frontend-install:
	cd frontend && npm install

frontend-dev:
	cd frontend && npm run dev
