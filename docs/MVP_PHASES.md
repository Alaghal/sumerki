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

# Phase 3: Frontend Skeleton

## Goal

Create minimal frontend app.

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

Frontend starts locally and shows basic pages.

---

# Phase 4: Database Migrations

## Goal

Create initial schema.

## Scope

Tables:

* `users`
* `kingdoms`

Rules:

* one user can have one kingdom
* UUID primary keys
* reversible migrations

## Result

Database can store users and kingdoms.

---

# Phase 5: Auth API

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

Player can register, login, and fetch current user.

---

# Phase 6: Kingdom Creation API

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

Player can create and fetch their kingdom.

---

# Phase 7: Frontend Auth Integration

## Goal

Connect frontend auth to backend.

## Scope

Implement:

* register form
* login form
* logout
* JWT storage
* protected routes
* session restore via `GET /api/me`

Routing rules:

* no token -> `/login`
* token but no kingdom -> `/create-kingdom`
* token and kingdom exists -> `/app`

## Result

Player can register, login, refresh page, and stay logged in.

---

# Phase 8: Frontend Kingdom Creation

## Goal

Connect kingdom creation UI to backend.

## Scope

Create kingdom form:

* kingdom name
* culture select
* culture descriptions

After creation:

* redirect to `/app`

Dashboard shows:

* kingdom name
* culture
* placeholder resource cards
* placeholder ruler card

## Result

Full flow works:

register -> create kingdom -> app dashboard.

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

Implement passive production calculation.

Important rule:

Do not use real-time workers yet.
Calculate resources based on `last_calculated_at`.

Endpoints:

* `GET /api/resources/me`

## Result

Player sees resources that grow over time.

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
* market
* walls
* shrine

Add table:

* `kingdom_buildings`

Endpoints:

* `GET /api/buildings/me`
* `POST /api/buildings/{type}/upgrade`

Rules:

* upgrades cost resources
* upgrade takes time
* only one upgrade per building at a time
* no premium queues

## Result

Player can spend resources to upgrade buildings.

---

# Phase 11: Ruler System

## Goal

Add ruler as flavor and light gameplay modifier.

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

Ruler is generated when kingdom is created.

Endpoint:

* `GET /api/ruler/me`

## Result

Player has a ruler card with stats and flavor.

---

# Phase 12: Army System

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
* training takes time
* units have simple stats:

    * attack
    * defense
    * speed
    * supply

## Result

Player can train basic army units.

---

# Phase 13: PvE Missions

## Goal

Allow player to send army to PvE locations.

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
* `mission_reports`

Endpoints:

* `POST /api/missions/start`
* `GET /api/missions/me`
* `GET /api/reports/me`

Rules:

* mission has start time and resolve time
* result is calculated after timer ends
* result may include resources, losses, and flavor text

## Result

Player can send units to PvE mission and receive report.

---

# Phase 14: Battle Report System

## Goal

Make mission results readable and atmospheric.

## Scope

Report includes:

* title
* result
* phases
* units sent
* units lost
* loot
* consequences
* short narrative text

Report types:

* pve_expedition
* pvp_raid later

## Result

Player receives readable mission reports, not just raw numbers.

---

# Phase 15: Simple PvP Raids

## Goal

Allow basic asynchronous PvP interaction.

## Scope

Endpoints:

* `GET /api/neighbors`
* `POST /api/raids/start`

Rules:

* player can raid another kingdom
* raid has travel time
* raid result is auto-calculated
* defender cannot be destroyed
* attacker can steal limited resources
* walls reduce raid success
* weak-player protection exists

Add:

* infamy/dread value for aggressive behavior

## Result

Player can raid another player, but cannot delete or fully cripple them.

---

# Phase 16: Patron System

## Goal

Add first political layer.

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

## Result

Player can choose a political path.

---

# Phase 17: Tribute and Pressure

## Goal

Make patron choice meaningful without making it oppressive.

## Scope

Add tribute logic for Empire of Dusk.

Rules:

* tribute is based on surplus, not total income
* tribute cannot block all progress
* if player cannot pay, create a crisis choice

Crisis choices:

* ask for delay
* accept debt
* give resources
* break patron relation
* request help from another side

## Result

Empire pressure exists, but does not create dead-end gameplay.

---

# Phase 18: Event System

## Goal

Add simple choice-based events.

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

MVP event count:

* 20 to 30 events

## Result

Player occasionally receives events and makes meaningful choices.

---

# Phase 19: Basic Dashboard Polish

## Goal

Make the MVP playable as one coherent flow.

## Scope

Frontend pages:

* city overview
* resources
* buildings
* ruler
* army
* missions
* reports
* patron
* events

No advanced graphics.

Use:

* React
* Tailwind
* simple cards
* simple SVG/icons if needed

## Result

Player can use the core MVP through the UI.

---

# Phase 20: Balance Pass

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

## Result

The game has a basic but playable economy loop.

---

# Phase 21: Smoke Tests and Seed Data

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
* view resources
* upgrade building
* train unit
* start mission
* receive report
* choose patron
* resolve event

## Result

Developer can verify MVP manually and with basic tests.

---

# Phase 22: First Playtest Build

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
