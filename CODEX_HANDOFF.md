# Codex Handoff

## Current Phase

Phase 4: Auth API.

## Status

Phase 4 is complete. The backend now supports user registration, login, JWT authentication middleware, and the current-user endpoint.

## Completed

- Created initial repository guidance in `AGENTS.md`.
- Added `.gitignore` for future Go and TypeScript development artifacts.
- Created project overview in `README.md`.
- Documented MVP scope in `docs/MVP_SCOPE.md`.
- Preserved and aligned phase plan in `docs/MVP_PHASES.md`.
- Created decision log in `docs/DECISIONS.md`.
- Drafted initial API contract in `docs/API_CONTRACT.md`.
- Drafted initial domain model in `docs/DOMAIN_MODEL.md`.
- Created phase documentation directory at `docs/phases/`.
- Added Phase 0 detail page at `docs/phases/phase-00-repository-bootstrap.md`.
- Added Docker Compose configuration for local PostgreSQL.
- Added `.env.example` with local PostgreSQL defaults.
- Added basic Makefile commands for local database management.
- Updated `README.md` with local database start, status, and stop commands.
- Updated `docs/MVP_PHASES.md` so the plan is easier to execute phase-by-phase and produces playable vertical slices earlier.
- Added minimal Go backend skeleton under `backend/`.
- Added config loading for `DATABASE_URL`, `JWT_SECRET`, and `BACKEND_PORT`.
- Added PostgreSQL connection helper using `database/sql` and pgx stdlib.
- Added Echo server setup with request logging, recover middleware, local CORS, and standard JSON errors.
- Added `GET /health` and `GET /ready`.
- Added backend run/tidy Makefile commands.
- Updated backend runtime instructions in `README.md`.
- Updated `docs/API_CONTRACT.md` with readiness endpoint documentation.
- Added Goose SQL migrations for `users` and `kingdoms`.
- Added Makefile migration commands for up, down, status, and local reset.
- Updated README migration instructions.
- Updated domain model field documentation for `updatedAt`, case-insensitive email uniqueness, password hash checks, and nullable patron validation.
- Added user domain model.
- Added user repository for creating users and finding users by normalized email or id.
- Added auth service for email normalization, password hashing, login, JWT generation, JWT validation, and current-user lookup.
- Added auth HTTP handlers for `POST /api/auth/register`, `POST /api/auth/login`, and `GET /api/me`.
- Added auth middleware for Bearer token extraction and JWT validation.
- Made `JWT_SECRET` required at backend startup.
- Added focused auth service and middleware tests.
- Updated README auth curl examples.
- Updated `docs/API_CONTRACT.md` with Phase 4 auth response and error details.

## Changed Files

- `README.md`
- `CODEX_HANDOFF.md`
- `docs/API_CONTRACT.md`
- `backend/go.mod`
- `backend/go.sum`
- `backend/cmd/server/main.go`
- `backend/internal/config/config.go`
- `backend/internal/domain/user.go`
- `backend/internal/repository/user_repository.go`
- `backend/internal/service/auth_service.go`
- `backend/internal/service/auth_service_test.go`
- `backend/internal/http/server.go`
- `backend/internal/http/apierror/apierror.go`
- `backend/internal/http/handlers/errors.go`
- `backend/internal/http/handlers/auth_handler.go`
- `backend/internal/http/handlers/me_handler.go`
- `backend/internal/http/middleware/auth_middleware.go`
- `backend/internal/http/middleware/auth_middleware_test.go`

## Constraints

- Auth API has been implemented.
- Users now have registration, login, and current-user API support.
- Kingdoms still have database tables only; no kingdom application APIs or repositories have been implemented.
- Frontend code has not been implemented.
- Gameplay systems have not been implemented.
- `.idea/` exists locally as an untracked editor directory and was left untouched.
  It is now ignored by `.gitignore`.

## Verification

- Ran `go mod tidy`.
- Ran `GOCACHE=/Users/andrey/Documents/pets/sumerki/.cache/go-build go test ./...`.
- Ran backend without `JWT_SECRET`; startup failed with `config error: JWT_SECRET is required`.
- Ran `POSTGRES_PORT=15432 docker compose up -d postgres` for Phase 4 verification because default port `5432` was still occupied by another Docker process.
- Ran `DATABASE_URL=postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable make migrate-up`; database was already at migration version 2.
- Ran backend with `DATABASE_URL=postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable`, `JWT_SECRET=test-secret`, and `BACKEND_PORT=18080`.
- Ran `curl -i http://localhost:18080/ready`; received HTTP 200 with `{"status":"ready","database":"ok"}`.
- Ran `POST /api/auth/register`; received HTTP 201 with normalized user email and JWT token, without `password_hash`.
- Re-ran register with the same email in different casing; received HTTP 409 with `email_already_exists`.
- Ran `POST /api/auth/login` with valid credentials and mixed-case email; received HTTP 200 with normalized user email and JWT token.
- Ran `POST /api/auth/login` with invalid credentials; received HTTP 401 with `invalid_credentials`.
- Ran `GET /api/me` with a valid Bearer token; received HTTP 200 with the current user.
- Ran `GET /api/me` without Authorization; received HTTP 401 with `missing_authorization_header`.
- Ran `GET /api/me` with an invalid token; received HTTP 401 with `invalid_token`.
- Covered expired token rejection with middleware unit test.
- Removed the temporary verification user from PostgreSQL.
- Stopped the backend with SIGINT and ran `POSTGRES_PORT=15432 docker compose down`.
- Ran final `docker compose ps` and confirmed no Sumerki Compose service was left running.

## What Works Now

- Local PostgreSQL can be started with `docker compose up -d postgres` or `make db-up`.
- Local PostgreSQL uses database `sumerki`, user `sumerki`, password `sumerki`, and host port `5432`.
- Backend can be started with `cd backend && DATABASE_URL="postgres://sumerki:sumerki@localhost:5432/sumerki?sslmode=disable" JWT_SECRET="dev-secret" go run ./cmd/server` when PostgreSQL is reachable.
- `GET /health` reports process health without requiring database connectivity.
- `GET /ready` reports database readiness and returns a standard JSON 503 error when PostgreSQL is unavailable.
- Backend supports request logging, panic recovery, local frontend CORS for `http://localhost:5173`, and graceful shutdown on SIGINT/SIGTERM.
- `make migrate-up` applies initial database schema migrations.
- `make migrate-status` shows Goose migration status.
- `make migrate-down` rolls back the latest migration.
- `make migrate-reset` resets local development migrations.
- Database constraints enforce case-insensitive user email uniqueness, non-empty email/password hash, one kingdom per user, valid culture, and valid nullable patron.
- `POST /api/auth/register` creates a user with a bcrypt password hash and returns a JWT.
- `POST /api/auth/login` authenticates normalized email and password and returns a JWT.
- `GET /api/me` returns the current user when a valid Bearer token is provided.
- Auth responses never expose `password_hash`.
- Auth middleware rejects missing, malformed, invalid, expired, and user-not-found tokens with standard JSON errors.

## Known Limitations

- No kingdom API, frontend, or gameplay code exists yet.
- Default PostgreSQL port `5432` was occupied by another Docker process during verification, so successful auth checks used temporary `POSTGRES_PORT=15432` without changing repository files.
- Goose `@latest` currently switches to Go toolchain `go1.25.11` during command execution because the latest Goose release requires Go >= 1.25.7.

## Next Recommended Step

Start Phase 5: Kingdom Creation API.
