# API Contract

This is the API contract for the implemented Playtest 001 MVP.

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

## Missions

### `GET /api/missions/available`

Returns configured PvE missions for the current authenticated user's kingdom.

Requires:

```http
Authorization: Bearer <token>
```

Response:

```json
{
  "missions": [
    {
      "key": "black_forest_expedition",
      "label": "Чёрный Лес",
      "type": "expedition",
      "description": "Охотники шепчут, что в Чёрном Лесу пропадают тропы, но старые склады всё ещё можно найти.",
      "durationSeconds": 120,
      "minimumRequirements": {
        "totalUnits": 5,
        "scouts": 0
      },
      "baseRewards": {
        "gold": 40,
        "food": 80,
        "wood": 120,
        "stone": 0,
        "population": 0
      },
      "risk": "medium"
    }
  ]
}
```

### `GET /api/missions/me`

Returns the current authenticated user's missions after applying lazy mission completion.

Requires:

```http
Authorization: Bearer <token>
```

Response:

```json
{
  "missions": [
    {
      "id": "uuid",
      "missionKey": "black_forest_expedition",
      "missionLabel": "Чёрный Лес",
      "missionType": "expedition",
      "status": "active",
      "startedAt": "2026-06-29T00:00:00Z",
      "finishesAt": "2026-06-29T00:02:00Z",
      "completedAt": null,
      "units": [
        {
          "unitType": "militia",
          "unitLabel": "Ополчение",
          "amountSent": 5,
          "amountLost": 0,
          "amountReturned": 0
        }
      ],
      "result": null
    }
  ]
}
```

### `POST /api/missions/start`

Starts a PvE mission and immediately removes sent units from the available army.

Requires:

```http
Authorization: Bearer <token>
```

Request:

```json
{
  "missionKey": "black_forest_expedition",
  "units": [
    {
      "unitType": "militia",
      "amount": 5
    },
    {
      "unitType": "scouts",
      "amount": 1
    }
  ]
}
```

Response:

```json
{
  "mission": {
    "id": "uuid",
    "missionKey": "black_forest_expedition",
    "missionLabel": "Чёрный Лес",
    "missionType": "expedition",
    "status": "active",
    "startedAt": "2026-06-29T00:00:00Z",
    "finishesAt": "2026-06-29T00:02:00Z",
    "completedAt": null,
    "units": [
      {
        "unitType": "militia",
        "unitLabel": "Ополчение",
        "amountSent": 5,
        "amountLost": 0,
        "amountReturned": 0
      }
    ],
    "result": null
  },
  "army": {
    "kingdomId": "uuid",
    "units": [],
    "trainingOrders": [],
    "summary": {
      "totalUnits": 7,
      "totalAttack": 14,
      "totalDefense": 21,
      "totalSupply": 7
    }
  }
}
```

Errors:

- `kingdom_not_found`
- `invalid_mission_key`
- `invalid_unit_type`
- `invalid_unit_amount`
- `insufficient_units`
- `mission_requirements_not_met`

## Raids

### `GET /api/neighbors`

Returns possible PvP raid targets. The current user's kingdom is excluded and exact defender resources or unit counts are not exposed.

Requires:

```http
Authorization: Bearer <token>
```

Response:

```json
{
  "neighbors": [
    {
      "kingdomId": "uuid",
      "name": "Серый Посад",
      "culture": "free_posad",
      "patron": "old_pact",
      "dread": 2,
      "powerEstimate": "similar",
      "canRaid": true,
      "blockedReason": null
    }
  ]
}
```

Blocked reasons:

- `target_newbie_protected`
- `target_too_weak`
- `raid_cooldown_active`
- `target_under_protection`

### `GET /api/raids/me`

Returns raids where the current user's kingdom is attacker or defender. Completed raids are resolved lazily before returning.

Requires:

```http
Authorization: Bearer <token>
```

Response:

