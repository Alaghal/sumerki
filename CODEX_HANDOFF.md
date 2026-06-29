# Codex Handoff

## Current Phase

Phase 5: Kingdom Creation API.

## Status

Phase 5 is complete. Authenticated users can create exactly one kingdom and fetch their current kingdom.

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
- Added user domain model, user repository, auth service, auth handlers, and auth middleware.
- Made `JWT_SECRET` required at backend startup.
- Added focused auth service and middleware tests.
- Added kingdom domain model.
- Added kingdom repository for creating a kingdom and finding a kingdom by user id.
- Added kingdom service with name trimming, name length validation, culture validation, and one-kingdom error mapping.
- Added authenticated `POST /api/kingdoms` and `GET /api/kingdoms/me`.
- Added focused kingdom service tests.
- Updated README with kingdom curl examples and fixed the backend run example to include `JWT_SECRET`.
- Updated `docs/API_CONTRACT.md` with exact kingdom API responses and errors.

## Changed Files

- `README.md`
- `CODEX_HANDOFF.md`
- `docs/API_CONTRACT.md`
- `backend/internal/http/server.go`
- `backend/internal/domain/kingdom.go`
- `backend/internal/repository/kingdom_repository.go`
- `backend/internal/service/kingdom_service.go`
- `backend/internal/service/kingdom_service_test.go`
- `backend/internal/http/handlers/kingdom_handler.go`

Phase 4 auth files are still present as untracked working-tree files in this checkout and were preserved while implementing Phase 5.

## Constraints

- Auth API has been implemented.
- Kingdom creation and current-kingdom lookup have been implemented.
- Users can have only one kingdom.
- Kingdom creation does not create resources, buildings, rulers, patron relations, armies, missions, events, tribute, or gameplay systems.
- Client-provided `user_id` and `patron` are not accepted by the create-kingdom request shape; authenticated user id comes from JWT context and patron remains `null`.
- Frontend code has not been implemented.
- `.idea/` exists locally as an untracked editor directory and was left untouched.
  It is now ignored by `.gitignore`.

## Verification

- Ran `GOCACHE=/Users/andrey/Documents/pets/sumerki/.cache/go-build go test ./...`.
- Ran `POSTGRES_PORT=15432 docker compose up -d postgres` for Phase 5 verification because sandboxed curl could not reach the elevated backend process and default port `5432` may be occupied by another Docker process.
- Ran `DATABASE_URL=postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable make migrate-up`; database was already at migration version 2.
- Ran backend with `DATABASE_URL=postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable`, `JWT_SECRET=test-secret`, and `BACKEND_PORT=18080`.
- Ran elevated `curl -i http://localhost:18080/ready`; received HTTP 200 with `{"status":"ready","database":"ok"}`.
- Registered temporary users through `POST /api/auth/register` and received JWT tokens.
- Ran unauthenticated `POST /api/kingdoms`; received HTTP 401 with `missing_authorization_header`.
- Ran authenticated `GET /api/kingdoms/me` before creation; received HTTP 200 with `{"kingdom":null}`.
- Ran authenticated `POST /api/kingdoms` with trimmed name input plus extra `user_id` and `patron` fields; received HTTP 201, trimmed kingdom name, authenticated `userId`, and `patron:null`.
- Ran second authenticated `POST /api/kingdoms`; received HTTP 409 with `kingdom_already_exists`.
- Ran authenticated `POST /api/kingdoms` with invalid culture; received HTTP 400 with `invalid_culture`.
- Ran authenticated `POST /api/kingdoms` with too-short name; received HTTP 400 with `kingdom_name_too_short`.
- Ran authenticated `POST /api/kingdoms` with too-long name; received HTTP 400 with `kingdom_name_too_long`.
- Ran authenticated `GET /api/kingdoms/me` after creation; received HTTP 200 with the created kingdom.
- Removed temporary verification users from PostgreSQL; related kingdoms were removed by cascade.
- Stopped the backend with SIGINT and ran `POSTGRES_PORT=15432 docker compose down`.
- Ran final `docker compose ps` and confirmed no Sumerki Compose service was left running.

## What Works Now

- Local PostgreSQL can be started with `docker compose up -d postgres` or `make db-up`.
- Local PostgreSQL uses database `sumerki`, user `sumerki`, password `sumerki`, and host port `5432`.
- Backend can be started with `cd backend && DATABASE_URL="postgres://sumerki:sumerki@localhost:5432/sumerki?sslmode=disable" JWT_SECRET="dev-secret" go run ./cmd/server` when PostgreSQL is reachable.
- `GET /health` reports process health without requiring database connectivity.
- `GET /ready` reports database readiness and returns a standard JSON 503 error when PostgreSQL is unavailable.
- `make migrate-up` applies initial database schema migrations.
- `POST /api/auth/register` creates a user with a bcrypt password hash and returns a JWT.
- `POST /api/auth/login` authenticates normalized email and password and returns a JWT.
- `GET /api/me` returns the current user when a valid Bearer token is provided.
- `POST /api/kingdoms` creates one kingdom for the authenticated user.
- `GET /api/kingdoms/me` returns `{"kingdom":null}` before creation and the user's kingdom after creation.
- Kingdom responses do not expose `password_hash`.

## Known Limitations

- No frontend, resources, buildings, ruler generation, army, missions, combat, events, patrons logic, tribute, alliances, map, or payments exist yet.
- Default PostgreSQL port `5432` may be occupied by another Docker process in this environment, so successful Phase 5 live checks used temporary `POSTGRES_PORT=15432` without changing repository files.
- Goose `@latest` currently switches to Go toolchain `go1.25.11` during command execution because the latest Goose release requires Go >= 1.25.7.

## Next Recommended Step

Start Phase 6: Ruler System.
