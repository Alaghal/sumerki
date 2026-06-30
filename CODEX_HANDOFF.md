# Codex Handoff

## Current Phase

Documentation Sync after Phase 21.

## Status

Documentation sync completed after Phase 21 Playtest 001. This was a documentation-only pass to align stale planning and reference docs with the implemented Playtest 001 state.

## Completed

- Synced `docs/MVP_PHASES.md` to the actual implemented phase order through Phase 21.
- Added top-level implementation status and phase-order correction notes to `docs/MVP_PHASES.md`.
- Marked Phases 0 through 21 as completed in `docs/MVP_PHASES.md`.
- Corrected Phase 15 and Phase 16 ordering:
  - Phase 15: Simple PvP Raids With Protection
  - Phase 16: Tribute And Pressure
- Removed stale Patron System mention of an unimplemented patron help endpoint.
- Updated phase descriptions for actual patron, raid, tribute/pressure, event, seed/smoke, and playtest deliverables.
- Updated the Final MVP Definition of Done as completed for Playtest 001.
- Updated `docs/API_CONTRACT.md` intro and report wording so it reflects the implemented Playtest 001 API and event report coverage.
- Updated `docs/DOMAIN_MODEL.md` intro and small stale phase-specific statements.
- Updated `docs/DOMAIN_MODEL.md` to show starting scouts as 3.
- Added current known limitations to `docs/DOMAIN_MODEL.md`.
- Added accepted decisions:
  - `0007. Playtest 001 Baseline`
  - `0008. Documentation Follows Implemented Phase Order`

## Changed Files

- `CODEX_HANDOFF.md`
- `docs/MVP_PHASES.md`
- `docs/API_CONTRACT.md`
- `docs/DOMAIN_MODEL.md`
- `docs/DECISIONS.md`

## Commands Run

- `rg` checks for stale patron help, old Phase 15/16 ordering, old draft wording, stale scout counts, old report type names, and stale report wording
- `rg -n "Status: Completed|Implementation Status|Phase Order Correction|Final MVP" docs/MVP_PHASES.md`
- `rg -n "patron/pressure|events/me|events/\\{id\\}/choose|neighbors|raids/me|raids/start|event" docs/API_CONTRACT.md`

No backend or frontend tests were run because this was a documentation-only sync and no code files were modified.

## What Works Now

- Documentation no longer claims Phases 16 or 21 are pending.
- `docs/MVP_PHASES.md` shows Playtest 001 as completed.
- `docs/MVP_PHASES.md` no longer contradicts `CODEX_HANDOFF.md` about Phase 15/16 ordering.
- Patron docs no longer list the unimplemented patron help endpoint.
- API contract and domain docs match README and the current implemented MVP feature set more closely.
- Decisions record the Playtest 001 baseline and the phase-order correction policy.

## Known Limitations

- Documentation sync is based on existing docs and handoff, not a fresh live API smoke run.
- Code was intentionally not changed.
- README was reviewed and not changed because it already points to the Playtest 001 docs and current Phase 21 state.

## Next Recommended Step

Begin post-playtest feedback collection or create issues from Playtest 001 feedback.
