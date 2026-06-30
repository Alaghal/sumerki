# Playtest Checklist

Use this checklist for the first internal manual playtest.

## Pre-Flight

- [ ] Database starts.
- [ ] Migrations apply.
- [ ] Seed works.
- [ ] Backend tests pass.
- [ ] Frontend typecheck/build passes.
- [ ] `smoke-api` passes, or the known warning is documented.

## Manual App Checklist

- [ ] Login works.
- [ ] Logout works.
- [ ] Refresh keeps the session.
- [ ] Create kingdom works.
- [ ] Ruler is visible.
- [ ] Resources are visible.
- [ ] Building upgrade works.
- [ ] Resources update after spending.
- [ ] Unit training works.
- [ ] Trained units appear after lazy refresh.
- [ ] Mission start works.
- [ ] Sent units become unavailable.
- [ ] Mission resolves.
- [ ] Report appears.
- [ ] Report can be read and marked read.
- [ ] Patron options are visible.
- [ ] Join patron works.
- [ ] Break patron works.
- [ ] Pressure is visible if implemented.
- [ ] Events are visible.
- [ ] Choosing an event applies effects.
- [ ] Event report appears.
- [ ] Neighbors are visible.
- [ ] Blocked raid reasons make sense.
- [ ] Raid can start when valid.
- [ ] Raid resolves.
- [ ] Attacker and defender reports appear.

## UX Checks

- [ ] No blank screens.
- [ ] Loading states show.
- [ ] Errors are readable.
- [ ] Buttons are disabled when action is impossible.
- [ ] Labels are understandable.
