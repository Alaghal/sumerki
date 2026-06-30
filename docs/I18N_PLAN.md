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
* report template titles by stable type/result keys

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
* event content by stable event and choice keys

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

## Event And Report Content

Phase 30 implemented frontend event/report content localization where stable keys are available.

Fallback order:

1. Frontend i18n by stable event, choice, report type, or report result key.
2. Backend-provided text.
3. Generic localized fallback.

Event content:

* Use `events.content.<eventKey>.title`.
* Use `events.content.<eventKey>.body`.
* Use `events.content.<eventKey>.choices.<choiceKey>.label`.
* Use `events.content.<eventKey>.choices.<choiceKey>.description`.
* Use `events.content.<eventKey>.choices.<choiceKey>.result.title`.
* Use `events.content.<eventKey>.choices.<choiceKey>.result.body`.

Report content:

* Use `reports.templates.<reportType>.title.<result>` for stable report title templates.
* Use `reports.phaseTitles.<backendPhaseTitle>` only for stable/common phase labels.
* Keep backend report body and detailed phase body as fallback unless the API exposes a stable content key.

Adding new event/report translations:

1. Add the backend `eventKey` and choice keys to `frontend/src/i18n/resources/ru/events.json`.
2. Add matching English content to `frontend/src/i18n/resources/en/events.json`.
3. Add report title templates to `reports.json` only when a stable type/result or future report key exists.
4. Run frontend typecheck/build.
