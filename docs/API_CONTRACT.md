# API Contract

This is the initial draft API contract for future backend phases. It is documentation only in Phase 0.

## General Conventions

- Base path: `/api`
- Request and response format: JSON
- Authentication: bearer JWT after login
- Error format: consistent JSON errors
- IDs: UUID strings
- Timestamps: ISO 8601 strings

## Error Shape

```json
{
  "error": {
    "code": "validation_error",
    "message": "Human-readable error message"
  }
}
```

## Health

### `GET /health`

Returns process health. This endpoint does not require database connectivity.

Response:

```json
{
  "status": "ok"
}
```

### `GET /ready`

Returns service readiness.

Response when the database is reachable:

```json
{
  "status": "ready",
  "database": "ok"
}
```

Response when the database is not reachable:

HTTP 503

```json
{
  "error": {
    "code": "database_unavailable",
    "message": "Database is not reachable"
  }
}
```

## Auth

### `POST /api/auth/register`

Registers a new user.

Request:

```json
{
  "email": "player@example.com",
  "password": "secret-password"
}
```

Success response:

HTTP 201

```json
{
  "user": {
    "id": "uuid",
    "email": "player@example.com"
  },
  "token": "jwt"
}
```

Validation errors:

- `invalid_email`
- `password_too_short`
- `email_already_exists`

### `POST /api/auth/login`

Logs in an existing user.

Request:

```json
{
  "email": "player@example.com",
  "password": "secret-password"
}
```

Success response:

HTTP 200

```json
{
  "user": {
    "id": "uuid",
    "email": "player@example.com"
  },
  "token": "jwt"
}
```

Login errors use the generic invalid credentials response and do not reveal whether the email exists:

```json
{
  "error": {
    "code": "invalid_credentials",
    "message": "Invalid email or password"
  }
}
```

### `GET /api/me`

Returns the current authenticated user.

Requires:

```http
Authorization: Bearer <token>
```

Response:

```json
{
  "user": {
    "id": "uuid",
    "email": "player@example.com"
  }
}
```

Auth errors:

- `missing_authorization_header`
- `invalid_authorization_header`
- `invalid_token`
- `expired_token`
- `user_not_found`

## Kingdoms

### `POST /api/kingdoms`

Creates the current user's kingdom.

Requires:

```http
Authorization: Bearer <token>
```

Request:

```json
{
  "name": "Воронья Сечь",
  "culture": "northern_principality"
}
```

Success response:

HTTP 201

```json
{
  "kingdom": {
    "id": "uuid",
    "userId": "uuid",
    "name": "Воронья Сечь",
    "culture": "northern_principality",
    "patron": null,
    "createdAt": "2026-06-29T00:00:00Z",
    "updatedAt": "2026-06-29T00:00:00Z"
  }
}
```

Validation errors:

- `kingdom_name_too_short`
- `kingdom_name_too_long`
- `invalid_culture`
- `kingdom_already_exists`

### `GET /api/kingdoms/me`

Returns the current user's kingdom.

Requires:

```http
Authorization: Bearer <token>
```

Response when the current user has a kingdom:

```json
{
  "kingdom": {
    "id": "uuid",
    "userId": "uuid",
    "name": "Воронья Сечь",
    "culture": "northern_principality",
    "patron": null,
    "createdAt": "2026-06-29T00:00:00Z",
    "updatedAt": "2026-06-29T00:00:00Z"
  }
}
```

Response when the current user has no kingdom:

```json
{
  "kingdom": null
}
```

## Ruler

### `GET /api/ruler/me`

Returns the current authenticated user's ruler.

Requires:

```http
Authorization: Bearer <token>
```

Response:

```json
{
  "ruler": {
    "id": "uuid",
    "kingdomId": "uuid",
    "name": "Боривой",
    "age": 42,
    "culture": "northern_principality",
    "authority": 61,
    "courage": 74,
    "cunning": 44,
    "honor": 68,
    "cruelty": 31,
    "ambition": 55,
    "paranoia": 29,
    "healthStatus": "healthy",
    "createdAt": "2026-06-29T00:00:00Z",
    "updatedAt": "2026-06-29T00:00:00Z"
  }
}
```

Response when the current user has no kingdom:

HTTP 404

```json
{
  "error": {
    "code": "kingdom_not_found",
    "message": "Create a kingdom before requesting a ruler"
  }
}
```

## Resources

### `GET /api/resources/me`

Returns the current authenticated user's resources after applying lazy production.

Requires:

```http
Authorization: Bearer <token>
```

Response:

```json
{
  "resources": {
    "kingdomId": "uuid",
    "gold": 520,
    "food": 330,
    "wood": 325,
    "stone": 215,
    "population": 101,
    "productionPerHour": {
      "gold": 20,
      "food": 30,
      "wood": 25,
      "stone": 15,
      "population": 1
    },
    "lastCalculatedAt": "2026-06-29T00:00:00Z",
    "updatedAt": "2026-06-29T00:00:00Z"
  }
}
```

Response when the current user has no kingdom:

HTTP 404

```json
{
  "error": {
    "code": "kingdom_not_found",
    "message": "Create a kingdom before requesting resources"
  }
}
```

## Buildings

### `GET /api/buildings/me`

Returns all buildings for the current authenticated user's kingdom after applying lazy upgrade completion.

Requires:

```http
Authorization: Bearer <token>
```

Response:

