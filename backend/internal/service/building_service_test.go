package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"sumerki/backend/internal/domain"
	"sumerki/backend/internal/gameconfig"
	"sumerki/backend/internal/repository"
)

type fakeBuildingRepository struct {
	buildingsByKingdomID map[string]map[string]domain.Building
	now                  func() time.Time
}

func newFakeBuildingRepository() *fakeBuildingRepository {
	return &fakeBuildingRepository{
		buildingsByKingdomID: map[string]map[string]domain.Building{},
		now:                  time.Now,
	}
}

func (r *fakeBuildingRepository) CreateInitial(_ context.Context, kingdomID string) error {
	if _, ok := r.buildingsByKingdomID[kingdomID]; !ok {
		r.buildingsByKingdomID[kingdomID] = map[string]domain.Building{}
	}
	now := r.now()
	for _, buildingType := range gameconfig.BuildingOrder {
		if _, exists := r.buildingsByKingdomID[kingdomID][buildingType]; exists {
			continue
		}
		r.buildingsByKingdomID[kingdomID][buildingType] = domain.Building{
			ID:        kingdomID + "-" + buildingType,
			KingdomID: kingdomID,
			Type:      buildingType,
			Level:     gameconfig.InitialBuildingLevels[buildingType],
			CreatedAt: now,
			UpdatedAt: now,
		}
	}
	return nil
}

func (r *fakeBuildingRepository) ListByKingdomID(_ context.Context, kingdomID string) ([]domain.Building, error) {
	buildingsByType, ok := r.buildingsByKingdomID[kingdomID]
	if !ok {
		return nil, nil
	}
	buildings := make([]domain.Building, 0, len(gameconfig.BuildingOrder))
	for _, buildingType := range gameconfig.BuildingOrder {
		if building, ok := buildingsByType[buildingType]; ok {
			buildings = append(buildings, building)
		}
	}
	return buildings, nil
}

func (r *fakeBuildingRepository) FindByKingdomIDAndType(_ context.Context, kingdomID string, buildingType string) (domain.Building, error) {
	buildingsByType, ok := r.buildingsByKingdomID[kingdomID]
	if !ok {
		return domain.Building{}, repository.ErrBuildingNotFound
	}
	building, ok := buildingsByType[buildingType]
	if !ok {
		return domain.Building{}, repository.ErrBuildingNotFound
	}
	return building, nil
}

func (r *fakeBuildingRepository) CompleteFinished(_ context.Context, kingdomID string, now time.Time) error {
	for buildingType, building := range r.buildingsByKingdomID[kingdomID] {
		if building.UpgradeFinishesAt != nil && !building.UpgradeFinishesAt.After(now) && building.Level < gameconfig.MaxBuildingLevel {
			building.Level++
			building.UpgradeStartedAt = nil
			building.UpgradeFinishesAt = nil
			building.UpdatedAt = now
			r.buildingsByKingdomID[kingdomID][buildingType] = building
		}
	}
	return nil
}

func (r *fakeBuildingRepository) StartUpgrade(_ context.Context, kingdomID string, buildingType string, startedAt time.Time, finishesAt time.Time) (domain.Building, error) {
	building, err := r.FindByKingdomIDAndType(context.Background(), kingdomID, buildingType)
	if err != nil {
		return domain.Building{}, err
	}
	building.UpgradeStartedAt = &startedAt
	building.UpgradeFinishesAt = &finishesAt
	building.UpdatedAt = startedAt
	r.buildingsByKingdomID[kingdomID][buildingType] = building
	return building, nil
}

func TestBuildingsCurrentReturnsInitialBuildings(t *testing.T) {
	kingdoms := newFakeKingdomRepository()
	if _, err := kingdoms.Create(context.Background(), "user-1", "Blackwater", CultureNorthernPrincipality); err != nil {
		t.Fatalf("create kingdom fixture failed: %v", err)
	}
	resources := NewResourcesService(kingdoms, newFakeResourcesRepository())
	service := NewBuildingService(kingdoms, newFakeBuildingRepository(), resources)

	buildings, err := service.Current(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("current buildings failed: %v", err)
	}
	if len(buildings) != len(gameconfig.BuildingOrder) {
		t.Fatalf("expected %d buildings, got %d", len(gameconfig.BuildingOrder), len(buildings))
	}
}

