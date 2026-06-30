# Balance Notes

Phase 19 is the first MVP balance pass. Values are tuned for short local playtests, not production.

## Early Loop

A new kingdom starts with enough resources to upgrade a basic economic building, build barracks level 1, train several militia or scouts, and send a low-risk mission in the first session.

Starting resources:

- gold: 600
- food: 400
- wood: 400
- stone: 300
- population: 120

Starting army:

- militia: 10
- scouts: 3

## Production

Base production plus level-1 economic buildings gives roughly:

- gold: 30 per hour
- food: 45 per hour
- wood: 37 per hour
- stone: 25 per hour
- population: 1 per hour

Resources resolve lazily during reads and commands. There are no background workers.

## Buildings And Units

Building upgrade duration remains `60 * target_level` seconds. Level 1 and 2 upgrades are intentionally short for MVP testing.

Militia and scouts are available immediately. Spearmen and archers require barracks level 1. Cavalry requires barracks level 2.

## Missions

Dry Ford is the shortest low-risk mission. Black Forest is the first general expedition. Old Kurgan is more dangerous and more rewarding.

Mission losses are intentionally mild in the first pass so the starting army can survive early experimentation.

## Raids

Raids remain limited by newbie protection, same-target cooldowns, defender protection, protected resource minimums, and capped loot. PvP should be worth trying, but not better than every PvE and economy action.

## Tribute And Pressure

Empire tribute and Old Pact contribution use lazy surplus-based resolution. Tribute cannot spend below protected minimums, and inactivity processing is capped.

## Events

The first event pack uses small resource, unit, honor, dread, and favor deltas. Event effects are flavor-forward and should not dominate missions or building progression.

## Known Rough Edges

- Balance is first-pass MVP tuning.
- Values are not production-ready.
- There is no analytics-driven balancing yet.
- There is no long-term economy simulation yet.
- There is no monetization balancing.
- Playtest feedback is needed.
