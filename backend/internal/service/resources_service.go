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
var ErrInsufficientResources = errors.New("insufficient resources")

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
	kingdoms           KingdomRepository
	resources          ResourcesRepository
	productionProvider ResourceProductionProvider
	now                func() time.Time
}

type ResourceProductionProvider interface {
	ProductionBonus(ctx context.Context, kingdomID string) (gameconfig.ResourceValues, error)
}

func NewResourcesService(kingdoms KingdomRepository, resources ResourcesRepository) *ResourcesService {
	return &ResourcesService{
		kingdoms:  kingdoms,
		resources: resources,
		now:       time.Now,
	}
}

func (s *ResourcesService) SetProductionProvider(provider ResourceProductionProvider) {
	s.productionProvider = provider
}

func (s *ResourcesService) Current(ctx context.Context, userID string) (ResourcesResult, error) {
	kingdom, err := s.kingdoms.FindByUserID(ctx, userID)
	if errors.Is(err, repository.ErrKingdomNotFound) {
		return ResourcesResult{}, ErrResourcesKingdomNotFound
	}
	if err != nil {
		return ResourcesResult{}, err
	}

	return s.CurrentForKingdom(ctx, kingdom.ID)
}

func (s *ResourcesService) CurrentForKingdom(ctx context.Context, kingdomID string) (ResourcesResult, error) {
	resources, err := s.ensureForKingdom(ctx, kingdomID)
	if err != nil {
		return ResourcesResult{}, err
	}

	production, err := s.productionForKingdom(ctx, kingdomID)
	if err != nil {
		return ResourcesResult{}, err
	}

	calculated, err := s.recalculate(ctx, resources, production)
	if err != nil {
		return ResourcesResult{}, err
	}

	return ResourcesResult{
		Resources:         calculated,
		ProductionPerHour: production,
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

func (s *ResourcesService) Spend(ctx context.Context, kingdomID string, cost gameconfig.ResourceValues) (ResourcesResult, error) {
	result, err := s.CurrentForKingdom(ctx, kingdomID)
	if err != nil {
		return ResourcesResult{}, err
	}

	resources := result.Resources
	if resources.Gold < cost.Gold ||
		resources.Food < cost.Food ||
		resources.Wood < cost.Wood ||
		resources.Stone < cost.Stone ||
		resources.Population < cost.Population {
		return ResourcesResult{}, ErrInsufficientResources
	}

	resources.Gold -= cost.Gold
	resources.Food -= cost.Food
	resources.Wood -= cost.Wood
	resources.Stone -= cost.Stone
	resources.Population -= cost.Population

	updated, err := s.resources.UpdateCalculated(ctx, resources)
	if err != nil {
		return ResourcesResult{}, err
	}

	return ResourcesResult{
		Resources:         updated,
		ProductionPerHour: result.ProductionPerHour,
	}, nil
}

func (s *ResourcesService) Grant(ctx context.Context, kingdomID string, reward gameconfig.ResourceValues) (ResourcesResult, error) {
	result, err := s.CurrentForKingdom(ctx, kingdomID)
	if err != nil {
		return ResourcesResult{}, err
	}

	resources := result.Resources
	resources.Gold += reward.Gold
	resources.Food += reward.Food
	resources.Wood += reward.Wood
	resources.Stone += reward.Stone
	resources.Population += reward.Population

	updated, err := s.resources.UpdateCalculated(ctx, resources)
	if err != nil {
		return ResourcesResult{}, err
	}

	return ResourcesResult{
		Resources:         updated,
		ProductionPerHour: result.ProductionPerHour,
	}, nil
}

func (s *ResourcesService) productionForKingdom(ctx context.Context, kingdomID string) (gameconfig.ResourceValues, error) {
	production := gameconfig.BaseProductionPerHour
	if s.productionProvider == nil {
		return production, nil
	}

	bonus, err := s.productionProvider.ProductionBonus(ctx, kingdomID)
	if err != nil {
		return gameconfig.ResourceValues{}, err
	}

	production.Gold += bonus.Gold
	production.Food += bonus.Food
	production.Wood += bonus.Wood
	production.Stone += bonus.Stone
	production.Population += bonus.Population

	return production, nil
}

func (s *ResourcesService) recalculate(ctx context.Context, resources domain.Resources, production gameconfig.ResourceValues) (domain.Resources, error) {
	now := s.now()
	if !now.After(resources.LastCalculatedAt) {
		return resources, nil
	}

	elapsedSeconds := int64(now.Sub(resources.LastCalculatedAt).Seconds())
	if elapsedSeconds <= 0 {
		return resources, nil
	}

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
