# UX And i18n Strategy

## Problem

The current Playtest 001 frontend is functionally complete but feels like a website/dashboard:

* too many cards on one page
* systems are stacked sequentially
* navigation does not switch meaningful modes
* mixed Russian/English copy
* technical labels leak into the player UI
* layout can overflow

## Target UX

The main screen should become:

* Top: resource HUD
* Left: mode navigation
* Center: map or game scene
* Right: context panel
* Bottom: activity feed

Example layout:

Top HUD:

* gold
* food
* wood
* stone
* population
* kingdom name
* patron
* alert count

Left navigation:

* Map
* City
* Army
* Missions
* Events
* Reports
* Patron
* Raids

Center:

* local map
* player settlement
* PvE locations
* neighbor kingdoms
* route/status markers

Right context:

* selected object details
* available actions
* forms and buttons

Bottom feed:

* active timers
* new reports
* active events
* raids/missions in progress

## Design Principle

The main screen should answer:

"What is happening around my settlement, and what can I do now?"

Not:

"Which backend systems have been implemented?"

## Localization Principle

* Russian is default.
* English is the second supported language.
* UI must be ready for more languages later.
* No player-facing string should be hardcoded in components after migration.
* Enums should be translated through i18n keys.
* API keys should not appear in the player UI.

## Technical i18n Direction

Use:

* i18next
* react-i18next

Suggested translation namespaces:

* common
* game
* resources
* buildings
* units
* missions
* reports
* patrons
* events
* raids
* errors

Suggested structure:

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

Default:

* `ru`

Fallback:

* `en`
