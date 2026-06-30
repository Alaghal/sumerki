# Playtest 002 Release Notes

## Focus

Playtest 002 focuses on UX, localization, and game feel after the Playtest 001 MVP systems were proven.

Main question:

Does Sumerki now feel more like a map-first strategy-lite game instead of a dashboard of systems?

## New Since Playtest 001

- Russian/English localization foundation
- language switcher
- current UI copy migrated to i18n
- dashboard split into focused panels
- Game Shell v1
- resource HUD
- local mode navigation
- Local SVG Map v1
- clickable map nodes:
  - settlement
  - Black Forest
  - Old Kurgan
  - Dry Ford
  - neighbors
  - patron road
  - omens/events
- node-specific context summaries
- improved activity feed
- responsive/overflow hardening
- localized event content where stable keys exist
- localized report shell/template text where stable keys exist

## Still Included From Playtest 001

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

- full province map
- territory capture
- city destruction
- alliances
- advanced diplomacy
- player trade
- market
- chat
- payments
- polished art
- mobile-first UX
- advanced combat simulation
- dark god avatar system
- NPC retaliation
- event chains
- patron armies
- production deployment
- analytics
- admin panel
- pathfinding
- map pan/zoom

## Known Rough Edges

- The map is symbolic and local.
- Some backend-generated report narrative may remain in its original language.
- Responsive hardening needs broader real-device testing.
- Dense panels may still need visual design polish.
- Final art direction is not implemented.
- Full mobile UX is not implemented.
- Seed data is local-only and uses public dev passwords.
- Smoke API mutates local data.
- Playtest 002 is local-development focused.

## How To Report Issues

Use `docs/FEEDBACK_TEMPLATE.md`.

Please include:

- browser
- OS
- screen width
- selected language
- seeded account or new account
- steps to reproduce
- screenshot/logs if useful
