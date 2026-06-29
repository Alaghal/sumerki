# Codex Handoff

## Current Phase

Phase 14: Patron System v1.

## Status

Phase 14 is complete according to the attached prompt: players can view patron options, view current patron status, join a patron, switch patron, break patron relation, and see patron state on the dashboard.

## Completed

- Added reversible Goose migration `00009_create_patron_relations.sql`.
- Added `patron_relations` table with patron, favor, standing, joined/left timestamps, and constraints.
- Added migration backfill from non-null `kingdoms.patron`.
- Added patron game config for:
  - `independent`
  - `empire_of_dusk`
  - `old_pact`
- Added patron domain model.
- Added patron repository with find, upsert, break, and backfill support.
- Added minimal kingdom repository update for `kingdoms.patron`.
- Added patron service with options, current status, join, switch, break, ownership scoping, and validation.
- Added patron HTTP handler.
- Added `GET /api/patron/options`.
- Added `GET /api/patron/me`.
- Added `POST /api/patron/join`.
- Added `POST /api/patron/break`.
- Added frontend patron API types and calls.
- Added dashboard Patron section with loading/error states, current patron status, available patron cards, join/select buttons, break relation button, and v1 wording.
- Updated dashboard kingdom patron display to use the Russian patron label.
- Updated README with patron curl examples and local manual flow.
- Updated API contract and domain model docs.
- Added patron service tests for options, no-kingdom handling, invalid patron, join, idempotent join, switching, kingdom patron updates, break idempotency, backfill, and ownership scoping.

## Phase Order Note

- The attached prompt defines Phase 14 as `Patron System v1`.
- `docs/MVP_PHASES.md` includes `POST /api/patron/request-help` in Phase 14, but the attached prompt explicitly excludes military help and patron requests. This session followed the prompt and did not implement request-help.
- This session did not start Phase 15.

## Changed Files

- `README.md`
- `CODEX_HANDOFF.md`
- `docs/API_CONTRACT.md`
- `docs/DOMAIN_MODEL.md`
- `backend/migrations/00009_create_patron_relations.sql`
- `backend/internal/domain/patron.go`
- `backend/internal/gameconfig/patrons.go`
- `backend/internal/repository/kingdom_repository.go`
- `backend/internal/repository/patron_repository.go`
- `backend/internal/service/patron_service.go`
- `backend/internal/service/patron_service_test.go`
- `backend/internal/http/handlers/patron_handler.go`
- `backend/internal/http/server.go`
- `frontend/src/api/client.ts`
- `frontend/src/pages/DashboardPage.tsx`

## Commands Run

- `gofmt -w backend/internal/domain/patron.go backend/internal/gameconfig/patrons.go backend/internal/repository/patron_repository.go backend/internal/repository/kingdom_repository.go backend/internal/service/patron_service.go backend/internal/service/patron_service_test.go backend/internal/http/handlers/patron_handler.go backend/internal/http/server.go`
- `cd backend && GOCACHE=/Users/andrey/Documents/pets/sumerki/.cache/go-build GOMODCACHE=/Users/andrey/Documents/pets/sumerki/.cache/go-mod go test ./...`
- `cd frontend && npm run typecheck`
- `cd frontend && npm run build`
- `POSTGRES_PORT=15432 docker compose up -d postgres`
- `DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable' make migrate-up`
- `DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable' make migrate-status`
- `lsof -nP -iTCP:18080 -sTCP:LISTEN`
- `cd backend && DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable' JWT_SECRET='dev-secret' BACKEND_PORT=18080 GOCACHE=/Users/andrey/Documents/pets/sumerki/.cache/go-build GOMODCACHE=/Users/andrey/Documents/pets/sumerki/.cache/go-mod go run ./cmd/server`
- `curl` smoke tests for unauthenticated patron access, register, patron status without kingdom, create kingdom, patron options, patron status, invalid join, valid join, idempotent join, patron switch, break, repeated break, and kingdom patron readback.

## Verification

- `go test ./...` completed successfully.
- `npm run typecheck` completed successfully.
- `npm run build` completed successfully.
- Goose applied `00009_create_patron_relations.sql` successfully.
- Goose status shows migrations `00001` through `00009` applied.
- Verified unauthenticated `GET /api/patron/options` returns 401.
- Verified authenticated user without kingdom gets `kingdom_not_found` from `GET /api/patron/me`.
- Verified `GET /api/patron/options` returns all 3 patrons.
- Verified `GET /api/patron/me` returns `patron: null` for a kingdom without patron.
- Verified `POST /api/patron/join` rejects invalid patron with `invalid_patron`.
- Verified `POST /api/patron/join` joins `old_pact` and updates `kingdom.patron`.
- Verified joining the current patron succeeds.
- Verified joining `empire_of_dusk` switches the existing relation and updates `kingdom.patron`.
- Verified `POST /api/patron/break` clears patron response and `kingdom.patron`.
- Verified repeated `POST /api/patron/break` succeeds and remains null.

Notes:

- Port `5432` was already allocated locally, so live verification used `POSTGRES_PORT=15432` and `DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable'`.
- The backend was run on `18080` for verification.
- Local API and Docker commands needed approved local-process/local-network execution.

## What Works Now

- Players can view the three MVP patron options.
- Players can view current patron status.
- Players with no selected patron see `patron: null`.
- Players can join `independent`, `empire_of_dusk`, or `old_pact`.
- Joining the current patron is idempotent.
- Joining a different patron switches the MVP relation.
- Joining updates `kingdoms.patron`.
- Players can break patron relation.
- Breaking clears `kingdoms.patron`.
- Dashboard shows current patron label, status, effects, future effects, and v1 limitations.

## Known Limitations

- Patron choice has no mechanical tribute yet.
- No military help is implemented yet.
- No PvP protection is implemented yet.
- No patron quests or patron orders are implemented yet.
- No cooldown or penalty for switching patrons exists yet.
- No debt, pressure, contribution, NPC raids, punitive raids, patron resources, or timers exist yet.
- No alliances, world map, dark gods, payments, chat, events, diplomacy between players, Redis, cron jobs, or background workers were implemented.
- Tribute and pressure are planned for a later phase.

## Next Recommended Step

Start Phase 15 only when explicitly requested. Based on `docs/MVP_PHASES.md`, the next likely phase is Tribute and Pressure.
