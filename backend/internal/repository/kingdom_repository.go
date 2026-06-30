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
		RETURNING id::text, user_id::text, name, culture, patron, dread, honor, created_at, updated_at
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
		SELECT id::text, user_id::text, name, culture, patron, dread, honor, created_at, updated_at
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

func (r *KingdomRepository) FindByID(ctx context.Context, kingdomID string) (domain.Kingdom, error) {
	const query = `
		SELECT id::text, user_id::text, name, culture, patron, dread, honor, created_at, updated_at
		FROM kingdoms
		WHERE id = $1
	`

	kingdom, err := scanKingdom(r.db.QueryRowContext(ctx, query, kingdomID))
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Kingdom{}, ErrKingdomNotFound
	}
	if err != nil {
		return domain.Kingdom{}, err
	}

	return kingdom, nil
}

func (r *KingdomRepository) ListNeighbors(ctx context.Context, currentKingdomID string, limit int) ([]domain.Kingdom, error) {
	const query = `
		SELECT id::text, user_id::text, name, culture, patron, dread, honor, created_at, updated_at
		FROM kingdoms
		WHERE id <> $1
		ORDER BY created_at ASC, name ASC
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, currentKingdomID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	kingdoms := []domain.Kingdom{}
	for rows.Next() {
		kingdom, err := scanKingdom(rows)
		if err != nil {
			return nil, err
		}
		kingdoms = append(kingdoms, kingdom)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return kingdoms, nil
}

func (r *KingdomRepository) AddDread(ctx context.Context, kingdomID string, amount int) error {
	const query = `
		UPDATE kingdoms
		SET dread = dread + $2,
			updated_at = now()
		WHERE id = $1
	`
	result, err := r.db.ExecContext(ctx, query, kingdomID, amount)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrKingdomNotFound
	}
	return nil
}

func (r *KingdomRepository) UpdatePatronByID(ctx context.Context, kingdomID string, patron *string) (domain.Kingdom, error) {
	const query = `
		UPDATE kingdoms
		SET patron = $2,
			updated_at = now()
		WHERE id = $1
		RETURNING id::text, user_id::text, name, culture, patron, dread, honor, created_at, updated_at
	`

	kingdom, err := scanKingdom(r.db.QueryRowContext(ctx, query, kingdomID, patron))
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
		&kingdom.Dread,
		&kingdom.Honor,
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
