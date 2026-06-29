package service

import (
	"context"
	"errors"
	"time"

	"sumerki/backend/internal/domain"
	"sumerki/backend/internal/gameconfig"
	"sumerki/backend/internal/repository"
)

var (
	ErrArmyKingdomNotFound  = errors.New("kingdom not found")
	ErrInvalidUnitType      = errors.New("invalid unit type")
	ErrInvalidTrainingCount = errors.New("invalid training amount")
	ErrBarracksLevelTooLow  = errors.New("barracks level too low")
)

type ArmyRepository interface {
	CreateInitial(ctx context.Context, kingdomID string) error
	ListUnitsByKingdomID(ctx context.Context, kingdomID string) ([]domain.Unit, error)
	FindUnitByKingdomIDAndType(ctx context.Context, kingdomID string, unitType string) (domain.Unit, error)
	AdjustUnitAmount(ctx context.Context, kingdomID string, unitType string, delta int64) error
	ListTrainingOrdersByKingdomID(ctx context.Context, kingdomID string) ([]domain.UnitTrainingOrder, error)
	CreateTrainingOrder(ctx context.Context, kingdomID string, unitType string, amount int64, startedAt time.Time, finishesAt time.Time) (domain.UnitTrainingOrder, error)
	CompleteFinishedTraining(ctx context.Context, kingdomID string, now time.Time) error
}

type BuildingLevelProvider interface {
	LevelForKingdom(ctx context.Context, kingdomID string, buildingType string) (int, error)
}

type ArmyView struct {
	KingdomID      string
	Units          []UnitView
	TrainingOrders []TrainingOrderView
	Summary        ArmySummary
}

type UnitView struct {
	Unit         domain.Unit
	Label        string
	Role         string
	Stats        gameconfig.UnitStats
	Cost         gameconfig.ResourceValues
	Seconds      int
	Requirements UnitRequirements
}

type UnitRequirements struct {
	BarracksLevel int
	IsMet         bool
}

type TrainingOrderView struct {
	Order     domain.UnitTrainingOrder
	UnitLabel string
}

type ArmySummary struct {
	TotalUnits   int64
	TotalAttack  int64
	TotalDefense int64
	TotalSupply  int64
}

type TrainUnitsResult struct {
	Order     TrainingOrderView
	Resources ResourcesResult
}

type ArmyService struct {
	kingdoms  KingdomRepository
	army      ArmyRepository
	resources *ResourcesService
	buildings BuildingLevelProvider
	now       func() time.Time
}

func NewArmyService(kingdoms KingdomRepository, army ArmyRepository, resources *ResourcesService, buildings BuildingLevelProvider) *ArmyService {
	return &ArmyService{
		kingdoms:  kingdoms,
		army:      army,
		resources: resources,
		buildings: buildings,
		now:       time.Now,
	}
}

func (s *ArmyService) AfterKingdomCreated(ctx context.Context, kingdom domain.Kingdom) error {
	return s.EnsureForKingdom(ctx, kingdom.ID)
}

func (s *ArmyService) EnsureForKingdom(ctx context.Context, kingdomID string) error {
	return s.army.CreateInitial(ctx, kingdomID)
}

func (s *ArmyService) Current(ctx context.Context, userID string) (ArmyView, error) {
	kingdom, err := s.kingdoms.FindByUserID(ctx, userID)
	if errors.Is(err, repository.ErrKingdomNotFound) {
		return ArmyView{}, ErrArmyKingdomNotFound
	}
	if err != nil {
		return ArmyView{}, err
	}

	return s.CurrentForKingdom(ctx, kingdom.ID)
}

func (s *ArmyService) CurrentForKingdom(ctx context.Context, kingdomID string) (ArmyView, error) {
	if err := s.EnsureForKingdom(ctx, kingdomID); err != nil {
		return ArmyView{}, err
	}
	if err := s.completeFinished(ctx, kingdomID); err != nil {
		return ArmyView{}, err
	}

	units, err := s.army.ListUnitsByKingdomID(ctx, kingdomID)
	if err != nil {
		return ArmyView{}, err
	}
	orders, err := s.army.ListTrainingOrdersByKingdomID(ctx, kingdomID)
	if err != nil {
		return ArmyView{}, err
	}
	barracksLevel, err := s.barracksLevel(ctx, kingdomID)
	if err != nil {
		return ArmyView{}, err
	}

	unitViews := s.unitViews(units, barracksLevel)
	return ArmyView{
		KingdomID:      kingdomID,
		Units:          unitViews,
		TrainingOrders: s.trainingOrderViews(orders),
		Summary:        armySummary(unitViews),
	}, nil
}

