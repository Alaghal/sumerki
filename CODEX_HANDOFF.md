# Codex Handoff

## Current Phase

Phase 12: PvE Missions with Basic Reports.

## Status

Phase 12 is complete according to the attached prompt: players can view configured PvE missions, send available units, have those units removed while the mission is active, resolve missions lazily, receive resource rewards, lose or recover units, and read basic mission reports in the dashboard.

## Completed

- Added reversible Goose migration `00007_create_missions_and_reports.sql`.
- Added `missions`, `mission_units`, and `mission_reports` tables with status/type/result constraints and indexes.
- Added mission configuration for:
  - `black_forest_expedition`
  - `old_kurgan_expedition`
  - `dry_ford_scouting`
- Added mission and report domain models.
- Added mission and report repositories.
- Added mission service with:
  - available mission listing
  - mission start validation
  - unit requirement checks
  - available-unit checks
  - sent-unit removal
  - lazy mission completion
  - survivor returns
  - deterministic MVP losses
  - resource rewards
  - basic report creation
- Added small army service/repository methods for mission unit subtraction and return.
- Added resource reward grant support after lazy recalculation.
- Added `GET /api/missions/available`.
- Added `GET /api/missions/me`.
- Added `POST /api/missions/start`.
- Added `GET /api/reports/me`.
- Added frontend API types and calls for missions and reports.
- Added dashboard Missions section with available mission cards, unit amount inputs, start action, current mission display, and refresh button.
- Added dashboard Reports section with latest mission reports.
- Updated README with mission/report curl examples and Phase 12 local flow.
- Updated API contract and domain model docs.
- Added mission service tests for invalid mission key, insufficient units, successful unit subtraction, lazy completion, rewards, unit return, and report creation.

## Phase Order Note

- The attached prompt defines Phase 12 as `PvE Missions with Basic Reports`.
- `docs/MVP_PHASES.md` currently defines Phase 12 as `PvE Missions`, which is aligned in substance.
- This session followed the attached prompt and did not modify `docs/MVP_PHASES.md`.

## Changed Files

- `README.md`
- `CODEX_HANDOFF.md`
- `docs/API_CONTRACT.md`
- `docs/DOMAIN_MODEL.md`
- `backend/migrations/00007_create_missions_and_reports.sql`
- `backend/internal/gameconfig/missions.go`
- `backend/internal/domain/mission.go`
- `backend/internal/domain/report.go`
- `backend/internal/repository/army_repository.go`
- `backend/internal/repository/mission_repository.go`
- `backend/internal/repository/report_repository.go`
- `backend/internal/service/army_service.go`
- `backend/internal/service/army_service_test.go`
- `backend/internal/service/resources_service.go`
- `backend/internal/service/mission_service.go`
- `backend/internal/service/mission_service_test.go`
- `backend/internal/http/handlers/mission_handler.go`
- `backend/internal/http/handlers/report_handler.go`
- `backend/internal/http/server.go`
- `frontend/src/api/client.ts`
- `frontend/src/api/errors.ts`
- `frontend/src/pages/DashboardPage.tsx`

## Commands Run

- `gofmt -w backend/internal/domain/mission.go backend/internal/domain/report.go backend/internal/gameconfig/missions.go backend/internal/repository/army_repository.go backend/internal/repository/mission_repository.go backend/internal/repository/report_repository.go backend/internal/service/army_service.go backend/internal/service/resources_service.go backend/internal/service/mission_service.go backend/internal/http/handlers/mission_handler.go backend/internal/http/handlers/report_handler.go backend/internal/http/server.go`
- `gofmt -w backend/internal/service/army_service_test.go backend/internal/service/mission_service_test.go`
- `cd backend && GOCACHE=/Users/andrey/Documents/pets/sumerki/.cache/go-build GOMODCACHE=/Users/andrey/Documents/pets/sumerki/.cache/go-mod go test ./...`
- `cd frontend && npm run typecheck`
- `cd frontend && npm run build`
- `DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable' make migrate-up`
- `DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable' make migrate-status`
- `cd backend && GOCACHE=/Users/andrey/Documents/pets/sumerki/.cache/go-build GOMODCACHE=/Users/andrey/Documents/pets/sumerki/.cache/go-mod DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable' JWT_SECRET='dev-secret' BACKEND_PORT=18080 go run ./cmd/server`
- `curl -i http://127.0.0.1:18080/ready`
- HTTP smoke test covering register, create kingdom, `GET /api/missions/available`, `POST /api/missions/start`, `GET /api/missions/me`, and `GET /api/army/me`
- Local DB timestamp adjustment to verify lazy mission completion
- HTTP smoke test covering `GET /api/reports/me`, completed mission state, returned units, and resource rewards
- Browser flow covering Missions/Reports dashboard rendering and starting a Black Forest mission from the UI
- `git diff --check`

## Verification

- `go test ./...` completed successfully.
- `npm run typecheck` completed successfully.
- `npm run build` completed successfully.
- Goose applied `00007_create_missions_and_reports.sql` successfully.
- Goose status shows migrations `00001` through `00007` applied.
- Verified backend readiness returned HTTP 200 with `{"status":"ready","database":"ok"}`.
- Verified `GET /api/missions/available` returns all 3 configured MVP missions.
- Verified `POST /api/missions/start` starts `black_forest_expedition` and immediately removes sent militia/scouts from available army.
- Verified `GET /api/missions/me` returns the active mission and sent unit allocation.
- Verified lazy mission completion by moving `finishes_at` into the past and reading reports.
- Verified completed mission has `completed` status, result payload, returned units, and reward data.
- Verified `GET /api/reports/me` creates and returns a basic `pve_mission` report after lazy completion.
- Verified resources increased by mission rewards.
- Verified army returned surviving units after completion.
- Verified the dashboard renders 3 mission cards and the reports section.
- Verified the dashboard can start a Black Forest mission and displays it as `В пути` with a finish time.
- Verified `git diff --check` has no whitespace errors.

Notes:

- Port `5432` was already in use locally, so live verification used `POSTGRES_PORT=15432` and `DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable'`.
- The backend was run on `18080` for verification, and the frontend was already running on `5173` with `VITE_API_BASE_URL=http://localhost:18080`.
- Combined shell-based HTTP smoke commands needed approved local-network execution.

## What Works Now

- Players can view available PvE missions.
- Players can start PvE missions with available units.
- Sent units are unavailable while a mission is active.
- Mission completion is resolved lazily by mission reads, report reads, or mission start.
- Completed missions return surviving units.
- Completed missions can permanently lose units according to simple MVP loss rules.
- Completed missions grant resource rewards after lazy resource recalculation.
- Completed missions create basic unread mission reports.
- Dashboard shows available missions, active/completed missions, and reports.

## Known Limitations

- Mission start, unit subtraction, mission row creation, mission unit creation, reward grants, report creation, and mission completion are not wrapped in one cross-repository database transaction.
- Mission outcomes use simple deterministic MVP rules; there is no advanced battle simulation.
- Failure is rare because mission start rejects below-minimum allocations.
- Reports are basic and cannot be marked read yet.
- Report text is intentionally simple; richer report polish is deferred.
- There is no separate locked-unit system beyond mission unit allocation records.
- No PvP raids, player-vs-player combat, route/pathfinding, heroes, equipment, unit XP, events, patrons, tribute, alliances, map, market, trading, payments, chat, or real-time systems were implemented.
- No background workers by design; resources, buildings, unit training, and missions use lazy resolution.

## Next Recommended Step

Start Phase 13 only when explicitly requested. Based on `docs/MVP_PHASES.md`, the next likely phase is Report Polish.
