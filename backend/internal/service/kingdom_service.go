package service

import (
	"context"
	"errors"
	"strings"

	"sumerki/backend/internal/domain"
	"sumerki/backend/internal/repository"
)

const (
	CultureNorthernPrincipality = "northern_principality"
	CultureLizardGrad           = "lizard_grad"
	CultureFreePosad            = "free_posad"
)

var (
	ErrKingdomNameTooShort  = errors.New("kingdom name too short")
	ErrKingdomNameTooLong   = errors.New("kingdom name too long")
	ErrInvalidCulture       = errors.New("invalid culture")
	ErrKingdomAlreadyExists = errors.New("kingdom already exists")
)

type KingdomRepository interface {
	Create(ctx context.Context, userID string, name string, culture string) (domain.Kingdom, error)
	FindByUserID(ctx context.Context, userID string) (domain.Kingdom, error)
}

type KingdomService struct {
	kingdoms KingdomRepository
}

func NewKingdomService(kingdoms KingdomRepository) *KingdomService {
	return &KingdomService{kingdoms: kingdoms}
}

func (s *KingdomService) Create(ctx context.Context, userID string, name string, culture string) (domain.Kingdom, error) {
	trimmedName := strings.TrimSpace(name)
	nameLength := len([]rune(trimmedName))
	if nameLength < 3 {
		return domain.Kingdom{}, ErrKingdomNameTooShort
	}
	if nameLength > 32 {
		return domain.Kingdom{}, ErrKingdomNameTooLong
	}
	if !validCulture(culture) {
		return domain.Kingdom{}, ErrInvalidCulture
	}

	kingdom, err := s.kingdoms.Create(ctx, userID, trimmedName, culture)
	if errors.Is(err, repository.ErrKingdomExists) {
		return domain.Kingdom{}, ErrKingdomAlreadyExists
	}
	if err != nil {
		return domain.Kingdom{}, err
	}

	return kingdom, nil
}

func (s *KingdomService) Current(ctx context.Context, userID string) (*domain.Kingdom, error) {
	kingdom, err := s.kingdoms.FindByUserID(ctx, userID)
	if errors.Is(err, repository.ErrKingdomNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &kingdom, nil
}

func validCulture(culture string) bool {
	switch culture {
	case CultureNorthernPrincipality, CultureLizardGrad, CultureFreePosad:
		return true
	default:
		return false
	}
}
