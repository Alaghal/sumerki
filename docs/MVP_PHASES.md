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
* receive battle/mission reports
* choose a patron: Independent, Empire of Dusk, or Old Pact
* receive simple events and make choices

## Out of scope for first MVP

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

## Goal

Prepare the empty repository.

## Scope

Create:

* `AGENTS.md`
* `CODEX_HANDOFF.md`
* `README.md`
* `docs/MVP_SCOPE.md`
* `docs/MVP_PHASES.md`
* `docs/DECISIONS.md`
* `docs/API_CONTRACT.md`
* `docs/DOMAIN_MODEL.md`
* `docs/phases/`

## Result

Repository has project rules, MVP scope, and phase documentation.

---

# Phase 1: Infrastructure Foundation

## Goal

Prepare local infrastructure.

## Scope

Add:

* Docker Compose
* PostgreSQL
* `.env.example`
* basic Makefile commands

Services:

* `postgres`

## Result

Developer can start PostgreSQL locally.

---

# Phase 1b: Infrastructure Verification

## Goal

Confirm local PostgreSQL startup before backend work begins.

## Scope

Verify:

* `docker compose up -d postgres`
* `docker compose ps`
* PostgreSQL container is running
* PostgreSQL accepts connections with database `sumerki`, user `sumerki`, password `sumerki`, and port `5432`
* `docker compose down`

Update:

* `CODEX_HANDOFF.md` with verification results

## Result

Infrastructure is proven locally and Phase 2 can rely on a working database.

---

# Phase 2: Backend Skeleton

## Goal

Create minimal backend service.

## Scope

Backend stack:

* Go
* Echo
* PostgreSQL

Create:

* backend folder structure
* config loading
* database connection
* health endpoint

Endpoint:

* `GET /health`

## Result

Backend starts locally and returns health status.

---

# Phase 3: Database Migration Foundation

## Goal

Create the first backend schema before frontend implementation.

## Scope

Add migration tooling and initial reversible migrations.

Tables:

* `users`
* `kingdoms`

Rules:

* one user can have one kingdom
* UUID primary keys
* reversible migrations
* local migration commands documented in README

## Result

Database can store users and kingdoms, and the backend/auth/kingdom flow can be implemented against real schema.

---

# Phase 4: Auth API

## Goal

Allow users to register and login.

## Scope

Endpoints:

* `POST /api/auth/register`
* `POST /api/auth/login`
* `GET /api/me`

Features:

* password hashing
* JWT auth
* auth middleware
* clean JSON errors

## Result

Player can register, login, and fetch current user through the backend.

---

# Phase 5: Kingdom Creation API

## Goal

Allow authenticated player to create a kingdom.

## Scope

Endpoints:

* `POST /api/kingdoms`
* `GET /api/kingdoms/me`

Cultures:

* `northern_principality`
* `lizard_grad`
* `free_posad`

Rules:

* one kingdom per user
* kingdom name length: 3 to 32
* culture must be valid

## Result

Player can create and fetch their kingdom through the backend.

---

# Phase 6: Ruler System

## Goal

Add ruler identity immediately after the account and kingdom loop.

## Scope

Add table:

* `rulers`

Ruler fields:

* name
* age
* culture
* authority
* courage
* cunning
* honor
* cruelty
* ambition
* paranoia
* health_status

Rules:

* ruler is generated when kingdom is created
* ruler belongs to exactly one kingdom
* initial ruler generation can be deterministic or simple random MVP logic

Endpoint:

* `GET /api/ruler/me`

## Result

Player has a ruler attached to the first playable account and kingdom loop.

---

# Phase 7: Frontend Skeleton

## Goal

Create minimal frontend app after the backend account, kingdom, and ruler APIs exist.

## Scope

Frontend stack:

* React
* TypeScript
* Vite
* Tailwind

Create pages:

* `/login`
* `/register`
* `/create-kingdom`
* `/app`

Create layout:

* top bar
* sidebar
* main content area

## Result

Frontend starts locally and shows basic pages ready to connect to backend APIs.

---

# Phase 8: Frontend Auth, Kingdom, and Ruler Integration

## Goal

Connect the first playable account and kingdom flow to the frontend.

## Scope

Implement:

* register form
* login form
* logout
* JWT storage
* protected routes
* session restore via `GET /api/me`
* kingdom creation form
* culture select with short descriptions
* ruler card on dashboard