```json
{
  "raids": [
    {
      "id": "uuid",
      "attackerKingdomId": "uuid",
      "attackerKingdomName": "Воронья Сечь",
      "defenderKingdomId": "uuid",
      "defenderKingdomName": "Серый Посад",
      "status": "active",
      "result": null,
      "startedAt": "2026-06-29T00:00:00Z",
      "arrivesAt": "2026-06-29T00:02:00Z",
      "completedAt": null,
      "units": [
        {
          "unitType": "militia",
          "unitLabel": "Ополчение",
          "amountSent": 5,
          "amountLost": 0,
          "amountReturned": 0
        }
      ],
      "loot": {
        "gold": 0,
        "food": 0,
        "wood": 0,
        "stone": 0,
        "population": 0
      }
    }
  ]
}
```

### `POST /api/raids/start`

Starts an asynchronous raid against another kingdom. Sent attacker units are removed immediately and returned when the raid resolves.

Requires:

```http
Authorization: Bearer <token>
```

Request:

```json
{
  "defenderKingdomId": "uuid",
  "units": [
    {
      "unitType": "militia",
      "amount": 5
    },
    {
      "unitType": "scouts",
      "amount": 1
    }
  ]
}
```

Response:

```json
{
  "raid": {
    "id": "uuid",
    "attackerKingdomId": "uuid",
    "attackerKingdomName": "Воронья Сечь",
    "defenderKingdomId": "uuid",
    "defenderKingdomName": "Серый Посад",
    "status": "active",
    "result": null,
    "startedAt": "2026-06-29T00:00:00Z",
    "arrivesAt": "2026-06-29T00:02:00Z",
    "completedAt": null,
    "units": [],
    "loot": {
      "gold": 0,
      "food": 0,
      "wood": 0,
      "stone": 0,
      "population": 0
    }
  },
  "army": {
    "kingdomId": "uuid",
    "units": [],
    "trainingOrders": [],
    "summary": {
      "totalUnits": 7,
      "totalAttack": 14,
      "totalDefense": 21,
      "totalSupply": 7
    }
  }
}
```

Errors:

- `kingdom_not_found`
- `target_not_found`
- `cannot_raid_self`
- `invalid_unit_type`
- `invalid_unit_amount`
- `insufficient_units`
- `raid_requirements_not_met`
- `target_newbie_protected`
- `target_too_weak`
- `raid_cooldown_active`
- `target_under_protection`

## Reports

### `GET /api/reports/me`

Returns the current authenticated user's mission, raid, and event reports after applying lazy completion.

Query parameters:

- `limit`: optional, defaults to 20, maximum 50
- `offset`: optional, defaults to 0

Requires:

```http
Authorization: Bearer <token>
```

Response:

```json
{
  "reports": [
    {
      "id": "uuid",
      "type": "pve_mission",
      "title": "Чёрный лес уступил дорогу",
      "body": "Отряд вышел из чащи до заката. За спинами пахло смолой, сырой землёй и чужим страхом.",
      "phases": [
        {
          "title": "На опушке",
          "body": "Разведчики нашли старую звериную тропу и повели людей мимо топей."
        }
      ],
      "result": "success",
      "rewards": {
        "gold": 40,
        "food": 80,
        "wood": 120,
        "stone": 0,
        "population": 0
      },
      "losses": {
        "militia": 1,
        "scouts": 0
      },
      "isRead": false,
      "createdAt": "2026-06-29T00:02:00Z"
    }
  ],
  "pagination": {
    "limit": 20,
    "offset": 0
  },
  "unreadCount": 1
}
```

### `GET /api/reports/{id}`

Returns one report owned by the current authenticated user's kingdom.

Requires:

```http
Authorization: Bearer <token>
```

Response:

