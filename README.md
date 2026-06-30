# Sumerki

Sumerki is a browser strategy game MVP about building a small kingdom in a dusk-haunted world.

The first milestone is a playable vertical slice where a player can register, create one kingdom, manage a settlement, gather resources, upgrade buildings, train basic units, run PvE missions, perform simple PvP raids, receive reports, choose a patron, and react to simple events.

## Current Phase

Phase 16: Tribute and Pressure.

This phase adds lazy patron pressure for the Empire of Dusk and Old Pact: tribute debt, soft contribution debt, crisis status, delay requests, tribute payment, and relation breaking. It does not include NPC retaliation, patron armies, alliances, map travel, real-time combat, payments, chat, WebSocket, or background workers.

## Documentation

- `AGENTS.md`: working instructions for coding agents.
- `CODEX_HANDOFF.md`: current implementation status and next step.
- `docs/MVP_SCOPE.md`: MVP goals, included features, and exclusions.
- `docs/MVP_PHASES.md`: phase-by-phase implementation plan.
- `docs/DECISIONS.md`: architecture and product decision log.
- `docs/API_CONTRACT.md`: draft backend API contract.
- `docs/DOMAIN_MODEL.md`: draft domain model.
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
8. Call `GET /api/neighbors` or view the dashboard raid neighbors.
9. Start a raid against kingdom B.
10. Call `GET /api/army/me` and confirm sent units are unavailable.
11. Wait for the raid timer.
12. Call `GET /api/raids/me` to resolve lazy completion.
13. View reports for attacker.
14. Login as user B.
15. View defender report.
16. Confirm defender was not destroyed and protected resources remain.
17. Continue the settlement loop with resources, buildings, army, missions, patrons, and reports.

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
