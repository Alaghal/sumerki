package service

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"sumerki/backend/internal/domain"
	"sumerki/backend/internal/repository"
)

const (
	HealthStatusHealthy = "healthy"
)

var ErrRulerKingdomNotFound = errors.New("kingdom not found")

type RulerRepository interface {
	Create(ctx context.Context, ruler domain.Ruler) (domain.Ruler, error)
	FindByKingdomID(ctx context.Context, kingdomID string) (domain.Ruler, error)
}

type RulerService struct {
	kingdoms KingdomRepository
	rulers   RulerRepository
	random   *rand.Rand
}

func NewRulerService(kingdoms KingdomRepository, rulers RulerRepository) *RulerService {
	return &RulerService{
		kingdoms: kingdoms,
		rulers:   rulers,
		random:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (s *RulerService) Current(ctx context.Context, userID string) (domain.Ruler, error) {
	kingdom, err := s.kingdoms.FindByUserID(ctx, userID)
	if errors.Is(err, repository.ErrKingdomNotFound) {
		return domain.Ruler{}, ErrRulerKingdomNotFound
	}
	if err != nil {
		return domain.Ruler{}, err
	}

	ruler, err := s.rulers.FindByKingdomID(ctx, kingdom.ID)
	if errors.Is(err, repository.ErrRulerNotFound) {
		return s.GenerateForKingdom(ctx, kingdom)
	}
	if err != nil {
		return domain.Ruler{}, err
	}

	return ruler, nil
}

func (s *RulerService) GenerateForKingdom(ctx context.Context, kingdom domain.Kingdom) (domain.Ruler, error) {
	ruler := s.generate(kingdom)
	created, err := s.rulers.Create(ctx, ruler)
	if errors.Is(err, repository.ErrRulerExists) {
		return s.rulers.FindByKingdomID(ctx, kingdom.ID)
	}
	if err != nil {
		return domain.Ruler{}, err
	}

	return created, nil
}

func (s *RulerService) generate(kingdom domain.Kingdom) domain.Ruler {
	return domain.Ruler{
		KingdomID:    kingdom.ID,
		Name:         s.pickName(kingdom.Culture),
		Age:          s.randomInRange(25, 60),
		Culture:      kingdom.Culture,
		Authority:    s.randomInRange(30, 80),
		Courage:      s.randomInRange(30, 80),
		Cunning:      s.randomInRange(30, 80),
		Honor:        s.randomInRange(30, 80),
		Cruelty:      s.randomInRange(30, 80),
		Ambition:     s.randomInRange(30, 80),
		Paranoia:     s.randomInRange(30, 80),
		HealthStatus: HealthStatusHealthy,
	}
}

func (s *RulerService) pickName(culture string) string {
	names := rulerNamePools[culture]
	if len(names) == 0 {
		names = rulerNamePools[CultureFreePosad]
	}

	return names[s.random.Intn(len(names))]
}

func (s *RulerService) randomInRange(minimum int, maximum int) int {
	return minimum + s.random.Intn(maximum-minimum+1)
}

var rulerNamePools = map[string][]string{
	CultureNorthernPrincipality: {
		"Боривой",
		"Радомир",
		"Милорад",
		"Всеслава",
		"Ярополк",
	},
	CultureLizardGrad: {
		"Шессар",
		"Иррах",
		"Ссавара",
		"Кхореш",
		"Тёплокам",
	},
	CultureFreePosad: {
		"Берест",
		"Лука",
		"Добрын",
		"Мирена",
		"Сребран",
	},
}