```json
{
  "report": {
    "id": "uuid",
    "type": "pve_mission",
    "title": "Чёрный лес уступил дорогу",
    "body": "Отряд вышел из чащи до заката. За спинами пахло смолой, сырой землёй и чужим страхом.",
    "phases": [],
    "result": "success",
    "rewards": {
      "gold": 40,
      "food": 80,
      "wood": 120,
      "stone": 0,
      "population": 0
    },
    "losses": {},
    "isRead": false,
    "createdAt": "2026-06-29T00:02:00Z"
  }
}
```

Errors:

- `kingdom_not_found`
- `report_not_found`

### `POST /api/reports/{id}/read`

Marks one report owned by the current authenticated user's kingdom as read. The operation is idempotent.

Requires:

```http
Authorization: Bearer <token>
```

Response:

```json
{
  "report": {
    "id": "uuid",
    "type": "pve_mission",
    "title": "Чёрный лес уступил дорогу",
    "body": "Отряд вышел из чащи до заката. За спинами пахло смолой, сырой землёй и чужим страхом.",
    "phases": [],
    "result": "success",
    "rewards": {
      "gold": 40,
      "food": 80,
      "wood": 120,
      "stone": 0,
      "population": 0
    },
    "losses": {},
    "isRead": true,
    "createdAt": "2026-06-29T00:02:00Z"
  }
}
```

Errors:

- `kingdom_not_found`
- `report_not_found`

## Events

### `GET /api/events/me`

Returns current active events and recent resolved or expired events for the authenticated user's kingdom. The read lazily expires stale events and lazily generates new active events up to the MVP active limit.

Requires:

```http
Authorization: Bearer <token>
```

Query params:

- `includeResolved`: optional boolean, default `true`
- `limit`: optional integer, default `20`, max `50`

Response:

```json
{
  "events": [
    {
      "id": "uuid",
      "eventKey": "found_old_idol",
      "category": "economy",
      "title": "Старый идол в лесу",
      "body": "Лесорубы нашли в корнях чёрный идол...",
      "status": "active",
      "generatedAt": "2026-06-30T08:00:00Z",
      "expiresAt": "2026-07-01T08:00:00Z",
      "resolvedAt": null,
      "selectedChoiceKey": null,
      "choices": [
        {
          "key": "sell_to_merchants",
          "label": "Продать купцам",
          "description": "Купцы не задают вопросов, если цена хорошая."
        }
      ],
      "result": null
    }
  ],
  "activeCount": 1
}
```

Errors:

- `kingdom_not_found`

### `POST /api/events/{id}/choose`

Resolves one active event choice for the authenticated user's kingdom. Effects are applied once, safely clamped so resources and units do not go below zero and population does not go below one. Resolved and expired events cannot be chosen.

Requires:

```http
Authorization: Bearer <token>
```

Request:

```json
{
  "choiceKey": "sell_to_merchants"
}
```

Response:

```json
{
  "event": {
    "id": "uuid",
    "eventKey": "found_old_idol",
    "category": "economy",
    "title": "Старый идол в лесу",
    "body": "Лесорубы нашли в корнях чёрный идол...",
    "status": "resolved",
    "generatedAt": "2026-06-30T08:00:00Z",
    "expiresAt": "2026-07-01T08:00:00Z",
    "resolvedAt": "2026-06-30T08:05:00Z",
    "selectedChoiceKey": "sell_to_merchants",
    "choices": [],
    "result": {
      "title": "Идол ушёл с купцами",
      "body": "Купцы забрали находку без лишних слов. В казне стало тяжелее, но старики смотрят на лес тревожнее.",
      "appliedEffects": {
        "resourceDelta": {
          "gold": 80
        },
        "kingdomDelta": {
          "honor": -1
        }
      }
    }
  },
  "resources": {
    "kingdomId": "uuid",
    "gold": 580,
    "food": 300,
    "wood": 300,
    "stone": 200,
    "population": 100,
    "productionPerHour": {
      "gold": 20,
      "food": 30,
      "wood": 25,
      "stone": 15,
      "population": 1
    },
    "lastCalculatedAt": "2026-06-30T08:05:00Z",
    "updatedAt": "2026-06-30T08:05:00Z"
  },
  "kingdom": {
    "id": "uuid",
    "userId": "uuid",
    "name": "Воронья Сечь",
    "culture": "northern_principality",
    "patron": null,
    "dread": 0,
    "honor": 0,
    "createdAt": "2026-06-29T00:00:00Z",
    "updatedAt": "2026-06-30T08:05:00Z"
  }
}
```