Routing rules:

* no token -> `/login`
* token but no kingdom -> `/create-kingdom`
* token and kingdom exists -> `/app`

Dashboard must show:

* kingdom name
* culture
* ruler summary
* placeholder resources area
* clear loading, error, and empty states

## Result

Playable vertical slice:

register -> login -> create kingdom -> view dashboard with ruler.

---

# Phase 9: Resources System

## Goal

Add basic resource economy.

## Scope

Resources:

* gold
* food
* wood
* stone
* population

Add table:

* `kingdom_resources`

Implement lazy passive production calculation.

Rules:

* do not use background workers
* calculate resource changes from `last_calculated_at`
* resolve resources during resource reads and commands that spend or grant resources
* population can grow or cap production, but must stay simple for MVP

Endpoint:

* `GET /api/resources/me`

Frontend:

* replace placeholder resources with live resource values
* show enough timing or production context for resource growth to be understandable

## Result

Player sees resources that grow over time through lazy resolution.

---

# Phase 10: Buildings System

## Goal

Allow player to upgrade basic buildings.

## Scope

Buildings:

* town_hall
* farm
* lumberyard
* quarry
* barracks
* walls
* shrine

Add table:

* `kingdom_buildings`

Endpoints:

* `GET /api/buildings/me`
* `POST /api/buildings/{type}/upgrade`

Rules:

* upgrades cost resources
* upgrade completion is resolved lazily from start and finish timestamps
* only one upgrade per building at a time
* no premium queues
* do not use background workers

Frontend:

* add a buildings view or dashboard section
* show building level, cost, status, and upgrade action
* refresh affected resources after upgrade commands

## Result

Player can spend resources to upgrade buildings and see progress without background workers.

---

# Phase 11: Army System

## Goal

Allow player to train basic units.

## Scope

Units:

* militia
* spearmen
* archers
* cavalry
* scouts

Add table:

* `kingdom_units`

Endpoints:

* `GET /api/army/me`
* `POST /api/army/train`

Rules:

* units cost resources
* training completion is resolved lazily from start and finish timestamps
* do not use background workers
* units have simple stats:

    * attack
    * defense
    * speed
    * supply

Frontend:

* add an army view or dashboard section
* show available unit counts, training costs, and training status
* refresh affected resources after training commands

## Result

Player can train basic army units and see their force grow through lazy resolution.

---

# Phase 12: PvE Missions

## Goal

Allow player to send army to PvE locations and receive basic reports.

## Scope

PvE locations:

* Black Forest
* Old Kurgan
* Dry Ford

Mission types:

* expedition
* scouting

Add tables:

* `missions`
* `reports`

Endpoints:

* `POST /api/missions/start`
* `GET /api/missions/me`
* `GET /api/reports/me`

Rules:

* mission has start time and resolve time
* mission completion is resolved lazily when missions or reports are read, or when a dependent command needs current state
* result may include resources, losses, and short flavor text
* do not use background workers
* basic report creation happens in this phase

Frontend:

* add missions view or dashboard section
* show available missions, sent units, mission status, and resolved report summaries
* add a basic reports list that can display PvE mission outcomes

## Result

Player can send units to PvE missions and receive basic reports in the UI.

---

# Phase 13: Report Polish

## Goal

Make mission and later raid results readable and atmospheric.

## Scope

Improve report content:

* title
* result
* phases
* units sent
* units lost
* loot
* consequences
* short narrative text

Report types:

* `pve_expedition`
* `pvp_raid`
* `event`

Frontend:

* improve report detail view or modal
* distinguish unread and read reports
* keep the reports UI simple and usable on the main dashboard flow

## Result

Player receives readable reports, not just raw numbers.

---

# Phase 14: Patron System

## Goal

Add first political layer before PvP raids.

## Scope

Patrons:

* independent
* empire_of_dusk
* old_pact

Add table:

* `patron_relations`

Endpoints:

* `GET /api/patron/me`
* `POST /api/patron/join`
* `POST /api/patron/break`
* `POST /api/patron/request-help`

Basic effects:

Independent:

* no tribute
* no protection

Empire of Dusk:

* pays tribute
* receives protection bonus
* gets trade/order flavor

Old Pact:

* can request defensive help
* owes contribution later
* gets oath/defense flavor

