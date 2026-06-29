# Codex Handoff

## Current Phase

Phase 10: Buildings API + UI.

## Status

Phase 10 is complete according to the attached prompt: kingdoms now have MVP building rows, the backend exposes authenticated building list and upgrade endpoints, building upgrades spend resources and complete through lazy resolution, building levels affect resource production, and the dashboard can display and start building upgrades.

## Completed

- Created initial repository guidance in `AGENTS.md`.
- Added `.gitignore` for future Go and TypeScript development artifacts.
- Created project overview in `README.md`.
- Documented MVP scope in `docs/MVP_SCOPE.md`.
- Preserved and aligned phase plan in `docs/MVP_PHASES.md`.
- Created decision log in `docs/DECISIONS.md`.
- Drafted and updated API contract in `docs/API_CONTRACT.md`.
- Drafted and updated domain model in `docs/DOMAIN_MODEL.md`.
- Created phase documentation directory at `docs/phases/`.
- Added Docker Compose configuration for local PostgreSQL.
- Added `.env.example` with local PostgreSQL defaults.
- Added basic Makefile commands for local database management and migrations.
- Added backend skeleton, health/readiness endpoints, auth API, kingdom creation API, frontend auth/kingdom flow, ruler system, and resources system in earlier phases.
- Added reversible Goose migration `00005_create_kingdom_buildings.sql`.
- Added `kingdom_buildings` table with one row per kingdom/building type, level constraints, upgrade timestamps, and safe backfill for existing kingdoms.
- Added MVP building configuration for:
  - `town_hall`
  - `farm`
  - `lumberyard`
  - `quarry`
  - `market`
  - `barracks`
  - `walls`
  - `shrine`
- Added building domain, repository, service, and HTTP handler layers.
- Added lazy upgrade completion when building endpoints or resource production are resolved.
- Added resource spending support for building upgrades.
- Added building production bonuses:
  - farm: `+15 food/hour` per level
  - lumberyard: `+12 wood/hour` per level
  - quarry: `+10 stone/hour` per level
  - market: `+10 gold/hour` per level
- Integrated initial building creation into kingdom creation.
- Integrated building production bonuses into `GET /api/resources/me`.
- Added `GET /api/buildings/me`.
- Added `POST /api/buildings/{type}/upgrade`.
- Added frontend building API types and calls.
- Replaced the dashboard buildings placeholder with real building cards, level/progress display, next-upgrade cost/duration, upgrade buttons, loading state, and error state.
- Updated README with building API curl examples and the manual Phase 10 flow.
- Added service tests for building listing, upgrade spending, insufficient resources, lazy completion, production bonuses, resource production bonuses, and kingdom building initialization.

## Phase Order Note

- The attached prompt defines Phase 10 as `Buildings API + UI`.
- `docs/MVP_PHASES.md` currently defines Phase 10 as `Buildings API + UI`, which is aligned.
- This session followed the attached prompt and did not modify `docs/MVP_PHASES.md`.

## Changed Files

- `README.md`
- `CODEX_HANDOFF.md`
- `docs/API_CONTRACT.md`
- `docs/DOMAIN_MODEL.md`
- `backend/migrations/00005_create_kingdom_buildings.sql`
- `backend/internal/gameconfig/buildings.go`
- `backend/internal/domain/building.go`
- `backend/internal/repository/building_repository.go`
- `backend/internal/service/building_service.go`
- `backend/internal/service/building_service_test.go`
- `backend/internal/service/resources_service.go`
- `backend/internal/service/resources_service_test.go`
- `backend/internal/service/kingdom_service_test.go`
- `backend/internal/http/handlers/building_handler.go`
- `backend/internal/http/server.go`
- `frontend/src/api/client.ts`
- `frontend/src/api/errors.ts`
- `frontend/src/pages/DashboardPage.tsx`

## Commands Run

