DATABASE_URL ?= postgres://sumerki:sumerki@localhost:5432/sumerki?sslmode=disable
JWT_SECRET ?= dev-secret
BACKEND_PORT ?= 8080

.PHONY: db-up db-down db-ps db-logs backend-run backend-tidy

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
