# Sumerki MVP Phases

## MVP Goal

Build the first playable vertical slice of Sumerki.

By the end of MVP, a player can:

* register and login
* create one kingdom
* see a basic settlement dashboard
* accumulate resources over time
* upgrade buildings
* train basic units
* send units to PvE missions
* perform simple PvP raids
* receive battle/mission/event reports
* choose a patron: Independent, Empire of Dusk, or Old Pact
* receive simple events and make choices

## Implementation Status

* Playtest 001 is complete.
* Phases 0 through 21 are implemented.
* The MVP is ready for first internal manual playtest.
* `CODEX_HANDOFF.md` remains the operational source of truth for the latest run details.
* Future phases after Playtest 001 are post-MVP or playtest-feedback driven.

## Phase Order Correction

The original plan placed Tribute and Pressure before Simple PvP Raids. Actual implementation followed the prompt-specified order where Simple PvP Raids with Protection landed before Tribute and Pressure.

This document now reflects the completed feature set and Playtest 001 order. Do not renumber code, migrations, or historical artifacts just to match this documentation cleanup.

## Out Of Scope For First MVP

Do not implement:

* alliances
* large province map
* seasonal world wars
* advanced diplomacy
* dark god avatar system
* trading between players
* market
* payments
* admin panel
* chat
* real-time combat
* mobile app
* complex animations
* Phaser/Pixi/Three.js

## Cross-Phase Rules

* Implement one phase at a time.
* Keep each phase small enough for Codex to execute and verify independently.
* Prefer earlier playable vertical slices over large late integration phases.
* Use incremental UI requirements inside each feature phase instead of saving dashboard work for a large polish phase.
* Resolve time-based resources, building upgrades, unit training, missions, raids, tribute, and events lazily during reads or commands.
* Do not add background workers for the MVP unless a later phase explicitly changes this plan.
* The resolved MVP resource list is: `gold`, `food`, `wood`, `stone`, `population`.

---

# Phase 0: Repository Bootstrap

Status: Completed

Created the initial repository documentation, agent rules, phase plan, API contract, domain model, decision log, handoff file, and phase notes folder.

---

# Phase 1: Infrastructure Foundation

Status: Completed

Added local Docker Compose PostgreSQL infrastructure, `.env.example`, Makefile database helpers, and README database commands.

---

# Phase 1b: Infrastructure Verification

Status: Completed

Verified local PostgreSQL startup, container status, connection settings, and local shutdown flow.

---

# Phase 2: Backend Skeleton

Status: Completed

Added the Go/Echo backend skeleton with config loading, PostgreSQL connection, request logging, recovery, local CORS, graceful shutdown, `GET /health`, and `GET /ready`.

---

# Phase 3: Database Migration Foundation

Status: Completed

Added Goose migration support and initial reversible migrations for `users` and `kingdoms`.

---

# Phase 4: Auth API

Status: Completed

Implemented email/password registration, login, JWT auth, `GET /api/me`, password hashing, auth middleware, and consistent JSON errors.

Endpoints:

* `POST /api/auth/register`
* `POST /api/auth/login`
* `GET /api/me`

---

# Phase 5: Kingdom Creation API

Status: Completed

Implemented one-kingdom-per-user creation and current kingdom lookup.

Endpoints:

* `POST /api/kingdoms`
* `GET /api/kingdoms/me`

---

# Phase 6: Ruler System

Status: Completed

Added rulers, generated a ruler during kingdom creation, and exposed the current ruler.

Endpoint:

* `GET /api/ruler/me`

---

# Phase 7: Frontend Skeleton

Status: Completed

Created the React/TypeScript/Vite/Tailwind frontend with login, register, kingdom creation, dashboard route, top bar, sidebar, and app shell.

---

# Phase 8: Frontend Auth, Kingdom, And Ruler Integration

Status: Completed

Connected register, login, logout, JWT session restore, protected routes, kingdom creation, and ruler dashboard display.

---

# Phase 9: Resources System

Status: Completed

Added `kingdom_resources`, lazy resource production, resource display, and the resource API.

Endpoint:

* `GET /api/resources/me`

---

# Phase 10: Buildings System

Status: Completed

Added MVP buildings, upgrade costs, lazy upgrade completion, production bonuses, building API, and dashboard building controls.

