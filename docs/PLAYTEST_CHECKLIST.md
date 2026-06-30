# Playtest 002 Checklist

Use this checklist for the second internal manual playtest.

## Pre-Flight

- [ ] `docker compose up -d postgres` starts PostgreSQL.
- [ ] `make playtest-setup` applies migrations and seeds data.
- [ ] `make test-backend` passes.
- [ ] `make test-frontend` passes.
- [ ] Backend starts with `make backend-run`.
- [ ] Frontend starts with `make frontend-dev`.
- [ ] `make smoke-api` passes if backend is running.

## Manual App Checklist

- [ ] Login works.
- [ ] Logout works.
- [ ] Refresh keeps the session.
- [ ] Create kingdom works for a new account.
- [ ] `/app` opens Game Shell.
- [ ] HUD is visible.
- [ ] Resources are visible.
- [ ] Mode navigation works.
- [ ] Map is visible.
- [ ] All map nodes are clickable.
- [ ] Settlement node opens city context.
- [ ] Mission nodes open missions context.
- [ ] Neighbor nodes open raids context.
- [ ] Patron road opens patron context.
- [ ] Omens/events opens events context.
- [ ] Activity feed is visible.
- [ ] Feed items switch modes where supported.
- [ ] RU/EN switcher works.
- [ ] Selected language persists after refresh.
- [ ] Event title/body/choices localize where stable keys exist.
- [ ] Report shell/template labels localize where stable keys exist.
- [ ] No raw i18n keys are visible.
- [ ] No raw enum labels are visible as player-facing labels.

## Responsive Checks

- [ ] No horizontal overflow at 1280px.
- [ ] No horizontal overflow at 768px.
- [ ] No horizontal overflow at 390px.
- [ ] HUD wraps safely.
- [ ] Mode navigation wraps safely.
- [ ] Map remains usable.
- [ ] Context panel remains usable.
- [ ] Dense forms/buttons remain usable.

## Core Gameplay Loop

- [ ] Ruler is visible.
- [ ] Building upgrade works.
- [ ] Resources update after spending.
- [ ] Unit training works.
- [ ] Trained units appear after lazy refresh.
- [ ] Mission start works.
- [ ] Sent units become unavailable.
- [ ] Mission resolves.
- [ ] Report appears.
- [ ] Report can be read and marked read.
- [ ] Patron options are visible.
- [ ] Join patron works.
- [ ] Break patron works.
- [ ] Pressure is visible if applicable.
- [ ] Patron pressure action works if available.
- [ ] Events are visible.
- [ ] Choosing an event applies effects.
- [ ] Event report appears.
- [ ] Neighbors are visible.
- [ ] Blocked raid reasons make sense.
- [ ] Raid can start when valid.
- [ ] Raid resolves.
- [ ] Attacker and defender reports appear.

## UX Checks

- [ ] No blank screens.
- [ ] Loading states show.
- [ ] Errors are readable.
- [ ] Buttons are disabled when action is impossible.
- [ ] Labels are understandable.
- [ ] Main screen feels more like a strategy game than a dashboard.
