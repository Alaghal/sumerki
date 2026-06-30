# UX And i18n Strategy

## Current Problem

Playtest 001 is functionally complete, but visually and structurally it still feels like a dashboard/site rather than a game.

Current UX issues:

* too many cards on one long page
* no true map-first game view
* navigation does not switch meaningful modes yet
* mixed Russian/English UI
* technical identifiers should not be shown to players
* some UI can overflow
* the game needs stronger spatial and contextual presentation

The interface currently answers:

"Which backend systems have been implemented?"

It should instead answer:

"What is happening around my settlement, and what can I do now?"

## Target UX

The target interface should be a strategy-lite game shell:

* top resource HUD
* left mode navigation
* central map or game scene
* right context panel
* bottom activity feed

The main screen should answer:

"What is happening around my settlement, and what can I do now?"

Not:

"Which backend systems have been implemented?"

## Game Shell Layout

Top HUD:

* gold
* food
* wood
* stone
* population
* kingdom name
* culture
* patron
* alert/report count

Left Navigation:

* Map
* City
* Army
* Missions
* Events
* Reports
* Patron
* Raids

Center Scene:

* local map
* player settlement
* PvE locations
* neighbor kingdoms
* route/status markers
* active mission/raid indicators

Right Context Panel:

* selected object details
* available actions
* forms and buttons
* selected city / mission / neighbor / patron / event / report context

Bottom Activity Feed:

* active building upgrades
* active unit training
* active missions
* active raids
* unread reports
* active events
* patron pressure warnings

## UX Principles

* show context, not all systems at once
* hide technical implementation details
* map-first where possible
* make next useful action visible
* use progressive disclosure
* keep MVP interactions simple
* avoid huge dashboards
* prefer cards inside context panels, not endless stacked sections
* keep timers and state visible
* do not add new gameplay mechanics during UX shell work

## Map v1 Direction

Use SVG/HTML inside React for the first map.

Do not use Phaser, Pixi, or Three.js for map v1.

The first map is not a full province system. It is a local node map with:

* player settlement
* Black Forest
* Old Kurgan
* Dry Ford
* neighbor kingdoms

Clicking a node opens the context panel.

## Localization Strategy

* Russian is default.
* English is the second supported language.
* More languages should be easy to add later.
* No player-facing string should remain hardcoded after migration.
* Enum/API keys should be translated through i18n.
* Backend may still return Russian event/report text until the later content localization phase.
* UI localization comes first; event/report content localization comes later.

## Technical i18n Direction

Recommended libraries:

* i18next
* react-i18next

Suggested translation folder structure:

```text
frontend/src/i18n/
  index.ts
  resources/
    ru/
      common.json
      game.json
      resources.json
      buildings.json
      units.json
      missions.json
      reports.json
      patrons.json
      events.json
      raids.json
      errors.json
    en/
      common.json
      game.json
      resources.json
      buildings.json
      units.json
      missions.json
      reports.json
      patrons.json
      events.json
      raids.json
      errors.json
```

Language persistence:

* localStorage key: `sumerki.ui.language`

Default language:

* `ru`

Fallback language:

* `en`

## Implementation Phases

* Phase 23: i18n Foundation ru/en
* Phase 24: UI Copy Migration
* Phase 25: Dashboard Decomposition
* Phase 26: Game Shell v1
* Phase 27: Local SVG Map v1
* Phase 28: Context Panels And Activity Feed
* Phase 29: Responsive And Overflow Hardening
* Phase 30: Localized Event And Report Content
* Phase 31: Playtest 002 UX Build
