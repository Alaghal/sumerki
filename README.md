# Sumerki

Sumerki is a browser strategy game MVP about building a small kingdom in a dusk-haunted world.

The first milestone is a playable vertical slice where a player can register, create one kingdom, manage a settlement, gather resources, upgrade buildings, train basic units, run PvE missions, perform simple PvP raids, receive reports, choose a patron, and react to simple events.

## Current Phase

Phase 3: Database Migration Foundation.

This phase adds Goose SQL migrations for the initial `users` and `kingdoms` schema. It does not include auth APIs, kingdom APIs, frontend code, resources, rulers, combat, events, patrons, or gameplay systems.

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

The backend requires `DATABASE_URL`. `BACKEND_PORT` defaults to `8080` when unset. `JWT_SECRET` is loaded for future auth work but is not used yet.

Run the backend:

```sh
cd backend
go mod tidy
DATABASE_URL="postgres://sumerki:sumerki@localhost:5432/sumerki?sslmode=disable" go run ./cmd/server
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
