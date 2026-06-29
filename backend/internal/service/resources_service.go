package service

import (
	"context"
	"errors"
	"time"

	"sumerki/backend/internal/domain"
	"sumerki/backend/internal/gameconfig"
	"sumerki/backend/internal/repository"
)

var ErrResourcesKingdomNotFound = errors.New("kingdom not found")

type ResourcesRepository interface {
	CreateInitial(ctx context.Context, kingdomID string) (domain.Resources, error)
	FindByKingdomID(ctx context.Context, kingdomID string) (domain.Resources, error)
	UpdateCalculated(ctx context.Context, resources domain.Resources) (domain.Resources, error)
}

type ResourcesResult struct {
	Resources         domain.Resources
	ProductionPerHour gameconfig.ResourceValues
}

type ResourcesService struct {
	kingdoms  KingdomRepository
	resources ResourcesRepository
	now       func() time.Time
}

func NewResourcesService(kingdoms KingdomRepository, resources ResourcesRepository) *ResourcesService {
	return &ResourcesService{
		kingdoms:  kingdoms,
		resources: resources,
		now:       time.Now,
	}
}

func (s *ResourcesService) Current(ctx context.Context, userID string) (ResourcesResult, error) {
	kingdom, err := s.kingdoms.FindByUserID(ctx, userID)
	if errors.Is(err, repository.ErrKingdomNotFound) {
		return ResourcesResult{}, ErrResourcesKingdomNotFound
	}
	if err != nil {
		return ResourcesResult{}, err
	}

	resources, err := s.ensureForKingdom(ctx, kingdom.ID)
	if err != nil {
		return ResourcesResult{}, err
	}

	calculated, err := s.recalculate(ctx, resources)
	if err != nil {
		return ResourcesResult{}, err
	}

	return ResourcesResult{
		Resources:         calculated,
		ProductionPerHour: gameconfig.BaseProductionPerHour,
	}, nil
}

func (s *ResourcesService) AfterKingdomCreated(ctx context.Context, kingdom domain.Kingdom) error {
	_, err := s.CreateForKingdom(ctx, kingdom.ID)
	return err
}

func (s *ResourcesService) CreateForKingdom(ctx context.Context, kingdomID string) (domain.Resources, error) {
	resources, err := s.resources.CreateInitial(ctx, kingdomID)
	if errors.Is(err, repository.ErrResourcesExist) {
		return s.resources.FindByKingdomID(ctx, kingdomID)
	}
	if err != nil {
		return domain.Resources{}, err
	}

	return resources, nil
}

func (s *ResourcesService) ensureForKingdom(ctx context.Context, kingdomID string) (domain.Resources, error) {
	resources, err := s.resources.FindByKingdomID(ctx, kingdomID)
	if errors.Is(err, repository.ErrResourcesNotFound) {
		return s.CreateForKingdom(ctx, kingdomID)
	}
	if err != nil {
		return domain.Resources{}, err
	}

	return resources, nil
}

func (s *ResourcesService) recalculate(ctx context.Context, resources domain.Resources) (domain.Resources, error) {
	now := s.now()
	if !now.After(resources.LastCalculatedAt) {
		return resources, nil
	}

	elapsedSeconds := int64(now.Sub(resources.LastCalculatedAt).Seconds())
	if elapsedSeconds <= 0 {
		return resources, nil
	}

	production := gameconfig.BaseProductionPerHour
	gained := gameconfig.ResourceValues{
		Gold:       elapsedSeconds * production.Gold / 3600,
		Food:       elapsedSeconds * production.Food / 3600,
		Wood:       elapsedSeconds * production.Wood / 3600,
		Stone:      elapsedSeconds * production.Stone / 3600,
		Population: elapsedSeconds * production.Population / 3600,
	}

	if gained.Gold == 0 && gained.Food == 0 && gained.Wood == 0 && gained.Stone == 0 && gained.Population == 0 {
		return resources, nil
	}

	resources.Gold += gained.Gold
	resources.Food += gained.Food
	resources.Wood += gained.Wood
	resources.Stone += gained.Stone
	resources.Population += gained.Population
	resources.LastCalculatedAt = now

	return s.resources.UpdateCalculated(ctx, resources)
}
