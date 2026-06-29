# Codex Handoff

## Current Phase

Phase 2: Backend Skeleton and Runtime Basics.

## Status

Phase 2 is complete. A minimal Go/Echo backend starts locally, loads runtime config, opens a PostgreSQL handle, exposes health/readiness endpoints, and shuts down gracefully.

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

## Changed Files

- `.env.example`
- `Makefile`
- `README.md`
- `CODEX_HANDOFF.md`
- `docs/API_CONTRACT.md`
- `backend/go.mod`
- `backend/go.sum`
- `backend/cmd/server/main.go`
- `backend/internal/config/config.go`
- `backend/internal/db/db.go`
- `backend/internal/http/server.go`
- `backend/internal/http/handlers/errors.go`
- `backend/internal/http/handlers/health.go`
- `backend/internal/http/middleware/cors.go`
- `backend/internal/domain/doc.go`
- `backend/internal/service/doc.go`
- `backend/internal/repository/doc.go`

## Constraints

- Backend is a skeleton only.
- Auth has not been implemented.
- Users and kingdoms have not been implemented.
- Frontend code has not been implemented.
- Database migrations have not been implemented.
- Gameplay systems have not been implemented.
- `.idea/` exists locally as an untracked editor directory and was left untouched.
  It is now ignored by `.gitignore`.

## Verification

- Ran `docker compose config` and confirmed the Compose file resolves with database `sumerki`, user `sumerki`, password `sumerki`, and port `5432`.
- Ran `docker compose up -d postgres`; Docker started pulling `postgres:16-alpine`, but the pull did not complete during the session and was interrupted.
- Ran `docker compose pull postgres`; it also started pulling `postgres:16-alpine`, did not complete during the session, and was interrupted.
- Ran `docker compose ps`; no container was running because the image pull had not completed.
- Reviewed `docs/MVP_PHASES.md` diff and confirmed Phase 1b, migration ordering, lazy resolution rules, resolved resources, earlier ruler and patron phases, split event phases, and incremental UI requirements are present.
- Ran `cd backend && go mod tidy`.
- Ran `GOCACHE=/Users/andrey/Documents/pets/sumerki/.cache/go-build go test ./...`.
- Ran backend without `DATABASE_URL`; startup failed with `config error: DATABASE_URL is required`.
- Ran backend with `DATABASE_URL=postgres://sumerki:sumerki@localhost:5432/sumerki?sslmode=disable`, `JWT_SECRET=dev-secret`, and `BACKEND_PORT=18080`; sandboxed port binding was denied, then the server started successfully with elevated local port binding.
- Ran `curl -i http://localhost:18080/health`; received HTTP 200 with `{"status":"ok"}`.
- Ran `curl -i http://localhost:18080/ready` without a reachable database; received HTTP 503 with standard `database_unavailable` JSON error.
- Ran `docker compose up -d postgres`; image pull completed, but default port `5432` could not be bound because another Docker process was already listening on that port.
- Ran `POSTGRES_PORT=15432 docker compose up -d postgres` for temporary local verification.
- Ran backend with `DATABASE_URL=postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable`; `GET /ready` returned HTTP 200 with `{"status":"ready","database":"ok"}`.
- Ran `POSTGRES_PORT=15432 docker compose down`.
- Ran final `docker compose ps` and confirmed no Sumerki Compose service was left running.

## What Works Now

- Local PostgreSQL can be started with `docker compose up -d postgres` or `make db-up`.
- Local PostgreSQL uses database `sumerki`, user `sumerki`, password `sumerki`, and host port `5432`.
- Backend can be started with `cd backend && DATABASE_URL="postgres://sumerki:sumerki@localhost:5432/sumerki?sslmode=disable" go run ./cmd/server` when PostgreSQL is reachable.
- `GET /health` reports process health without requiring database connectivity.
- `GET /ready` reports database readiness and returns a standard JSON 503 error when PostgreSQL is unavailable.
- Backend supports request logging, panic recovery, local frontend CORS for `http://localhost:5173`, and graceful shutdown on SIGINT/SIGTERM.

## Known Limitations

- No auth, user, kingdom, migration, frontend, or gameplay code exists yet.
- No migrations or schema exist yet.
- `JWT_SECRET` is loaded but intentionally unused until the auth phase.
- Default PostgreSQL port `5432` was occupied by another Docker process during verification, so the successful readiness check used temporary `POSTGRES_PORT=15432` without changing repository files.

## Next Recommended Step

Start Phase 3: Database Migration Foundation.
