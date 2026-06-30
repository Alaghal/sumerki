package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"sumerki/backend/internal/domain"
)

var ErrEventNotFound = errors.New("event not found")

type EventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) ListActiveGameEvents(ctx context.Context) ([]domain.GameEvent, error) {
	const query = `
		SELECT id::text, event_key, category, title, body, trigger_type, weight, is_active,
			cooldown_seconds, expires_after_seconds, conditions_json, created_at, updated_at
		FROM game_events
		WHERE is_active = true
		ORDER BY category ASC, event_key ASC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := []domain.GameEvent{}
	for rows.Next() {
		event, err := scanGameEvent(rows)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return events, nil
}

func (r *EventRepository) ListChoicesByGameEventID(ctx context.Context, gameEventID string) ([]domain.EventChoice, error) {
	const query = `
		SELECT id::text, game_event_id::text, choice_key, label, description, effects_json, result_title, result_body, created_at, updated_at
		FROM event_choices
		WHERE game_event_id = $1
		ORDER BY created_at ASC, choice_key ASC
	`
	rows, err := r.db.QueryContext(ctx, query, gameEventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	choices := []domain.EventChoice{}
	for rows.Next() {
		choice, err := scanEventChoice(rows)
		if err != nil {
			return nil, err
		}
		choices = append(choices, choice)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return choices, nil
}

func (r *EventRepository) ListActiveByKingdomID(ctx context.Context, kingdomID string) ([]domain.KingdomEvent, error) {
	const query = kingdomEventSelect + `
		WHERE ke.kingdom_id = $1 AND ke.status = 'active'
		ORDER BY ke.generated_at ASC, ge.category ASC, ge.event_key ASC
	`
	return r.listKingdomEvents(ctx, query, kingdomID)
}

func (r *EventRepository) ListRecentResolvedByKingdomID(ctx context.Context, kingdomID string, limit int) ([]domain.KingdomEvent, error) {
	const query = kingdomEventSelect + `
		WHERE ke.kingdom_id = $1 AND ke.status IN ('resolved', 'expired')
		ORDER BY COALESCE(ke.resolved_at, ke.updated_at) DESC
		LIMIT $2
	`
	return r.listKingdomEvents(ctx, query, kingdomID, limit)
}

func (r *EventRepository) CreateKingdomEvent(ctx context.Context, kingdomID string, gameEventID string, generatedAt time.Time, expiresAt time.Time) (domain.KingdomEvent, error) {
	const query = `
		INSERT INTO kingdom_events (kingdom_id, game_event_id, generated_at, expires_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id::text
	`
	var id string
	if err := r.db.QueryRowContext(ctx, query, kingdomID, gameEventID, generatedAt, expiresAt).Scan(&id); err != nil {
		return domain.KingdomEvent{}, err
	}
	return r.FindByIDAndKingdomID(ctx, id, kingdomID)
}

func (r *EventRepository) ExpireStale(ctx context.Context, kingdomID string, now time.Time) error {
	const query = `
		UPDATE kingdom_events
		SET status = 'expired',
			updated_at = now()
		WHERE kingdom_id = $1
		  AND status = 'active'
		  AND expires_at <= $2
	`
	_, err := r.db.ExecContext(ctx, query, kingdomID, now)
	return err
}

func (r *EventRepository) FindByIDAndKingdomID(ctx context.Context, eventID string, kingdomID string) (domain.KingdomEvent, error) {
	const query = kingdomEventSelect + `
		WHERE ke.id = $1 AND ke.kingdom_id = $2
	`
	events, err := r.listKingdomEvents(ctx, query, eventID, kingdomID)
	if err != nil {
		return domain.KingdomEvent{}, err
	}
	if len(events) == 0 {
		return domain.KingdomEvent{}, ErrEventNotFound
	}
	return events[0], nil
}

func (r *EventRepository) MarkResolved(ctx context.Context, eventID string, kingdomID string, choiceKey string, resolvedAt time.Time, resultJSON []byte) (domain.KingdomEvent, error) {
	const query = `
		UPDATE kingdom_events
		SET status = 'resolved',
			resolved_at = $3,
			selected_choice_key = $4,
			result_json = $5,
			updated_at = now()
		WHERE id = $1 AND kingdom_id = $2 AND status = 'active'
		RETURNING id::text
	`
	var id string
	if err := r.db.QueryRowContext(ctx, query, eventID, kingdomID, resolvedAt, choiceKey, resultJSON).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.KingdomEvent{}, ErrEventNotFound
		}
		return domain.KingdomEvent{}, err
	}
	return r.FindByIDAndKingdomID(ctx, id, kingdomID)
}

func (r *EventRepository) HasRecentByEventKey(ctx context.Context, kingdomID string, eventKey string, since time.Time) (bool, error) {
	const query = `
		SELECT EXISTS (
			SELECT 1
			FROM kingdom_events ke
			JOIN game_events ge ON ge.id = ke.game_event_id
			WHERE ke.kingdom_id = $1
			  AND ge.event_key = $2
			  AND ke.generated_at >= $3
		)
	`
	var exists bool
	if err := r.db.QueryRowContext(ctx, query, kingdomID, eventKey, since).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

const kingdomEventSelect = `
	SELECT
		ke.id::text, ke.kingdom_id::text, ke.game_event_id::text, ke.status, ke.generated_at, ke.expires_at,
		ke.resolved_at, ke.selected_choice_key, ke.result_json, ke.created_at, ke.updated_at,
		ge.id::text, ge.event_key, ge.category, ge.title, ge.body, ge.trigger_type, ge.weight, ge.is_active,
		ge.cooldown_seconds, ge.expires_after_seconds, ge.conditions_json, ge.created_at, ge.updated_at
	FROM kingdom_events ke
	JOIN game_events ge ON ge.id = ke.game_event_id
`

func (r *EventRepository) listKingdomEvents(ctx context.Context, query string, args ...interface{}) ([]domain.KingdomEvent, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := []domain.KingdomEvent{}
	for rows.Next() {
		event, err := scanKingdomEvent(rows)
		if err != nil {
			return nil, err
		}
		choices, err := r.ListChoicesByGameEventID(ctx, event.GameEventID)
		if err != nil {
			return nil, err
		}
		event.Choices = choices
		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return events, nil
}

func scanGameEvent(row scanner) (domain.GameEvent, error) {
	var event domain.GameEvent
	err := row.Scan(
		&event.ID,
		&event.Key,
		&event.Category,
		&event.Title,
		&event.Body,
		&event.TriggerType,
		&event.Weight,
		&event.IsActive,
		&event.CooldownSeconds,
		&event.ExpiresAfterSeconds,
		&event.ConditionsJSON,
		&event.CreatedAt,
		&event.UpdatedAt,
	)
	if err != nil {
		return domain.GameEvent{}, err
	}
	return event, nil
}

func scanEventChoice(row scanner) (domain.EventChoice, error) {
	var choice domain.EventChoice
	err := row.Scan(
		&choice.ID,
		&choice.GameEventID,
		&choice.Key,
		&choice.Label,
		&choice.Description,
		&choice.EffectsJSON,
		&choice.ResultTitle,
		&choice.ResultBody,
		&choice.CreatedAt,
		&choice.UpdatedAt,
	)
	if err != nil {
		return domain.EventChoice{}, err
	}
	return choice, nil
}

func scanKingdomEvent(row scanner) (domain.KingdomEvent, error) {
	var event domain.KingdomEvent
	var resolvedAt sql.NullTime
	var selectedChoiceKey sql.NullString
	err := row.Scan(
		&event.ID,
		&event.KingdomID,
		&event.GameEventID,
		&event.Status,
		&event.GeneratedAt,
		&event.ExpiresAt,
		&resolvedAt,
		&selectedChoiceKey,
		&event.ResultJSON,
		&event.CreatedAt,
		&event.UpdatedAt,
		&event.GameEvent.ID,
		&event.GameEvent.Key,
		&event.GameEvent.Category,
		&event.GameEvent.Title,
		&event.GameEvent.Body,
		&event.GameEvent.TriggerType,
		&event.GameEvent.Weight,
		&event.GameEvent.IsActive,
		&event.GameEvent.CooldownSeconds,
		&event.GameEvent.ExpiresAfterSeconds,
		&event.GameEvent.ConditionsJSON,
		&event.GameEvent.CreatedAt,
		&event.GameEvent.UpdatedAt,
	)
	if err != nil {
		return domain.KingdomEvent{}, err
	}
	if resolvedAt.Valid {
		event.ResolvedAt = &resolvedAt.Time
	}
	if selectedChoiceKey.Valid {
		event.SelectedChoiceKey = &selectedChoiceKey.String
	}
	return event, nil
}
