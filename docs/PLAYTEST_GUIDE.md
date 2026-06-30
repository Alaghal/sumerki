# Sumerki Playtest 002 Guide

## Purpose

Playtest 002 evaluates UX and game feel after Game Shell, Local SVG Map, context panels, activity feed, responsive hardening, and Russian/English localization.

Main question:

Does Sumerki now feel more like a map-first strategy-lite game instead of a dashboard of systems?

## What To Test

- Does the main screen feel like a strategy game?
- Is the resource HUD useful?
- Is the local map understandable?
- Do map nodes make sense?
- Does clicking nodes show useful context?
- Does mode navigation make sense?
- Is the activity feed useful?
- Does RU/EN switching work?
- Does English feel readable?
- Does Russian still feel atmospheric?
- Does the UI avoid horizontal overflow on your screen?
- Are reports/events understandable?
- Can you still complete the core gameplay loop?

## What Not To Judge Yet

- final art direction
- full mobile UX
- province strategy
- alliances
- diplomacy depth
- production deployment
- economy balance as final

## Recommended Playtest Length

Use a 20 to 40 minute session.

An optional second check after mission, raid, or training timers complete is useful. Do not assume the game is running on production hosting or available in the background.

## Local Run Instructions

The full local setup is documented in `README.md` and `docs/SMOKE_TESTS.md`.

Compact setup:

```sh
docker compose up -d postgres
make playtest-setup
make backend-run
```

In another terminal:

```sh
cd frontend
npm install
npm run dev
```

Open the frontend URL printed by Vite, usually:

```text
http://localhost:5173
```

## Test Accounts

All seeded accounts use password `password123`.

- `northern@example.com`
- `lizard@example.com`
- `posad@example.com`
- `raider@example.com`

## Suggested Playtest Route

1. Start local environment.
2. Login with a seeded account or create a new account.
3. Open `/app`.
4. Switch language to Russian, then English, then back.
5. Look at the HUD and say what each resource means.
6. Click the settlement node.
7. Click each mission node.
8. Start a mission if possible.
9. Click neighbor nodes.
10. Try to understand whether a raid is possible.
11. Click patron road.
12. Check patron pressure.
13. Click omens/events.
14. Resolve an event.
15. Open reports.
16. Read at least one report.
17. Use activity feed to switch modes.
18. Resize the browser or test a narrow screen if practical.
19. Continue the normal loop: building, training, mission, event, report, raid.

## Feedback Format

Use `docs/FEEDBACK_TEMPLATE.md` for structured feedback.
