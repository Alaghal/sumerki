# Codex Handoff

## Current Phase

Phase 19: Balance Pass.

## Status

Phase 19 is implemented according to the attached prompt. The first-session MVP loop is tuned so a new player can upgrade an economic building, train units, run a short PvE mission, resolve an event, and inspect patron/raid sections without manual database edits.

## Completed

- Increased starting resources to:
  - gold: 600
  - food: 400
  - wood: 400
  - stone: 300
  - population: 120
- Increased starting scouts from 2 to 3 so `dry_ford_scouting` can be sent at its recommended unit count immediately.
- Shortened `dry_ford_scouting` from 90 seconds to 75 seconds.
- Reduced early mission loss percentages:
  - `black_forest_expedition`: 8% base, 12% max
  - `old_kurgan_expedition`: 14% base, 18% max
  - `dry_ford_scouting`: 3% base, 5% max
- Kept building costs, building timers, unit costs, unit timers, raid limits, tribute pressure, and event mechanics unchanged because they already fit the MVP target ranges closely.
- Added concise balance documentation in `docs/BALANCE.md`.
- Updated README current phase and documentation list.
- Updated the army service test expectation for the new starting unit count.

## Balance Areas Changed

- Starting resources.
- Starting scout count.
- PvE mission duration for the first scouting mission.
- PvE mission loss rates and caps.
- Balance documentation.

## Phase Scope Note

- No new gameplay systems were added.
- No schema or migration changes were made.
- No background workers, cron jobs, queues, Redis, WebSocket, or real-time combat were added.
- No auth, kingdom creation, API response shape, frontend route, raid mechanic, patron mechanic, event mechanic, or report mechanic was changed.
- Event seed data was reviewed but not modified; current event effects are already MVP-small enough for this pass.
- This session did not start Phase 20.

## Changed Files

- `README.md`
- `CODEX_HANDOFF.md`
- `docs/BALANCE.md`
- `backend/internal/gameconfig/resources.go`
- `backend/internal/gameconfig/units.go`
- `backend/internal/gameconfig/missions.go`
- `backend/internal/service/army_service_test.go`

## Commands Run

- `gofmt -w backend/internal/gameconfig/resources.go backend/internal/gameconfig/units.go backend/internal/gameconfig/missions.go backend/internal/service/army_service_test.go`
- `cd backend && GOCACHE=/Users/andrey/Documents/pets/sumerki/.cache/go-build go test ./...`
- `cd frontend && npm run typecheck`
- `cd frontend && npm run build`
- `DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable' make migrate-status`
- `DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable' JWT_SECRET='dev-secret' BACKEND_PORT=18080 GOCACHE=/Users/andrey/Documents/pets/sumerki/.cache/go-build go run ./cmd/server`
- `curl -s -X POST http://localhost:18080/api/auth/register ...`
- `curl -s -X POST http://localhost:18080/api/kingdoms ...`
- `curl -s http://localhost:18080/api/resources/me ...`
- `curl -s http://localhost:18080/api/army/me ...`
- `curl -s http://localhost:18080/api/missions/available ...`
- `curl -s http://localhost:18080/api/events/me ...`
- `curl -s -X POST http://localhost:18080/api/buildings/farm/upgrade ...`
- `curl -s -X POST http://localhost:18080/api/army/train ...`
- `curl -s -X POST http://localhost:18080/api/missions/start ...`
- `curl -s -X POST http://localhost:18080/api/events/<EVENT_ID>/choose ...`
- `curl -s http://localhost:18080/api/missions/me ...`
- `curl -s http://localhost:18080/api/reports/me ...`
- `curl -s http://localhost:18080/api/patron/options ...`
- `curl -s http://localhost:18080/api/neighbors ...`

## Verification

- Backend tests passed with local Go build cache.
- Frontend typecheck passed.
- Frontend production build passed.
- Goose status shows migrations `00001` through `00013` applied.
- Live API smoke test confirmed:
  - registered a new smoke-test user
  - created a new smoke-test kingdom
  - starting resources returned as `600/400/400/300/120`
  - production returned as `30/45/37/25/1` per hour with level-1 economic buildings
  - starting army returned as 10 militia and 3 scouts
  - available missions returned `dry_ford_scouting` with 75-second duration
  - farm upgrade started successfully
  - 3 militia training started successfully
  - `dry_ford_scouting` started with 3 scouts
  - event choice resolved and applied effects
  - mission lazily completed after its timer
  - mission report was created and returned by reports API
  - patron options were visible
  - raid neighbors were visible and protected by newbie protection

## What Works Now

- A new player can take multiple meaningful actions immediately after kingdom creation.
- Starting resources allow at least one basic economic upgrade and still leave room for unit training or barracks.
- Starting army can run the first scouting mission at the recommended scout count.
- The first scouting mission resolves quickly enough for local MVP testing.
- Early mission losses are less likely to feel devastating.
- Raid loot and raid protection remain limited and safe.
- Tribute and pressure remain surplus-based and protected by minimum reserves.
- Balance assumptions are documented in `docs/BALANCE.md`.

## Known Limitations

- Balance is first-pass MVP tuning.
- Values are not production-ready.
- No analytics-driven balancing exists yet.
- No long-term economy simulation exists yet.
- No monetization balancing exists yet.
- Playtest feedback is needed.
- Event effects were reviewed but not re-seeded in this phase.
- Patron copy still says some obligations are future-facing even though pressure endpoints exist; this is messaging polish, not balance logic.
- During live smoke, parallel reads of completed missions/reports produced duplicate mission reports. This appears to be an existing lazy-resolution race outside the Phase 19 balance scope.
- The local Docker database is currently published on `15432`, so verification used `DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable'`.

## Next Recommended Phase

Start Phase 20 only when explicitly requested.
