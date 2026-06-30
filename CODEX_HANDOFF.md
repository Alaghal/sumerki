# Codex Handoff

## Current Phase

Phase 21: First Playtest Build.

## Status

Phase 21 is implemented according to the attached prompt. The MVP is now packaged for the first internal manual playtest with focused tester instructions, structured feedback, known limitations, release notes, playtest helper targets, and a visible Playtest 001 label in the frontend.

## Completed

- Added `docs/PLAYTEST_GUIDE.md` as the main first-playtest guide.
- Added `docs/FEEDBACK_TEMPLATE.md` for structured tester feedback.
- Added `docs/KNOWN_LIMITATIONS.md` to separate MVP limitations from bugs.
- Updated `docs/PLAYTEST_CHECKLIST.md` into a pre-flight, manual app, and UX checklist.
- Added `docs/RELEASE_NOTES_PLAYTEST_001.md`.
- Updated README with:
  - Phase 21 status
  - compact first playtest quick start
  - seed account list
  - playtest documentation links
- Added Makefile playtest helpers:
  - `playtest-setup`
  - `playtest-check`
  - `playtest-reset`
- Added a tiny frontend clarity label: `Playtest 001` in the top bar.

## Phase Scope Note

- No new gameplay systems were added.
- No schema or migration changes were made.
- No API response shapes were changed.
- No gameplay balance values were changed.
- No backend runtime systems were added.
- No production deployment, payments, analytics, admin panel, chat, alliances, map/province systems, dark god systems, NPC retaliation, background workers, cron jobs, Redis, or WebSocket were added.
- Frontend change was limited to a small playtest build label.
- This session did not start Phase 22.

## Changed Files

- `README.md`
- `CODEX_HANDOFF.md`
- `Makefile`
- `docs/PLAYTEST_GUIDE.md`
- `docs/FEEDBACK_TEMPLATE.md`
- `docs/KNOWN_LIMITATIONS.md`
- `docs/PLAYTEST_CHECKLIST.md`
- `docs/RELEASE_NOTES_PLAYTEST_001.md`
- `frontend/src/components/layout/TopBar.tsx`

## Commands Run

- `make test-backend`
- `make test-frontend`
- `make playtest-check`
- `DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable' make playtest-setup`
- `DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable' JWT_SECRET='dev-secret' BACKEND_PORT=18080 GOCACHE=/Users/andrey/Documents/pets/sumerki/.cache/go-build go run ./cmd/server`
- `API_BASE_URL='http://localhost:18080' make smoke-api`

## Verification

- Backend tests passed.
- Frontend typecheck passed.
- Frontend production build passed.
- `make playtest-check` passed and printed the expected smoke-api instruction.
- `make playtest-setup` ran successfully against the local Docker database on `15432`.
- Goose reported no pending migrations and current version `13`.
- `make seed-dev` ran as part of `playtest-setup` and refreshed all four dev accounts.
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

- A tester can use `docs/PLAYTEST_GUIDE.md` to understand the purpose and route of Playtest 001.
- A developer can run `make playtest-setup` to migrate and seed local playtest data.
- A developer can run `make playtest-check` for backend and frontend checks.
- A developer can run `make playtest-reset` for a documented local destructive reset and reseed.
- README has a compact First Playtest section and links to the detailed playtest docs.
- Testers have a structured feedback template.
- Current limitations are clearly documented.
- The dashboard top bar identifies the build as `Playtest 001`.

## Known Limitations

- Playtest 001 is local-development only.
- Seed data is local-only and uses public dev passwords.
- `reset-db` and `playtest-reset` are destructive and assume the local Docker PostgreSQL service/user/database from this repository.
- `scripts/smoke-api.sh` requires `jq`.
- Smoke API mutates local data by starting upgrades, training, missions, raids, and event choices.
- Smoke API checks happy paths but does not wait for mission or raid timers to resolve.
- Existing lazy mission/report duplicate behavior noted in Phase 19 was not changed.
- The local Docker database is currently published on `15432`, so verification used `DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable'`.

## Next Recommended Phase

Start Phase 22 only when explicitly requested.