Frontend:

* add patron choice UI
* show current patron status and basic effects
* make patron state visible before raids are introduced

## Result

Player can choose a political path before interacting with PvP.

---

# Phase 15: Tribute and Pressure

## Goal

Make patron choice meaningful without making it oppressive.

## Scope

Add tribute logic for Empire of Dusk.

Rules:

* tribute is based on surplus, not total income
* tribute resolves lazily when patron state, resources, or events are read or changed
* tribute cannot block all progress
* if player cannot pay, create a crisis choice
* do not use background workers

Crisis choices:

* ask for delay
* accept debt
* give resources
* break patron relation
* request help from another side

Frontend:

* show tribute status and next pressure point in patron UI
* show crisis choices when they exist

## Result

Empire pressure exists, but does not create dead-end gameplay.

---

# Phase 16: Simple PvP Raids

## Goal

Allow basic asynchronous PvP interaction after patron choice exists.

## Scope

Endpoints:

* `GET /api/neighbors`
* `POST /api/raids/start`

Rules:

* player can raid another kingdom
* raid has travel time
* raid result is resolved lazily from timestamps
* defender cannot be destroyed
* attacker can steal limited resources
* walls reduce raid success
* patron state can influence simple protection or flavor
* weak-player protection exists
* do not use background workers

Add:

* infamy/dread value for aggressive behavior

Frontend:

* add simple neighbors and raid UI
* show raid limits, travel status, and basic outcomes
* show raid reports through the reports UI

## Result

Player can raid another player, but cannot delete or fully cripple them.

---

# Phase 17: Event Engine

## Goal

Add reusable choice-based event mechanics.

## Scope

Add tables:

* `game_events`
* `kingdom_events`
* `event_choices`

Event categories:

* economy
* ruler
* military
* patron
* dark_omen

Endpoints:

* `GET /api/events/me`
* `POST /api/events/{id}/choose`

Rules:

* event availability and expiry resolve lazily
* event choices can grant resources, affect units, affect patron state, or create reports
* do not use background workers

Frontend:

* add an events view or dashboard section
* show available choices and results clearly

## Result

The game can present and resolve simple narrative choices.

---

# Phase 18: First Event Content Pack

## Goal

Add enough event content to make the MVP loop feel alive.

## Scope

Create first MVP events:

* 20 to 30 events
* at least one event per category
* short, readable event text
* clear mechanical effects

Frontend:

* verify event text, choices, and outcomes fit the existing UI

## Result

Player occasionally receives events and makes meaningful choices.

---

# Phase 19: Balance Pass

## Goal

Make the first loop playable.

## Scope

Tune:

* resource production
* building costs
* building times
* unit costs
* mission rewards
* raid limits
* tribute values
* event rewards and penalties

Add config files if useful:

* `backend/internal/gameconfig/`

Frontend:

* adjust labels, empty states, disabled states, and visible feedback where balance changes affect player understanding

## Result

The game has a basic but playable economy loop.

---

# Phase 20: Smoke Tests and Seed Data

## Goal

Make MVP easy to verify.

## Scope

Add:

* seed data
* test users if useful
* basic backend tests
* API smoke checklist
* frontend smoke checklist

Required checks:

* register
* login
* create kingdom
* view ruler
* view resources
* upgrade building
* train unit
* start mission
* receive report
* choose patron
* start raid
* resolve event

## Result

Developer can verify MVP manually and with basic tests.

---

# Phase 21: First Playtest Build

## Goal

Prepare first internal test.

## Scope

Add:

* README playtest instructions
* known limitations
* reset database instructions
* feedback checklist

Playtest questions:

* Is the first 10 minutes clear?
* Is resource growth understandable?
* Are building choices meaningful?
* Are mission reports interesting?
* Is PvP annoying or exciting?
* Does patron choice feel meaningful?

## Result

MVP is ready for first manual playtest.

---

# Final MVP Definition of Done

MVP is complete when:

* user can register and login
* user can create one kingdom
* kingdom has a ruler
* kingdom has resources
* resources grow over time
* player can upgrade buildings
* player can train units
* player can send PvE missions
* player can receive reports
* player can raid another player
* player can choose patron
* player can receive and resolve events
* player cannot be fully destroyed by another player
* README explains how to run project
* CODEX_HANDOFF.md is up to date
