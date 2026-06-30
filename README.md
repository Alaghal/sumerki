# Sumerki

Sumerki is a browser strategy game MVP about building a small kingdom in a dusk-haunted world.

The first milestone is a playable vertical slice where a player can register, create one kingdom, manage a settlement, gather resources, upgrade buildings, train basic units, run PvE missions, perform simple PvP raids, receive reports, choose a patron, and react to simple events.

## Current Phase

Phase 29: Responsive And Overflow Hardening.

Playtest 001 is complete. The frontend now has a Russian/English localization foundation, most current UI labels are localized, and `/app` uses a first game-shell layout with a symbolic local SVG map. The game shell now has responsive overflow hardening for narrow screens, dense panels, the HUD, mode navigation, context panels, the activity feed, and the local map.

## Documentation

- `AGENTS.md`: working instructions for coding agents.
- `CODEX_HANDOFF.md`: current implementation status and next step.
- `docs/MVP_SCOPE.md`: MVP goals, included features, and exclusions.
- `docs/MVP_PHASES.md`: phase-by-phase implementation plan.
- `docs/DECISIONS.md`: architecture and product decision log.
- `docs/API_CONTRACT.md`: draft backend API contract.
- `docs/DOMAIN_MODEL.md`: draft domain model.
- `docs/BALANCE.md`: first-pass MVP balance assumptions.
- `docs/SMOKE_TESTS.md`: local smoke test flow.
- `docs/PLAYTEST_GUIDE.md`: first internal playtest guide.
- `docs/FEEDBACK_TEMPLATE.md`: structured playtest feedback template.
- `docs/KNOWN_LIMITATIONS.md`: current MVP limitations.
- `docs/PLAYTEST_CHECKLIST.md`: manual browser playtest checklist.
- `docs/RELEASE_NOTES_PLAYTEST_001.md`: first playtest release notes.
- `docs/POST_PLAYTEST_ROADMAP.md`: planned UX/i18n roadmap after Playtest 001.
- `docs/UX_I18N_STRATEGY.md`: target game-shell UX and localization strategy.
- `docs/I18N_PLAN.md`: practical ru/en localization implementation plan.
- `docs/UI_COPY_RULES.md`: player-facing copy and enum display rules.
- `docs/phases/`: detailed phase notes.

## Planned Stack

- Backend: Go, Echo, PostgreSQL.
- Frontend: React, TypeScript, Vite, Tailwind.
- Local infrastructure: Docker Compose with PostgreSQL.

## Local Database

The local PostgreSQL service uses these defaults:

- database: `sumerki`
- user: `sumerki`
- password: `sumerki`
- port: `5432`

Start PostgreSQL:

```sh
docker compose up -d postgres
```

Check service status:

```sh
docker compose ps
```

Stop local infrastructure:

```sh
docker compose down
```

Equivalent Makefile shortcuts are available:

```sh
make db-up
make db-ps
make db-down
```

Reset the local Docker database:

```sh
make reset-db
```

`reset-db` is destructive and intended only for local development.

## Backend

The backend requires `DATABASE_URL` and `JWT_SECRET`. `BACKEND_PORT` defaults to `8080` when unset.

Run the backend:

```sh
cd backend
go mod tidy
DATABASE_URL="postgres://sumerki:sumerki@localhost:5432/sumerki?sslmode=disable" JWT_SECRET="dev-secret" go run ./cmd/server
```

Equivalent Makefile shortcuts are available:

```sh
make backend-tidy
make backend-run
```

Check process health:

```sh
curl http://localhost:8080/health
```

Check database readiness:

```sh
curl http://localhost:8080/ready
```

## Local MVP Verification

Start infrastructure, migrate, and seed dev data:

```sh
docker compose up -d postgres
make migrate-up
make seed-dev
```

Run checks:

```sh
make test-backend
make test-frontend
```

Run the backend:

```sh
make backend-run
```

In another terminal, run the API smoke test:

```sh
make smoke-api
```

Dev accounts use password `password123`:

- `northern@example.com`
- `lizard@example.com`
- `posad@example.com`
- `raider@example.com`

Dev seed data is local only. Do not use dev passwords in production.

See `docs/SMOKE_TESTS.md` and `docs/PLAYTEST_CHECKLIST.md` for the full local verification flow.

