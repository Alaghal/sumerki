# Codex Handoff

## Current Phase

Phase 29: Responsive And Overflow Hardening.

## Status

Phase 29 completed.

## Completed

- Hardened `/app` against horizontal overflow in the app shell, top bar, shared cards/buttons, game shell, HUD, mode navigation, context panel, activity feed, local SVG map, map legend, and node context summaries.
- Updated dense dashboard rows to wrap instead of forcing single-line label/value layouts.
- Made unit training, mission, and raid input grids stack on narrow screens.
- Adjusted local map sizing and node labels so the SVG map scales inside its container on narrow viewports.
- Changed activity feed text from truncation to wrapping for readability.
- Added the missing `navigation.missions` label in English and Russian so mode navigation no longer falls back to a raw i18n key.
- Preserved existing gameplay actions, API calls, backend behavior, balance, and data contracts.
- Kept new gameplay mechanics, province systems, territory capture, pathfinding, pan/zoom, canvas, Phaser, Pixi, and Three.js out of scope.
- Updated README, MVP phases, known limitations, post-playtest roadmap, and handoff.

## Changed Files

- `CODEX_HANDOFF.md`
- `README.md`
- `docs/MVP_PHASES.md`
- `docs/KNOWN_LIMITATIONS.md`
- `docs/POST_PLAYTEST_ROADMAP.md`
- `frontend/src/components/layout/AppShell.tsx`
- `frontend/src/components/layout/TopBar.tsx`
- `frontend/src/components/ui/Button.tsx`
- `frontend/src/components/ui/Card.tsx`
- `frontend/src/features/dashboard/DashboardPanels.tsx`
- `frontend/src/features/game/ActivityFeed.tsx`
- `frontend/src/features/game/GameContextPanel.tsx`
- `frontend/src/features/game/GameHud.tsx`
- `frontend/src/features/game/GameModeNavigation.tsx`
- `frontend/src/features/game/GameShell.tsx`
- `frontend/src/features/map/LocalMap.tsx`
- `frontend/src/features/map/MapLegend.tsx`
- `frontend/src/features/map/MapNodeContextSummary.tsx`
- `frontend/src/i18n/resources/en/game.json`
- `frontend/src/i18n/resources/ru/game.json`

## Commands Run

- `cd frontend && npm run typecheck`
- `cd frontend && npm run build`
- `rg -n "overflow|break-words|min-w-0|max-w-full|flex-wrap|whitespace" frontend/src/features frontend/src/components frontend/src/pages`
- `rg -n "Phaser|Pixi|Three|canvas" frontend/src`
- `rg -n "sm:grid-cols-5|sm:grid-cols-\\[1fr_8rem_auto\\]|max-w-\\[14rem\\] truncate|flex justify-between gap-4|flex items-center justify-between gap-4|truncate" frontend/src/features frontend/src/components frontend/src/pages`
- `git diff --check`
- `git status --short`
- Browser check: `/app` at 390px, 768px, and 1280px with `documentElement.scrollWidth === clientWidth`.
- Browser check: clicked all 8 mode navigation buttons at 390px with no horizontal overflow.
- Browser check: language switcher changed nav labels between English and Russian at 390px with no horizontal overflow.

## What Works Now

- `/app` has no detected horizontal overflow at 390px, 768px, or 1280px in the in-app browser.
- HUD resources wrap into responsive columns.
- Mode navigation wraps safely and works on narrow screens.
- The symbolic SVG map scales inside its container, with smaller node labels on narrow screens.
- Context panel content wraps safely and remains scrollable on large viewports.
- Activity feed items wrap instead of truncating important text.
- Dense panels wrap label/value rows and stack unit input forms on narrow screens.
- Long event/report text inherits safe wrapping from shared card/context containers.
- Existing gameplay actions remain usable.
- Russian/English i18n still works.
- Frontend typecheck and production build pass.

## Known Limitations

- Map is still symbolic and local.
- There is no province system.
- There is no territory capture.
- There is no pathfinding or pan/zoom.
- Backend-provided event/report narrative text remains in its original language.
- Responsive hardening is first-pass and still needs broader real-device playtesting.
- Some dense panels may still need visual design polish.

## Next Recommended Step

Phase 30: Localized Event And Report Content.
