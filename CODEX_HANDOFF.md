# Codex Handoff

## Current Phase

Phase 9: Resources API + UI.

## Status

Phase 9 is complete according to the attached prompt: kingdoms now have stored resources, resources grow through lazy calculation, the backend exposes `GET /api/resources/me`, and the dashboard displays real resource data with manual refresh.

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
- Added basic Makefile commands for local database management.
- Added backend skeleton, health/readiness endpoints, auth API, kingdom creation API, frontend auth/kingdom flow, and ruler system in earlier phases.
- Added reversible Goose migration `00004_create_kingdom_resources.sql`.
- Added `kingdom_resources` table with one row per kingdom, nonnegative constraints, timestamps, and safe backfill for existing kingdoms.
- Added resource configuration in `backend/internal/gameconfig/resources.go`.
- Added resources domain, repository, service, and HTTP handler layers.
- Added simple base production per hour:
  - gold: 20
  - food: 30
  - wood: 25
  - stone: 15
  - population: 1
- Added lazy resource calculation from `last_calculated_at`.
- Added safe lazy resource row creation when an existing kingdom is missing resources.
- Integrated initial resource creation into kingdom creation.
- Added `GET /api/resources/me`.
- Added frontend `Resources` type and `getMyResources()` API call.
- Replaced the dashboard placeholder resources card with real resource data, loading state, error state, production text, and manual refresh button.
- Updated README with the resources curl example and local manual flow.
- Added minimal service tests for resource creation, lazy production, missing kingdom behavior, and kingdom creation integration.

## Phase Order Note

- The attached prompt defines Phase 9 as `Resources API + UI`.
- `docs/MVP_PHASES.md` currently defines Phase 9 as `Resources System`, which is aligned in substance.
- This session followed the attached prompt and did not modify `docs/MVP_PHASES.md`.

## Changed Files

- `README.md`
- `CODEX_HANDOFF.md`
- `docs/API_CONTRACT.md`
- `docs/DOMAIN_MODEL.md`
- `backend/migrations/00004_create_kingdom_resources.sql`
- `backend/internal/gameconfig/resources.go`
- `backend/internal/domain/resources.go`
- `backend/internal/repository/resources_repository.go`
- `backend/internal/service/resources_service.go`
- `backend/internal/service/resources_service_test.go`
- `backend/internal/service/kingdom_service.go`
- `backend/internal/service/kingdom_service_test.go`
- `backend/internal/service/ruler_service.go`
- `backend/internal/http/handlers/resources_handler.go`
- `backend/internal/http/server.go`
- `frontend/src/api/client.ts`
- `frontend/src/pages/DashboardPage.tsx`

## Commands Run

- `gofmt -w backend/internal/gameconfig/resources.go backend/internal/domain/resources.go backend/internal/repository/resources_repository.go backend/internal/service/resources_service.go backend/internal/service/ruler_service.go backend/internal/service/kingdom_service.go backend/internal/http/handlers/resources_handler.go backend/internal/http/server.go backend/internal/service/kingdom_service_test.go backend/internal/service/resources_service_test.go`
- `npm run typecheck`
- `npm run build`
- `GOCACHE=/Users/andrey/Documents/pets/sumerki/.cache/go-build GOMODCACHE=/Users/andrey/Documents/pets/sumerki/.cache/go-mod go test ./...`
- `DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable' make migrate-up`
- `GOCACHE=/Users/andrey/Documents/pets/sumerki/.cache/go-build GOMODCACHE=/Users/andrey/Documents/pets/sumerki/.cache/go-mod DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable' JWT_SECRET='dev-secret' BACKEND_PORT=18080 go run ./cmd/server`
- `curl -sS -i http://localhost:18080/ready`
- HTTP smoke test covering unauthenticated `GET /api/resources/me`, authenticated no-kingdom `GET /api/resources/me`, `POST /api/kingdoms`, and successful `GET /api/resources/me`
- Local DB timestamp adjustment to verify lazy one-hour production
- Browser flow covering register, create kingdom, dashboard resource card, manual refresh, and dashboard refresh
- `DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable' make migrate-status`
- `git diff --check`

## Verification

- `npm run typecheck` completed successfully.
- `npm run build` completed successfully.
- `go test ./...` completed successfully.
- Goose applied `00004_create_kingdom_resources.sql` successfully.
- Goose status shows migrations `00001`, `00002`, `00003`, and `00004` applied.
- Verified backend readiness returned HTTP 200 with `{"status":"ready","database":"ok"}`.
- Verified `GET /api/resources/me` without auth returns HTTP 401.
- Verified authenticated `GET /api/resources/me` without a kingdom returns HTTP 404 with `kingdom_not_found`.
- Verified creating a new kingdom creates initial resources.
- Verified `GET /api/resources/me` returns the current user's resources with production per hour.
- Verified lazy production by moving a test row's `last_calculated_at` one hour back and confirming resources increased by one hour of production.
- Verified the dashboard shows real resource labels, values, and `+N / час` production text.
- Verified the manual `Обновить ресурсы` button calls the resources endpoint successfully.
- Verified dashboard still shows resources after page refresh.
- Verified `git diff --check` has no whitespace errors.

Notes:

- Port `5432` was already in use locally, so live verification used `POSTGRES_PORT=15432` and `DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable'`.
- The backend was run on `18080` for verification, and the frontend was already running on `5173` with `VITE_API_BASE_URL=http://localhost:18080`.
- Kingdom creation, ruler creation, and resource creation are not wrapped in one database transaction. The service creates them sequentially; missing rulers and resources are safely created lazily by their read endpoints.

## What Works Now

- New migrations create and backfill `kingdom_resources`.
- Existing kingdoms receive resources through migration backfill or lazy creation.
- New kingdoms get initial resources during creation.
- `GET /api/resources/me` requires authentication.
- `GET /api/resources/me` returns only the authenticated user's kingdom resources.
- Users without kingdoms receive the standard `kingdom_not_found` JSON error.
- Resources are recalculated lazily from `last_calculated_at` when `GET /api/resources/me` is called.
- Dashboard loads and displays real resources after session and kingdom load.
- Dashboard resource card can be refreshed manually.
- Dashboard keeps working after refresh.

## Known Limitations

- Resources use simple base production only.
- No buildings modify production yet.
- No resource caps yet.
- No resource spending yet.
- No background resource workers by design.
- Resource calculation uses simple integer flooring for MVP and does not preserve fractional progress when any whole-unit resource gain is applied.
- No buildings, army, missions, combat, raids, events, patrons, tribute, market, trading, payments, or gameplay spending systems were implemented.
- Kingdom, ruler, and resource creation are not transactional in this phase.

## Next Recommended Step

Start Phase 10 only when explicitly requested. Based on the prompt-driven order, the next likely phase is the buildings system.