## First Playtest

Quick start:

```sh
docker compose up -d postgres
make playtest-setup
make test-backend
make test-frontend
make backend-run
```

In another terminal:

```sh
make smoke-api
cd frontend
npm install
npm run dev
```

Open the frontend URL printed by Vite, usually `http://localhost:5173`.

Seed accounts use password `password123`:

- `northern@example.com`
- `lizard@example.com`
- `posad@example.com`
- `raider@example.com`

Playtest docs:

- `docs/PLAYTEST_GUIDE.md`
- `docs/FEEDBACK_TEMPLATE.md`
- `docs/KNOWN_LIMITATIONS.md`
- `docs/PLAYTEST_CHECKLIST.md`
- `docs/RELEASE_NOTES_PLAYTEST_001.md`

Post-Playtest 001 UX/i18n docs:

- `docs/POST_PLAYTEST_ROADMAP.md`
- `docs/UX_I18N_STRATEGY.md`
- `docs/I18N_PLAN.md`
- `docs/UI_COPY_RULES.md`

## Auth API

Register:

```sh
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"player@example.com","password":"password123"}'
```

Login:

```sh
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"player@example.com","password":"password123"}'
```

Fetch current user:

```sh
curl http://localhost:8080/api/me \
  -H "Authorization: Bearer <token>"
```

## Kingdom API

Create a kingdom:

```sh
curl -X POST http://localhost:8080/api/kingdoms \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <TOKEN>" \
  -d '{"name":"Воронья Сечь","culture":"northern_principality"}'
```

Fetch current kingdom:

```sh
curl http://localhost:8080/api/kingdoms/me \
  -H "Authorization: Bearer <TOKEN>"
```

## Patron API

Get patron options:

```sh
curl http://localhost:8080/api/patron/options \
  -H "Authorization: Bearer <TOKEN>"
```

Get current patron:

```sh
curl http://localhost:8080/api/patron/me \
  -H "Authorization: Bearer <TOKEN>"
```

Join Old Pact:

```sh
curl -X POST http://localhost:8080/api/patron/join \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <TOKEN>" \
  -d '{"patron":"old_pact"}'
```

Break patron:

```sh
curl -X POST http://localhost:8080/api/patron/break \
  -H "Authorization: Bearer <TOKEN>"
```

Get patron pressure:

```sh
curl http://localhost:8080/api/patron/pressure \
  -H "Authorization: Bearer <TOKEN>"
```

Pay available tribute or contribution above protected reserves:

```sh
curl -X POST http://localhost:8080/api/patron/pay-tribute \
  -H "Authorization: Bearer <TOKEN>"
```

Ask for a delay during a pressure crisis:

```sh
curl -X POST http://localhost:8080/api/patron/crisis-choice \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <TOKEN>" \
  -d '{"choice":"ask_delay"}'
```

## Ruler API

Fetch current ruler:

```sh
curl http://localhost:8080/api/ruler/me \
  -H "Authorization: Bearer <TOKEN>"
```

## Resources API

Fetch current resources:

```sh
curl http://localhost:8080/api/resources/me \
  -H "Authorization: Bearer <TOKEN>"
```

## Buildings API

Fetch current buildings:

```sh
curl http://localhost:8080/api/buildings/me \
  -H "Authorization: Bearer <TOKEN>"
```

Upgrade a farm:

```sh
curl -X POST http://localhost:8080/api/buildings/farm/upgrade \
  -H "Authorization: Bearer <TOKEN>"
```

## Army API

Fetch current army:

```sh
curl http://localhost:8080/api/army/me \
  -H "Authorization: Bearer <TOKEN>"
```

Train militia:

```sh
curl -X POST http://localhost:8080/api/army/train \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <TOKEN>" \
  -d '{"unitType":"militia","amount":5}'
```

## Missions API

Fetch available missions:

```sh
curl http://localhost:8080/api/missions/available \
  -H "Authorization: Bearer <TOKEN>"
```

Start a mission:

```sh
curl -X POST http://localhost:8080/api/missions/start \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <TOKEN>" \
  -d '{"missionKey":"black_forest_expedition","units":[{"unitType":"militia","amount":5},{"unitType":"scouts","amount":1}]}'
```

