# Smoke Tests

Phase 20 makes the local MVP loop easy to verify without manual database edits.

## Prerequisites

- Docker Compose
- Go
- Node.js and npm
- `curl`
- `jq`

## Local Flow

Start PostgreSQL and apply migrations:

```sh
docker compose up -d postgres
make migrate-up
```

Seed predictable dev data:

```sh
make seed-dev
```

Run the backend:

```sh
make backend-run
```

In another terminal, run the API smoke test:

```sh
make smoke-api
```

Run automated checks:

```sh
make test-backend
make test-frontend
```

## Dev Accounts

All dev accounts use password `password123`.

- `northern@example.com`: Воронья Сечь, northern principality, Old Pact
- `lizard@example.com`: Тёплый Камень, lizard grad, Independent
- `posad@example.com`: Серый Посад, free posad, Empire of Dusk
- `raider@example.com`: Чёрный Брод, northern principality, no patron

Dev seed data is local only. Do not use these passwords outside local development.

## Resetting Local Data

`make reset-db` drops and recreates the local Docker database. It is destructive and intended only for local development.

After reset:

```sh
make migrate-up
make seed-dev
```

## API Smoke Script

`scripts/smoke-api.sh` uses `API_BASE_URL` and defaults to `http://localhost:8080`.

```sh
API_BASE_URL=http://localhost:8080 make smoke-api
```

The script verifies:

- auth
- kingdom
- ruler
- resources
- buildings
- army
- missions
- reports
- patron
- events
- raids

Optional actions blocked by game state print `NOTE` or `WARN` and continue. Broken required endpoints fail the script.

## Frontend Manual Check

Run the frontend:

```sh
cd frontend
npm install
npm run dev
```

Open the app and login as:

```text
northern@example.com / password123
```

Use `docs/PLAYTEST_CHECKLIST.md` for the manual browser pass.

## Common Issues

- If `make smoke-api` cannot connect, confirm the backend is running and `API_BASE_URL` matches the backend port.
- If `jq` is missing, install it before running the smoke script.
- If raid targets are unavailable, run `make seed-dev`; seeded kingdoms are aged beyond newbie protection for local raid smoke tests.
- If Goose download fails in an offline environment, rerun after Go module cache is available.
