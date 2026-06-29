# Domain Model

This is the initial domain model for the Sumerki MVP. It is documentation only in Phase 0.

## User

Represents a player account.

Fields:

- `id`: UUID
- `email`: unique email address
- `passwordHash`: hashed password
- `createdAt`: account creation timestamp
- `updatedAt`: last update timestamp

Rules:

- A user can have one kingdom in the MVP.
- Email uniqueness is case-insensitive.
- Passwords are never stored in plain text.
- Password hash must not be empty.

## Kingdom

Represents the player's settlement and strategic identity.

Fields:

- `id`: UUID
- `userId`: owner user UUID
- `name`: kingdom name
- `culture`: selected starting culture
- `patron`: selected patron, nullable before patron choice
- `createdAt`: kingdom creation timestamp
- `updatedAt`: last update timestamp

Rules:

- One kingdom per user.
- Every kingdom has one ruler.
- Every kingdom has one resource row.
- Every kingdom has all MVP building rows.
- Every kingdom has all MVP unit rows.
- Name length is 3 to 32 characters.
- Culture must be one of the supported culture values.
- Patron can be null until selected, then must be one of the supported patron values.

## Ruler

Represents the current leader of a kingdom.

Fields:

- `id`: UUID
- `kingdomId`: kingdom UUID
- `name`: generated ruler name
- `age`: ruler age, 18 to 80
- `culture`: inherited from the kingdom culture
- `authority`: flavor stat, 1 to 100
- `courage`: flavor stat, 1 to 100
- `cunning`: flavor stat, 1 to 100
- `honor`: flavor stat, 1 to 100
- `cruelty`: flavor stat, 1 to 100
- `ambition`: flavor stat, 1 to 100
- `paranoia`: flavor stat, 1 to 100
- `healthStatus`: ruler health state
- `createdAt`: ruler creation timestamp
- `updatedAt`: last update timestamp

Rules:

- Each kingdom has exactly one ruler.
- A ruler is generated when a kingdom is created.
- Existing kingdoms can receive a default or lazily generated ruler.
- Ruler stats are flavor only in v1 and do not modify gameplay.
- No heirs, dynasties, intrigue, ruler death, ruler actions, or gameplay modifiers are included in Phase 8.

## Culture

Starting cultural identity for a kingdom.

Values:

- `northern_principality`
- `lizard_grad`
- `free_posad`

## Patron

Political or supernatural alignment chosen during progression.

Values:

- `independent`
- `empire_of_dusk`
- `old_pact`

## Resource

Represents stored kingdom materials.

MVP resources:

- gold
- food
- wood
- stone
- population

Post-MVP resources may add other materials such as iron or silver, but they are out of scope for the Phase 9 MVP resource system.

Fields:

- `kingdomId`: kingdom UUID
- `gold`: stored gold
- `food`: stored food
- `wood`: stored wood
- `stone`: stored stone
- `population`: stored population
- `lastCalculatedAt`: timestamp used for lazy production
- `createdAt`: resource row creation timestamp
- `updatedAt`: last update timestamp

Rules:

- Resources are integer values and cannot be negative.
- Resource production is calculated lazily when resources are read or later commands need current resources.
- Phase 9 uses simple base production per hour only.
- Phase 9 has no resource caps.
- Resource spending is performed by later feature commands such as building upgrades and unit training.
- No background workers, cron jobs, or real-time server ticking are used for resources.

## Building

Represents an upgradeable structure in a kingdom.

MVP buildings:

- `town_hall`: starts at level 1, no gameplay effect in Phase 10
- `farm`: starts at level 1, increases food production
- `lumberyard`: starts at level 1, increases wood production
- `quarry`: starts at level 1, increases stone production
- `market`: starts at level 1, increases gold production
- `barracks`: starts at level 0, future unit training
- `walls`: starts at level 0, future raid defense
- `shrine`: starts at level 0, future events and patrons

Fields:

- `id`: UUID
- `kingdomId`: kingdom UUID
- `type`: building type
- `level`: current building level, 0 to 5
- `upgradeStartedAt`: nullable upgrade start timestamp
- `upgradeFinishesAt`: nullable upgrade finish timestamp
- `createdAt`: building row creation timestamp
- `updatedAt`: last update timestamp

Rules:

- Each kingdom has one row for each MVP building type.
- Max level is 5.
- Upgrade costs are deterministic and use gold, wood, and stone only in Phase 10.
- Upgrade duration is `60 * target_level` seconds.
- Upgrade completion is resolved lazily when buildings or resources are read, or when an upgrade command is issued.
- Farm, lumberyard, quarry, and market modify resource production in Phase 10.
- Town hall, barracks, walls, and shrine have no gameplay effect yet or only placeholder future purpose.
- No building destruction, cancellation, premium queues, speedups, or multiple construction queues are included in Phase 10.

## Unit

Represents trainable military forces.

MVP units:

