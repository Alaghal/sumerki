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

type fakeResourcesRepository struct {
	resourcesByKingdomID map[string]domain.Resources
	now                  func() time.Time
}

type fakeProductionProvider struct {
	bonus gameconfig.ResourceValues
}

func (p fakeProductionProvider) ProductionBonus(_ context.Context, _ string) (gameconfig.ResourceValues, error) {
	return p.bonus, nil
}

func newFakeResourcesRepository() *fakeResourcesRepository {
	return &fakeResourcesRepository{
		resourcesByKingdomID: map[string]domain.Resources{},
		now:                  time.Now,
	}
}

func (r *fakeResourcesRepository) CreateInitial(_ context.Context, kingdomID string) (domain.Resources, error) {
	if _, exists := r.resourcesByKingdomID[kingdomID]; exists {
		return domain.Resources{}, repository.ErrResourcesExist
	}

	now := r.now()
	initial := gameconfig.StartingResources
	resources := domain.Resources{
		KingdomID:        kingdomID,
		Gold:             initial.Gold,
		Food:             initial.Food,
		Wood:             initial.Wood,
		Stone:            initial.Stone,
		Population:       initial.Population,
		LastCalculatedAt: now,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	r.resourcesByKingdomID[kingdomID] = resources

	return resources, nil
}

func (r *fakeResourcesRepository) FindByKingdomID(_ context.Context, kingdomID string) (domain.Resources, error) {
	resources, ok := r.resourcesByKingdomID[kingdomID]
	if !ok {
		return domain.Resources{}, repository.ErrResourcesNotFound
	}

	return resources, nil
}

func (r *fakeResourcesRepository) UpdateCalculated(_ context.Context, resources domain.Resources) (domain.Resources, error) {
	resources.UpdatedAt = r.now()
	r.resourcesByKingdomID[resources.KingdomID] = resources
	return resources, nil
}

func TestResourcesCurrentCreatesMissingResources(t *testing.T) {
	kingdoms := newFakeKingdomRepository()
	kingdom, err := kingdoms.Create(context.Background(), "user-1", "Blackwater", CultureNorthernPrincipality)
	if err != nil {
		t.Fatalf("create kingdom fixture failed: %v", err)
	}

	service := NewResourcesService(kingdoms, newFakeResourcesRepository())

	result, err := service.Current(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("current resources failed: %v", err)
	}

	if result.Resources.KingdomID != kingdom.ID {
		t.Fatalf("expected kingdom id %q, got %q", kingdom.ID, result.Resources.KingdomID)
	}
	if result.Resources.Gold != gameconfig.StartingResources.Gold {
		t.Fatalf("expected starting gold %d, got %d", gameconfig.StartingResources.Gold, result.Resources.Gold)
	}
	if result.ProductionPerHour.Gold != gameconfig.BaseProductionPerHour.Gold {
		t.Fatalf("expected production in response")
	}
}

func TestResourcesCurrentReturnsKingdomNotFound(t *testing.T) {
	service := NewResourcesService(newFakeKingdomRepository(), newFakeResourcesRepository())

	_, err := service.Current(context.Background(), "user-1")
	if !errors.Is(err, ErrResourcesKingdomNotFound) {
		t.Fatalf("expected kingdom not found, got %v", err)
	}
}

func TestResourcesCurrentAppliesLazyProduction(t *testing.T) {
	now := time.Date(2026, 6, 29, 12, 0, 0, 0, time.UTC)
	kingdoms := newFakeKingdomRepository()
	kingdom, err := kingdoms.Create(context.Background(), "user-1", "Blackwater", CultureNorthernPrincipality)
	if err != nil {
		t.Fatalf("create kingdom fixture failed: %v", err)
	}

	resources := newFakeResourcesRepository()
	resources.now = func() time.Time { return now.Add(-1 * time.Hour) }
	if _, err := resources.CreateInitial(context.Background(), kingdom.ID); err != nil {
		t.Fatalf("create resources fixture failed: %v", err)
	}

	service := NewResourcesService(kingdoms, resources)
	service.now = func() time.Time { return now }
	resources.now = func() time.Time { return now }

	result, err := service.Current(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("current resources failed: %v", err)
	}

	expected := gameconfig.StartingResources.Gold + gameconfig.BaseProductionPerHour.Gold
	if result.Resources.Gold != expected {
		t.Fatalf("expected gold %d, got %d", expected, result.Resources.Gold)
	}
	if result.Resources.Population < 0 {
		t.Fatalf("population must never be negative")
	}
}

func TestResourcesCurrentIncludesProductionBonus(t *testing.T) {
	kingdoms := newFakeKingdomRepository()
	if _, err := kingdoms.Create(context.Background(), "user-1", "Blackwater", CultureNorthernPrincipality); err != nil {
		t.Fatalf("create kingdom fixture failed: %v", err)
	}

	service := NewResourcesService(kingdoms, newFakeResourcesRepository())
	service.SetProductionProvider(fakeProductionProvider{bonus: gameconfig.ResourceValues{Food: 15}})

	result, err := service.Current(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("current resources failed: %v", err)
	}

	expected := gameconfig.BaseProductionPerHour.Food + 15
	if result.ProductionPerHour.Food != expected {
		t.Fatalf("expected food production %d, got %d", expected, result.ProductionPerHour.Food)
	}
}

func TestResourcesCurrentSkipsUpdateWhenNoWholeUnitsGained(t *testing.T) {
	now := time.Date(2026, 6, 29, 12, 0, 0, 0, time.UTC)
	kingdoms := newFakeKingdomRepository()
	kingdom, err := kingdoms.Create(context.Background(), "user-1", "Blackwater", CultureNorthernPrincipality)
	if err != nil {
		t.Fatalf("create kingdom fixture failed: %v", err)
	}

	resources := newFakeResourcesRepository()
	resources.now = func() time.Time { return now.Add(-1 * time.Second) }
	created, err := resources.CreateInitial(context.Background(), kingdom.ID)
	if err != nil {
		t.Fatalf("create resources fixture failed: %v", err)
	}

	service := NewResourcesService(kingdoms, resources)
	service.now = func() time.Time { return now }

	result, err := service.Current(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("current resources failed: %v", err)
	}

	if !result.Resources.LastCalculatedAt.Equal(created.LastCalculatedAt) {
		t.Fatalf("expected last calculated time to remain unchanged")
	}
}
