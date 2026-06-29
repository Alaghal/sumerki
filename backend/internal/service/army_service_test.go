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

type fakeArmyRepository struct {
	unitsByKingdomID  map[string]map[string]domain.Unit
	ordersByKingdomID map[string][]domain.UnitTrainingOrder
	now               func() time.Time
}

type fakeBuildingLevelProvider struct {
	level int
}

func (p fakeBuildingLevelProvider) LevelForKingdom(_ context.Context, _ string, _ string) (int, error) {
	return p.level, nil
}

func newFakeArmyRepository() *fakeArmyRepository {
	return &fakeArmyRepository{
		unitsByKingdomID:  map[string]map[string]domain.Unit{},
		ordersByKingdomID: map[string][]domain.UnitTrainingOrder{},
		now:               time.Now,
	}
}

func (r *fakeArmyRepository) CreateInitial(_ context.Context, kingdomID string) error {
	if _, ok := r.unitsByKingdomID[kingdomID]; !ok {
		r.unitsByKingdomID[kingdomID] = map[string]domain.Unit{}
	}
	now := r.now()
	for _, unitType := range gameconfig.UnitOrder {
		if _, exists := r.unitsByKingdomID[kingdomID][unitType]; exists {
			continue
		}
		r.unitsByKingdomID[kingdomID][unitType] = domain.Unit{
			ID:        kingdomID + "-" + unitType,
			KingdomID: kingdomID,
			Type:      unitType,
			Amount:    gameconfig.InitialUnitAmounts[unitType],
			CreatedAt: now,
			UpdatedAt: now,
		}
	}
	return nil
}

func (r *fakeArmyRepository) ListUnitsByKingdomID(_ context.Context, kingdomID string) ([]domain.Unit, error) {
	unitsByType := r.unitsByKingdomID[kingdomID]
	units := make([]domain.Unit, 0, len(gameconfig.UnitOrder))
	for _, unitType := range gameconfig.UnitOrder {
		if unit, ok := unitsByType[unitType]; ok {
			units = append(units, unit)
		}
	}
	return units, nil
}

func (r *fakeArmyRepository) FindUnitByKingdomIDAndType(_ context.Context, kingdomID string, unitType string) (domain.Unit, error) {
	unit, ok := r.unitsByKingdomID[kingdomID][unitType]
	if !ok {
		return domain.Unit{}, repository.ErrUnitNotFound
	}
	return unit, nil
}

func (r *fakeArmyRepository) AdjustUnitAmount(_ context.Context, kingdomID string, unitType string, delta int64) error {
	unit, ok := r.unitsByKingdomID[kingdomID][unitType]
	if !ok || unit.Amount+delta < 0 {
		return repository.ErrUnitNotFound
	}
	unit.Amount += delta
	unit.UpdatedAt = r.now()
	r.unitsByKingdomID[kingdomID][unitType] = unit
	return nil
}

func (r *fakeArmyRepository) ListTrainingOrdersByKingdomID(_ context.Context, kingdomID string) ([]domain.UnitTrainingOrder, error) {
	orders := []domain.UnitTrainingOrder{}
	for _, order := range r.ordersByKingdomID[kingdomID] {
		if order.Status == "training" {
			orders = append(orders, order)
		}
	}
	return orders, nil
}

func (r *fakeArmyRepository) CreateTrainingOrder(_ context.Context, kingdomID string, unitType string, amount int64, startedAt time.Time, finishesAt time.Time) (domain.UnitTrainingOrder, error) {
	order := domain.UnitTrainingOrder{
		ID:         kingdomID + "-order",
		KingdomID:  kingdomID,
		UnitType:   unitType,
		Amount:     amount,
		Status:     "training",
		StartedAt:  startedAt,
		FinishesAt: finishesAt,
		CreatedAt:  startedAt,
		UpdatedAt:  startedAt,
	}
	r.ordersByKingdomID[kingdomID] = append(r.ordersByKingdomID[kingdomID], order)
	return order, nil
}

func (r *fakeArmyRepository) CompleteFinishedTraining(_ context.Context, kingdomID string, now time.Time) error {
	orders := r.ordersByKingdomID[kingdomID]
	for index, order := range orders {
		if order.Status != "training" || order.FinishesAt.After(now) {
			continue
		}
		unit := r.unitsByKingdomID[kingdomID][order.UnitType]
		unit.Amount += order.Amount
		unit.UpdatedAt = now
		r.unitsByKingdomID[kingdomID][order.UnitType] = unit
		order.Status = "completed"
		order.CompletedAt = &now
		order.UpdatedAt = now
		orders[index] = order
	}
	r.ordersByKingdomID[kingdomID] = orders
	return nil
}

