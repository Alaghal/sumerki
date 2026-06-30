# Codex Handoff

## Current Phase

Phase 25: Dashboard Decomposition.

## Status

Phase 25 completed.

## Completed

- Split the large DashboardPage render into focused feature panels.
- Added dashboard feature components for header, refresh, kingdom, patron, resources, ruler, buildings, army, missions, raids, events, and reports.
- Kept DashboardPage as the orchestration layer for session access, API loading, state ownership, refresh functions, and action handlers.
- Preserved existing gameplay behavior, API calls, visual order, and i18n label handling.
- Kept Game Shell, SVG map, mode navigation, backend changes, API shape changes, gameplay mechanics, and balance changes out of scope.
- Updated README, MVP phases, known limitations, post-playtest roadmap, and handoff.

## Changed Files

- `CODEX_HANDOFF.md`
- `README.md`
- `docs/MVP_PHASES.md`
- `docs/KNOWN_LIMITATIONS.md`
- `docs/POST_PLAYTEST_ROADMAP.md`
- `frontend/src/features/dashboard/DashboardPanels.tsx`
- `frontend/src/features/dashboard/shared.ts`
- `frontend/src/pages/DashboardPage.tsx`

## Commands Run

- `cd frontend && npm run typecheck`
- `cd frontend && npm run build`
- `rg -n "DashboardPage" frontend/src/features frontend/src/pages`
- `rg -n "\"Kingdom\"|\"Culture\"|\"Patron\"|\"Player\"|\"Favor\"|\"Logout\"|\"No data\"" frontend/src`
- `rg -n "northern_principality|black_forest_expedition|target_too_weak" frontend/src`
- `rg -n "cultureLabels|patronLabels|standingLabels|pressureStatusLabels|healthLabels|missionTypeLabels|missionStatusLabels|missionResultLabels|reportTypeLabels|eventCategoryLabels|eventStatusLabels|powerEstimateLabels|raidBlockedReasonLabels" frontend/src/pages/DashboardPage.tsx frontend/src/features/dashboard`
- `git diff --check`
- `git status --short`
- `git diff --stat`

## What Works Now

- DashboardPage no longer contains inline JSX for every major gameplay section.
- Ruler, resources, buildings, army, missions, reports, patron, events, and raids render through focused panel components.
- Dashboard state and API behavior remain centralized in DashboardPage.
- Russian/English i18n namespaces continue to drive current UI labels.
- Frontend typecheck and production build pass.

## Known Limitations

- The UI still feels like a dashboard/site rather than a map-first game shell.
- Game Shell is not implemented yet.
- SVG map is not implemented yet.
- Mode navigation is not implemented yet.
- Backend-provided event/report narrative text remains in its original language.
- Some internal enum strings remain in code/types/translation keys but should not appear as normal player UI.

## Next Recommended Step

Phase 26: Game Shell v1.
