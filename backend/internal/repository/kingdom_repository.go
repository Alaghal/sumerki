package repository

import (
	"context"
	"database/sql"
	"errors"

	"sumerki/backend/internal/domain"
)

var (
	ErrKingdomExists   = errors.New("kingdom already exists")
	ErrKingdomNotFound = errors.New("kingdom not found")
)

type KingdomRepository struct {
	db *sql.DB
}

func NewKingdomRepository(db *sql.DB) *KingdomRepository {
	return &KingdomRepository{db: db}
}

func (r *KingdomRepository) Create(ctx context.Context, userID string, name string, culture string) (domain.Kingdom, error) {
	const query = `
		INSERT INTO kingdoms (user_id, name, culture)
		VALUES ($1, $2, $3)
		RETURNING id::text, user_id::text, name, culture, patron, created_at, updated_at
	`

	kingdom, err := scanKingdom(r.db.QueryRowContext(ctx, query, userID, name, culture))
	if err != nil {
		if isUniqueViolation(err) {
			return domain.Kingdom{}, ErrKingdomExists
		}
		return domain.Kingdom{}, err
	}

	return kingdom, nil
}

func (r *KingdomRepository) FindByUserID(ctx context.Context, userID string) (domain.Kingdom, error) {
	const query = `
		SELECT id::text, user_id::text, name, culture, patron, created_at, updated_at
		FROM kingdoms
		WHERE user_id = $1
	`

	kingdom, err := scanKingdom(r.db.QueryRowContext(ctx, query, userID))
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Kingdom{}, ErrKingdomNotFound
	}
	if err != nil {
		return domain.Kingdom{}, err
	}

	return kingdom, nil
}

type scanner interface {
	Scan(dest ...interface{}) error
}

func scanKingdom(row scanner) (domain.Kingdom, error) {
	var kingdom domain.Kingdom
	var patron sql.NullString

	err := row.Scan(
		&kingdom.ID,
		&kingdom.UserID,
		&kingdom.Name,
		&kingdom.Culture,
		&patron,
		&kingdom.CreatedAt,
		&kingdom.UpdatedAt,
	)
	if err != nil {
		return domain.Kingdom{}, err
	}

	if patron.Valid {
		kingdom.Patron = &patron.String
	}

	return kingdom, nil
}