Errors:

- `kingdom_not_found`
- `event_not_found`
- `event_expired`
- `event_already_resolved`
- `invalid_event_choice`
- `event_choice_not_available`

## Patron

### `GET /api/patron/options`

Returns all patron choices.

Requires:

```http
Authorization: Bearer <token>
```

Response:

```json
{
  "patrons": [
    {
      "key": "independent",
      "label": "Независимость",
      "shortDescription": "Ты никому не служишь и никому не платишь. Свобода полная, защита только своя.",
      "flavor": "Свободные владетели живут без печатей и клятв. Но когда приходит беда, никто не обязан идти им на помощь.",
      "currentEffects": [
        "Нет дани",
        "Нет защиты",
        "Полная свобода решений"
      ],
      "futureEffects": [
        "Безопасность зависит только от собственных сил"
      ]
    }
  ]
}
```

### `GET /api/patron/me`

Returns the current authenticated user's patron status.

Requires:

```http
Authorization: Bearer <token>
```

Response when no patron is selected:

```json
{
  "patron": null,
  "availablePatrons": [
    "independent",
    "empire_of_dusk",
    "old_pact"
  ]
}
```

Response when a patron is selected:

```json
{
  "patron": {
    "id": "uuid",
    "kingdomId": "uuid",
    "key": "empire_of_dusk",
    "label": "Империя Заката",
    "favor": 0,
    "standing": "neutral",
    "joinedAt": "2026-06-29T00:00:00Z",
    "leftAt": null,
    "currentEffects": [
      "Дань ещё не активна в этой версии",
      "Защита ещё не активна в этой версии"
    ],
    "futureEffects": [
      "Позже Империя сможет требовать дань",
      "Позже Империя сможет давать защиту"
    ]
  },
  "availablePatrons": [
    "independent",
    "empire_of_dusk",
    "old_pact"
  ]
}
```

Errors:

- `kingdom_not_found`

### `POST /api/patron/join`

Joins or switches the current authenticated user's patron. Joining the current patron is idempotent.

Requires:

```http
Authorization: Bearer <token>
```

Request:

```json
{
  "patron": "old_pact"
}
```

Response:

```json
{
  "patron": {
    "id": "uuid",
    "kingdomId": "uuid",
    "key": "old_pact",
    "label": "Старый Договор",
    "favor": 0,
    "standing": "neutral",
    "joinedAt": "2026-06-29T00:00:00Z",
    "leftAt": null,
    "currentEffects": [
      "Обязательства ещё не активны в этой версии",
      "Помощь ещё не активна в этой версии"
    ],
    "futureEffects": [
      "Позже Старый Договор сможет требовать вклад",
      "Позже Старый Договор сможет помогать в обороне"
    ]
  },
  "kingdom": {
    "id": "uuid",
    "patron": "old_pact"
  }
}
```

Errors:

- `kingdom_not_found`
- `invalid_patron`

### `POST /api/patron/break`

Breaks the current patron relation if one exists. The operation is idempotent and clears `kingdom.patron`.

Requires:

```http
Authorization: Bearer <token>
```

Response:

```json
{
  "patron": null,
  "kingdom": {
    "id": "uuid",
    "patron": null
  }
}
```

Errors:

- `kingdom_not_found`

### `GET /api/patron/pressure`

Returns the authenticated kingdom's current patron pressure state after lazy resolution. If no patron is selected, `pressure` is `null`.

Requires:

```http
Authorization: Bearer <token>
```

Response:

