# i18n Plan

## Phase 23: Foundation

* Install i18next and react-i18next.
* Configure default language ru.
* Configure fallback language en.
* Add language switcher.
* Persist language in localStorage.
* Add translation namespace structure.
* Translate shell/topbar/sidebar/auth/create-kingdom/basic dashboard labels.

## Phase 24: Copy Migration

Move these from hardcoded constants to i18n:

* culture labels
* patron labels
* standing labels
* pressure status labels
* health labels
* ruler stat labels
* resource labels
* building labels
* unit labels
* mission status labels
* mission type labels
* report type labels
* event category labels
* event status labels
* raid blocked reason labels
* power estimate labels
* API error messages
* buttons
* loading states
* empty states

## Phase 30: Event And Report Content

Migrate event/report content to localization.

Preferred long-term direction:

* Backend returns stable event keys and choice keys.
* Frontend translates event title/body/choice/result text.
* Backend remains responsible for effects and resolution.
* Existing Russian backend text can remain as fallback during transition.
