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

