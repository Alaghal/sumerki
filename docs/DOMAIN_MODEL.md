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
- `dread`: simple aggression/infamy value
- `honor`: reserved simple reputation value
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

MVP patron choices:

- `independent`: no tribute, no protection, full freedom
- `empire_of_dusk`: order, tribute debt, and pressure
- `old_pact`: oath and soft contribution obligation

Rules:

- Patron choice updates `kingdom.patron`.
- `kingdom.patron` is null before first choice.
- Choosing `independent` stores `independent`.
- Breaking a patron relation clears `kingdom.patron` to null.
- Empire tribute and Old Pact contribution pressure are resolved lazily.
- No patron military help, patron quests, patron orders, cooldowns, switching penalties, NPC raids, or retaliation are included in the MVP pressure system.

## PatronRelation

Represents the current MVP patron relation for one kingdom.

Fields:

- `id`: UUID
- `kingdomId`: kingdom UUID
- `patron`: patron key
- `favor`: integer from -100 to 100, defaults to 0
- `standing`: `hostile`, `cold`, `neutral`, `warm`, or `loyal`
- `joinedAt`: relation join timestamp
- `leftAt`: nullable timestamp set when the relation is broken
- `createdAt`: row creation timestamp
- `updatedAt`: last update timestamp

Rules:

- Each kingdom has at most one patron relation row in the MVP.
- Joining a patron creates or updates the existing relation row.
- Joining the current patron is idempotent.
- Joining a different patron switches the relation without cooldown or penalty in Phase 14.
- Breaking a relation sets `leftAt` and clears `kingdom.patron`.
- Existing kingdoms with non-null `kingdom.patron` can be backfilled into a matching relation.
- Patron state is owned by the authenticated user's kingdom; clients never submit `kingdomId`.

## PatronPressureState

Represents lazy tribute pressure for one kingdom.

Fields:

- `id`: UUID
- `kingdomId`: kingdom UUID
- `patron`: patron key
- `tributeDebtGold`: unpaid Empire gold tribute
- `tributeDebtFood`: unpaid Empire food tribute
- `contributionDebtFood`: unpaid Old Pact food contribution
- `pressureLevel`: integer from 0 to 100
- `crisisStatus`: `none`, `warning`, `active`, or `delayed`
- `crisisStartedAt`: nullable crisis start timestamp reserved for later polish
- `nextTributeAt`: next lazy tribute or contribution resolution timestamp
- `lastResolvedAt`: last pressure resolution timestamp
- `delayUntil`: nullable delay timestamp after asking for more time
- `createdAt`: row creation timestamp
- `updatedAt`: last update timestamp

Rules:

- Each kingdom has at most one pressure state row.
- `independent` has no tribute, contribution debt, or pressure.
- `empire_of_dusk` periodically takes tribute from surplus gold and food above protected minimums; unpaid amounts become debt and can raise pressure.
- `old_pact` tracks a soft food contribution debt and caps pressure at a low level.
- Paying tribute or contribution spends only resources above protected minimums.
- Asking for delay sets `crisisStatus` to `delayed`, sets `delayUntil`, and increases pressure slightly.
- Breaking the patron relation clears pressure and debt.
- Tribute, contribution, and pressure are resolved lazily on patron pressure reads, tribute payments, crisis choices, patron joins, and patron status reads.
- No background workers, cron jobs, Redis, queues, real-time ticking, NPC retaliation, patron armies, or dark god events are used.

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

Represents asynchronous PvP activity between kingdoms.

Fields:

- `id`: UUID
- `attackerKingdomId`: attacking kingdom UUID
- `defenderKingdomId`: defending kingdom UUID
- `status`: `active` or `completed`
- `startedAt`: raid start timestamp
- `arrivesAt`: lazy completion timestamp
- `completedAt`: nullable completion timestamp
- `result`: nullable raid result
- `lootJson`: resource loot payload
- `attackerLossesJson`: attacker unit losses payload
- `defenderLossesJson`: defender pressure/loss payload
- `resultJson`: full result payload
- `createdAt`: row creation timestamp
- `updatedAt`: last update timestamp

Rules:

- A kingdom sends available units to raid another kingdom.
- Sent attacker units are subtracted immediately and unavailable until the raid resolves.
- Defender units are snapshotted for scoring but are not removed before the raid.
- Raid completion is resolved lazily when raids or reports are read, or when another raid starts.
- Result is deterministic in Phase 15: `attacker_success`, `defender_success`, or `bloody_stalemate`.
- Attacker survivors are returned on completion.
- Successful raids steal limited gold, food, wood, and stone only.
- Population is never stolen.
- Defender resources cannot be stolen below protected minimums.
- Defender city and buildings cannot be destroyed or captured.
- Weak-player protection blocks raids against new kingdoms, much weaker targets, repeated same-target raids, and over-farmed defenders.
- Starting and resolving raids increases attacker dread.
- Patron state can add small defensive flavor modifiers, but no patron armies or costs are implemented.
- No territory capture, sieges, alliances, diplomacy, map routes, real-time combat, revenge, bounties, NPC retaliation, or patron military help are included in Phase 15.

