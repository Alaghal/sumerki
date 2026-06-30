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

* Move hardcoded player-facing labels to i18n.
* Remove mixed Russian/English UI.
* Hide enum/API keys from player-facing UI.
* Translate shell, dashboard labels, errors, buttons, resources, units, buildings, missions, reports, patrons, events, raids.

### Phase 25: Dashboard Decomposition

* Split the huge DashboardPage into feature panels.
* Extract resources, buildings, army, missions, reports, patron, events, raids into separate feature components.
* Keep behavior the same.
* Prepare for game shell.

### Phase 26: Game Shell v1

* Replace the pure dashboard layout with game shell:
  * top resource HUD
  * left mode navigation
  * central scene area
  * right context panel
  * bottom activity feed
* Do not add new gameplay mechanics.

### Phase 27: Local SVG Map v1

* Add simple local map:
  * player settlement
  * Black Forest
  * Old Kurgan
  * Dry Ford
  * neighbor kingdoms
* Clicking map nodes opens context panels.
* Use SVG/HTML, not Phaser/Pixi/Three.js.

### Phase 28: Context Panels And Activity Feed

* Move major actions into contextual panels:
  * city context
  * mission context
  * raid target context
  * patron context
  * event context
* Add bottom feed for:
  * active building upgrades
  * active training
  * active missions
  * active raids
  * unread reports
  * active events

### Phase 29: Responsive And Overflow Hardening

* Fix layout overflow.
* Ensure long text wraps.
* Ensure buttons wrap.
* Ensure panels scroll when needed.
* Test desktop, laptop, tablet-width, and narrow mobile-ish screens.
* No new game mechanics.

### Phase 30: Localized Event And Report Content

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