Endpoints:

* `GET /api/buildings/me`
* `POST /api/buildings/{type}/upgrade`

---

# Phase 11: Army System

Status: Completed

Added MVP units, unit training costs, barracks requirements, lazy training completion, army API, and dashboard training controls.

Endpoints:

* `GET /api/army/me`
* `POST /api/army/train`

---

# Phase 12: PvE Missions

Status: Completed

Added configured PvE missions, sent-unit tracking, lazy mission completion, resource rewards, losses, and basic reports.

Endpoints:

* `GET /api/missions/available`
* `GET /api/missions/me`
* `POST /api/missions/start`
* `GET /api/reports/me`

---

# Phase 13: Report Polish

Status: Completed

Improved reports with narrative phases, unread/read state, report detail, read marking, and support for later raid and event report types.

Endpoints:

* `GET /api/reports/me`
* `GET /api/reports/{id}`
* `POST /api/reports/{id}/read`

Report types:

* `pve_mission`
* `pvp_raid_attacker`
* `pvp_raid_defender`
* `event`

---

# Phase 14: Patron System

Status: Completed

Added basic patron selection and relation state for Independent, Empire of Dusk, and Old Pact. This phase intentionally did not include tribute, pressure, or military help.

Endpoints:

* `GET /api/patron/options`
* `GET /api/patron/me`
* `POST /api/patron/join`
* `POST /api/patron/break`

---

# Phase 15: Simple PvP Raids With Protection

Status: Completed

Added asynchronous PvP raids with lazy completion, neighbors, weak-player protection, newbie protection, same-target cooldown, defender protection, limited loot, protected resource minimums, dread gain, and attacker/defender reports.

Endpoints:

* `GET /api/neighbors`
* `GET /api/raids/me`
* `POST /api/raids/start`

Rules:

* defender cities and buildings cannot be destroyed
* population cannot be stolen
* loot is limited and cannot reduce defenders below protected minimums
* no territory capture, alliance war, NPC retaliation, or real-time combat

---

# Phase 16: Tribute And Pressure

Status: Completed

Added lazy patron pressure, surplus-based Empire tribute, Old Pact contribution pressure, protected minimums, pressure thresholds, payment, delay crisis choice, and patron break handling.

Endpoints:

* `GET /api/patron/pressure`
* `POST /api/patron/pay-tribute`
* `POST /api/patron/crisis-choice`

Rules:

* tribute and contribution are resolved lazily
* tribute spends only resources above protected minimums
* crisis choice supports asking for delay or breaking patron
* no NPC retaliation, patron armies, or background workers

---

# Phase 17: Event Engine

Status: Completed

Added reusable event templates, event choices, generated kingdom events, lazy event generation/expiry, choice resolution, event effects, and event reports.

Endpoints:

* `GET /api/events/me`
* `POST /api/events/{id}/choose`

Report type:

* `event`

---

# Phase 18: First Event Content Pack

Status: Completed

Added the first MVP-sized event content pack with 25 active event templates: five economy, five ruler, five military, five patron, and five dark omen events.

---

# Phase 19: Balance Pass

Status: Completed

Tuned first-session resources, starting army, PvE mission pacing, and early mission loss rates. Added `docs/BALANCE.md`.

---

# Phase 20: Smoke Tests And Seed Data

Status: Completed

Added local dev seed data, `seed-dev`, `smoke-api`, reset/test Makefile helpers, `docs/SMOKE_TESTS.md`, and a local playtest checklist.

Key helpers:

* `make seed-dev`
* `make smoke-api`
* `make test-backend`
* `make test-frontend`
* `make test-all`
* `make reset-db`

---

# Phase 21: First Playtest Build

Status: Completed

Prepared Playtest 001 with playtest instructions, feedback template, known limitations, release notes, playtest checklist updates, playtest Makefile helpers, and a small frontend `Playtest 001` label.

Key docs:

* `docs/PLAYTEST_GUIDE.md`
* `docs/FEEDBACK_TEMPLATE.md`
* `docs/KNOWN_LIMITATIONS.md`
* `docs/PLAYTEST_CHECKLIST.md`
* `docs/RELEASE_NOTES_PLAYTEST_001.md`

Key helpers:

