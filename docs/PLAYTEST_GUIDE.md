# Sumerki First Playtest Guide

## Purpose

This is the first internal MVP playtest for Sumerki. It is meant to test the core loop, not judge a finished game.

Core loop:

- login/register
- create a kingdom
- inspect ruler
- watch resources
- upgrade a building
- train units
- send a mission
- read a report
- choose a patron
- resolve an event
- try a raid if available

## What To Test

Focus on:

- clarity
- pacing
- whether you know what to do next
- whether timers feel too slow or too fast
- whether resource costs feel fair
- whether mission reports are interesting
- whether raids feel threatening but not account-killing
- whether patron pressure feels understandable
- whether events feel atmospheric and meaningful

## What Not To Judge Yet

These areas are intentionally unfinished or out of scope for the first MVP:

- no alliances
- no province map
- no diplomacy
- no market
- no chat
- no payments
- no polished art
- no mobile optimization
- no advanced combat
- no dark god avatar system
- no NPC retaliation

## Recommended Playtest Length

Use a 20 to 40 minute first session.

An optional second check after mission, raid, or training timers complete is useful. Do not assume the game is running on production hosting or available in the background.

## Local Run Instructions

The full local setup is documented in `README.md` and `docs/SMOKE_TESTS.md`.

Compact setup:

```sh
docker compose up -d postgres
make migrate-up
make seed-dev
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

## Playtest Route

1. Login or register.
2. Create a kingdom if needed.
3. Look at the dashboard.
4. Upgrade one building.
5. Train militia or scouts.
6. Send a PvE mission.
7. Resolve an event.
8. Choose a patron.
9. Check patron pressure if available.
10. Try viewing neighbors.
11. Start a raid only if there is a valid target.
12. Read reports.
13. Refresh the page and check state persists.

## Feedback Format

Use `docs/FEEDBACK_TEMPLATE.md` for structured feedback.
