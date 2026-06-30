# i18n Plan

## Goals

* support Russian and English
* keep Russian as default
* make adding new languages easy
* eliminate mixed-language UI
* hide enum/API keys from player-facing UI
* avoid hardcoded strings in components after migration

## Non-Goals For Phase 23

* no backend localization
* no event/report content localization yet
* no automatic machine translation
* no dynamic locale loading unless simple
* no browser language auto-detection unless simple and safe
* no RTL support yet

## Phase 23: i18n Foundation ru/en

Scope:

* install i18next and react-i18next
* configure i18n provider
* default language ru
* fallback language en
* persist language in localStorage key `sumerki.ui.language`
* add language switcher
* add namespace structure
* translate shell/topbar/sidebar/auth/create-kingdom/basic dashboard labels

## Phase 24: UI Copy Migration

Move hardcoded labels to translation files:

* culture labels
* patron labels
* standing labels
* pressure status labels
* ruler health labels
* ruler stat labels
* resource labels
* building labels
* unit labels
* unit stat labels
* mission status labels
* mission type labels
* report type labels
* event category labels
* event status labels
* raid blocked reason labels
* raid power estimate labels
* API error messages
* buttons
* loading states
* empty states
* form labels
* validation messages

## Translation Namespaces

common:

* buttons
* loading states
* empty states
* generic errors
* language switcher

game:

* app name
* navigation
* HUD
* game shell
* activity feed
* dashboard labels

resources:

* gold
* food
* wood
* stone
* population

buildings:

* town_hall
* farm
* lumberyard
* quarry
* barracks
* market
* walls
* shrine

units:

* militia
* spearmen
* archers
* cavalry
* scouts
* unit stats

missions:

* mission labels
* mission statuses
* mission types
* mission form copy

reports:

* report types
* report statuses
* read/unread
* report actions

patrons:

* independent
* empire_of_dusk
* old_pact
* pressure
* tribute
* crisis states

events:

* event categories
* event states
* event UI shell copy
* later: event content

raids:

* raid statuses
* raid results
* blocked reasons
* power estimates

errors:

* backend error code mapping
* frontend validation errors

## Key Naming Rules

Examples:

* `resources.gold.name`
* `units.militia.name`
* `units.militia.description`
* `buildings.farm.name`
* `patrons.empire_of_dusk.name`
* `raids.blockedReasons.target_too_weak`
* `missions.status.active`
* `events.categories.dark_omen`
* `errors.invalid_credentials`

## Adding A New Language

1. Create `frontend/src/i18n/resources/<language-code>/`.
2. Copy namespace JSON files from `en` or `ru`.
3. Translate values.
4. Register language in i18n config.
5. Add language option to switcher.
6. Run frontend typecheck/build.
7. Manually check main flows.

## Event And Report Content Later

Current Phase 23/24 focus is UI localization.

Event/report content localization is Phase 30.

Preferred future approach:

* backend returns `eventKey` and `choiceKey`
* frontend translates title/body/choice/result by keys
* backend keeps effects and resolution logic
* old backend text can be fallback during transition