* `make playtest-setup`
* `make playtest-check`
* `make playtest-reset`

---

# Final MVP Definition Of Done

Playtest 001 completion status:

* [x] user can register and login
* [x] user can create one kingdom
* [x] kingdom has a ruler
* [x] kingdom has resources
* [x] resources grow over time
* [x] player can upgrade buildings
* [x] player can train units
* [x] player can send PvE missions
* [x] player can receive reports
* [x] player can raid another player
* [x] player can choose patron
* [x] patron pressure exists for Empire and Old Pact paths
* [x] player can receive and resolve events
* [x] player cannot be fully destroyed by another player
* [x] README explains how to run and playtest the project
* [x] `CODEX_HANDOFF.md` is up to date

---

# Post-Playtest 001 Planned Phases

Phases 0 through 21 remain completed. Phases 22 through 31 are planned post-MVP UX/i18n phases. No new gameplay systems are planned until UX/i18n work is evaluated.

## Phase 22: UX And i18n Strategy Docs

Status: Completed

Documented the target map-first game shell, localization strategy, i18n implementation plan, and UI copy rules before frontend implementation.

## Phase 23: i18n Foundation ru/en

Status: Completed

Added i18next/react-i18next, Russian and English resources, default Russian language with English fallback, language persistence in `localStorage`, a TopBar language switcher, and initial shell/auth/kingdom/basic dashboard translation coverage.

Broader UI copy migration was completed in Phase 24.

## Phase 24: UI Copy Migration

Status: Completed

Migrated hardcoded player-facing UI copy to i18n, added gameplay translation namespaces, and moved DashboardPage gameplay labels for resources, buildings, units, missions, reports, patrons, events, raids, and ruler details to translations.

Full backend event/report narrative content localization remains Phase 30.

## Phase 25: Dashboard Decomposition

Status: Completed

Split the large DashboardPage render into focused feature panels for ruler, resources, buildings, army, missions, reports, patron, events, and raids.

DashboardPage remains the container for session access, API loading, state ownership, refresh functions, and action handlers. Gameplay behavior and i18n label handling were preserved.

Game Shell v1 remains Phase 26.

## Phase 26: Game Shell v1

Status: Completed

Added Game Shell v1 around the existing gameplay panels.

The `/app` screen now has a top resource HUD, local mode navigation, central scene placeholder, right context panel, and bottom activity feed. Existing gameplay behavior, API calls, and i18n handling were preserved.

The real local SVG map remains Phase 27.

## Phase 27: Local SVG Map v1

Status: Completed

Added Local SVG Map v1 to the Game Shell center scene.

The map shows the player settlement, Black Forest, Old Kurgan, Dry Ford, neighbor kingdoms, patron road, and omens/events. Map node clicks switch the local mode and right context panel using existing frontend state.

This is a local symbolic map only. It does not add a province system, territory capture, pathfinding, pan/zoom, new backend endpoints, or new gameplay mechanics.

Phase 28 remains context panel and activity feed refinement.

## Phase 28: Context Panels And Activity Feed

Status: Completed

Added node-specific context summaries above the existing mode panels.

The right context panel now responds to selected map nodes with compact city, mission, neighbor, patron, event, army, raid, and report summaries while preserving the existing gameplay panels and actions.

Improved the bottom activity feed with grouped actionable items for construction, training, missions, raids, events, reports, and patron pressure using existing frontend data. Feed items can switch local modes.

Responsive and overflow hardening remains Phase 29.

## Phase 29: Responsive And Overflow Hardening

Status: Planned

* Fix layout overflow.
* Ensure long text wraps.
* Ensure buttons wrap.
* Ensure panels scroll when needed.
* Test desktop, laptop, tablet-width, and narrow mobile-ish screens.
* No new game mechanics.

## Phase 30: Localized Event And Report Content

Status: Planned

* Localize event and report text.
* Prefer frontend translation by event/report keys where practical.
* Keep backend compatible.
* Add English text for current MVP event pack and report templates.
* Keep Russian as default.

## Phase 31: Playtest 002 UX Build

Status: Planned

* Package the UX/i18n improvements for a second playtest.
* Update playtest guide, release notes, known limitations, checklist, and feedback template.
* Focus on whether the game now feels like a strategy game rather than a dashboard.