Fetch current missions:

```sh
curl http://localhost:8080/api/missions/me \
  -H "Authorization: Bearer <TOKEN>"
```

Fetch reports:

```sh
curl "http://localhost:8080/api/reports/me?limit=20&offset=0" \
  -H "Authorization: Bearer <TOKEN>"
```

Fetch one report:

```sh
curl http://localhost:8080/api/reports/<REPORT_ID> \
  -H "Authorization: Bearer <TOKEN>"
```

Mark a report as read:

```sh
curl -X POST http://localhost:8080/api/reports/<REPORT_ID>/read \
  -H "Authorization: Bearer <TOKEN>"
```

## Raids API

Get neighbors:

```sh
curl http://localhost:8080/api/neighbors \
  -H "Authorization: Bearer <TOKEN>"
```

Start raid:

```sh
curl -X POST http://localhost:8080/api/raids/start \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <TOKEN>" \
  -d '{"defenderKingdomId":"<TARGET_KINGDOM_ID>","units":[{"unitType":"militia","amount":5},{"unitType":"scouts","amount":1}]}'
```

Get raids:

```sh
curl http://localhost:8080/api/raids/me \
  -H "Authorization: Bearer <TOKEN>"
```

## Events API

Phase 18 adds the first 25 event content templates. The same Phase 17 event curl flow can be used to test them.

Get current events:

```sh
curl http://localhost:8080/api/events/me \
  -H "Authorization: Bearer <TOKEN>"
```

Choose an event option:

```sh
curl -X POST http://localhost:8080/api/events/<EVENT_ID>/choose \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <TOKEN>" \
  -d '{"choiceKey":"sell_to_merchants"}'
```

## Frontend

The frontend reads `VITE_API_BASE_URL` and defaults to `http://localhost:8080`.

Run the frontend:

```sh
cd frontend
npm install
VITE_API_BASE_URL=http://localhost:8080 npm run dev
```

For a complete local account flow, use separate terminals:

```sh
docker compose up -d postgres
make migrate-up
make backend-run
```

Then start the frontend:

```sh
cd frontend
VITE_API_BASE_URL=http://localhost:8080 npm run dev
```

Open:

```sh
http://localhost:5173
```

Expected flow:

1. Register a new account.
2. Create a kingdom.
3. Create user B and kingdom B.
4. Login as user A.
5. Open `/app`.
6. View patron options and choose a patron if desired.
7. View patron pressure and pay tribute or ask for a delay when available.
8. Call `GET /api/events/me` or view dashboard events.
9. Choose an event option.
10. Verify resources, army, patron, or kingdom effects if applicable.
11. Verify an event report appears in reports.
12. Refresh events and confirm the resolved event is not applied again.
13. Call `GET /api/neighbors` or view the dashboard raid neighbors.
14. Start a raid against kingdom B.
15. Call `GET /api/army/me` and confirm sent units are unavailable.
16. Wait for the raid timer.
17. Call `GET /api/raids/me` to resolve lazy completion.
18. View reports for attacker.
19. Login as user B.
20. View defender report.
21. Confirm defender was not destroyed and protected resources remain.
22. Continue the settlement loop with resources, buildings, army, missions, patrons, events, and reports.

The frontend stores the MVP JWT in `localStorage` under `sumerki.auth.token`. Refreshing the page restores the session through `GET /api/me`.

If you only need to start the Vite server:

```sh
cd frontend
npm run dev
```

Expected local URL:

```sh
http://localhost:5173
```

Equivalent Makefile shortcuts are available:

```sh
make frontend-install
make frontend-dev
```

## Database Migrations

Start PostgreSQL before running migrations:

```sh
docker compose up -d postgres
```

Apply migrations:

```sh
make migrate-up
```

Check migration status:

```sh
make migrate-status
```

Rollback the latest migration:

```sh
make migrate-down
```

Reset local development migrations:

```sh
make migrate-reset
```

The migration commands use this default local database URL unless `DATABASE_URL` is provided:

```sh
postgres://sumerki:sumerki@localhost:5432/sumerki?sslmode=disable
```

## Phase Discipline

Each phase should be implemented independently and should update `CODEX_HANDOFF.md` before handoff.