- `gofmt -w backend/internal/domain/building.go backend/internal/gameconfig/buildings.go backend/internal/repository/building_repository.go backend/internal/service/building_service.go backend/internal/service/building_service_test.go backend/internal/service/resources_service.go backend/internal/service/resources_service_test.go backend/internal/http/handlers/building_handler.go backend/internal/http/server.go backend/internal/service/kingdom_service_test.go`
- `cd frontend && npm run typecheck`
- `cd frontend && npm run build`
- `cd backend && GOCACHE=/Users/andrey/Documents/pets/sumerki/.cache/go-build GOMODCACHE=/Users/andrey/Documents/pets/sumerki/.cache/go-mod go test ./...`
- `DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable' make migrate-up`
- `GOCACHE=/Users/andrey/Documents/pets/sumerki/.cache/go-build GOMODCACHE=/Users/andrey/Documents/pets/sumerki/.cache/go-mod DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable' JWT_SECRET='dev-secret' BACKEND_PORT=18080 go run ./cmd/server`
- `curl -i http://localhost:18080/ready`
- HTTP smoke test covering unauthenticated `GET /api/buildings/me`, authenticated no-kingdom `GET /api/buildings/me`, `POST /api/kingdoms`, successful `GET /api/buildings/me`, successful `GET /api/resources/me`, and successful `POST /api/buildings/farm/upgrade`
- Local DB timestamp adjustment to verify lazy building completion
- Browser flow covering dashboard building cards and starting a farm upgrade from the UI
- `DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable' make migrate-status`
- `git diff --check`

## Verification

- `npm run typecheck` completed successfully.
- `npm run build` completed successfully.
- `go test ./...` completed successfully.
- Goose applied `00005_create_kingdom_buildings.sql` successfully.
- Goose status shows migrations `00001`, `00002`, `00003`, `00004`, and `00005` applied.
- Verified backend readiness returned HTTP 200 with `{"status":"ready","database":"ok"}`.
- Verified `GET /api/buildings/me` without auth returns HTTP 401.
- Verified authenticated `GET /api/buildings/me` without a kingdom returns HTTP 404 with `kingdom_not_found`.
- Verified creating a new kingdom creates all 8 MVP buildings.
- Verified `GET /api/buildings/me` returns all 8 buildings with levels, effects, upgrade state, and next-upgrade data.
- Verified `POST /api/buildings/farm/upgrade` spends resources and starts an upgrade timer.
- Verified lazy building completion by moving `upgrade_finishes_at` into the past and confirming the farm reached level 2 on the next read.
- Verified `GET /api/resources/me` includes building production bonuses.
- Verified the dashboard shows all 8 building cards.
- Verified the dashboard can start a farm upgrade and then displays `Улучшается`.
- Verified `git diff --check` has no whitespace errors.

Notes:

- Port `5432` was already in use locally, so live verification used `POSTGRES_PORT=15432` and `DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable'`.
- The backend was run on `18080` for verification, and the frontend was already running on `5173` with `VITE_API_BASE_URL=http://localhost:18080`.
- The first sandboxed `make migrate-status` attempt failed because `go run github.com/pressly/goose/v3/cmd/goose@latest` could not resolve `proxy.golang.org`; the same command succeeded with approved network/module access.

## What Works Now

- New migrations create and backfill `kingdom_buildings`.
- Existing kingdoms receive buildings through migration backfill or lazy creation.
- New kingdoms get initial buildings during creation.
- `GET /api/buildings/me` requires authentication.
- `GET /api/buildings/me` returns only the authenticated user's kingdom buildings.
- Users without kingdoms receive the standard `kingdom_not_found` JSON error.
- Building upgrades validate type, max level, existing upgrade state, and resource affordability.
- Building upgrades spend resources and set `upgrade_started_at` / `upgrade_finishes_at`.
- Finished upgrades are resolved lazily by building reads, upgrade attempts, and resource reads.
- Farm, lumberyard, quarry, and market levels modify resource production.
- Dashboard loads and displays real buildings after session and kingdom load.
- Dashboard can start building upgrades.

## Known Limitations

- Building upgrades and resource spending are not wrapped in one database transaction yet; a failure after spending but before setting the upgrade timer could leave resources spent without the upgrade.
- Kingdom, ruler, resource, and building initialization are not transactional in this phase.
- Building costs and durations are simple MVP constants.
- Building upgrade queues are single-upgrade-per-building only; no cancellation, instant completion, premium queue, or multi-queue behavior.
- Production modifiers only cover farm, lumberyard, quarry, and market.
- Town hall, barracks, walls, and shrine have no gameplay effect yet.
- Resource caps are still not implemented.
- No background workers by design; resources and buildings use lazy resolution.
- No armies, missions, combat, raids, events, patrons, tribute, market, trading, payments, chat, alliances, large map, or real-time systems were implemented.

## Next Recommended Step

Start Phase 11 only when explicitly requested. Based on `docs/MVP_PHASES.md`, the next likely phase is Army System.
