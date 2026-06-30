# Codex Handoff

## Current Phase

Phase 17: Event Engine.

## Status

Phase 17 is implemented according to the attached prompt: players can receive lazily generated events, choose one option, apply simple safe effects, receive event reports, and see resolved event state. The full event content pack is not implemented.

## Completed

- Added reversible Goose migration `00012_create_events.sql`.
- Added `game_events`, `event_choices`, and `kingdom_events` tables.
- Extended `mission_reports.type` to support `event`.
- Seeded 5 smoke-test events, one per category:
  - economy: `found_old_idol`
  - ruler: `ruler_bad_dream`
  - military: `volunteers_at_gate`
  - patron: `patron_envoy`
  - dark omen: `black_birds_over_walls`
- Added event domain, repository, service, and HTTP handler.
- Added lazy event expiry and deterministic lazy generation up to 3 active events per kingdom.
- Added MVP condition support for `requiresPatron`.
- Added event choice resolution with safe effects:
  - resource deltas
  - unit deltas
  - kingdom dread/honor deltas
  - patron favor deltas
- Added event reports with phases:
  - `Событие`
  - `Выбор`
  - `Последствия`
- Added endpoints:
  - `GET /api/events/me`
  - `POST /api/events/:id/choose`
- Added frontend event API types and calls.
- Added dashboard Events UI with active events, choice buttons, resolved events, loading, error, and empty states.
- Updated reports UI to label `event` reports.
- Updated README, API contract, and domain model docs.

## Phase Order Note

- `docs/MVP_PHASES.md` already lists Phase 17 as Event Engine.
- The prior phase-order mismatch remains historical: raids and tribute were implemented according to prior prompt order.
- This session followed the attached prompt and did not start Phase 18.

## Changed Files

- `README.md`
- `CODEX_HANDOFF.md`
- `docs/API_CONTRACT.md`
- `docs/DOMAIN_MODEL.md`
- `backend/migrations/00012_create_events.sql`
- `backend/internal/domain/event.go`
- `backend/internal/repository/event_repository.go`
- `backend/internal/repository/kingdom_repository.go`
- `backend/internal/repository/patron_repository.go`
- `backend/internal/service/event_service.go`
- `backend/internal/service/army_service.go`
- `backend/internal/service/resources_service.go`
- `backend/internal/http/handlers/event_handler.go`
- `backend/internal/http/handlers/kingdom_handler.go`
- `backend/internal/http/server.go`
- `frontend/src/api/client.ts`
- `frontend/src/api/errors.ts`
- `frontend/src/pages/DashboardPage.tsx`

## Commands Run

- `gofmt -w backend/internal/domain/event.go backend/internal/repository/event_repository.go backend/internal/service/event_service.go backend/internal/service/resources_service.go backend/internal/service/army_service.go backend/internal/repository/kingdom_repository.go backend/internal/repository/patron_repository.go backend/internal/http/handlers/event_handler.go backend/internal/http/handlers/kingdom_handler.go backend/internal/http/server.go`
- `cd frontend && npm run typecheck`
- `cd backend && GOCACHE=/Users/andrey/Documents/pets/sumerki/.cache/go-build go test ./...`
- `cd frontend && npm run build`
- `docker compose ps`
- `DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable' make migrate-up`
- `DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable' make migrate-down`
- `DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable' make migrate-up`
- `docker compose exec -T postgres psql -U sumerki -d sumerki -c "select category, event_key, count(c.id) as choices from game_events g join event_choices c on c.game_event_id = g.id group by category, event_key order by category, event_key;"`
- `docker compose exec -T postgres psql -U sumerki -d sumerki -c '\d kingdom_events'`
- `DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable' make migrate-status`
- `DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable' JWT_SECRET='dev-secret' BACKEND_PORT=18080 GOCACHE=/Users/andrey/Documents/pets/sumerki/.cache/go-build go run ./cmd/server`
- `curl -i http://localhost:18080/api/events/me`
- `curl -i http://localhost:18080/health`

## Verification

- `npm run typecheck` completed successfully.
- `GOCACHE=/Users/andrey/Documents/pets/sumerki/.cache/go-build go test ./...` completed successfully.
- `npm run build` completed successfully.
- Docker shows `sumerki-postgres-1` running and healthy on local port `15432`.
- Goose applied `00012_create_events.sql` successfully.
- Goose rolled back and reapplied `00012_create_events.sql` successfully.
- Goose status shows migrations `00001` through `00012` applied.
- SQL check confirmed 5 seeded events with 2 choices each.
- `psql \d kingdom_events` confirmed the instance table, status constraint, indexes, and foreign keys.
- Live smoke test:
  - `GET /health` returned 200 with `{"status":"ok"}`.
  - unauthenticated `GET /api/events/me` returned 401 with `missing_authorization_header`.

## What Works Now

- `GET /api/events/me` requires authentication.
- Event reads lazily expire stale active events.
- Event reads lazily generate up to 3 active events.
- Duplicate active event keys are avoided for the same kingdom.
- Event cooldown history prevents immediate regeneration.
- Patron-only event eligibility supports `requiresPatron`.
- `POST /api/events/:id/choose` requires authentication and scopes by authenticated user's kingdom.
- Choosing a valid active event applies configured effects once.
- Resources and units are clamped so they do not go below zero.
- Population is clamped so it does not go below one.
- Dread and honor do not go below zero.
- Patron favor changes clamp to -100..100 and are ignored when no active patron relation exists.
- Event reports are created with report type `event`.
- Dashboard shows active and resolved events and can resolve event choices.
- Dashboard refreshes events, resources, army, patron, and reports after an event choice.

## Known Limitations

- Only 5 smoke-test events are seeded.
- Full 20 to 30 event content pack is Phase 18.
- No event chains.
- No recurring event chains.
- No complex branching questlines.
- No advanced conditions beyond `requiresPatron`.
- No map/province events.
- No alliance events.
- No dark god avatar system.
- No NPC retaliation or patron military help.
- No scheduled/background event generation by design.
- Event reports are simple.
- Event effects and event resolution are not wrapped in one cross-repository database transaction.
- Report rewards/losses still use the existing simple report payloads; full applied effects are stored in the event result JSON.
- The local Docker database is currently published on `15432`, so verification used `DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable'`.

## Next Recommended Phase

Start Phase 18: First Event Content Pack only when explicitly requested.
