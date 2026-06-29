package repository

import (
	"context"
	"database/sql"
	"errors"

	"sumerki/backend/internal/domain"

	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrEmailExists  = errors.New("email already exists")
	ErrUserNotFound = errors.New("user not found")
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, email string, passwordHash string) (domain.User, error) {
	const query = `
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id::text, email, password_hash, created_at, updated_at
	`

	var user domain.User
	err := r.db.QueryRowContext(ctx, query, email, passwordHash).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return domain.User{}, ErrEmailExists
		}
		return domain.User{}, err
	}

	return user, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	const query = `
		SELECT id::text, email, password_hash, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	return r.findOne(ctx, query, email)
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (domain.User, error) {
	const query = `
		SELECT id::text, email, password_hash, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	return r.findOne(ctx, query, id)
}

func (r *UserRepository) findOne(ctx context.Context, query string, arg string) (domain.User, error) {
	var user domain.User
	err := r.db.QueryRowContext(ctx, query, arg).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.User{}, ErrUserNotFound
	}
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
