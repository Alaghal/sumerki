# Decisions

This file records product and technical decisions that should remain stable across phases.

## 0001. Build a Browser MVP First

Status: Accepted

Sumerki will start as a browser-based MVP. Native mobile apps are out of scope for the first MVP.

## 0002. Use Go, Echo, and PostgreSQL for Backend

Status: Accepted

The backend will use Go with Echo and PostgreSQL in later phases. This is documented in Phase 0 but not implemented yet.

## 0003. Use React, TypeScript, Vite, and Tailwind for Frontend

Status: Accepted

The frontend will use React, TypeScript, Vite, and Tailwind in later phases. This is documented in Phase 0 but not implemented yet.

## 0004. Keep One Kingdom Per User in MVP

Status: Accepted

The MVP limits each user account to one kingdom. This keeps early domain logic and UI simpler.

## 0005. Avoid Advanced Game Systems in MVP

Status: Accepted

Large province maps, alliances, markets, chat, real-time combat, advanced diplomacy, and complex animations are intentionally deferred.

## 0006. Implement by Explicit Phases

Status: Accepted

Work should follow `docs/MVP_PHASES.md`. A requested phase should not pull in code or infrastructure from later phases unless explicitly approved.

## 0007. Playtest 001 Baseline

Status: Accepted

Playtest 001 is the first local-development playtest baseline. It includes auth, kingdom, ruler, resources, buildings, army, PvE missions, reports, patron selection, patron pressure, events, PvP raids with protection, seed data, smoke tests, and playtest documentation.

## 0008. Documentation Follows Implemented Phase Order

Status: Accepted

The phase documentation should reflect the implemented Playtest 001 order and actual feature set. Some earlier prompt order drift placed Tribute and Pressure before Raids; the repository now documents the completed order without renumbering migrations or code history.

## 0009. Map-First Game Shell After Playtest 001

Status: Accepted

After Playtest 001, the frontend should move from a dashboard-style layout to a map-first game shell with HUD, navigation, central map/scene, context panel, and activity feed.

## 0010. i18n Foundation Before Further UI Polish

Status: Accepted

Before deeper UI polish, the frontend should add a localization foundation for Russian and English with a structure that allows more languages later.

## 0011. Do Not Use Phaser/Pixi/Three.js For First Map UX

Status: Accepted

The first game map should use simple SVG/HTML inside React. Phaser, Pixi, and Three.js remain out of scope until a later phase proves they are necessary.
