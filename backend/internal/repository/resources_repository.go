package repository

import (
	"context"
	"database/sql"
	"errors"

	"sumerki/backend/internal/domain"
	"sumerki/backend/internal/gameconfig"
)

var (
	ErrResourcesExist    = errors.New("resources already exist")
	ErrResourcesNotFound = errors.New("resources not found")
)

type ResourcesRepository struct {
	db *sql.DB
}

func NewResourcesRepository(db *sql.DB) *ResourcesRepository {
	return &ResourcesRepository{db: db}
}

func (r *ResourcesRepository) CreateInitial(ctx context.Context, kingdomID string) (domain.Resources, error) {
	initial := gameconfig.StartingResources
	const query = `
		INSERT INTO kingdom_resources (kingdom_id, gold, food, wood, stone, population)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING kingdom_id::text, gold, food, wood, stone, population, last_calculated_at, created_at, updated_at
	`

	resources, err := scanResources(r.db.QueryRowContext(
		ctx,
		query,
		kingdomID,
		initial.Gold,
		initial.Food,
		initial.Wood,
		initial.Stone,
		initial.Population,
	))
	if err != nil {
		if isUniqueViolation(err) {
			return domain.Resources{}, ErrResourcesExist
		}
		return domain.Resources{}, err
	}

	return resources, nil
}

func (r *ResourcesRepository) FindByKingdomID(ctx context.Context, kingdomID string) (domain.Resources, error) {
	const query = `
		SELECT kingdom_id::text, gold, food, wood, stone, population, last_calculated_at, created_at, updated_at
		FROM kingdom_resources
		WHERE kingdom_id = $1
	`

	resources, err := scanResources(r.db.QueryRowContext(ctx, query, kingdomID))
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Resources{}, ErrResourcesNotFound
	}
	if err != nil {
		return domain.Resources{}, err
	}

	return resources, nil
}

func (r *ResourcesRepository) UpdateCalculated(ctx context.Context, resources domain.Resources) (domain.Resources, error) {
	const query = `
		UPDATE kingdom_resources
		SET gold = $2,
			food = $3,
			wood = $4,
			stone = $5,
			population = $6,
			last_calculated_at = $7,
			updated_at = now()
		WHERE kingdom_id = $1
		RETURNING kingdom_id::text, gold, food, wood, stone, population, last_calculated_at, created_at, updated_at
	`

	return scanResources(r.db.QueryRowContext(
		ctx,
		query,
		resources.KingdomID,
		resources.Gold,
		resources.Food,
		resources.Wood,
		resources.Stone,
		resources.Population,
		resources.LastCalculatedAt,
	))
}

func scanResources(row scanner) (domain.Resources, error) {
	var resources domain.Resources
	err := row.Scan(
		&resources.KingdomID,
		&resources.Gold,
		&resources.Food,
		&resources.Wood,
		&resources.Stone,
		&resources.Population,
		&resources.LastCalculatedAt,
		&resources.CreatedAt,
		&resources.UpdatedAt,
	)
	if err != nil {
		return domain.Resources{}, err
	}

	return resources, nil
}
