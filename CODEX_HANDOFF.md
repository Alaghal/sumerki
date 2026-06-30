# Codex Handoff

## Current Phase

Phase 18: First Event Content Pack.

## Status

Phase 18 is implemented according to the attached prompt: the event catalog now has an MVP-sized content pack with 25 active event templates distributed evenly across the five event categories.

## Completed

- Added reversible Goose migration `00013_seed_first_event_content_pack.sql`.
- Kept the 5 Phase 17 smoke-test events.
- Added 20 new event templates:
  - 4 economy events
  - 4 ruler events
  - 4 military events
  - 4 patron events
  - 4 dark omen events
- Final active event catalog:
  - 5 economy events
  - 5 ruler events
  - 5 military events
  - 5 patron events
  - 5 dark omen events
  - 25 active events total
- Added two choices for every new event.
- Used only the existing Phase 17 effect schema:
  - `resourceDelta`
  - `unitDelta`
  - `kingdomDelta`
  - `patronFavorDelta`
- Used only existing simple patron conditions.
- Kept seeding idempotent with `ON CONFLICT` updates for events and choices.
- Updated README with Phase 18 status and event content note.
- Updated domain model docs with Phase 18 content-pack count and distribution.

## Phase Scope Note

- This phase added content only.
- No new event mechanics were added.
- No API response shapes were changed.
- No raid, mission, tribute, army, resource, building, auth, kingdom, or ruler logic was changed.
- This session did not start Phase 19.

## Changed Files

- `README.md`
- `CODEX_HANDOFF.md`
- `docs/DOMAIN_MODEL.md`
- `backend/migrations/00013_seed_first_event_content_pack.sql`

## Commands Run

- `DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable' make migrate-up`
- `docker compose exec -T postgres psql -U sumerki -d sumerki -c "select category, count(*) from game_events where is_active group by category order by category;"`
- `docker compose exec -T postgres psql -U sumerki -d sumerki -c "select count(*) as active_events from game_events where is_active;"`
- `docker compose exec -T postgres psql -U sumerki -d sumerki -c "select min(choice_count), max(choice_count), count(*) from (select g.event_key, count(c.id) as choice_count from game_events g left join event_choices c on c.game_event_id = g.id where g.is_active group by g.event_key) x;"`
- `docker compose exec -T postgres psql -U sumerki -d sumerki -c "select count(*) as duplicate_keys from (select event_key from game_events group by event_key having count(*) > 1) d;"`
- `DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable' make migrate-down`
- `DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable' make migrate-up`
- `cd backend && GOCACHE=/Users/andrey/Documents/pets/sumerki/.cache/go-build go test ./...`
- `cd frontend && npm run typecheck`
- `DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable' make migrate-status`
- `cd frontend && npm run build`
- `DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable' JWT_SECRET='dev-secret' BACKEND_PORT=18080 GOCACHE=/Users/andrey/Documents/pets/sumerki/.cache/go-build go run ./cmd/server`
- `curl -i http://localhost:18080/health`
- `curl -s -X POST http://localhost:18080/api/auth/register ...`
- `curl -s -X POST http://localhost:18080/api/kingdoms ...`
- `curl -s http://localhost:18080/api/events/me ...`
- `curl -s -X POST http://localhost:18080/api/events/<EVENT_ID>/choose ...`
- `curl -s http://localhost:18080/api/reports/me ...`

## Verification

- Goose applied `00013_seed_first_event_content_pack.sql` successfully.
- Goose rolled back and reapplied `00013_seed_first_event_content_pack.sql` successfully.
- Goose status shows migrations `00001` through `00013` applied.
- SQL count confirmed exactly 25 active events.
- SQL category count confirmed 5 active events in each category:
  - `dark_omen`
  - `economy`
  - `military`
  - `patron`
  - `ruler`
- SQL choice count confirmed every active event has exactly 2 choices.
- SQL duplicate check confirmed no duplicate `event_key` values.
- `GOCACHE=/Users/andrey/Documents/pets/sumerki/.cache/go-build go test ./...` completed successfully.
- `npm run typecheck` completed successfully.
- `npm run build` completed successfully.
- Live API smoke test:
  - `GET /health` returned 200.
  - Registered a smoke-test user.
  - Created a smoke-test kingdom.
  - `GET /api/events/me` returned 3 active events from the larger pool.
  - `POST /api/events/{id}/choose` resolved a new event and applied effects.
  - `GET /api/reports/me` returned an `event` report.

## What Works Now

- The event engine has enough content to feel alive in the MVP loop.
- There are 25 active event templates total.
- All five event categories are represented evenly.
- Every active event has two choices.
- New content uses concise Russian text with Sumerki frontier tone.
- New content uses only existing MVP-safe effects.
- Patron events use only existing simple patron conditions.
- Existing `GET /api/events/me` works with the larger event pool.
- Existing `POST /api/events/{id}/choose` works with the new events.
- Existing reports UI can display event reports.
- Existing dashboard Events UI can render the larger event pool.

## Known Limitations

- Event pack is still MVP-sized.
- No event chains.
- No recurring event chains.
- No advanced conditions.
- No map/province events.
- No alliance events.
- No dark god avatar mechanics.
- No NPC retaliation.
- No patron military help.
- No scheduled/background event generation by design.
- Effects remain intentionally small and not fully balanced.
- The local Docker database is currently published on `15432`, so verification used `DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable'`.

## Next Recommended Phase

Start Phase 19 only when explicitly requested.