## RaidUnit

Represents units allocated or snapshotted for a raid.

Fields:

- `id`: UUID
- `raidId`: raid UUID
- `side`: `attacker` or `defender`
- `unitType`: unit type
- `amountSent`: units sent or snapshotted
- `amountLost`: units lost in raid result
- `amountReturned`: surviving units returned or remaining in snapshot
- `createdAt`: row creation timestamp
- `updatedAt`: last update timestamp

Rules:

- Attacker raid unit rows represent unavailable sent units.
- Defender raid unit rows are a snapshot used for scoring and reports.
- Defender losses are intentionally very limited and not applied to defender `kingdom_units` in Phase 15.

## Report

Represents the result of a mission, raid, battle, or event.

Phase 15 report types:

- `pve_mission`
- `pvp_raid_attacker`
- `pvp_raid_defender`
- `event`

Fields:

- `id`: UUID
- `kingdomId`: kingdom UUID
- `missionId`: nullable mission UUID
- `type`: report type
- `title`: report title
- `body`: readable report body
- `phasesJson`: ordered narrative phase payload
- `result`: mission or raid result
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
- Raid reports are created for both attacker and defender when a raid resolves.
- Event reports are created when event choices resolve.

## GameEvent

Catalog template for a choice-based event.

Fields:

- `id`: UUID
- `eventKey`: stable event key
- `category`: `economy`, `ruler`, `military`, `patron`, or `dark_omen`
- `title`: event title
- `body`: event setup text
- `triggerType`: currently `pool`
- `weight`: positive integer reserved for future weighted selection
- `isActive`: whether the template can be generated
- `cooldownSeconds`: per-kingdom cooldown before the same event key can appear again
- `expiresAfterSeconds`: lifetime for generated event instances
- `conditionsJson`: simple generation conditions
- `createdAt`: row creation timestamp
- `updatedAt`: last update timestamp

## EventChoice

Choice available for a `GameEvent`.

Fields:

- `id`: UUID
- `gameEventId`: game event UUID
- `choiceKey`: stable choice key within the event
- `label`: choice button text
- `description`: short choice description
- `effectsJson`: configured effects payload
- `resultTitle`: report/result title after choosing
- `resultBody`: report/result body after choosing
- `createdAt`: row creation timestamp
- `updatedAt`: last update timestamp

Rules:

- Choice keys are unique per game event.
- Unknown effects are ignored by the MVP service after JSON decoding.

## KingdomEvent

Generated event instance for one kingdom.

Fields:

- `id`: UUID
- `kingdomId`: kingdom UUID
- `gameEventId`: game event UUID
- `status`: `active`, `resolved`, or `expired`
- `generatedAt`: generation timestamp
- `expiresAt`: expiry timestamp
- `resolvedAt`: nullable resolution timestamp
- `selectedChoiceKey`: nullable selected choice key
- `resultJson`: nullable result payload with applied effects
- `createdAt`: row creation timestamp
- `updatedAt`: last update timestamp

Rules:

- Events are generated and expired lazily on event reads and choice commands.
- Each kingdom can have up to three active events.
- The same event key cannot be active twice for one kingdom.
- The same event key is not regenerated during its cooldown window after generation.
- Expired events cannot be chosen and are not deleted.
- Resolved events cannot be chosen again, preventing duplicate effect application.

## Event Effects

Optional `effectsJson` keys:

- `resourceDelta`: `gold`, `food`, `wood`, `stone`, `population`
- `unitDelta`: `militia`, `spearmen`, `archers`, `cavalry`, `scouts`
- `kingdomDelta`: `dread`, `honor`
- `patronFavorDelta`: integer favor change

Rules:

- Resource deltas can be positive or negative.
- Resources and units never go below zero from event effects.
- Population never goes below one.
- Unit deltas affect only currently available units, not units away on missions or raids.
- Dread and honor never go below zero.
- Patron favor is clamped from -100 to 100.
- Patron favor effects are ignored if the kingdom has no active patron relation.
- Tribute, debt, raid, and mission effects are not included in Phase 17.

## Event Conditions

Supported `conditionsJson` key:

- `requiresPatron`: `any`, `none`, `independent`, `empire_of_dusk`, or `old_pact`

Rules:

- Missing conditions mean the event is eligible.
- `any` requires a non-null kingdom patron.
- `none` requires no selected patron.
- A specific patron value requires an exact match.
- No complex conditions, map/province events, chains, recurring chains, or dark god avatar logic are included in Phase 17.

## Event Content

Represents a simple narrative or strategic choice.

Rules:

- Phase 17 seeds five smoke-test events, one per category.
- Phase 18 expands the active catalog to 25 event templates.
- Phase 18 distribution is five active events per category: economy, ruler, military, patron, and dark omen.
- Phase 18 content uses the existing Phase 17 effect and condition schemas.
- Event reports use simple local text and three phases: `Событие`, `Выбор`, and `Последствия`.
- Event chains, advanced conditions, map/province events, alliance events, and dark god avatar mechanics remain out of scope.
