# Codex Handoff

## Current Phase

Phase 7: Frontend Auth + Kingdom Flow.

## Status

Phase 7 is complete according to the attached prompt: the frontend now connects to the existing backend auth and kingdom APIs for the first account flow.

## Completed

- Created initial repository guidance in `AGENTS.md`.
- Added `.gitignore` for future Go and TypeScript development artifacts.
- Created project overview in `README.md`.
- Documented MVP scope in `docs/MVP_SCOPE.md`.
- Preserved and aligned phase plan in `docs/MVP_PHASES.md`.
- Created decision log in `docs/DECISIONS.md`.
- Drafted initial API contract in `docs/API_CONTRACT.md`.
- Drafted initial domain model in `docs/DOMAIN_MODEL.md`.
- Created phase documentation directory at `docs/phases/`.
- Added Phase 0 detail page at `docs/phases/phase-00-repository-bootstrap.md`.
- Added Docker Compose configuration for local PostgreSQL.
- Added `.env.example` with local PostgreSQL defaults.
- Added basic Makefile commands for local database management.
- Added backend skeleton, health/readiness endpoints, auth API, and kingdom creation API in earlier phases.
- Added frontend app structure under `frontend/`.
- Added Vite, React, TypeScript, Tailwind, and React Router setup.
- Connected register and login forms to the backend.
- Added typed frontend API functions for auth, session, and kingdom calls.
- Added MVP JWT storage in `localStorage` under `sumerki.auth.token`.
- Added React session context for token, user, kingdom, loading, and error state.
- Added protected route behavior for `/app` and `/create-kingdom`.
- Added public route redirects for authenticated users.
- Connected the create kingdom form to the backend.
- Added logout from the top bar.
- Updated the dashboard to show the real kingdom name, culture, patron state, and user email.
- Kept resources, ruler, buildings, army, and reports as placeholders.
- Updated README with the full local account flow.

## Phase Order Note

- The attached prompt defines Phase 7 as `Frontend Auth + Kingdom Flow`.
- `docs/MVP_PHASES.md` currently defines Phase 7 as `Frontend Skeleton` and Phase 8 as `Frontend Auth, Kingdom, and Ruler Integration`.
- This session followed the attached prompt and did not modify `docs/MVP_PHASES.md`.
- Ruler integration was not implemented because the attached prompt explicitly excludes it.

## Changed Files

- `README.md`
- `CODEX_HANDOFF.md`
- `frontend/src/App.tsx`
- `frontend/src/api/client.ts`
- `frontend/src/api/errors.ts`
- `frontend/src/context/SessionContext.tsx`
- `frontend/src/components/layout/TopBar.tsx`
- `frontend/src/pages/CreateKingdomPage.tsx`
- `frontend/src/pages/DashboardPage.tsx`
- `frontend/src/pages/LoginPage.tsx`
- `frontend/src/pages/RegisterPage.tsx`
- `frontend/src/routes/AppRoutes.tsx`

No backend files, migrations, Docker Compose files, or project docs outside README and this handoff were modified.

## Commands Run

- `npm run typecheck`
- `npm run build`
- `docker compose up -d postgres`
- `POSTGRES_PORT=15432 docker compose up -d postgres`
- `DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable' make migrate-up`
- `DATABASE_URL='postgres://sumerki:sumerki@localhost:15432/sumerki?sslmode=disable' JWT_SECRET='dev-secret' BACKEND_PORT=18080 go run ./cmd/server`
- `curl -sS http://localhost:18080/health`
- `curl -sS -i http://localhost:18080/ready`
- `VITE_API_BASE_URL=http://localhost:18080 npm run dev -- --host 127.0.0.1 --port 5173`
- `npm install`

## Verification

- `npm install` completed successfully.
- `npm run typecheck` completed successfully.
- `npm run build` completed successfully.
- Verified backend health returned `{"status":"ok"}`.
- Verified backend readiness returned HTTP 200 with `{"status":"ready","database":"ok"}`.
- Verified in the browser:
  - register redirects to `/create-kingdom`
  - create kingdom redirects to `/app`
  - dashboard displays the real kingdom name
  - dashboard displays the real culture label
  - dashboard displays the user email
  - dashboard displays `Без покровителя` when patron is null
  - logout redirects to `/login`
  - `/app` without a session redirects to `/login`
  - login redirects to `/app`
  - page refresh restores the session and dashboard
- Verified no forbidden paths were changed with `git diff --name-only -- backend backend/migrations docker-compose.yml docs/MVP_SCOPE.md docs/MVP_PHASES.md docs/API_CONTRACT.md docs/DOMAIN_MODEL.md`.

Notes:

- Port `5432` was already in use locally, so live verification used `POSTGRES_PORT=15432`.
- The backend was run on `18080` for verification, and the frontend used `VITE_API_BASE_URL=http://localhost:18080`.
- Browser policy blocked setting an intentionally invalid token through a `javascript:` URL, so invalid-token behavior was not browser-forced. The frontend code clears session state for `invalid_token`, `expired_token`, and related auth errors from `/api/me` or `/api/kingdoms/me`.

## What Works Now

- `cd frontend && npm install` installs frontend dependencies.
- `cd frontend && npm run dev` starts the Vite development server.
- `cd frontend && npm run typecheck` typechecks the frontend.
- `cd frontend && npm run build` builds the frontend.
- Registering through `/register` calls `POST /api/auth/register` and stores the JWT.
- Logging in through `/login` calls `POST /api/auth/login` and stores the JWT.
- Refreshing the app restores the session through `GET /api/me` and fetches the kingdom through `GET /api/kingdoms/me`.
- Authenticated users without a kingdom are sent to `/create-kingdom`.
- Creating a kingdom calls `POST /api/kingdoms` and redirects to `/app`.
- Authenticated users with a kingdom are redirected away from `/login`, `/register`, and `/create-kingdom` to `/app`.
- Logout clears the local session and returns the user to `/login`.

## Known Limitations

- Dashboard still uses placeholder resources, ruler, buildings, army, and reports.
- No resources system exists yet.
- No gameplay systems exist yet.
- No ruler integration was implemented in this phase.
- `localStorage` token storage is MVP-only and should be revisited before production hardening.
- The frontend assumes the existing backend error shape documented in `docs/API_CONTRACT.md`.

## Next Recommended Step

Start the next explicitly requested phase only. Before more phase work, reconcile the prompt-driven phase order with `docs/MVP_PHASES.md` or continue noting the mismatch in handoffs.
