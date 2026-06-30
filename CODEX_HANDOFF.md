# Codex Handoff

## Current Phase

Phase 16: Tribute and Pressure.

## Status

Phase 16 is implemented according to the attached prompt: patron tribute and pressure now resolve lazily for the Empire of Dusk and Old Pact. Independent kingdoms have no tribute, debt, or pressure.

## Completed

- Added reversible Goose migration `00011_create_patron_pressure_states.sql`.
- Added `patron_pressure_states` table for tribute debt, Old Pact contribution debt, pressure level, crisis status, next tribute time, and delay state.
- Added patron pressure domain model, repository, game config, and service.
- Added protected-reserve resource spending helper for tribute and contribution payment.
- Integrated patron pressure lifecycle with patron join, patron break, and patron status reads.
- Added lazy pressure resolution for:
  - `GET /api/patron/pressure`
  - `POST /api/patron/pay-tribute`
  - `POST /api/patron/crisis-choice`
  - `GET /api/patron/me`
  - patron join and break flows
- Added `GET /api/patron/pressure`.
- Added `POST /api/patron/pay-tribute`.
- Added `POST /api/patron/crisis-choice`.
- Added dashboard patron pressure UI with debt, crisis state, protected reserves, pay action, delay request, and crisis break action.
- Updated README with patron pressure curl examples and local flow note.
- Updated API contract and domain model docs.

## Phase Order Note

- `docs/MVP_PHASES.md` lists Tribute and Pressure before Simple PvP Raids.
- The prior handoff recorded Simple PvP Raids with Protection as Phase 15 because that earlier prompt requested raids first.
- This session treated the requested Phase 16 as Tribute and Pressure and did not reimplement raids.
- This session did not start the next phase.

## Changed Files

- `README.md`
- `CODEX_HANDOFF.md`
- `docs/API_CONTRACT.md`
- `docs/DOMAIN_MODEL.md`
- `backend/migrations/00011_create_patron_pressure_states.sql`
- `backend/internal/domain/patron_pressure.go`
- `backend/internal/gameconfig/patron_pressure.go`
- `backend/internal/repository/patron_pressure_repository.go`
- `backend/internal/service/patron_pressure_service.go`
- `backend/internal/service/patron_service.go`
- `backend/internal/service/resources_service.go`
- `backend/internal/http/handlers/patron_pressure_handler.go`
- `backend/internal/http/server.go`
- `frontend/src/api/client.ts`
- `frontend/src/pages/DashboardPage.tsx`

## Commands Run

- `gofmt -w backend/internal/domain/patron_pressure.go backend/internal/gameconfig/patron_pressure.go backend/internal/repository/patron_pressure_repository.go backend/internal/service/patron_service.go backend/internal/service/patron_pressure_service.go backend/internal/service/resources_service.go backend/internal/http/handlers/patron_pressure_handler.go backend/internal/http/server.go`
- `cd frontend && npm run typecheck`
- `cd backend && go test ./...` (failed because the sandbox could not write to `/Users/andrey/Library/Caches/go-build`)
- `cd backend && GOCACHE=/Users/andrey/Documents/pets/sumerki/.cache/go-build go test ./...`
- `cd frontend && npm run build`
- `docker compose ps`
- `DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable' make migrate-up`
- `docker compose exec -T postgres psql -U sumerki -d sumerki -c '\d patron_pressure_states'`
- `DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable' make migrate-down`
- `DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable' make migrate-up`
- `DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable' make migrate-status`

## Verification

- `npm run typecheck` completed successfully.
- `GOCACHE=/Users/andrey/Documents/pets/sumerki/.cache/go-build go test ./...` completed successfully.
- `npm run build` completed successfully.
- Docker shows `sumerki-postgres-1` running and healthy on local port `15432`.
- Goose applied `00011_create_patron_pressure_states.sql` successfully.
- Goose rolled back and reapplied `00011_create_patron_pressure_states.sql` successfully.
- Goose status shows migrations `00001` through `00011` applied.
- `psql \d patron_pressure_states` confirmed the table, unique kingdom row, foreign key to `kingdoms(id)`, nonnegative debt constraints, patron constraint, pressure range constraint, and crisis status constraint.

## What Works Now

- Empire of Dusk pressure resolves lazily over time.
- Empire tribute is paid from surplus gold and food above protected reserves.
- Unpaid Empire tribute becomes debt and can raise pressure.
- Old Pact tracks a soft food contribution debt with low capped pressure.
- Independent kingdoms have no pressure, debt, tribute, or next tribute time.
- Players can view current pressure, debt, available actions, protected minimums, and next tribute time.
- Players can pay tribute/contribution when debt exists.
- Players can ask for a delay during warning/active pressure.
- Players can break patron relation through the pressure crisis endpoint.
- Dashboard exposes the basic pressure loop without background workers.

## Known Limitations

- No NPC raids or retaliation.
- No patron military help or Old Pact troop support.
- No tribute events or event engine integration.
- No vassalage, city seizure, land transfer, alliances, large map, market, trading, payments, chat, WebSocket, or real-time combat.
- Patron pressure updates are not wrapped in cross-repository database transactions.
- The local Docker database is currently published on `15432`, so verification used `DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable'` instead of the Makefile default `5432`.
- Live authenticated API smoke tests were not run in this turn; verification covered compile/tests/build and database migration/schema.

## Next Recommended Phase

Start the next prompt-specified phase only when explicitly requested. Based on the current implemented features, the next likely area is Event Engine or Report Polish, depending on the revised phase plan.