func TestArmyCurrentReturnsInitialUnits(t *testing.T) {
	kingdoms := newFakeKingdomRepository()
	if _, err := kingdoms.Create(context.Background(), "user-1", "Blackwater", CultureNorthernPrincipality); err != nil {
		t.Fatalf("create kingdom fixture failed: %v", err)
	}
	resources := NewResourcesService(kingdoms, newFakeResourcesRepository())
	service := NewArmyService(kingdoms, newFakeArmyRepository(), resources, fakeBuildingLevelProvider{})

	army, err := service.Current(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("current army failed: %v", err)
	}

	if len(army.Units) != len(gameconfig.UnitOrder) {
		t.Fatalf("expected %d unit rows, got %d", len(gameconfig.UnitOrder), len(army.Units))
	}
	if army.Summary.TotalUnits != 12 {
		t.Fatalf("expected 12 starting units, got %d", army.Summary.TotalUnits)
	}
}

func TestArmyTrainMilitiaSpendsResourcesAndCreatesOrder(t *testing.T) {
	now := time.Date(2026, 6, 29, 12, 0, 0, 0, time.UTC)
	kingdoms := newFakeKingdomRepository()
	if _, err := kingdoms.Create(context.Background(), "user-1", "Blackwater", CultureNorthernPrincipality); err != nil {
		t.Fatalf("create kingdom fixture failed: %v", err)
	}
	resourcesRepo := newFakeResourcesRepository()
	resourcesRepo.now = func() time.Time { return now }
	resources := NewResourcesService(kingdoms, resourcesRepo)
	resources.now = func() time.Time { return now }
	armyRepo := newFakeArmyRepository()
	armyRepo.now = func() time.Time { return now }
	service := NewArmyService(kingdoms, armyRepo, resources, fakeBuildingLevelProvider{})
	service.now = func() time.Time { return now }

	result, err := service.Train(context.Background(), "user-1", "militia", 5)
	if err != nil {
		t.Fatalf("train failed: %v", err)
	}

	if result.Order.Order.Amount != 5 {
		t.Fatalf("expected order amount 5, got %d", result.Order.Order.Amount)
	}
	if result.Resources.Resources.Gold != gameconfig.StartingResources.Gold-75 {
		t.Fatalf("expected spent gold, got %d", result.Resources.Resources.Gold)
	}
	if result.Order.Order.FinishesAt.Sub(result.Order.Order.StartedAt) != 25*time.Second {
		t.Fatalf("expected 25 second training duration")
	}
}

func TestArmyTrainSpearmenRequiresBarracks(t *testing.T) {
	kingdoms := newFakeKingdomRepository()
	if _, err := kingdoms.Create(context.Background(), "user-1", "Blackwater", CultureNorthernPrincipality); err != nil {
		t.Fatalf("create kingdom fixture failed: %v", err)
	}
	resources := NewResourcesService(kingdoms, newFakeResourcesRepository())
	service := NewArmyService(kingdoms, newFakeArmyRepository(), resources, fakeBuildingLevelProvider{})

	_, err := service.Train(context.Background(), "user-1", "spearmen", 1)
	if !errors.Is(err, ErrBarracksLevelTooLow) {
		t.Fatalf("expected barracks level error, got %v", err)
	}
}

func TestArmyTrainRejectsInvalidAmount(t *testing.T) {
	service := NewArmyService(newFakeKingdomRepository(), newFakeArmyRepository(), nil, fakeBuildingLevelProvider{})

	_, err := service.Train(context.Background(), "user-1", "militia", gameconfig.MaxTrainingAmount+1)
	if !errors.Is(err, ErrInvalidTrainingCount) {
		t.Fatalf("expected invalid amount, got %v", err)
	}
}

func TestArmyLazyCompletionAddsUnits(t *testing.T) {
	now := time.Date(2026, 6, 29, 12, 0, 0, 0, time.UTC)
	kingdoms := newFakeKingdomRepository()
	kingdom, err := kingdoms.Create(context.Background(), "user-1", "Blackwater", CultureNorthernPrincipality)
	if err != nil {
		t.Fatalf("create kingdom fixture failed: %v", err)
	}
	resources := NewResourcesService(kingdoms, newFakeResourcesRepository())
	armyRepo := newFakeArmyRepository()
	service := NewArmyService(kingdoms, armyRepo, resources, fakeBuildingLevelProvider{})
	service.now = func() time.Time { return now }
	if err := service.EnsureForKingdom(context.Background(), kingdom.ID); err != nil {
		t.Fatalf("ensure army failed: %v", err)
	}
	started := now.Add(-1 * time.Minute)
	finished := now.Add(-30 * time.Second)
	if _, err := armyRepo.CreateTrainingOrder(context.Background(), kingdom.ID, "militia", 3, started, finished); err != nil {
		t.Fatalf("create order fixture failed: %v", err)
	}

	army, err := service.Current(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("current army failed: %v", err)
	}

	if len(army.TrainingOrders) != 0 {
		t.Fatalf("expected no active training orders, got %d", len(army.TrainingOrders))
	}
	for _, unit := range army.Units {
		if unit.Unit.Type == "militia" && unit.Unit.Amount != 13 {
			t.Fatalf("expected 13 militia, got %d", unit.Unit.Amount)
		}
	}
}
