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

func (s *ResourcesService) TransferRaidLoot(ctx context.Context, attackerKingdomID string, defenderKingdomID string, percent int64) (gameconfig.ResourceValues, error) {
	return s.transferRaidLoot(ctx, attackerKingdomID, defenderKingdomID, percent, true)
}

func (s *ResourcesService) TransferRaidStalemateLoot(ctx context.Context, attackerKingdomID string, defenderKingdomID string, percent int64) (gameconfig.ResourceValues, error) {
	return s.transferRaidLoot(ctx, attackerKingdomID, defenderKingdomID, percent, false)
}

func (s *ResourcesService) SpendAboveProtected(ctx context.Context, kingdomID string, due gameconfig.ResourceValues, protected gameconfig.ResourceValues) (gameconfig.ResourceValues, ResourcesResult, error) {
	result, err := s.CurrentForKingdom(ctx, kingdomID)
	if err != nil {
		return gameconfig.ResourceValues{}, ResourcesResult{}, err
	}

	resources := result.Resources
	paid := gameconfig.ResourceValues{
		Gold:  payableAboveProtected(resources.Gold, protected.Gold, due.Gold),
		Food:  payableAboveProtected(resources.Food, protected.Food, due.Food),
		Wood:  payableAboveProtected(resources.Wood, protected.Wood, due.Wood),
		Stone: payableAboveProtected(resources.Stone, protected.Stone, due.Stone),
	}

	resources.Gold -= paid.Gold
	resources.Food -= paid.Food
	resources.Wood -= paid.Wood
	resources.Stone -= paid.Stone

	updated, err := s.resources.UpdateCalculated(ctx, resources)
	if err != nil {
		return gameconfig.ResourceValues{}, ResourcesResult{}, err
	}

	return paid, ResourcesResult{
		Resources:         updated,
		ProductionPerHour: result.ProductionPerHour,
	}, nil
}

func (s *ResourcesService) transferRaidLoot(ctx context.Context, attackerKingdomID string, defenderKingdomID string, percent int64, includeStone bool) (gameconfig.ResourceValues, error) {
	defenderResult, err := s.CurrentForKingdom(ctx, defenderKingdomID)
	if err != nil {
		return gameconfig.ResourceValues{}, err
	}
	attackerResult, err := s.CurrentForKingdom(ctx, attackerKingdomID)
	if err != nil {
		return gameconfig.ResourceValues{}, err
	}

	defender := defenderResult.Resources
	attacker := attackerResult.Resources
	loot := raidLoot(defender, gameconfig.ProtectedRaidResources, percent)
	if !includeStone {
		loot.Stone = 0
	}

	defender.Gold -= loot.Gold
	defender.Food -= loot.Food
	defender.Wood -= loot.Wood
	defender.Stone -= loot.Stone
	attacker.Gold += loot.Gold
	attacker.Food += loot.Food
	attacker.Wood += loot.Wood
	attacker.Stone += loot.Stone

	if _, err := s.resources.UpdateCalculated(ctx, defender); err != nil {
		return gameconfig.ResourceValues{}, err
	}
	if _, err := s.resources.UpdateCalculated(ctx, attacker); err != nil {
		return gameconfig.ResourceValues{}, err
	}

	return loot, nil
}

func raidLoot(resources domain.Resources, protected gameconfig.ResourceValues, percent int64) gameconfig.ResourceValues {
	return gameconfig.ResourceValues{
		Gold:  lootForResource(resources.Gold, protected.Gold, percent),
		Food:  lootForResource(resources.Food, protected.Food, percent),
		Wood:  lootForResource(resources.Wood, protected.Wood, percent),
		Stone: lootForResource(resources.Stone, protected.Stone, percent),
	}
}

func lootForResource(value int64, protected int64, percent int64) int64 {
	if value <= protected {
		return 0
	}
	loot := value * percent / 100
	maxLoot := value - protected
	if loot > maxLoot {
		return maxLoot
	}
	return loot
}

func payableAboveProtected(value int64, protected int64, due int64) int64 {
	if due <= 0 || value <= protected {
		return 0
	}
	available := value - protected
	if available < due {
		return available
	}
	return due
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
