# Codex Handoff

## Current Phase

Phase 24: UI Copy Migration.

## Status

Phase 24 completed.

## Completed

- Added gameplay translation namespaces for resources, buildings, units, missions, reports, patrons, events, and raids.
- Registered the new namespaces in the i18n setup.
- Migrated DashboardPage player-facing gameplay labels to i18n.
- Reduced mixed Russian/English UI in the current screens.
- Replaced known resource/building/unit/mission/patron/event/raid enum labels with translated labels.
- Kept backend event/report narrative localization out of scope for Phase 30.
- Updated README, MVP phases, known limitations, post-playtest roadmap, and handoff.

## Changed Files

- `CODEX_HANDOFF.md`
- `README.md`
- `docs/MVP_PHASES.md`
- `docs/KNOWN_LIMITATIONS.md`
- `docs/POST_PLAYTEST_ROADMAP.md`
- `frontend/src/api/errors.ts`
- `frontend/src/components/i18n/LanguageSwitcher.tsx`
- `frontend/src/i18n/index.ts`
- `frontend/src/i18n/resources/en/common.json`
- `frontend/src/i18n/resources/en/game.json`
- `frontend/src/i18n/resources/en/buildings.json`
- `frontend/src/i18n/resources/en/events.json`
- `frontend/src/i18n/resources/en/missions.json`
- `frontend/src/i18n/resources/en/patrons.json`
- `frontend/src/i18n/resources/en/raids.json`
- `frontend/src/i18n/resources/en/reports.json`
- `frontend/src/i18n/resources/en/resources.json`
- `frontend/src/i18n/resources/en/units.json`
- `frontend/src/i18n/resources/ru/common.json`
- `frontend/src/i18n/resources/ru/game.json`
- `frontend/src/i18n/resources/ru/buildings.json`
- `frontend/src/i18n/resources/ru/events.json`
- `frontend/src/i18n/resources/ru/missions.json`
- `frontend/src/i18n/resources/ru/patrons.json`
- `frontend/src/i18n/resources/ru/raids.json`
- `frontend/src/i18n/resources/ru/reports.json`
- `frontend/src/i18n/resources/ru/resources.json`
- `frontend/src/i18n/resources/ru/units.json`
- `frontend/src/pages/DashboardPage.tsx`

## Commands Run

- `npm run typecheck`
- `npm run build`
- `rg -n "cultureLabels|patronLabels|standingLabels|pressureStatusLabels|healthLabels|missionTypeLabels|missionStatusLabels|missionResultLabels|reportTypeLabels|eventCategoryLabels|eventStatusLabels|powerEstimateLabels|raidBlockedReasonLabels" frontend/src/pages/DashboardPage.tsx`
- `rg -n "\"Kingdom\"|\"Culture\"|\"Patron\"|\"Player\"|\"Favor\"|\"Logout\"|\"No data\"|northern_principality|black_forest_expedition|target_too_weak" frontend/src`
- `git diff --check`
- `git status --short`

## What Works Now

- Frontend initializes i18next before rendering the app.
- Russian is the default UI language and English is the fallback.
- Selected language persists in `localStorage` under `sumerki.ui.language`.
- TopBar includes a language switcher.
- Current UI labels switch between Russian and English more consistently.
- Gameplay labels are mostly translated through i18n namespaces.
- Additional languages can add equivalent namespace files.

## Known Limitations

- DashboardPage is still structurally large and should be decomposed in Phase 25.
- Game Shell is not implemented yet.
- SVG map is not implemented yet.
- Backend-provided event/report narrative text remains in its original language.
- Some internal enum strings remain in code/types/translation keys but should not appear as normal player UI.

## Next Recommended Step

Phase 25: Dashboard Decomposition.
