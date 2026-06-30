# Codex Handoff

## Current Phase

Phase 27: Local SVG Map v1.

## Status

Phase 27 completed.

## Completed

- Added Local SVG Map v1 for the Game Shell center scene.
- Added map nodes for the player settlement, Black Forest, Old Kurgan, Dry Ford, two neighbor slots, patron road, and omens/events.
- Added map node clicks that switch local mode/context to city, missions, raids, patron, or events.
- Selecting known neighbor nodes also selects the corresponding raid target using existing frontend state.
- Added map status indicators using existing frontend data for active missions, available missions, active raids, active events, unread reports, patron pressure, and patron relation.
- Added a compact map legend.
- Added a new `map` i18n namespace for Russian and English map copy.
- Preserved existing gameplay behavior, API calls, action handlers, and i18n behavior.
- Kept backend changes, API shape changes, province systems, territory capture, pathfinding, pan/zoom, canvas, Phaser, Pixi, Three.js, and new gameplay mechanics out of scope.
- Updated README, MVP phases, known limitations, post-playtest roadmap, and handoff.

## Changed Files

- `CODEX_HANDOFF.md`
- `README.md`
- `docs/MVP_PHASES.md`
- `docs/KNOWN_LIMITATIONS.md`
- `docs/POST_PLAYTEST_ROADMAP.md`
- `frontend/src/features/game/GameScenePlaceholder.tsx`
- `frontend/src/features/map/LocalMap.tsx`
- `frontend/src/features/map/MapLegend.tsx`
- `frontend/src/features/map/mapNodes.ts`
- `frontend/src/features/map/types.ts`
- `frontend/src/i18n/index.ts`
- `frontend/src/i18n/resources/en/map.json`
- `frontend/src/i18n/resources/ru/map.json`
- `frontend/src/pages/DashboardPage.tsx`

## Commands Run

- `cd frontend && npm run typecheck`
- `cd frontend && npm run build`
- `rg -n "LocalMap|mapNodes|selectedMapNode|MapLegend" frontend/src`
- `rg -n "Phaser|Pixi|Three|canvas" frontend/src`
- `git diff --check`
- `git status --short`

## What Works Now

- `/app` has a central local SVG map instead of a placeholder.
- Players can click map nodes to navigate to city, missions, raids, patron, and events.
- Existing right context panels and actions remain usable.
- Map indicators reflect current frontend state where available.
- Russian/English i18n still works through the registered namespace setup.
- Frontend typecheck and production build pass.

## Known Limitations

- The map is symbolic and local, not a province system.
- There is no territory capture.
- There is no pathfinding.
- There is no pan/zoom.
- Neighbor nodes use existing neighbor data and fall back to generic labels.
- Mission nodes provide mode switching and basic status only.
- Deeper node-specific context behavior remains Phase 28.
- Backend-provided event/report narrative text remains in its original language.

## Next Recommended Step

Phase 28: Context Panels And Activity Feed.
