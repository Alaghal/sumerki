package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"sumerki/backend/internal/domain"
	"sumerki/backend/internal/gameconfig"
)

var (
	ErrBuildingExists   = errors.New("building already exists")
	ErrBuildingNotFound = errors.New("building not found")
)

type BuildingRepository struct {
	db *sql.DB
}

func NewBuildingRepository(db *sql.DB) *BuildingRepository {
	return &BuildingRepository{db: db}
}

func (r *BuildingRepository) CreateInitial(ctx context.Context, kingdomID string) error {
	for _, buildingType := range gameconfig.BuildingOrder {
		level := gameconfig.InitialBuildingLevels[buildingType]
		const query = `
			INSERT INTO kingdom_buildings (kingdom_id, type, level)
			VALUES ($1, $2, $3)
			ON CONFLICT (kingdom_id, type) DO NOTHING
		`
		if _, err := r.db.ExecContext(ctx, query, kingdomID, buildingType, level); err != nil {
			return err
		}
	}

	return nil
}

func (r *BuildingRepository) ListByKingdomID(ctx context.Context, kingdomID string) ([]domain.Building, error) {
	const query = `
		SELECT id::text, kingdom_id::text, type, level, upgrade_started_at, upgrade_finishes_at, created_at, updated_at
		FROM kingdom_buildings
		WHERE kingdom_id = $1
		ORDER BY array_position(ARRAY[
			'town_hall',
			'farm',
			'lumberyard',
			'quarry',
			'market',
			'barracks',
			'walls',
			'shrine'
		]::text[], type)
	`

	rows, err := r.db.QueryContext(ctx, query, kingdomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	buildings := []domain.Building{}
	for rows.Next() {
		building, err := scanBuilding(rows)
		if err != nil {
			return nil, err
		}
		buildings = append(buildings, building)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return buildings, nil
}

func (r *BuildingRepository) FindByKingdomIDAndType(ctx context.Context, kingdomID string, buildingType string) (domain.Building, error) {
	const query = `
		SELECT id::text, kingdom_id::text, type, level, upgrade_started_at, upgrade_finishes_at, created_at, updated_at
		FROM kingdom_buildings
		WHERE kingdom_id = $1 AND type = $2
	`

	building, err := scanBuilding(r.db.QueryRowContext(ctx, query, kingdomID, buildingType))
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Building{}, ErrBuildingNotFound
	}
	if err != nil {
		return domain.Building{}, err
	}

	return building, nil
}

func (r *BuildingRepository) CompleteFinished(ctx context.Context, kingdomID string, now time.Time) error {
	const query = `
		UPDATE kingdom_buildings
		SET level = level + 1,
			upgrade_started_at = NULL,
			upgrade_finishes_at = NULL,
			updated_at = now()
		WHERE kingdom_id = $1
		  AND upgrade_finishes_at IS NOT NULL
		  AND upgrade_finishes_at <= $2
		  AND level < $3
	`

	_, err := r.db.ExecContext(ctx, query, kingdomID, now, gameconfig.MaxBuildingLevel)
	return err
}

func (r *BuildingRepository) StartUpgrade(ctx context.Context, kingdomID string, buildingType string, startedAt time.Time, finishesAt time.Time) (domain.Building, error) {
	const query = `
		UPDATE kingdom_buildings
		SET upgrade_started_at = $3,
			upgrade_finishes_at = $4,
			updated_at = now()
		WHERE kingdom_id = $1 AND type = $2
		RETURNING id::text, kingdom_id::text, type, level, upgrade_started_at, upgrade_finishes_at, created_at, updated_at
	`

	building, err := scanBuilding(r.db.QueryRowContext(ctx, query, kingdomID, buildingType, startedAt, finishesAt))
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Building{}, ErrBuildingNotFound
	}
	if err != nil {
		return domain.Building{}, err
	}

	return building, nil
}

func scanBuilding(row scanner) (domain.Building, error) {
	var building domain.Building
	var upgradeStartedAt sql.NullTime
	var upgradeFinishesAt sql.NullTime
	err := row.Scan(
		&building.ID,
		&building.KingdomID,
		&building.Type,
		&building.Level,
		&upgradeStartedAt,
		&upgradeFinishesAt,
		&building.CreatedAt,
		&building.UpdatedAt,
	)
	if err != nil {
		return domain.Building{}, err
	}

	if upgradeStartedAt.Valid {
		building.UpgradeStartedAt = &upgradeStartedAt.Time
	}
	if upgradeFinishesAt.Valid {
		building.UpgradeFinishesAt = &upgradeFinishesAt.Time
	}

	return building, nil
}
