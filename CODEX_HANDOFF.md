# Codex Handoff

## Current Phase

Phase 28: Context Panels And Activity Feed.

## Status

Phase 28 completed.

## Completed

- Added selected context state derived from map node clicks and local mode changes.
- Added node-specific context summaries for settlement, missions, neighbors, patron road, events, army, raids, and reports.
- Updated the right context area so the compact context summary appears above the existing mode panel.
- Improved the bottom activity feed with grouped actionable items for construction, training, missions, raids, events, reports, and patron pressure.
- Feed items switch local modes when clicked.
- Preserved existing gameplay panels, action handlers, API calls, and i18n behavior.
- Kept backend changes, API shape changes, new gameplay mechanics, province systems, territory capture, pathfinding, pan/zoom, canvas, Phaser, Pixi, and Three.js out of scope.
- Updated README, MVP phases, known limitations, post-playtest roadmap, and handoff.

## Changed Files

- `CODEX_HANDOFF.md`
- `README.md`
- `docs/MVP_PHASES.md`
- `docs/KNOWN_LIMITATIONS.md`
- `docs/POST_PLAYTEST_ROADMAP.md`
- `frontend/src/features/game/ActivityFeed.tsx`
- `frontend/src/features/map/MapNodeContextSummary.tsx`
- `frontend/src/features/map/types.ts`
- `frontend/src/i18n/resources/en/game.json`
- `frontend/src/i18n/resources/ru/game.json`
- `frontend/src/pages/DashboardPage.tsx`

## Commands Run

- `cd frontend && npm run typecheck`
- `cd frontend && npm run build`
- `rg -n "ContextSummary|MapNodeContext|selectedContext|ActivityFeed" frontend/src`
- `rg -n "Phaser|Pixi|Three|canvas" frontend/src`
- `git diff --check`
- `git status --short`

## What Works Now

- Selecting map nodes gives more useful context in the right panel.
- The right context panel still shows the appropriate existing mode panel below the summary.
- Activity feed shows clearer current activity using existing frontend data.
- Activity feed items can switch local modes.
- Existing gameplay actions remain usable.
- Russian/English i18n still works.
- Frontend typecheck and production build pass.

## Known Limitations

- Responsive and overflow polish remains Phase 29.
- Map is still symbolic and local.
- There is no province system.
- There is no territory capture.
- There is no pathfinding or pan/zoom.
- Backend-provided event/report narrative text remains in its original language.
- Some dense panels may still need visual polish.

## Next Recommended Step

Phase 29: Responsive And Overflow Hardening.
