# Codex Handoff

## Current Phase

Phase 8: Ruler System v1.

## Status

Phase 8 is complete according to the attached prompt: every kingdom now has a simple ruler, the backend exposes `GET /api/ruler/me`, and the dashboard displays real ruler data.

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
- Added backend skeleton, health/readiness endpoints, auth API, kingdom creation API, and frontend auth/kingdom flow in earlier phases.
- Added reversible Goose migration `00003_create_rulers.sql`.
- Added `rulers` table with UUID primary key, one ruler per kingdom, stat constraints, culture constraints, health status constraints, timestamps, and safe backfill for existing kingdoms.
- Added ruler domain, repository, service, and HTTP handler layers.
- Added simple culture-specific ruler generation.
- Added safe lazy ruler creation when an existing kingdom is missing a ruler.
- Integrated ruler creation into kingdom creation.
- Added `GET /api/ruler/me`.
- Added frontend `Ruler` type and `getMyRuler()` API call.
- Replaced the dashboard placeholder ruler card with real ruler data, loading state, and error state.
- Updated README with the ruler curl example and local manual flow.
- Added minimal service tests for ruler generation and kingdom creation integration.

## Phase Order Note

- The attached prompt defines Phase 8 as `Ruler System v1`.
- `docs/MVP_PHASES.md` currently defines Phase 8 as `Frontend Auth, Kingdom, and Ruler Integration`, while its older Phase 6 is `Ruler System`.
- This session followed the attached prompt and did not modify `docs/MVP_PHASES.md`.

## Changed Files

- `README.md`
- `CODEX_HANDOFF.md`
- `docs/API_CONTRACT.md`
- `docs/DOMAIN_MODEL.md`
- `backend/migrations/00003_create_rulers.sql`
- `backend/internal/domain/ruler.go`
- `backend/internal/repository/ruler_repository.go`
- `backend/internal/service/ruler_service.go`
- `backend/internal/service/ruler_service_test.go`
- `backend/internal/service/kingdom_service.go`
- `backend/internal/service/kingdom_service_test.go`
- `backend/internal/http/handlers/ruler_handler.go`
- `backend/internal/http/server.go`
- `frontend/src/api/client.ts`
- `frontend/src/pages/DashboardPage.tsx`

## Commands Run

- `gofmt -w backend/internal/domain/ruler.go backend/internal/repository/ruler_repository.go backend/internal/service/ruler_service.go backend/internal/service/kingdom_service.go backend/internal/http/handlers/ruler_handler.go backend/internal/http/server.go backend/internal/service/kingdom_service_test.go backend/internal/service/ruler_service_test.go`
- `npm run typecheck`
- `npm run build`
- `GOCACHE=/Users/andrey/Documents/pets/sumerki/.cache/go-build GOMODCACHE=/Users/andrey/Documents/pets/sumerki/.cache/go-mod go test ./...`
- `DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable' make migrate-up`
- `GOCACHE=/Users/andrey/Documents/pets/sumerki/.cache/go-build GOMODCACHE=/Users/andrey/Documents/pets/sumerki/.cache/go-mod DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable' JWT_SECRET='dev-secret' BACKEND_PORT=18080 go run ./cmd/server`
- `curl -sS -i http://localhost:18080/ready`
- HTTP smoke test covering unauthenticated `GET /api/ruler/me`, authenticated no-kingdom `GET /api/ruler/me`, `POST /api/kingdoms`, and successful `GET /api/ruler/me`
- Browser flow covering register, create kingdom, dashboard ruler card, and dashboard refresh
- `DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable' make migrate-status`
- `git diff --check`

## Verification

- `npm run typecheck` completed successfully.
- `npm run build` completed successfully.
- `go test ./...` completed successfully.
- Goose applied `00003_create_rulers.sql` successfully.
- Goose status shows migrations `00001`, `00002`, and `00003` applied.
- Verified backend readiness returned HTTP 200 with `{"status":"ready","database":"ok"}`.
- Verified `GET /api/ruler/me` without auth returns HTTP 401.
- Verified authenticated `GET /api/ruler/me` without a kingdom returns HTTP 404 with `kingdom_not_found`.
- Verified creating a new kingdom creates a ruler.
- Verified `GET /api/ruler/me` returns the current user's ruler with kingdom id, name, age, culture, stats, health status, and timestamps.
- Verified ruler age and stat generation ranges through service tests.
- Verified dashboard shows the real ruler card with Russian stat labels.
- Verified dashboard still shows ruler data after page refresh.
- Verified `git diff --check` has no whitespace errors.

Notes:

- Port `5432` was already in use locally, so live verification used `POSTGRES_PORT=15432` and `DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable'`.
- The backend was run on `18080` for verification, and the frontend was already running on `5173` with `VITE_API_BASE_URL=http://localhost:18080`.
- Kingdom creation and ruler creation are not wrapped in one database transaction. The service creates the kingdom first and then creates the ruler immediately after; missing rulers are safely created lazily by `GET /api/ruler/me`.

## What Works Now

- New migrations create and backfill `rulers`.
- Existing kingdoms receive rulers through migration backfill or lazy creation.
- New kingdoms get a generated ruler during creation.
- `GET /api/ruler/me` requires authentication.
- `GET /api/ruler/me` returns only the authenticated user's ruler.
- Users without kingdoms receive the standard `kingdom_not_found` JSON error.
- Dashboard loads and displays real ruler data after session and kingdom load.
- Dashboard keeps working after refresh.

## Known Limitations

- Ruler stats are flavor only for now.
- No ruler traits yet.
- No ruler actions yet.
- No ruler death, heirs, or dynasties yet.
- No intrigue system yet.
- No gameplay modifiers from ruler stats yet.
- No resources, buildings, army, missions, combat, events, patrons, tribute, alliances, or map systems were implemented.
- Kingdom and ruler creation are not transactional in this phase.

## Next Recommended Step

Start Phase 9 only when explicitly requested. Based on the prompt-driven order, the next likely phase is the resources system.
