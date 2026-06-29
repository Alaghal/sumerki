package repository

import (
	"context"
	"database/sql"
	"errors"

	"sumerki/backend/internal/domain"
)

var ErrPatronRelationNotFound = errors.New("patron relation not found")

type PatronRepository struct {
	db *sql.DB
}

func NewPatronRepository(db *sql.DB) *PatronRepository {
	return &PatronRepository{db: db}
}

func (r *PatronRepository) FindByKingdomID(ctx context.Context, kingdomID string) (domain.PatronRelation, error) {
	const query = `
		SELECT id::text, kingdom_id::text, patron, favor, standing, joined_at, left_at, created_at, updated_at
		FROM patron_relations
		WHERE kingdom_id = $1 AND left_at IS NULL
	`

	relation, err := scanPatronRelation(r.db.QueryRowContext(ctx, query, kingdomID))
	if errors.Is(err, sql.ErrNoRows) {
		return domain.PatronRelation{}, ErrPatronRelationNotFound
	}
	if err != nil {
		return domain.PatronRelation{}, err
	}
	return relation, nil
}

func (r *PatronRepository) UpsertForKingdom(ctx context.Context, kingdomID string, patron string) (domain.PatronRelation, error) {
	const query = `
		INSERT INTO patron_relations (kingdom_id, patron, favor, standing, joined_at, left_at)
		VALUES ($1, $2, 0, 'neutral', now(), NULL)
		ON CONFLICT (kingdom_id) DO UPDATE
		SET patron = EXCLUDED.patron,
			favor = 0,
			standing = 'neutral',
			joined_at = now(),
			left_at = NULL,
			updated_at = now()
		RETURNING id::text, kingdom_id::text, patron, favor, standing, joined_at, left_at, created_at, updated_at
	`

	return scanPatronRelation(r.db.QueryRowContext(ctx, query, kingdomID, patron))
}

func (r *PatronRepository) BreakForKingdom(ctx context.Context, kingdomID string) error {
	const query = `
		UPDATE patron_relations
		SET left_at = COALESCE(left_at, now()),
			updated_at = now()
		WHERE kingdom_id = $1
	`

	_, err := r.db.ExecContext(ctx, query, kingdomID)
	return err
}

func (r *PatronRepository) BackfillForKingdom(ctx context.Context, kingdomID string, patron string) (domain.PatronRelation, error) {
	const query = `
		INSERT INTO patron_relations (kingdom_id, patron)
		VALUES ($1, $2)
		ON CONFLICT (kingdom_id) DO NOTHING
	`
	if _, err := r.db.ExecContext(ctx, query, kingdomID, patron); err != nil {
		return domain.PatronRelation{}, err
	}
	return r.FindByKingdomID(ctx, kingdomID)
}

func scanPatronRelation(row scanner) (domain.PatronRelation, error) {
	var relation domain.PatronRelation
	var leftAt sql.NullTime
	err := row.Scan(
		&relation.ID,
		&relation.KingdomID,
		&relation.Patron,
		&relation.Favor,
		&relation.Standing,
		&relation.JoinedAt,
		&leftAt,
		&relation.CreatedAt,
		&relation.UpdatedAt,
	)
	if err != nil {
		return domain.PatronRelation{}, err
	}
	if leftAt.Valid {
		relation.LeftAt = &leftAt.Time
	}
	return relation, nil
}