func TestBuildingUpgradeSpendsResourcesAndStartsTimer(t *testing.T) {
	now := time.Date(2026, 6, 29, 12, 0, 0, 0, time.UTC)
	kingdoms := newFakeKingdomRepository()
	kingdom, err := kingdoms.Create(context.Background(), "user-1", "Blackwater", CultureNorthernPrincipality)
	if err != nil {
		t.Fatalf("create kingdom fixture failed: %v", err)
	}
	resourcesRepo := newFakeResourcesRepository()
	resourcesRepo.now = func() time.Time { return now }
	resources := NewResourcesService(kingdoms, resourcesRepo)
	resources.now = func() time.Time { return now }
	buildingsRepo := newFakeBuildingRepository()
	buildingsRepo.now = func() time.Time { return now }
	service := NewBuildingService(kingdoms, buildingsRepo, resources)
	service.now = func() time.Time { return now }
	if err := service.EnsureForKingdom(context.Background(), kingdom.ID); err != nil {
		t.Fatalf("ensure buildings failed: %v", err)
	}

	result, err := service.Upgrade(context.Background(), "user-1", "farm")
	if err != nil {
		t.Fatalf("upgrade failed: %v", err)
	}

	if !result.Building.Building.IsUpgrading() {
		t.Fatal("expected building to be upgrading")
	}
	if result.Resources.Resources.Gold != gameconfig.StartingResources.Gold-160 {
		t.Fatalf("expected spent gold, got %d", result.Resources.Resources.Gold)
	}
}

func TestBuildingUpgradeFailsWithInsufficientResources(t *testing.T) {
	kingdoms := newFakeKingdomRepository()
	kingdom, err := kingdoms.Create(context.Background(), "user-1", "Blackwater", CultureNorthernPrincipality)
	if err != nil {
		t.Fatalf("create kingdom fixture failed: %v", err)
	}
	resourcesRepo := newFakeResourcesRepository()
	now := time.Date(2026, 6, 29, 12, 0, 0, 0, time.UTC)
	resourcesRepo.resourcesByKingdomID[kingdom.ID] = domain.Resources{
		KingdomID:        kingdom.ID,
		LastCalculatedAt: now,
	}
	resources := NewResourcesService(kingdoms, resourcesRepo)
	resources.now = func() time.Time { return now }
	service := NewBuildingService(kingdoms, newFakeBuildingRepository(), resources)

	_, err = service.Upgrade(context.Background(), "user-1", "farm")
	if !errors.Is(err, ErrInsufficientResources) {
		t.Fatalf("expected insufficient resources, got %v", err)
	}
}

func TestBuildingLazyCompletionIncreasesLevel(t *testing.T) {
	now := time.Date(2026, 6, 29, 12, 0, 0, 0, time.UTC)
	kingdoms := newFakeKingdomRepository()
	kingdom, err := kingdoms.Create(context.Background(), "user-1", "Blackwater", CultureNorthernPrincipality)
	if err != nil {
		t.Fatalf("create kingdom fixture failed: %v", err)
	}
	resources := NewResourcesService(kingdoms, newFakeResourcesRepository())
	buildings := newFakeBuildingRepository()
	service := NewBuildingService(kingdoms, buildings, resources)
	service.now = func() time.Time { return now }
	if err := service.EnsureForKingdom(context.Background(), kingdom.ID); err != nil {
		t.Fatalf("ensure buildings failed: %v", err)
	}
	started := now.Add(-2 * time.Minute)
	finished := now.Add(-1 * time.Minute)
	if _, err := buildings.StartUpgrade(context.Background(), kingdom.ID, "farm", started, finished); err != nil {
		t.Fatalf("start upgrade fixture failed: %v", err)
	}

	current, err := service.Current(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("current buildings failed: %v", err)
	}

	for _, building := range current {
		if building.Building.Type == "farm" && building.Building.Level != 2 {
			t.Fatalf("expected farm level 2, got %d", building.Building.Level)
		}
	}
}

func TestBuildingProductionBonusUsesLevels(t *testing.T) {
	kingdoms := newFakeKingdomRepository()
	kingdom, err := kingdoms.Create(context.Background(), "user-1", "Blackwater", CultureNorthernPrincipality)
	if err != nil {
		t.Fatalf("create kingdom fixture failed: %v", err)
	}
	resources := NewResourcesService(kingdoms, newFakeResourcesRepository())
	service := NewBuildingService(kingdoms, newFakeBuildingRepository(), resources)
	if err := service.EnsureForKingdom(context.Background(), kingdom.ID); err != nil {
		t.Fatalf("ensure buildings failed: %v", err)
	}

	bonus, err := service.ProductionBonus(context.Background(), kingdom.ID)
	if err != nil {
		t.Fatalf("production bonus failed: %v", err)
	}
	if bonus.Food != 15 || bonus.Wood != 12 || bonus.Stone != 10 || bonus.Gold != 10 {
		t.Fatalf("unexpected production bonus: %#v", bonus)
	}
}
