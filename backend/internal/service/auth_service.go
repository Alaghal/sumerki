package service

import (
	"context"
	"errors"
	"net/mail"
	"strings"
	"time"

	"sumerki/backend/internal/domain"
	"sumerki/backend/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const tokenTTL = 24 * time.Hour

var (
	ErrInvalidEmail       = errors.New("invalid email")
	ErrPasswordTooShort   = errors.New("password too short")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrExpiredToken       = errors.New("expired token")
	ErrUserNotFound       = errors.New("user not found")
)

type AuthUserRepository interface {
	Create(ctx context.Context, email string, passwordHash string) (domain.User, error)
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindByID(ctx context.Context, id string) (domain.User, error)
}

type AuthService struct {
	users     AuthUserRepository
	jwtSecret []byte
	now       func() time.Time
}

type AuthResult struct {
	User  domain.User
	Token string
}

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func NewAuthService(users AuthUserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		users:     users,
		jwtSecret: []byte(jwtSecret),
		now:       time.Now,
	}
}

func (s *AuthService) Register(ctx context.Context, email string, password string) (AuthResult, error) {
	normalizedEmail, err := normalizeEmail(email)
	if err != nil {
		return AuthResult{}, err
	}
	if len(password) < 8 {
		return AuthResult{}, ErrPasswordTooShort
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return AuthResult{}, err
	}

	user, err := s.users.Create(ctx, normalizedEmail, string(passwordHash))
	if errors.Is(err, repository.ErrEmailExists) {
		return AuthResult{}, ErrEmailAlreadyExists
	}
	if err != nil {
		return AuthResult{}, err
	}

	token, err := s.GenerateToken(user.ID)
	if err != nil {
		return AuthResult{}, err
	}

	return AuthResult{User: user, Token: token}, nil
}

func (s *AuthService) Login(ctx context.Context, email string, password string) (AuthResult, error) {
	normalizedEmail, err := normalizeEmail(email)
	if err != nil {
		return AuthResult{}, ErrInvalidCredentials
	}

	user, err := s.users.FindByEmail(ctx, normalizedEmail)
	if errors.Is(err, repository.ErrUserNotFound) {
		return AuthResult{}, ErrInvalidCredentials
	}
	if err != nil {
		return AuthResult{}, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return AuthResult{}, ErrInvalidCredentials
	}

	token, err := s.GenerateToken(user.ID)
	if err != nil {
		return AuthResult{}, err
	}

	return AuthResult{User: user, Token: token}, nil
}

func (s *AuthService) CurrentUser(ctx context.Context, userID string) (domain.User, error) {
	user, err := s.users.FindByID(ctx, userID)
	if errors.Is(err, repository.ErrUserNotFound) {
		return domain.User{}, ErrUserNotFound
	}
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (s *AuthService) GenerateToken(userID string) (string, error) {
	now := s.now().UTC()
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(tokenTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *AuthService) ValidateToken(tokenString string) (string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, ErrInvalidToken
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", ErrExpiredToken
		}
		return "", ErrInvalidToken
	}
	if !token.Valid || claims.UserID == "" {
		return "", ErrInvalidToken
	}

	return claims.UserID, nil
}

func normalizeEmail(email string) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(email))
	if normalized == "" || strings.Count(normalized, "@") != 1 {
		return "", ErrInvalidEmail
	}

	address, err := mail.ParseAddress(normalized)
	if err != nil || address.Address != normalized {
		return "", ErrInvalidEmail
	}

	return normalized, nil
}
