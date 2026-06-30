# Codex Handoff

## Current Phase

Phase 23: i18n Foundation ru/en.

## Status

Phase 23 completed with a minimal Russian/English frontend localization foundation.

## Completed

- Installed `i18next` and `react-i18next`.
- Added `frontend/src/i18n/index.ts` with default Russian language, English fallback, supported language metadata, and `localStorage` persistence at `sumerki.ui.language`.
- Added initial translation namespaces for `common`, `game`, `auth`, `kingdom`, and `errors`.
- Imported the i18n setup before app render.
- Added a TopBar language switcher.
- Migrated TopBar, Sidebar, Login, Register, Create Kingdom, and basic Dashboard header/kingdom-card copy to i18n.
- Kept full DashboardPage/gameplay copy migration out of scope for Phase 24.
- Updated README, MVP phases, and known limitations for Phase 23.
- Updated post-playtest roadmap with Phase 23 completed status.

## Changed Files

- `CODEX_HANDOFF.md`
- `README.md`
- `docs/MVP_PHASES.md`
- `docs/KNOWN_LIMITATIONS.md`
- `docs/POST_PLAYTEST_ROADMAP.md`
- `frontend/package.json`
- `frontend/package-lock.json`
- `frontend/src/api/errors.ts`
- `frontend/src/components/i18n/LanguageSwitcher.tsx`
- `frontend/src/components/layout/Sidebar.tsx`
- `frontend/src/components/layout/TopBar.tsx`
- `frontend/src/i18n/index.ts`
- `frontend/src/i18n/resources/en/auth.json`
- `frontend/src/i18n/resources/en/common.json`
- `frontend/src/i18n/resources/en/errors.json`
- `frontend/src/i18n/resources/en/game.json`
- `frontend/src/i18n/resources/en/kingdom.json`
- `frontend/src/i18n/resources/ru/auth.json`
- `frontend/src/i18n/resources/ru/common.json`
- `frontend/src/i18n/resources/ru/errors.json`
- `frontend/src/i18n/resources/ru/game.json`
- `frontend/src/i18n/resources/ru/kingdom.json`
- `frontend/src/main.tsx`
- `frontend/src/pages/CreateKingdomPage.tsx`
- `frontend/src/pages/DashboardPage.tsx`
- `frontend/src/pages/LoginPage.tsx`
- `frontend/src/pages/RegisterPage.tsx`

## Commands Run

- `npm install i18next react-i18next`
- `npm run typecheck`
- `npm run build`
- `rg "sumerki.ui.language|fallbackLng|LanguageSwitcher|useTranslation" frontend/src`
- `rg "Phase 23|i18n Foundation|UI Copy Migration" README.md docs/MVP_PHASES.md docs/KNOWN_LIMITATIONS.md CODEX_HANDOFF.md`
- `git diff --check`
- `git status --short`

## What Works Now

- Frontend initializes i18next before rendering the app.
- Russian is the default UI language and English is the fallback.
- Selected language persists in `localStorage` under `sumerki.ui.language`.
- TopBar includes a language switcher.
- Shell/auth/create-kingdom/basic dashboard copy can switch between Russian and English.
- Existing API error mapping continues through the new `errors` namespace.

## Known Limitations

- Not all dashboard copy is migrated yet.
- Event/report content is not localized yet.
- Backend-provided text remains in its original language.
- Additional languages are not added yet, but the structure supports them.
- No game shell code is implemented yet.
- No map UI is implemented yet.

## Next Recommended Step

Phase 24: UI Copy Migration.
