# Post-Playtest Roadmap

## Current State

Playtest 001 is complete.

The MVP has working systems:

* auth
* kingdom creation
* ruler
* resources
* buildings
* army
* PvE missions
* reports
* patrons
* tribute and pressure
* events
* PvP raids with protection
* seed data
* smoke tests
* first playtest documentation

However, the current UI is still structurally closer to a dashboard than a game.

## Main Post-Playtest Goal

Transform Sumerki from a working systems dashboard into a map-first strategy-lite game interface.

The player should feel:

* "I am ruling a small settlement in Sumerki"

not:

* "I am looking at an admin panel of implemented APIs"

## Planned Phases

### Phase 22: UX And i18n Strategy Docs

Status: Completed

* Write the UX target.
* Write the localization strategy.
* Define the game shell layout.
* Define UI copy rules.
* No code changes.

### Phase 23: i18n Foundation ru/en

Status: Completed

* Add i18n infrastructure.
* Support Russian and English.
* Add language switcher.
* Persist selected language.
* Create translation namespace structure.
* Do not migrate all game text yet.

### Phase 24: UI Copy Migration

Status: Completed

* Move hardcoded player-facing labels to i18n.
* Remove mixed Russian/English UI.
* Hide enum/API keys from player-facing UI.
* Translate shell, dashboard labels, errors, buttons, resources, units, buildings, missions, reports, patrons, events, raids.

### Phase 25: Dashboard Decomposition

Status: Completed

* Split the large DashboardPage render into feature panels.
* Extracted ruler, resources, buildings, army, missions, reports, patron, events, and raids into focused components.
* Kept behavior and i18n handling the same.
* Prepared the frontend structure for game shell work.

### Phase 26: Game Shell v1

Status: Completed

* Replaced the pure dashboard layout with Game Shell v1.
* Added a top resource HUD, local mode navigation, central scene placeholder, right context panel, and bottom activity feed.
* Kept existing gameplay behavior and API calls unchanged.
* Kept the real SVG map out of scope for Phase 27.

### Phase 27: Local SVG Map v1

Status: Completed

* Added a symbolic local SVG map.
* Added settlement, PvE locations, neighbors, patron road, and omens/events nodes.
* Clicking map nodes switches local modes and right context panels.
* Kept province systems, territory capture, pathfinding, and new mechanics out of scope.

### Phase 28: Context Panels And Activity Feed

Status: Completed

* Added selected-node context summaries to the right panel.
* Kept existing mode panels and actions usable below the summaries.
* Improved the activity feed with grouped actionable items.
* Feed items can switch local modes.

### Phase 29: Responsive And Overflow Hardening

Status: Completed

* Hardened game shell, HUD, mode navigation, context panel, activity feed, local map, and dense panels against horizontal overflow.
* Ensured long labels and event/report text wrap more safely.
* Made unit training, mission, and raid input grids stack on narrow screens.
* Browser-checked `/app` at 390px, 768px, and 1280px with no horizontal overflow.
* Kept new game mechanics out of scope.

### Phase 30: Localized Event And Report Content

Status: Planned

* Localize event and report text.
* Prefer frontend translation by event/report keys where practical.
* Keep backend compatible.
* Add English text for current MVP event pack and report templates.
* Keep Russian as default.

### Phase 31: Playtest 002 UX Build

* Package the UX/i18n improvements for a second playtest.
* Update playtest guide, release notes, known limitations, checklist, and feedback template.
* Focus on whether the game now feels like a strategy game rather than a dashboard.

## Future Big Tracks After UX/i18n

These are optional post-Playtest tracks, not committed phases:

* Retention / LiveOps-lite
* Alliances v1
* Province map v1
* Dark Path v1
* Production readiness
* Analytics and deployment