- `militia`: starts with 10, no barracks requirement
- `scouts`: starts with 2, no barracks requirement
- `spearmen`: starts with 0, requires barracks level 1
- `archers`: starts with 0, requires barracks level 1
- `cavalry`: starts with 0, requires barracks level 2

Fields:

- `id`: UUID
- `kingdomId`: kingdom UUID
- `type`: unit type
- `amount`: available trained units
- `createdAt`: row creation timestamp
- `updatedAt`: last update timestamp

Rules:

- Each kingdom has one row for each MVP unit type.
- Unit amount cannot be negative.
- Units have simple attack, defense, speed, and supply stats.
- Units have deterministic per-unit training costs and seconds-per-unit values.
- Training spends resources immediately, including population.
- Training completion is resolved lazily when army is read or another army command needs current units.
- No missions, combat, raids, unit death, healing, equipment, heroes, upkeep, cancellation, premium speedups, or multiple army formations are included in Phase 11.

## UnitTrainingOrder

Represents asynchronous unit training in progress.

Fields:

- `id`: UUID
- `kingdomId`: kingdom UUID
- `unitType`: unit type being trained
- `amount`: number of units being trained
- `status`: `training` or `completed`
- `startedAt`: training start timestamp
- `finishesAt`: training finish timestamp
- `completedAt`: nullable completion timestamp
- `createdAt`: row creation timestamp
- `updatedAt`: last update timestamp

Rules:

- Training order amount must be 1 to 50 in the Phase 11 service.
- Training duration is `seconds_per_unit * amount`.
- Finished training orders add their amount to the matching unit row.
- Finished training orders are marked completed with `completedAt`.
- Completion uses lazy resolution only; no background workers, cron jobs, Redis, queues, or server ticking.

## Mission

Represents asynchronous PvE activity.

MVP missions:

- `black_forest_expedition`
- `old_kurgan_expedition`
- `dry_ford_scouting`

Fields:

- `id`: UUID
- `kingdomId`: kingdom UUID
- `missionKey`: configured mission key
- `missionType`: `expedition` or `scouting`
- `status`: `active` or `completed`
- `startedAt`: mission start timestamp
- `finishesAt`: mission finish timestamp
- `completedAt`: nullable completion timestamp
- `resultJson`: stored result payload with rewards and losses
- `createdAt`: row creation timestamp
- `updatedAt`: last update timestamp

Rules:

- A kingdom sends available units to a PvE mission.
- Sent units are subtracted from available `kingdom_units` immediately.
- Sent units are unavailable until the mission resolves.
- Mission completion is resolved lazily when missions or reports are read, or when mission start needs current state.
- Completed missions return surviving units, permanently lose killed units, grant resource rewards, and create a basic report.
- Mission outcome and losses use simple deterministic MVP rules.
- No background workers, cron jobs, Redis, queues, or real-time ticking are used.
- PvP raids, player-vs-player combat, route/pathfinding, heroes, equipment, unit XP, events, patrons, tribute, and advanced battle simulation are not implemented in Phase 12.

## MissionUnit

Represents units allocated to a PvE mission.

Fields:

- `id`: UUID
- `missionId`: mission UUID
- `unitType`: unit type sent
- `amountSent`: number of units sent
- `amountLost`: number of units lost after completion
- `amountReturned`: number of units returned after completion
- `createdAt`: row creation timestamp
- `updatedAt`: last update timestamp

Rules:

- `amountSent` must be positive.
- `amountLost` and `amountReturned` cannot be negative.
- `amountLost + amountReturned` cannot exceed `amountSent`.
- Mission unit rows do not represent a separate locked-unit system; they are the mission allocation record.

## Raid

Represents simple PvP activity.

Rules:

- A kingdom sends units to raid another kingdom.
- The raid resolves into rewards, losses, and a report.
- Real-time combat is out of scope.

## Report

Represents the result of a mission, raid, battle, or event.

Phase 13 report type:

- `pve_mission`

Fields:

- `id`: UUID
- `kingdomId`: kingdom UUID
- `missionId`: nullable mission UUID
- `type`: report type
- `title`: report title
- `body`: readable report body
- `phasesJson`: ordered narrative phase payload
- `result`: `success`, `partial_success`, or `failure`
- `rewardsJson`: resource rewards payload
- `lossesJson`: unit losses payload
- `isRead`: read/unread flag
- `createdAt`: report creation timestamp

Rules:

- Mission reports are created when PvE missions resolve.
- PvE mission reports use local deterministic narrative templates keyed by mission and result.
- Report phases are ordered title/body sections for the existing report detail view.
- Existing reports without phase data are treated as having an empty phase list.
- Reports can be marked read idempotently.
- Report ownership is scoped to the authenticated user's kingdom.
- Report delivery is pull-based through the API; there are no comments, notifications, WebSocket, or background jobs in Phase 13.
- Raid reports and event reports are later phases.

## Event

Represents a simple narrative or strategic choice.

Rules:

- Events offer a small set of choices.
- Choices can grant resources, modify production, affect units, or influence patron alignment.
