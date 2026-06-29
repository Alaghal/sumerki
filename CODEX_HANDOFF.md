# Codex Handoff

## Current Phase

Phase 11: Army API + UI.

## Status

Phase 11 is complete according to the attached prompt: kingdoms now have MVP unit rows, the backend exposes authenticated army list and training endpoints, unit training spends resources and completes through lazy resolution, barracks requirements are enforced through the buildings service, and the dashboard can display units and start training orders.

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
- Added backend skeleton, health/readiness endpoints, auth API, kingdom creation API, frontend auth/kingdom flow, ruler system, resources system, and buildings system in earlier phases.
- Added reversible Goose migration `00006_create_army.sql`.
- Added `kingdom_units` table with one row per kingdom/unit type, nonnegative amount constraints, and safe backfill for existing kingdoms.
- Added `unit_training_orders` table with training/completed status, finish timestamps, and indexes by kingdom, status, and finish time.
- Added MVP unit configuration for:
  - `militia`
  - `spearmen`
  - `archers`
  - `cavalry`
  - `scouts`
- Added army domain, repository, service, and HTTP handler layers.
- Added lazy unit training completion when army endpoints are read or training starts.
- Added resource spending support for unit training, including population spending.
- Added barracks requirement checks through the buildings service:
  - militia: no barracks requirement
  - scouts: no barracks requirement
  - spearmen: barracks level 1
  - archers: barracks level 1
  - cavalry: barracks level 2
- Integrated initial unit creation into kingdom creation.
- Added `GET /api/army/me`.
- Added `POST /api/army/train`.
- Added frontend army API types and calls.
- Replaced the dashboard army placeholder with real unit cards, stats, costs, requirements, training form, active training orders, loading state, and error state.
- Updated README with army API curl examples and the manual Phase 11 flow.
- Added service tests for initial army rows, training resource spending, barracks requirements, invalid amount rejection, and lazy completion.

## Phase Order Note

- The attached prompt defines Phase 11 as `Army API + UI`.
- `docs/MVP_PHASES.md` currently defines Phase 11 as `Army System`, which is aligned in substance.
- This session followed the attached prompt and did not modify `docs/MVP_PHASES.md`.

## Changed Files

- `README.md`
- `CODEX_HANDOFF.md`
- `docs/API_CONTRACT.md`
- `docs/DOMAIN_MODEL.md`
- `backend/migrations/00006_create_army.sql`
- `backend/internal/gameconfig/units.go`
- `backend/internal/domain/army.go`
- `backend/internal/repository/army_repository.go`
- `backend/internal/service/army_service.go`
- `backend/internal/service/army_service_test.go`
- `backend/internal/service/building_service.go`
- `backend/internal/http/handlers/army_handler.go`
- `backend/internal/http/server.go`
- `frontend/src/api/client.ts`
- `frontend/src/api/errors.ts`
- `frontend/src/pages/DashboardPage.tsx`

## Commands Run

- `gofmt -w backend/internal/domain/army.go backend/internal/gameconfig/units.go backend/internal/repository/army_repository.go backend/internal/service/army_service.go backend/internal/service/building_service.go backend/internal/http/handlers/army_handler.go backend/internal/http/server.go`
- `gofmt -w backend/internal/service/army_service_test.go`
- `cd frontend && npm run typecheck`
- `cd frontend && npm run build`
- `cd backend && GOCACHE=/Users/andrey/Documents/pets/sumerki/.cache/go-build GOMODCACHE=/Users/andrey/Documents/pets/sumerki/.cache/go-mod go test ./...`
- `DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable' make migrate-up`
- `DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable' make migrate-status`
- `cd backend && GOCACHE=/Users/andrey/Documents/pets/sumerki/.cache/go-build GOMODCACHE=/Users/andrey/Documents/pets/sumerki/.cache/go-mod DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable' JWT_SECRET='dev-secret' BACKEND_PORT=18080 go run ./cmd/server`
- `curl -i http://127.0.0.1:18080/ready`
- HTTP smoke test covering `POST /api/auth/register`, `POST /api/kingdoms`, successful `GET /api/army/me`, successful `POST /api/army/train`, and successful `GET /api/resources/me`
- Local DB timestamp adjustment to verify lazy unit training completion
- Browser flow covering dashboard army cards and starting militia training from the UI
- `git diff --check`

## Verification

- `npm run typecheck` completed successfully.
- `npm run build` completed successfully.
- `go test ./...` completed successfully.
- Goose applied `00006_create_army.sql` successfully.
- Goose status shows migrations `00001`, `00002`, `00003`, `00004`, `00005`, and `00006` applied.
- Verified backend readiness returned HTTP 200 with `{"status":"ready","database":"ok"}`.
- Verified creating a new kingdom creates all 5 MVP unit rows.
- Verified `GET /api/army/me` returns all 5 units, stats, costs, barracks requirements, active training orders, and army summary.
- Verified `POST /api/army/train` for 5 militia spends gold, food, and population, then creates a training order.
- Verified `GET /api/resources/me` returns the reduced resources after training starts.
- Verified lazy training completion by moving `finishes_at` into the past and confirming militia increased from 10 to 15 and active training orders became empty.
- Verified the dashboard shows the Army section with all 5 unit cards.
- Verified the dashboard can start militia training and displays a `Завершится` training order.
- Verified `git diff --check` has no whitespace errors.

Notes:

- Port `5432` was already in use locally, so live verification used `POSTGRES_PORT=15432` and `DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable'`.
- The backend was run on `18080` for verification, and the frontend was already running on `5173` with `VITE_API_BASE_URL=http://localhost:18080`.
- Combined shell-based HTTP smoke commands needed approved local-network execution; individual `curl` readiness checks worked against `127.0.0.1`.

## What Works Now

- New migrations create and backfill `kingdom_units`.
- New migrations create `unit_training_orders`.
- Existing kingdoms receive units through migration backfill or lazy creation.
- New kingdoms get initial units during creation.
- `GET /api/army/me` requires authentication.
- `GET /api/army/me` returns only the authenticated user's kingdom army.
- Users without kingdoms receive the standard `kingdom_not_found` JSON error.
- Unit training validates type, amount, barracks level, and resource affordability.
- Unit training spends resources immediately and creates a training order.
- Finished training orders are resolved lazily by army reads and training commands.
- Dashboard loads and displays real army data after session and kingdom load.
- Dashboard can start unit training and refresh affected resources.

## Known Limitations

- Unit training and resource spending are not wrapped in one database transaction yet; a failure after spending but before creating the training order could leave resources spent without an order.
- Kingdom, ruler, resource, building, and unit initialization are not transactional in this phase.
- Unit costs, stats, and training durations are simple MVP constants.
- Training orders cannot be cancelled.
- There are no premium speedups, queues, or formation systems.
- There is no unit upkeep.
- There is no unit death, healing, equipment, or heroes leading armies.
- Armies cannot be sent anywhere yet.
- No missions, PvE expeditions, combat, raids, reports, events, patrons, tribute, alliances, map, market, trading, payments, chat, or real-time systems were implemented.
- No background workers by design; resources, buildings, and unit training use lazy resolution.

## Next Recommended Step

Start Phase 12 only when explicitly requested. Based on `docs/MVP_PHASES.md`, the next likely phase is PvE Missions.
