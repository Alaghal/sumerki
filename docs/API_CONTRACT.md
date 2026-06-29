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

## Enumerations

Cultures:

- `northern_principality`
- `lizard_grad`
- `free_posad`

Patrons:

- `independent`
- `empire_of_dusk`
- `old_pact`
