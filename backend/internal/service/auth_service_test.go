package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"sumerki/backend/internal/domain"
	"sumerki/backend/internal/repository"
)

type fakeUserRepository struct {
	usersByID    map[string]domain.User
	usersByEmail map[string]domain.User
	nextID       string
}

func newFakeUserRepository() *fakeUserRepository {
	return &fakeUserRepository{
		usersByID:    map[string]domain.User{},
		usersByEmail: map[string]domain.User{},
		nextID:       "user-1",
	}
}

func (r *fakeUserRepository) Create(_ context.Context, email string, passwordHash string) (domain.User, error) {
	if _, exists := r.usersByEmail[email]; exists {
		return domain.User{}, repository.ErrEmailExists
	}

	user := domain.User{
		ID:           r.nextID,
		Email:        email,
		PasswordHash: passwordHash,
	}
	r.usersByID[user.ID] = user
	r.usersByEmail[user.Email] = user

	return user, nil
}

func (r *fakeUserRepository) FindByEmail(_ context.Context, email string) (domain.User, error) {
	user, ok := r.usersByEmail[email]
	if !ok {
		return domain.User{}, repository.ErrUserNotFound
	}

	return user, nil
}

func (r *fakeUserRepository) FindByID(_ context.Context, id string) (domain.User, error) {
	user, ok := r.usersByID[id]
	if !ok {
		return domain.User{}, repository.ErrUserNotFound
	}

	return user, nil
}

func TestRegisterNormalizesEmailAndHashesPassword(t *testing.T) {
	repo := newFakeUserRepository()
	auth := NewAuthService(repo, "test-secret")

	result, err := auth.Register(context.Background(), " Player@Example.COM ", "password123")
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	if result.User.Email != "player@example.com" {
		t.Fatalf("expected normalized email, got %q", result.User.Email)
	}
	if result.User.PasswordHash == "password123" {
		t.Fatal("password hash must not equal plaintext password")
	}
	if result.Token == "" {
		t.Fatal("expected token")
	}
}

func TestLoginRejectsInvalidPassword(t *testing.T) {
	repo := newFakeUserRepository()
	auth := NewAuthService(repo, "test-secret")

	if _, err := auth.Register(context.Background(), "player@example.com", "password123"); err != nil {
		t.Fatalf("register failed: %v", err)
	}

	_, err := auth.Login(context.Background(), "player@example.com", "wrong-password")
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected invalid credentials, got %v", err)
	}
}

func TestGenerateAndValidateToken(t *testing.T) {
	auth := NewAuthService(newFakeUserRepository(), "test-secret")
	auth.now = func() time.Time {
		return time.Date(2026, 6, 29, 12, 0, 0, 0, time.UTC)
	}

	token, err := auth.GenerateToken("user-1")
	if err != nil {
		t.Fatalf("generate token failed: %v", err)
	}

	userID, err := auth.ValidateToken(token)
	if err != nil {
		t.Fatalf("validate token failed: %v", err)
	}
	if userID != "user-1" {
		t.Fatalf("expected user-1, got %q", userID)
	}
}

func TestValidateTokenRejectsExpiredToken(t *testing.T) {
	auth := NewAuthService(newFakeUserRepository(), "test-secret")
	auth.now = func() time.Time {
		return time.Now().Add(-25 * time.Hour)
	}

	token, err := auth.GenerateToken("user-1")
	if err != nil {
		t.Fatalf("generate token failed: %v", err)
	}

	_, err = auth.ValidateToken(token)
	if !errors.Is(err, ErrExpiredToken) {
		t.Fatalf("expected expired token, got %v", err)
	}
}
