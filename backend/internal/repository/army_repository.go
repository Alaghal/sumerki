package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"sumerki/backend/internal/domain"
	"sumerki/backend/internal/gameconfig"
)

var ErrUnitNotFound = errors.New("unit not found")

type ArmyRepository struct {
	db *sql.DB
}

func NewArmyRepository(db *sql.DB) *ArmyRepository {
	return &ArmyRepository{db: db}
}

func (r *ArmyRepository) CreateInitial(ctx context.Context, kingdomID string) error {
	for _, unitType := range gameconfig.UnitOrder {
		amount := gameconfig.InitialUnitAmounts[unitType]
		const query = `
			INSERT INTO kingdom_units (kingdom_id, unit_type, amount)
			VALUES ($1, $2, $3)
			ON CONFLICT (kingdom_id, unit_type) DO NOTHING
		`
		if _, err := r.db.ExecContext(ctx, query, kingdomID, unitType, amount); err != nil {
			return err
		}
	}

	return nil
}

func (r *ArmyRepository) ListUnitsByKingdomID(ctx context.Context, kingdomID string) ([]domain.Unit, error) {
	const query = `
		SELECT id::text, kingdom_id::text, unit_type, amount, created_at, updated_at
		FROM kingdom_units
		WHERE kingdom_id = $1
		ORDER BY array_position(ARRAY[
			'militia',
			'spearmen',
			'archers',
			'cavalry',
			'scouts'
		]::text[], unit_type)
	`

	rows, err := r.db.QueryContext(ctx, query, kingdomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	units := []domain.Unit{}
	for rows.Next() {
		unit, err := scanUnit(rows)
		if err != nil {
			return nil, err
		}
		units = append(units, unit)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return units, nil
}

func (r *ArmyRepository) FindUnitByKingdomIDAndType(ctx context.Context, kingdomID string, unitType string) (domain.Unit, error) {
	const query = `
		SELECT id::text, kingdom_id::text, unit_type, amount, created_at, updated_at
		FROM kingdom_units
		WHERE kingdom_id = $1 AND unit_type = $2
	`

	unit, err := scanUnit(r.db.QueryRowContext(ctx, query, kingdomID, unitType))
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Unit{}, ErrUnitNotFound
	}
	if err != nil {
		return domain.Unit{}, err
	}

	return unit, nil
}

func (r *ArmyRepository) AdjustUnitAmount(ctx context.Context, kingdomID string, unitType string, delta int64) error {
	const query = `
		UPDATE kingdom_units
		SET amount = amount + $3,
			updated_at = now()
		WHERE kingdom_id = $1
		  AND unit_type = $2
		  AND amount + $3 >= 0
	`

	result, err := r.db.ExecContext(ctx, query, kingdomID, unitType, delta)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrUnitNotFound
	}

	return nil
}

func (r *ArmyRepository) ListTrainingOrdersByKingdomID(ctx context.Context, kingdomID string) ([]domain.UnitTrainingOrder, error) {
	const query = `
		SELECT id::text, kingdom_id::text, unit_type, amount, status, started_at, finishes_at, completed_at, created_at, updated_at
		FROM unit_training_orders
		WHERE kingdom_id = $1 AND status = 'training'
		ORDER BY finishes_at ASC, created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, kingdomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := []domain.UnitTrainingOrder{}
	for rows.Next() {
		order, err := scanTrainingOrder(rows)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *ArmyRepository) CreateTrainingOrder(ctx context.Context, kingdomID string, unitType string, amount int64, startedAt time.Time, finishesAt time.Time) (domain.UnitTrainingOrder, error) {
	const query = `
		INSERT INTO unit_training_orders (kingdom_id, unit_type, amount, started_at, finishes_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id::text, kingdom_id::text, unit_type, amount, status, started_at, finishes_at, completed_at, created_at, updated_at
	`

	order, err := scanTrainingOrder(r.db.QueryRowContext(ctx, query, kingdomID, unitType, amount, startedAt, finishesAt))
	if err != nil {
		return domain.UnitTrainingOrder{}, err
	}

	return order, nil
}

func (r *ArmyRepository) CompleteFinishedTraining(ctx context.Context, kingdomID string, now time.Time) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	const selectQuery = `
		SELECT id::text, unit_type, amount
		FROM unit_training_orders
		WHERE kingdom_id = $1
		  AND status = 'training'
		  AND finishes_at <= $2
		FOR UPDATE
	`

	rows, err := tx.QueryContext(ctx, selectQuery, kingdomID, now)
	if err != nil {
		return err
	}

	type finishedOrder struct {
		id       string
		unitType string
		amount   int64
	}

	orders := []finishedOrder{}
	for rows.Next() {
		var order finishedOrder
		if err := rows.Scan(&order.id, &order.unitType, &order.amount); err != nil {
			rows.Close()
			return err
		}
		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return err
	}
	rows.Close()

	if len(orders) == 0 {
		return tx.Commit()
	}

	const updateUnitsQuery = `
		UPDATE kingdom_units
		SET amount = amount + $3,
			updated_at = now()
		WHERE kingdom_id = $1 AND unit_type = $2
	`
	const completeOrderQuery = `
		UPDATE unit_training_orders
		SET status = 'completed',
			completed_at = $2,
			updated_at = now()
		WHERE id = $1
	`

	for _, order := range orders {
		if _, err := tx.ExecContext(ctx, updateUnitsQuery, kingdomID, order.unitType, order.amount); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, completeOrderQuery, order.id, now); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func scanUnit(row scanner) (domain.Unit, error) {
	var unit domain.Unit
	err := row.Scan(
		&unit.ID,
		&unit.KingdomID,
		&unit.Type,
		&unit.Amount,
		&unit.CreatedAt,
		&unit.UpdatedAt,
	)
	if err != nil {
		return domain.Unit{}, err
	}

	return unit, nil
}

func scanTrainingOrder(row scanner) (domain.UnitTrainingOrder, error) {
	var order domain.UnitTrainingOrder
	var completedAt sql.NullTime
	err := row.Scan(
		&order.ID,
		&order.KingdomID,
		&order.UnitType,
		&order.Amount,
		&order.Status,
		&order.StartedAt,
		&order.FinishesAt,
		&completedAt,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		return domain.UnitTrainingOrder{}, err
	}

	if completedAt.Valid {
		order.CompletedAt = &completedAt.Time
	}

	return order, nil
}
