package repository

import (
	"context"
	"database/sql"
	"errors"

	"sumerki/backend/internal/domain"
)

var (
	ErrRulerExists   = errors.New("ruler already exists")
	ErrRulerNotFound = errors.New("ruler not found")
)

type RulerRepository struct {
	db *sql.DB
}

func NewRulerRepository(db *sql.DB) *RulerRepository {
	return &RulerRepository{db: db}
}

func (r *RulerRepository) Create(ctx context.Context, ruler domain.Ruler) (domain.Ruler, error) {
	const query = `
		INSERT INTO rulers (
			kingdom_id,
			name,
			age,
			culture,
			authority,
			courage,
			cunning,
			honor,
			cruelty,
			ambition,
			paranoia,
			health_status
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING
			id::text,
			kingdom_id::text,
			name,
			age,
			culture,
			authority,
			courage,
			cunning,
			honor,
			cruelty,
			ambition,
			paranoia,
			health_status,
			created_at,
			updated_at
	`

	created, err := scanRuler(r.db.QueryRowContext(
		ctx,
		query,
		ruler.KingdomID,
		ruler.Name,
		ruler.Age,
		ruler.Culture,
		ruler.Authority,
		ruler.Courage,
		ruler.Cunning,
		ruler.Honor,
		ruler.Cruelty,
		ruler.Ambition,
		ruler.Paranoia,
		ruler.HealthStatus,
	))
	if err != nil {
		if isUniqueViolation(err) {
			return domain.Ruler{}, ErrRulerExists
		}
		return domain.Ruler{}, err
	}

	return created, nil
}

func (r *RulerRepository) FindByKingdomID(ctx context.Context, kingdomID string) (domain.Ruler, error) {
	const query = `
		SELECT
			id::text,
			kingdom_id::text,
			name,
			age,
			culture,
			authority,
			courage,
			cunning,
			honor,
			cruelty,
			ambition,
			paranoia,
			health_status,
			created_at,
			updated_at
		FROM rulers
		WHERE kingdom_id = $1
	`

	ruler, err := scanRuler(r.db.QueryRowContext(ctx, query, kingdomID))
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Ruler{}, ErrRulerNotFound
	}
	if err != nil {
		return domain.Ruler{}, err
	}

	return ruler, nil
}

func scanRuler(row scanner) (domain.Ruler, error) {
	var ruler domain.Ruler
	err := row.Scan(
		&ruler.ID,
		&ruler.KingdomID,
		&ruler.Name,
		&ruler.Age,
		&ruler.Culture,
		&ruler.Authority,
		&ruler.Courage,
		&ruler.Cunning,
		&ruler.Honor,
		&ruler.Cruelty,
		&ruler.Ambition,
		&ruler.Paranoia,
		&ruler.HealthStatus,
		&ruler.CreatedAt,
		&ruler.UpdatedAt,
	)
	if err != nil {
		return domain.Ruler{}, err
	}

	return ruler, nil
}
