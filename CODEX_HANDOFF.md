# Codex Handoff

## Current Phase

Phase 20: Smoke Tests and Seed Data.

## Status

Phase 20 is implemented according to the attached prompt. Local development now has predictable seed data, Makefile verification helpers, an API smoke script, smoke-test documentation, and a manual playtest checklist.

## Completed

- Added `backend/cmd/seed-dev`.
- Added idempotent local dev seed data for:
  - `northern@example.com`
  - `lizard@example.com`
  - `posad@example.com`
  - `raider@example.com`
- Dev password for all seeded users is `password123`.
- Seeded kingdoms include rulers, resources, buildings, units, and patron state where selected.
- Seeded resources are:
  - gold: 2000
  - food: 1500
  - wood: 1500
  - stone: 1200
  - population: 250
- Seeded units are:
  - militia: 40
  - scouts: 10
  - spearmen: 15
  - archers: 15
  - cavalry: 5
- Seeded kingdoms are aged beyond newbie protection for local raid smoke tests.
- `northern@example.com` has boosted local buildings:
  - farm level 2
  - lumberyard level 2
  - quarry level 2
  - market level 2
  - barracks level 2
  - walls level 1
- Added `scripts/smoke-api.sh`.
- Added Makefile helpers:
  - `seed-dev`
  - `reset-db`
  - `smoke-api`
  - `test-backend`
  - `test-frontend`
  - `test-all`
- Added local Go build cache support in Makefile with `GOCACHE ?= $(CURDIR)/.cache/go-build`.
- Added `docs/SMOKE_TESTS.md`.
- Added `docs/PLAYTEST_CHECKLIST.md`.
- Updated README with local MVP verification commands and dev account notes.

## Phase Scope Note

- No new gameplay systems were added.
- No schema or migration changes were made.
- No API response shapes were changed.
- No frontend UI changes were made.
- No gameplay balance values were changed.
- No background workers, cron jobs, queues, Redis, WebSocket, analytics, production deployment, payments, admin panel, chat, alliances, or map/province systems were added.
- `reset-db` is intentionally destructive and documented as local Docker development only.
- This session did not start Phase 21.

## Changed Files

- `Makefile`
- `README.md`
- `CODEX_HANDOFF.md`
- `backend/cmd/seed-dev/main.go`
- `docs/SMOKE_TESTS.md`
- `docs/PLAYTEST_CHECKLIST.md`
- `scripts/smoke-api.sh`

## Commands Run

- `gofmt -w backend/cmd/seed-dev/main.go`
- `cd backend && GOCACHE=/Users/andrey/Documents/pets/sumerki/.cache/go-build go test ./...`
- `bash -n scripts/smoke-api.sh`
- `DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable' make seed-dev`
- `DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable' make seed-dev`
- `docker compose exec -T postgres psql -U sumerki -d sumerki -c "select email, count(*) ..."`
- `docker compose exec -T postgres psql -U sumerki -d sumerki -c "select k.name, k.patron, r.gold, ..."`
- `docker compose exec -T postgres psql -U sumerki -d sumerki -c "select u.email, unit_type, amount ..."`
- `docker compose exec -T postgres psql -U sumerki -d sumerki -c "select type, level ..."`
- `DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable' JWT_SECRET='dev-secret' BACKEND_PORT=18080 GOCACHE=/Users/andrey/Documents/pets/sumerki/.cache/go-build go run ./cmd/server`
- `API_BASE_URL='http://localhost:18080' make smoke-api`
- `make test-backend`
- `make test-frontend`
- `DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable' make migrate-status`

## Verification

- `backend/cmd/seed-dev` compiles.
- Backend tests passed.
- Frontend typecheck passed.
- Frontend production build passed.
- `scripts/smoke-api.sh` passed shell syntax check.
- `make seed-dev` completed successfully.
- Re-running `make seed-dev` returned the same dev user and kingdom IDs, confirming idempotent behavior for the dev accounts.
- SQL confirmed exactly one user row for each dev account.
- SQL confirmed seeded kingdoms have the expected resource values.
- SQL confirmed `northern@example.com` has the expected seeded unit amounts.
- SQL confirmed `northern@example.com` has boosted building levels.
- Goose status shows migrations `00001` through `00013` applied.
- Live API smoke test against `http://localhost:18080` completed successfully:
  - auth login
  - kingdom fetch
  - ruler fetch
  - resources fetch
  - building upgrade
  - army fetch
  - unit training
  - mission start
  - missions fetch
  - reports fetch
  - patron options/status/pressure
  - event choice
  - raids fetch
  - raid start

## What Works Now

- A developer can seed predictable local users and kingdoms with `make seed-dev`.
- Running `make seed-dev` twice does not duplicate the four dev users or kingdoms.
- Dev users can login with `password123`.
- Seeded kingdoms have enough resources and units for missions and raids.
- At least one seeded kingdom has patron pressure state.
- At least one seeded account can start a mission.
- At least one seeded account can view and choose events.
- At least one seeded raid target is available after seeding.
- `make smoke-api` verifies the core MVP API flow.
- `make test-backend`, `make test-frontend`, and `make test-all` provide convenient local checks.
- `make reset-db` provides a documented local destructive reset path.
- README and smoke docs explain the full local verification flow.

## Known Limitations

- Seed data is local-only and uses public dev passwords.
- `reset-db` is destructive and assumes the local Docker PostgreSQL service/user/database from this repository.
- `scripts/smoke-api.sh` requires `jq`.
- Smoke API checks endpoint availability and basic happy paths; it does not wait for mission or raid timers to resolve.
- Smoke API mutates local data by starting upgrades, training, missions, raids, and event choices.
- Existing lazy mission/report duplicate behavior noted in Phase 19 was not changed.
- The local Docker database is currently published on `15432`, so verification used `DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable'`.

## Next Recommended Phase

Start Phase 21 only when explicitly requested.
