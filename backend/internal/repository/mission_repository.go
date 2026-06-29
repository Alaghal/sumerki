package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"sumerki/backend/internal/domain"
)

var ErrMissionNotFound = errors.New("mission not found")

type MissionRepository struct {
	db *sql.DB
}

func NewMissionRepository(db *sql.DB) *MissionRepository {
	return &MissionRepository{db: db}
}

func (r *MissionRepository) CreateMission(ctx context.Context, kingdomID string, missionKey string, missionType string, startedAt time.Time, finishesAt time.Time) (domain.Mission, error) {
	const query = `
		INSERT INTO missions (kingdom_id, mission_key, mission_type, started_at, finishes_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id::text, kingdom_id::text, mission_key, mission_type, status, started_at, finishes_at, completed_at, result_json, created_at, updated_at
	`

	return scanMission(r.db.QueryRowContext(ctx, query, kingdomID, missionKey, missionType, startedAt, finishesAt))
}

func (r *MissionRepository) CreateMissionUnit(ctx context.Context, missionID string, unitType string, amount int64) (domain.MissionUnit, error) {
	const query = `
		INSERT INTO mission_units (mission_id, unit_type, amount_sent)
		VALUES ($1, $2, $3)
		RETURNING id::text, mission_id::text, unit_type, amount_sent, amount_lost, amount_returned, created_at, updated_at
	`

	return scanMissionUnit(r.db.QueryRowContext(ctx, query, missionID, unitType, amount))
}

func (r *MissionRepository) ListByKingdomID(ctx context.Context, kingdomID string) ([]domain.Mission, error) {
	const query = `
		SELECT id::text, kingdom_id::text, mission_key, mission_type, status, started_at, finishes_at, completed_at, result_json, created_at, updated_at
		FROM missions
		WHERE kingdom_id = $1
		ORDER BY created_at DESC
	`

	return r.list(ctx, query, kingdomID)
}

func (r *MissionRepository) ListActiveByKingdomID(ctx context.Context, kingdomID string) ([]domain.Mission, error) {
	const query = `
		SELECT id::text, kingdom_id::text, mission_key, mission_type, status, started_at, finishes_at, completed_at, result_json, created_at, updated_at
		FROM missions
		WHERE kingdom_id = $1 AND status = 'active'
		ORDER BY finishes_at ASC
	`

	return r.list(ctx, query, kingdomID)
}

func (r *MissionRepository) ListUnitsByMissionID(ctx context.Context, missionID string) ([]domain.MissionUnit, error) {
	const query = `
		SELECT id::text, mission_id::text, unit_type, amount_sent, amount_lost, amount_returned, created_at, updated_at
		FROM mission_units
		WHERE mission_id = $1
		ORDER BY array_position(ARRAY[
			'militia',
			'spearmen',
			'archers',
			'cavalry',
			'scouts'
		]::text[], unit_type)
	`

	rows, err := r.db.QueryContext(ctx, query, missionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	units := []domain.MissionUnit{}
	for rows.Next() {
		unit, err := scanMissionUnit(rows)
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

func (r *MissionRepository) CompleteMission(ctx context.Context, missionID string, completedAt time.Time, resultJSON []byte) error {
	const query = `
		UPDATE missions
		SET status = 'completed',
			completed_at = $2,
			result_json = $3,
			updated_at = now()
		WHERE id = $1 AND status = 'active'
	`

	result, err := r.db.ExecContext(ctx, query, missionID, completedAt, resultJSON)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrMissionNotFound
	}

	return nil
}

func (r *MissionRepository) UpdateMissionUnitResult(ctx context.Context, missionUnitID string, lost int64, returned int64) error {
	const query = `
		UPDATE mission_units
		SET amount_lost = $2,
			amount_returned = $3,
			updated_at = now()
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, missionUnitID, lost, returned)
	return err
}

func (r *MissionRepository) list(ctx context.Context, query string, kingdomID string) ([]domain.Mission, error) {
	rows, err := r.db.QueryContext(ctx, query, kingdomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	missions := []domain.Mission{}
	for rows.Next() {
		mission, err := scanMission(rows)
		if err != nil {
			return nil, err
		}
		missions = append(missions, mission)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return missions, nil
}

func scanMission(row scanner) (domain.Mission, error) {
	var mission domain.Mission
	var completedAt sql.NullTime
	var resultJSON []byte
	err := row.Scan(
		&mission.ID,
		&mission.KingdomID,
		&mission.Key,
		&mission.Type,
		&mission.Status,
		&mission.StartedAt,
		&mission.FinishesAt,
		&completedAt,
		&resultJSON,
		&mission.CreatedAt,
		&mission.UpdatedAt,
	)
	if err != nil {
		return domain.Mission{}, err
	}
	if completedAt.Valid {
		mission.CompletedAt = &completedAt.Time
	}
	if resultJSON != nil {
		mission.ResultJSON = resultJSON
	}

	return mission, nil
}

func scanMissionUnit(row scanner) (domain.MissionUnit, error) {
	var unit domain.MissionUnit
	err := row.Scan(
		&unit.ID,
		&unit.MissionID,
		&unit.UnitType,
		&unit.AmountSent,
		&unit.AmountLost,
		&unit.AmountReturned,
		&unit.CreatedAt,
		&unit.UpdatedAt,
	)
	if err != nil {
		return domain.MissionUnit{}, err
	}

	return unit, nil
}
