# Codex Handoff

## Current Phase

Phase 26: Game Shell v1.

## Status

Phase 26 completed.

## Completed

- Added Game Shell v1 for the `/app` screen.
- Added a top resource HUD with kingdom, culture, patron, resources, unread reports, and active events.
- Added local mode navigation for map, city, army, missions, events, reports, patron, and raids.
- Added a central scene placeholder that prepares the screen for Phase 27 Local SVG Map v1.
- Added a right context panel that shows the existing mode-specific dashboard panels.
- Added a bottom activity feed summarizing active upgrades, training, missions, raids, events, and unread reports.
- Preserved existing dashboard/gameplay behavior, API calls, action handlers, and i18n behavior.
- Kept backend changes, API shape changes, gameplay mechanics, balance changes, and the real SVG map out of scope.
- Updated README, MVP phases, known limitations, post-playtest roadmap, and handoff.

## Changed Files

- `CODEX_HANDOFF.md`
- `README.md`
- `docs/MVP_PHASES.md`
- `docs/KNOWN_LIMITATIONS.md`
- `docs/POST_PLAYTEST_ROADMAP.md`
- `frontend/src/components/layout/AppShell.tsx`
- `frontend/src/features/game/ActivityFeed.tsx`
- `frontend/src/features/game/GameContextPanel.tsx`
- `frontend/src/features/game/GameHud.tsx`
- `frontend/src/features/game/GameModeNavigation.tsx`
- `frontend/src/features/game/GameScenePlaceholder.tsx`
- `frontend/src/features/game/GameShell.tsx`
- `frontend/src/features/game/gameModes.ts`
- `frontend/src/features/game/types.ts`
- `frontend/src/i18n/resources/en/game.json`
- `frontend/src/i18n/resources/ru/game.json`
- `frontend/src/pages/DashboardPage.tsx`

## Commands Run

- `cd frontend && npm run typecheck`
- `cd frontend && npm run build`
- `rg -n "GameShell|GameHud|GameModeNavigation|ActivityFeed|GameScenePlaceholder|GameContextPanel" frontend/src`
- `rg -n "\"Kingdom\"|\"Culture\"|\"Patron\"|\"Player\"|\"Favor\"|\"Logout\"|\"No data\"" frontend/src`
- `git diff --check`
- `git status --short`

## What Works Now

- `/app` uses a game-style shell instead of the pure stacked dashboard layout.
- Players can switch local game modes without changing routes.
- The right context panel keeps current gameplay actions usable.
- The resource HUD and activity feed summarize current settlement state using existing frontend data.
- Russian/English i18n still works through the existing namespace setup.
- Frontend typecheck and production build pass.

## Known Limitations

- The center scene is still a placeholder.
- Real SVG map nodes are not implemented yet.
- Mode changes are frontend-local and do not represent real world navigation.
- Backend-provided event/report narrative text remains in its original language.
- Some panels may still feel dense until Phase 28 context-panel refinement.

## Next Recommended Step

Phase 27: Local SVG Map v1.
