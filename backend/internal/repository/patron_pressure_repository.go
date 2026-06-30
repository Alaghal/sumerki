package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"sumerki/backend/internal/domain"
)

var ErrPatronPressureNotFound = errors.New("patron pressure not found")

type PatronPressureRepository struct {
	db *sql.DB
}

func NewPatronPressureRepository(db *sql.DB) *PatronPressureRepository {
	return &PatronPressureRepository{db: db}
}

func (r *PatronPressureRepository) FindByKingdomID(ctx context.Context, kingdomID string) (domain.PatronPressureState, error) {
	const query = `
		SELECT id::text, kingdom_id::text, patron, tribute_debt_gold, tribute_debt_food, contribution_debt_food,
			pressure_level, crisis_status, crisis_started_at, next_tribute_at, last_resolved_at, delay_until, created_at, updated_at
		FROM patron_pressure_states
		WHERE kingdom_id = $1
	`
	state, err := scanPatronPressureState(r.db.QueryRowContext(ctx, query, kingdomID))
	if errors.Is(err, sql.ErrNoRows) {
		return domain.PatronPressureState{}, ErrPatronPressureNotFound
	}
	if err != nil {
		return domain.PatronPressureState{}, err
	}
	return state, nil
}

func (r *PatronPressureRepository) UpsertForPatron(ctx context.Context, kingdomID string, patron string, nextTributeAt time.Time, resetDebt bool) (domain.PatronPressureState, error) {
	if resetDebt {
		return r.upsertReset(ctx, kingdomID, patron, nextTributeAt)
	}
	return r.upsertKeepDebt(ctx, kingdomID, patron, nextTributeAt)
}

func (r *PatronPressureRepository) upsertReset(ctx context.Context, kingdomID string, patron string, nextTributeAt time.Time) (domain.PatronPressureState, error) {
	const query = `
		INSERT INTO patron_pressure_states (kingdom_id, patron, next_tribute_at, last_resolved_at)
		VALUES ($1, $2, $3, now())
		ON CONFLICT (kingdom_id) DO UPDATE
		SET patron = EXCLUDED.patron,
			tribute_debt_gold = 0,
			tribute_debt_food = 0,
			contribution_debt_food = 0,
			pressure_level = 0,
			crisis_status = 'none',
			crisis_started_at = NULL,
			next_tribute_at = EXCLUDED.next_tribute_at,
			last_resolved_at = now(),
			delay_until = NULL,
			updated_at = now()
		RETURNING id::text, kingdom_id::text, patron, tribute_debt_gold, tribute_debt_food, contribution_debt_food,
			pressure_level, crisis_status, crisis_started_at, next_tribute_at, last_resolved_at, delay_until, created_at, updated_at
	`
	return scanPatronPressureState(r.db.QueryRowContext(ctx, query, kingdomID, patron, nextTributeAt))
}

func (r *PatronPressureRepository) upsertKeepDebt(ctx context.Context, kingdomID string, patron string, nextTributeAt time.Time) (domain.PatronPressureState, error) {
	const query = `
		INSERT INTO patron_pressure_states (kingdom_id, patron, next_tribute_at, last_resolved_at)
		VALUES ($1, $2, $3, now())
		ON CONFLICT (kingdom_id) DO UPDATE
		SET patron = EXCLUDED.patron,
			crisis_status = 'none',
			crisis_started_at = NULL,
			next_tribute_at = EXCLUDED.next_tribute_at,
			last_resolved_at = now(),
			delay_until = NULL,
			updated_at = now()
		RETURNING id::text, kingdom_id::text, patron, tribute_debt_gold, tribute_debt_food, contribution_debt_food,
			pressure_level, crisis_status, crisis_started_at, next_tribute_at, last_resolved_at, delay_until, created_at, updated_at
	`
	return scanPatronPressureState(r.db.QueryRowContext(ctx, query, kingdomID, patron, nextTributeAt))
}

func (r *PatronPressureRepository) Save(ctx context.Context, state domain.PatronPressureState) (domain.PatronPressureState, error) {
	const query = `
		UPDATE patron_pressure_states
		SET patron = $2,
			tribute_debt_gold = $3,
			tribute_debt_food = $4,
			contribution_debt_food = $5,
			pressure_level = $6,
			crisis_status = $7,
			crisis_started_at = $8,
			next_tribute_at = $9,
			last_resolved_at = $10,
			delay_until = $11,
			updated_at = now()
		WHERE kingdom_id = $1
		RETURNING id::text, kingdom_id::text, patron, tribute_debt_gold, tribute_debt_food, contribution_debt_food,
			pressure_level, crisis_status, crisis_started_at, next_tribute_at, last_resolved_at, delay_until, created_at, updated_at
	`
	return scanPatronPressureState(r.db.QueryRowContext(
		ctx,
		query,
		state.KingdomID,
		state.Patron,
		state.TributeDebtGold,
		state.TributeDebtFood,
		state.ContributionDebtFood,
		state.PressureLevel,
		state.CrisisStatus,
		state.CrisisStartedAt,
		state.NextTributeAt,
		state.LastResolvedAt,
		state.DelayUntil,
	))
}

func (r *PatronPressureRepository) ClearForKingdom(ctx context.Context, kingdomID string) error {
	const query = `
		UPDATE patron_pressure_states
		SET tribute_debt_gold = 0,
			tribute_debt_food = 0,
			contribution_debt_food = 0,
			pressure_level = 0,
			crisis_status = 'none',
			crisis_started_at = NULL,
			delay_until = NULL,
			updated_at = now()
		WHERE kingdom_id = $1
	`
	_, err := r.db.ExecContext(ctx, query, kingdomID)
	return err
}

func scanPatronPressureState(row scanner) (domain.PatronPressureState, error) {
	var state domain.PatronPressureState
	var crisisStartedAt sql.NullTime
	var delayUntil sql.NullTime
	err := row.Scan(
		&state.ID,
		&state.KingdomID,
		&state.Patron,
		&state.TributeDebtGold,
		&state.TributeDebtFood,
		&state.ContributionDebtFood,
		&state.PressureLevel,
		&state.CrisisStatus,
		&crisisStartedAt,
		&state.NextTributeAt,
		&state.LastResolvedAt,
		&delayUntil,
		&state.CreatedAt,
		&state.UpdatedAt,
	)
	if err != nil {
		return domain.PatronPressureState{}, err
	}
	if crisisStartedAt.Valid {
		state.CrisisStartedAt = &crisisStartedAt.Time
	}
	if delayUntil.Valid {
		state.DelayUntil = &delayUntil.Time
	}
	return state, nil
}
