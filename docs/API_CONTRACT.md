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

Response:

```json
{
  "user": {
    "id": "uuid",
    "email": "player@example.com"
  },
  "token": "jwt"
}
```

### `POST /api/auth/login`

Logs in an existing user.

Request:

```json
{
  "email": "player@example.com",
  "password": "secret-password"
}
```

Response:

```json
{
  "user": {
    "id": "uuid",
    "email": "player@example.com"
  },
  "token": "jwt"
}
```

### `GET /api/me`

Returns the current authenticated user.

Response:

```json
{
  "user": {
    "id": "uuid",
    "email": "player@example.com"
  }
}
```

## Kingdoms

### `POST /api/kingdoms`

Creates the current user's kingdom.

Request:

```json
{
  "name": "Blackwater",
  "culture": "northern_principality"
}
```

Response:

```json
{
  "kingdom": {
    "id": "uuid",
    "userId": "uuid",
    "name": "Blackwater",
    "culture": "northern_principality",
    "patron": null,
    "createdAt": "2026-06-29T00:00:00Z"
  }
}
```

### `GET /api/kingdoms/me`

Returns the current user's kingdom.

Response:

```json
{
  "kingdom": {
    "id": "uuid",
    "userId": "uuid",
    "name": "Blackwater",
    "culture": "northern_principality",
    "patron": null,
    "createdAt": "2026-06-29T00:00:00Z"
  }
}
```

## Enumerations

Cultures:

- `northern_principality`
- `lizard_grad`
- `free_posad`

Patrons:

- `independent`
- `empire_of_dusk`
- `old_pact`
