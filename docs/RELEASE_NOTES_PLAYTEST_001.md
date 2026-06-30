# Playtest 001 Release Notes

For the next internal UX-focused build, see `docs/RELEASE_NOTES_PLAYTEST_002.md`.

## Included

- auth
- kingdom creation
- ruler display
- resources
- buildings
- army
- PvE missions
- reports
- patrons
- tribute and pressure
- events
- PvP raids with protection
- seed data
- smoke tests

## Not Included

- alliances
- province map
- territory capture
- city destruction
- advanced diplomacy
- player trade
- market
- chat
- payments
- polished art
- mobile optimization
- advanced combat simulation
- dark god avatar system
- NPC retaliation
- event chains
- patron armies
- production deployment
- analytics
- admin panel

## Known Rough Edges

- This is a local development playtest build, not production hosting.
- Seed data is local-only and uses public dev passwords.
- Balance values are first-pass MVP values and need playtest feedback.
- The event pack is intentionally small.
- Reports use local templates.
- Patron copy and pressure mechanics may need clearer wording after testing.
- Smoke API checks happy paths but does not wait for mission or raid timers to resolve.
- Existing lazy mission/report duplicate behavior noted in `CODEX_HANDOFF.md` has not been changed.
- Visual design is functional MVP UI, not final art direction.

## How To Report Issues

Use `docs/FEEDBACK_TEMPLATE.md`.

For bugs, include:

- what happened
- what you expected
- steps to reproduce
- browser and OS
- screenshot or logs if available