func (s *ArmyService) Train(ctx context.Context, userID string, unitType string, amount int64) (TrainUnitsResult, error) {
	if !gameconfig.IsUnitType(unitType) {
		return TrainUnitsResult{}, ErrInvalidUnitType
	}
	if amount < 1 || amount > gameconfig.MaxTrainingAmount {
		return TrainUnitsResult{}, ErrInvalidTrainingCount
	}

	kingdom, err := s.kingdoms.FindByUserID(ctx, userID)
	if errors.Is(err, repository.ErrKingdomNotFound) {
		return TrainUnitsResult{}, ErrArmyKingdomNotFound
	}
	if err != nil {
		return TrainUnitsResult{}, err
	}

	if err := s.EnsureForKingdom(ctx, kingdom.ID); err != nil {
		return TrainUnitsResult{}, err
	}
	if err := s.completeFinished(ctx, kingdom.ID); err != nil {
		return TrainUnitsResult{}, err
	}

	cfg := gameconfig.Units[unitType]
	barracksLevel, err := s.barracksLevel(ctx, kingdom.ID)
	if err != nil {
		return TrainUnitsResult{}, err
	}
	if barracksLevel < cfg.RequiredBarracksLevel {
		return TrainUnitsResult{}, ErrBarracksLevelTooLow
	}

	if _, err := s.army.FindUnitByKingdomIDAndType(ctx, kingdom.ID, unitType); err != nil {
		return TrainUnitsResult{}, err
	}

	resources, err := s.resources.Spend(ctx, kingdom.ID, gameconfig.UnitTrainingCost(unitType, amount))
	if err != nil {
		return TrainUnitsResult{}, err
	}

	startedAt := s.now()
	finishesAt := startedAt.Add(time.Duration(gameconfig.UnitTrainingDurationSeconds(unitType, amount)) * time.Second)
	order, err := s.army.CreateTrainingOrder(ctx, kingdom.ID, unitType, amount, startedAt, finishesAt)
	if err != nil {
		return TrainUnitsResult{}, err
	}

	return TrainUnitsResult{
		Order:     s.trainingOrderView(order),
		Resources: resources,
	}, nil
}

func (s *ArmyService) PrepareForMission(ctx context.Context, kingdomID string) (ArmyView, error) {
	return s.CurrentForKingdom(ctx, kingdomID)
}

func (s *ArmyService) SubtractForMission(ctx context.Context, kingdomID string, units map[string]int64) error {
	for unitType, amount := range units {
		if amount <= 0 {
			continue
		}
		if !gameconfig.IsUnitType(unitType) {
			return ErrInvalidUnitType
		}
		if err := s.army.AdjustUnitAmount(ctx, kingdomID, unitType, -amount); err != nil {
			return err
		}
	}
	return nil
}

func (s *ArmyService) ReturnFromMission(ctx context.Context, kingdomID string, units map[string]int64) error {
	if err := s.EnsureForKingdom(ctx, kingdomID); err != nil {
		return err
	}
	for unitType, amount := range units {
		if amount <= 0 {
			continue
		}
		if !gameconfig.IsUnitType(unitType) {
			return ErrInvalidUnitType
		}
		if err := s.army.AdjustUnitAmount(ctx, kingdomID, unitType, amount); err != nil {
			return err
		}
	}
	return nil
}

func (s *ArmyService) completeFinished(ctx context.Context, kingdomID string) error {
	return s.army.CompleteFinishedTraining(ctx, kingdomID, s.now())
}

func (s *ArmyService) barracksLevel(ctx context.Context, kingdomID string) (int, error) {
	if s.buildings == nil {
		return 0, nil
	}
	return s.buildings.LevelForKingdom(ctx, kingdomID, "barracks")
}

func (s *ArmyService) unitViews(units []domain.Unit, barracksLevel int) []UnitView {
	views := make([]UnitView, 0, len(units))
	for _, unit := range units {
		cfg := gameconfig.Units[unit.Type]
		views = append(views, UnitView{
			Unit:    unit,
			Label:   cfg.Label,
			Role:    cfg.Role,
			Stats:   cfg.Stats,
			Cost:    cfg.Cost,
			Seconds: cfg.SecondsPerUnit,
			Requirements: UnitRequirements{
				BarracksLevel: cfg.RequiredBarracksLevel,
				IsMet:         barracksLevel >= cfg.RequiredBarracksLevel,
			},
		})
	}
	return views
}

func (s *ArmyService) trainingOrderViews(orders []domain.UnitTrainingOrder) []TrainingOrderView {
	views := make([]TrainingOrderView, 0, len(orders))
	for _, order := range orders {
		views = append(views, s.trainingOrderView(order))
	}
	return views
}

func (s *ArmyService) trainingOrderView(order domain.UnitTrainingOrder) TrainingOrderView {
	return TrainingOrderView{
		Order:     order,
		UnitLabel: gameconfig.Units[order.UnitType].Label,
	}
}

func armySummary(units []UnitView) ArmySummary {
	var summary ArmySummary
	for _, unit := range units {
		amount := unit.Unit.Amount
		summary.TotalUnits += amount
		summary.TotalAttack += amount * int64(unit.Stats.Attack)
		summary.TotalDefense += amount * int64(unit.Stats.Defense)
		summary.TotalSupply += amount * int64(unit.Stats.Supply)
	}
	return summary
}
