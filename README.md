# Sumerki

Sumerki is a browser strategy game MVP about building a small kingdom in a dusk-haunted world.

The first milestone is a playable vertical slice where a player can register, create one kingdom, manage a settlement, gather resources, upgrade buildings, train basic units, run PvE missions, perform simple PvP raids, receive reports, choose a patron, and react to simple events.

## Current Phase

Phase 9: Resources API + UI.

This phase adds stored kingdom resources, lazy base production, `GET /api/resources/me`, and real resource values on the dashboard. It does not include buildings, armies, missions, combat, spending, resource caps, or background workers.

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
3. Call `GET /api/resources/me` with the returned token if you want to inspect the resources API directly.
4. Wait a short time.
5. Call `GET /api/resources/me` again to see lazy production apply once whole units have accrued.
6. Open `/app` and see the real kingdom, ruler, and resource cards.
7. Logout.
8. Login again with the same account.
9. Return to the dashboard.

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
