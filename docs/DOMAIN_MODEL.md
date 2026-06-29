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

Initial likely resources:

- food
- wood
- stone
- iron
- silver

Exact resource names and production rules should be finalized before implementation.

## Building

Represents an upgradeable structure in a kingdom.

Possible MVP examples:

- town hall
- lumber camp
- quarry
- barracks
- shrine

Buildings should have levels and resource costs.

## Unit

Represents trainable military forces.

Possible MVP examples:

- militia
- scout
- spearman

Units should have training costs, training time or instant MVP creation rules, and combat stats.

## Mission

Represents PvE activity.

Rules:

- A kingdom sends units to a mission.
- The mission resolves into rewards, losses, and a report.
- MVP resolution can be simple and deterministic if needed.

## Raid

Represents simple PvP activity.

Rules:

- A kingdom sends units to raid another kingdom.
- The raid resolves into rewards, losses, and a report.
- Real-time combat is out of scope.

## Report

Represents the result of a mission, raid, battle, or event.

Fields should include:

- report type
- title
- body
- rewards
- losses
- created timestamp
- read/unread state

## Event

Represents a simple narrative or strategic choice.

Rules:

- Events offer a small set of choices.
- Choices can grant resources, modify production, affect units, or influence patron alignment.