```json
{
  "pressure": {
    "patron": "empire_of_dusk",
    "patronLabel": "Империя Заката",
    "pressureLevel": 35,
    "crisisStatus": "warning",
    "tributeDebt": {
      "gold": 40,
      "food": 10
    },
    "contributionDebt": {
      "food": 0
    },
    "nextTributeAt": "2026-06-30T08:00:00Z",
    "delayUntil": null,
    "summary": "Империя ждёт дань, но не может забрать неприкосновенный запас.",
    "availableActions": [
      "pay_tribute",
      "break_patron",
      "ask_delay"
    ],
    "protectedMinimums": {
      "gold": 150,
      "food": 150,
      "wood": 100,
      "stone": 75,
      "population": null
    }
  }
}
```

Errors:

- `kingdom_not_found`

### `POST /api/patron/pay-tribute`

Resolves patron pressure lazily, then pays available tribute or Old Pact contribution from resources above protected minimums.

Requires:

```http
Authorization: Bearer <token>
```

Response:

```json
{
  "pressure": {
    "patron": "empire_of_dusk",
    "patronLabel": "Империя Заката",
    "pressureLevel": 20,
    "crisisStatus": "none",
    "tributeDebt": {
      "gold": 0,
      "food": 0
    },
    "contributionDebt": {
      "food": 0
    },
    "nextTributeAt": "2026-06-30T09:00:00Z",
    "delayUntil": null,
    "summary": "Империя ждёт дань, но не может забрать неприкосновенный запас.",
    "availableActions": [
      "break_patron"
    ],
    "protectedMinimums": {
      "gold": 150,
      "food": 150,
      "wood": 100,
      "stone": 75,
      "population": null
    }
  },
  "resources": {
    "kingdomId": "uuid",
    "gold": 150,
    "food": 150,
    "wood": 100,
    "stone": 75,
    "population": 50,
    "productionPerHour": {
      "gold": 10,
      "food": 15,
      "wood": 8,
      "stone": 5,
      "population": 0
    },
    "lastCalculatedAt": "2026-06-30T08:00:00Z",
    "updatedAt": "2026-06-30T08:00:00Z"
  }
}
```

Errors:

- `kingdom_not_found`
- `no_patron_selected`
- `no_tribute_due`

### `POST /api/patron/crisis-choice`

Applies a simple pressure crisis choice. Supported choices are `ask_delay` and `break_patron`.

Requires:

```http
Authorization: Bearer <token>
```

Request:

```json
{
  "choice": "ask_delay"
}
```

Response for `ask_delay`:

```json
{
  "pressure": {
    "patron": "empire_of_dusk",
    "patronLabel": "Империя Заката",
    "pressureLevel": 40,
    "crisisStatus": "delayed",
    "tributeDebt": {
      "gold": 40,
      "food": 10
    },
    "contributionDebt": {
      "food": 0
    },
    "nextTributeAt": "2026-06-30T08:00:00Z",
    "delayUntil": "2026-06-30T10:00:00Z",
    "summary": "Имперский сборщик получил отсрочку, но запомнил слабость двора.",
    "availableActions": [
      "pay_tribute",
      "break_patron"
    ],
    "protectedMinimums": {
      "gold": 150,
      "food": 150,
      "wood": 100,
      "stone": 75,
      "population": null
    }
  }
}
```

Response for `break_patron`:

```json
{
  "pressure": null,
  "kingdom": {
    "id": "uuid",
    "patron": null
  }
}
```

Errors:

- `kingdom_not_found`
- `no_patron_selected`
- `invalid_crisis_choice`
- `crisis_choice_not_available`

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

Mission types:

- `expedition`
- `scouting`

Mission statuses:

- `active`
- `completed`

Mission results:

- `success`
- `partial_success`
- `failure`

Report types:

- `pve_mission`
- `pvp_raid_attacker`
- `pvp_raid_defender`
- `event`

Raid results:

- `attacker_success`
- `defender_success`
- `bloody_stalemate`
- `repelled_by_protection`
