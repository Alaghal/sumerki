# Codex Handoff

## Current Phase

Phase 13: Report Polish and Narrative Templates.

## Status

Phase 13 is complete according to the attached prompt: PvE mission reports now have deterministic local narrative templates, ordered report phases, unread counts, report detail, and idempotent mark-as-read support.

## Completed

- Added reversible Goose migration `00008_add_report_phases.sql`.
- Added `mission_reports.phases_json` with default empty JSON array for existing reports.
- Added local deterministic report templates for:
  - `black_forest_expedition`
  - `old_kurgan_expedition`
  - `dry_ford_scouting`
- Added report templates for `success`, `partial_success`, and `failure`.
- Added report phase data with title/body sections.
- Extended mission report domain and repository access.
- Added paginated report listing with unread count.
- Added report detail lookup scoped to the authenticated user's kingdom.
- Added idempotent mark-as-read support scoped to the authenticated user's kingdom.
- Added `GET /api/reports/:id`.
- Added `POST /api/reports/:id/read`.
- Updated `GET /api/reports/me` response with `phases`, `pagination`, and `unreadCount`.
- Updated frontend API types and helpers for report detail and mark-read.
- Updated dashboard Reports UI with unread count, read/unread state, report detail expansion, phases, refresh, and mark-as-read action.
- Updated README report curl examples and local flow.
- Updated API contract and domain model docs.
- Added service tests for report unread count, legacy empty phases, ownership scoping, idempotent read, and template phases.

## Phase Order Note

- The attached prompt defines Phase 13 as `Report Polish and Narrative Templates`.
- This session followed the attached prompt and did not start Phase 14.
- No PvP raids, events, patrons, tribute, dark gods, payments, chat, comments, notifications, WebSocket, background jobs, alliances, or map systems were implemented.

## Changed Files

- `README.md`
- `CODEX_HANDOFF.md`
- `docs/API_CONTRACT.md`
- `docs/DOMAIN_MODEL.md`
- `backend/migrations/00008_add_report_phases.sql`
- `backend/internal/domain/report.go`
- `backend/internal/gameconfig/report_templates.go`
- `backend/internal/repository/report_repository.go`
- `backend/internal/service/mission_service.go`
- `backend/internal/service/mission_service_test.go`
- `backend/internal/http/handlers/report_handler.go`
- `backend/internal/http/server.go`
- `frontend/src/api/client.ts`
- `frontend/src/pages/DashboardPage.tsx`

## Commands Run

- `gofmt -w backend/internal/domain/report.go backend/internal/gameconfig/report_templates.go backend/internal/repository/report_repository.go backend/internal/service/mission_service.go backend/internal/service/mission_service_test.go backend/internal/http/handlers/report_handler.go backend/internal/http/server.go`
- `cd backend && GOCACHE=/Users/andrey/Documents/pets/sumerki/.cache/go-build GOMODCACHE=/Users/andrey/Documents/pets/sumerki/.cache/go-mod go test ./...`
- `cd frontend && npm run typecheck`
- `cd frontend && npm run build`
- `docker compose up -d postgres`
- `docker compose ps`
- `POSTGRES_PORT=15432 docker compose up -d postgres`
- `POSTGRES_PORT=15432 docker compose ps`
- `DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable' make migrate-up`
- `DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable' make migrate-status`
- `cd backend && DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable' JWT_SECRET='dev-secret' BACKEND_PORT=18080 GOCACHE=/Users/andrey/Documents/pets/sumerki/.cache/go-build GOMODCACHE=/Users/andrey/Documents/pets/sumerki/.cache/go-mod go run ./cmd/server`
- `curl` smoke tests for register, create kingdom, mission start, report list, report detail, and report mark-read.
- `POSTGRES_PORT=15432 docker compose exec -T postgres psql ...` to move test mission `finishes_at` timestamps into the past for lazy report resolution.

## Verification

- `go test ./...` completed successfully.
- `npm run typecheck` completed successfully.
- `npm run build` completed successfully.
- Goose applied `00008_add_report_phases.sql` successfully.
- Goose status shows migrations `00001` through `00008` applied.
- Verified existing reports return `phases: []` after the migration.
- Verified a newly generated `dry_ford_scouting` report includes narrative phases.
- Verified `GET /api/reports/me?limit=20&offset=0` returns `reports`, `pagination`, and `unreadCount`.
- Verified `GET /api/reports/:id` returns one owned report with phases.
- Verified `POST /api/reports/:id/read` returns `isRead: true`.
- Verified repeated `POST /api/reports/:id/read` succeeds and remains `isRead: true`.
- Verified unread count decreases after marking one report read.

Notes:

- Port `5432` was already allocated locally, so live verification used `POSTGRES_PORT=15432` and `DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable'`.
- A stale backend process was listening on `18080`; it was stopped before running the updated backend on `18080`.
- Combined local API and Docker commands needed approved local-process/local-network execution.

## What Works Now

- Completed PvE missions create reports with local narrative templates and ordered phases.
- Existing pre-Phase-13 reports remain readable with an empty phases list.
- Players can list reports with pagination metadata and unread count.
- Players can open a single owned report.
- Players can mark an owned report as read.
- Mark-read is idempotent.
- Dashboard shows report unread count, read/unread state, details, phases, rewards, losses, and refresh/read actions.

## Known Limitations

- Report templates are deterministic and local; there is no AI-generated text.
- Report content is only implemented for the three current PvE mission keys.
- Mission resolution still is not wrapped in one cross-repository database transaction.
- Mission outcomes use simple deterministic MVP rules; there is no advanced battle simulation.
- There are no comments, notifications, WebSocket updates, or background jobs for reports.
- No PvP raids, player-vs-player combat, events, patrons, tribute, dark gods, alliances, map, market, trading, payments, chat, or real-time systems were implemented.
- No background workers by design; resources, buildings, unit training, and missions use lazy resolution.

## Next Recommended Step

Start Phase 14 only when explicitly requested. Based on `docs/MVP_PHASES.md`, the next likely phase is Patron System.
