.PHONY: db-up db-down db-ps db-logs

db-up:
	docker compose up -d postgres

db-down:
	docker compose down

db-ps:
	docker compose ps

db-logs:
	docker compose logs -f postgres
