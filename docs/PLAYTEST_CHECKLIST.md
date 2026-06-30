# Playtest Checklist

Use this checklist after running `make seed-dev` and starting the backend and frontend.

## Account

- Login works.
- Logout works.
- Refresh keeps the session.

## Kingdom

- Kingdom is visible.
- Ruler is visible.
- Resources are visible.

## Economy

- Resources grow after refresh.
- Building upgrade starts.
- Resources update after spending.

## Army

- Unit training starts.
- Trained units appear after lazy refresh.
- Unavailable unit requirements are understandable.

## Missions

- Mission start works.
- Sent units become unavailable.
- Mission resolves after its timer.
- Report appears.

## Raids

- Neighbors are visible.
- Blocked reasons make sense.
- Raid can start when a valid target exists.
- Defender is not destroyed.
- Protected resources remain.

## Patron

- Patron options are visible.
- Join works.
- Break works.
- Pressure is visible.
- Tribute does not reduce resources below protected minimums.

## Events

- Events are visible.
- Choosing an event applies effects.
- Event report appears.
- Resolved event is not applied twice.

## General

- No blank pages.
- Errors are readable.
- Loading states are understandable.