```json
{
  "buildings": [
    {
      "id": "uuid",
      "kingdomId": "uuid",
      "type": "town_hall",
      "label": "Ратуша",
      "level": 1,
      "maxLevel": 5,
      "isUpgrading": false,
      "upgradeStartedAt": null,
      "upgradeFinishesAt": null,
      "nextUpgrade": {
        "targetLevel": 2,
        "cost": {
          "gold": 300,
          "food": 0,
          "wood": 240,
          "stone": 200,
          "population": 0
        },
        "durationSeconds": 120,
        "canUpgrade": true,
        "blockedReason": null
      },
      "effects": [
        "Нет эффекта в текущей версии"
      ],
      "createdAt": "2026-06-29T00:00:00Z",
      "updatedAt": "2026-06-29T00:00:00Z"
    }
  ]
}
```

Response when the current user has no kingdom:

HTTP 404

```json
{
  "error": {
    "code": "kingdom_not_found",
    "message": "Create a kingdom before requesting buildings"
  }
}
```

### `POST /api/buildings/{type}/upgrade`

Starts an upgrade for one building and spends resources immediately.

Requires:

```http
Authorization: Bearer <token>
```

Response:

```json
{
  "building": {
    "id": "uuid",
    "kingdomId": "uuid",
    "type": "farm",
    "label": "Ферма",
    "level": 1,
    "maxLevel": 5,
    "isUpgrading": true,
    "upgradeStartedAt": "2026-06-29T00:00:00Z",
    "upgradeFinishesAt": "2026-06-29T00:02:00Z",
    "nextUpgrade": null,
    "effects": [
      "+15 еды/час за уровень"
    ],
    "createdAt": "2026-06-29T00:00:00Z",
    "updatedAt": "2026-06-29T00:00:00Z"
  },
  "resources": {
    "kingdomId": "uuid",
    "gold": 420,
    "food": 300,
    "wood": 220,
    "stone": 180,
    "population": 100,
    "productionPerHour": {
      "gold": 30,
      "food": 45,
      "wood": 37,
      "stone": 25,
      "population": 1
    },
    "lastCalculatedAt": "2026-06-29T00:00:00Z",
    "updatedAt": "2026-06-29T00:00:00Z"
  }
}
```

Errors:

- `kingdom_not_found`
- `invalid_building_type`
- `building_not_found`
- `building_already_upgrading`
- `building_max_level`
- `insufficient_resources`

## Army

### `GET /api/army/me`

Returns the current authenticated user's units after applying lazy training completion.

Requires:

```http
Authorization: Bearer <token>
```

Response:

```json
{
  "army": {
    "kingdomId": "uuid",
    "units": [
      {
        "type": "militia",
        "label": "Ополчение",
        "amount": 10,
        "stats": {
          "attack": 2,
          "defense": 3,
          "speed": 2,
          "supply": 1
        },
        "cost": {
          "gold": 15,
          "food": 10,
          "wood": 0,
          "stone": 0,
          "population": 1
        },
        "secondsPerUnit": 5,
        "requirements": {
          "barracksLevel": 0,
          "isMet": true
        }
      }
    ],
    "trainingOrders": [
      {
        "id": "uuid",
        "unitType": "spearmen",
        "unitLabel": "Копейщики",
        "amount": 5,
        "status": "training",
        "startedAt": "2026-06-29T00:00:00Z",
        "finishesAt": "2026-06-29T00:01:00Z",
        "completedAt": null
      }
    ],
    "summary": {
      "totalUnits": 12,
      "totalAttack": 24,
      "totalDefense": 34,
      "totalSupply": 12
    }
  }
}
```

Response when the current user has no kingdom:

HTTP 404

```json
{
  "error": {
    "code": "kingdom_not_found",
    "message": "Create a kingdom before requesting army"
  }
}
```

### `POST /api/army/train`

Starts a unit training order and spends resources immediately.

Requires:

```http
Authorization: Bearer <token>
```

Request:

```json
{
  "unitType": "militia",
  "amount": 5
}
```

Response:

```json
{
  "trainingOrder": {
    "id": "uuid",
    "unitType": "militia",
    "unitLabel": "Ополчение",
    "amount": 5,
    "status": "training",
    "startedAt": "2026-06-29T00:00:00Z",
    "finishesAt": "2026-06-29T00:00:25Z",
    "completedAt": null
  },
  "resources": {
    "kingdomId": "uuid",
    "gold": 425,
    "food": 250,
    "wood": 300,
    "stone": 200,
    "population": 95,
    "productionPerHour": {
      "gold": 30,
      "food": 45,
      "wood": 37,
      "stone": 25,
      "population": 1
    },
    "lastCalculatedAt": "2026-06-29T00:00:00Z",
    "updatedAt": "2026-06-29T00:00:00Z"
  }
}
```

Errors:

- `kingdom_not_found`
- `invalid_unit_type`
- `invalid_training_amount`
- `insufficient_resources`
- `barracks_level_too_low`

## Enumerations

Cultures:

- `northern_principality`
- `lizard_grad`
- `free_posad`

Patrons:

- `independent`
- `empire_of_dusk`
- `old_pact`

Ruler health statuses:

- `healthy`
- `wounded`
- `sick`

Resources:

- `gold`
- `food`
- `wood`
- `stone`
- `population`

Building types:

- `town_hall`
- `farm`
- `lumberyard`
- `quarry`
- `barracks`
- `market`
- `walls`
- `shrine`

Unit types:

- `militia`
- `spearmen`
- `archers`
- `cavalry`
- `scouts`
