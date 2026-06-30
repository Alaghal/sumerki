# Codex Handoff

## Current Phase

Phase 15: Simple PvP Raids with Protection.

## Status

Phase 15 is implemented according to the attached prompt: players can view neighbors, start asynchronous raids, have sent units removed until resolution, resolve raids lazily, steal limited protected resources, gain dread, and receive attacker/defender raid reports.

## Completed

- Added reversible Goose migration `00010_create_raids.sql`.
- Added `kingdoms.dread` and `kingdoms.honor`.
- Added `raids` and `raid_units` tables.
- Extended `mission_reports` type/result constraints for PvP raid reports.
- Added raid game config for duration, unit limits, loot caps, protected minimums, cooldowns, and protection windows.
- Added raid domain model.
- Added raid repository.
- Added raid service with:
  - neighbor discovery
  - weak-player/newbie/repeat/global protection checks
  - raid start validation
  - attacker unit subtraction
  - defender unit snapshot
  - lazy raid completion
  - deterministic score/result calculation
  - walls and patron defensive modifiers
  - limited loot transfer
  - attacker survivor return
  - dread gain
  - attacker and defender reports
- Added `GET /api/neighbors`.
- Added `GET /api/raids/me`.
- Added `POST /api/raids/start`.
- Existing report listing now also resolves completed raids before returning reports.
- Added frontend raid API types and calls.
- Added dashboard Neighbors/Raids UI with target selection, unit inputs, start action, active/completed raid display, and refresh.
- Updated reports UI labels to support PvP raid results.
- Updated README with raid curl examples and manual flow.
- Updated API contract and domain model docs.
- Added raid service unit tests for result thresholds, protected loot, and unit validation.

## Phase Order Note

- The attached prompt defines Phase 15 as `Simple PvP Raids with Protection`.
- `docs/MVP_PHASES.md` currently places `Tribute and Pressure` before `Simple PvP Raids`; this session followed the attached prompt and did not implement tribute/pressure.
- This session did not start Phase 16.

## Changed Files

- `README.md`
- `CODEX_HANDOFF.md`
- `docs/API_CONTRACT.md`
- `docs/DOMAIN_MODEL.md`
- `backend/migrations/00010_create_raids.sql`
- `backend/internal/domain/kingdom.go`
- `backend/internal/domain/raid.go`
- `backend/internal/gameconfig/raids.go`
- `backend/internal/repository/kingdom_repository.go`
- `backend/internal/repository/raid_repository.go`
- `backend/internal/repository/report_repository.go`
- `backend/internal/service/mission_service.go`
- `backend/internal/service/raid_service.go`
- `backend/internal/service/raid_service_test.go`
- `backend/internal/service/resources_service.go`
- `backend/internal/http/handlers/raid_handler.go`
- `backend/internal/http/server.go`
- `frontend/src/api/client.ts`
- `frontend/src/pages/DashboardPage.tsx`

## Commands Run

- `gofmt -w backend/internal/domain/raid.go backend/internal/gameconfig/raids.go backend/internal/repository/kingdom_repository.go backend/internal/repository/report_repository.go backend/internal/repository/raid_repository.go backend/internal/service/resources_service.go backend/internal/service/raid_service.go backend/internal/service/raid_service_test.go backend/internal/service/mission_service.go backend/internal/http/handlers/raid_handler.go backend/internal/http/server.go`
- `cd backend && GOCACHE=/Users/andrey/Documents/pets/sumerki/.cache/go-build GOMODCACHE=/Users/andrey/Documents/pets/sumerki/.cache/go-mod go test ./...`
- `cd frontend && npm run typecheck`
- `cd frontend && npm run build`
- `POSTGRES_PORT=15432 docker compose up -d postgres`
- `DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable' make migrate-up`
- `DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable' make migrate-status`
- Attempted to run local backend on `18080`; escalated approval timed out twice, and the sandboxed attempt failed with `listen tcp :18080: bind: operation not permitted`.

## Verification

- `go test ./...` completed successfully.
- `npm run typecheck` completed successfully.
- `npm run build` completed successfully.
- Goose applied `00010_create_raids.sql` successfully.
- Goose status shows migrations `00001` through `00010` applied.
- Live API smoke tests were not completed because the local backend process could not be started in this turn after approval timeouts and sandbox bind denial.

Notes:

- Port `5432` was already allocated locally in earlier work, so migration verification used `POSTGRES_PORT=15432` and `DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable'`.
- Defender unit losses are calculated and reported but not subtracted from defender `kingdom_units` in Phase 15.
- Raid completion idempotency is guarded by completing only `active` raids; reports are created after the completion update succeeds.

## What Works Now

- Players can view up to 20 neighboring kingdoms.
- Neighbor results expose culture, patron, dread, power estimate, `canRaid`, and blocked reason without exposing exact resources or unit counts.
- Players can start raids against valid targets.
- Sent attacker units are unavailable while the raid is active.
- Raids resolve lazily when raids or reports are read or when another raid starts.
- Completed raids return attacker survivors.
- Successful raids steal limited resources while respecting protected minimums.
- Population is never stolen.
- Defender city/buildings are never destroyed.
- Newbie protection, too-weak protection, same-target cooldown, and defender global protection are enforced on raid start.
- Raid completion creates attacker and defender reports.
- Dashboard shows neighbors, raid start controls, and current/completed raids.

## Known Limitations

- Raid combat is simple deterministic score comparison.
- Defender unit losses are report-only and are not applied to defender army yet.
- No territory capture.
- No city destruction.
- No alliance wars.
- No NPC retaliation.
- No patron military help beyond a small defensive score modifier.
- No revenge or bounty system.
- No map distance or travel routes.
- Raid start and completion are not wrapped in one cross-repository database transaction.
- No background workers by design; raid completion is lazy.

## Next Recommended Step

Start Phase 16 only when explicitly requested. Based on the attached execution order, the next likely phase is Tribute and Pressure or the next prompt-specified phase.
