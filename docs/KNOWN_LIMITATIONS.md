# Known Limitations

This document describes current MVP and playtest limitations. It is not a roadmap promise.

## Gameplay Limitations

- no alliances
- no province map
- no territory capture
- no city destruction
- no advanced diplomacy
- no player trade
- no market
- no chat
- no advanced combat simulation
- no hero or ruler actions beyond display
- no event chains
- no dark god avatar system
- no NPC retaliation
- no patron armies

## Technical Limitations

- local development setup only
- no production deployment
- no analytics
- no admin panel
- no email verification
- JWT and localStorage are MVP-only
- no rate limiting
- no production-grade anti-cheat
- no browser automation tests

## Balance Limitations

- values are first-pass MVP values
- no long-term economy simulation yet
- no analytics-driven balance yet
- PvP and tribute need playtest feedback

## Content Limitations

- event pack is small
- reports use local templates
- report narrative text may still fall back to backend-provided text when no stable report content key is exposed
- no final art direction

## UX Limitations

- Game Shell v1 and Local SVG Map v1 exist
- the map is local and symbolic only
- mode navigation is local UI state, not route-level world navigation
- node-specific context summaries exist, but deeper node workflows are still basic
- activity feed is more useful, but dense panels still need visual polish
- no full province system
- no territory capture
- no pathfinding
- no pan/zoom
- some UI may still be too dense on very small screens
- responsive hardening is first-pass and still needs broader device playtesting
- final art direction is not implemented

## Localization Limitations

- Russian/English i18n foundation exists
- most current UI shell and gameplay labels are migrated
- current MVP event content has Russian/English frontend translations where stable event and choice keys exist
- report shell, type, result, and template titles are localized where stable report type/result keys exist
- backend-generated report narrative may remain in its original language when no stable content key is available
- additional languages still require translation files, event/report namespace content, and registration

## Playtest 002 Limitations

- Playtest 002 is local-development focused.
- Seed data is local-only and uses public dev passwords.
- Smoke API checks mutate local data.
- No new gameplay systems were added for Playtest 002.
