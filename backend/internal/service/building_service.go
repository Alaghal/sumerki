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
	ErrBuildingKingdomNotFound  = errors.New("kingdom not found")
	ErrInvalidBuildingType      = errors.New("invalid building type")
	ErrBuildingNotFound         = errors.New("building not found")
	ErrBuildingAlreadyUpgrading = errors.New("building already upgrading")
	ErrBuildingMaxLevel         = errors.New("building max level")
)

type BuildingRepository interface {
	CreateInitial(ctx context.Context, kingdomID string) error
	ListByKingdomID(ctx context.Context, kingdomID string) ([]domain.Building, error)
	FindByKingdomIDAndType(ctx context.Context, kingdomID string, buildingType string) (domain.Building, error)
	CompleteFinished(ctx context.Context, kingdomID string, now time.Time) error
	StartUpgrade(ctx context.Context, kingdomID string, buildingType string, startedAt time.Time, finishesAt time.Time) (domain.Building, error)
}

type BuildingView struct {
	Building domain.Building
	Label    string
	MaxLevel int
	Effects  []string
	Next     *BuildingNextUpgrade
}

type BuildingNextUpgrade struct {
	TargetLevel     int
	Cost            gameconfig.ResourceValues
	DurationSeconds int
	CanUpgrade      bool
	BlockedReason   *string
}

type BuildingUpgradeResult struct {
	Building  BuildingView
	Resources ResourcesResult
}

type BuildingService struct {
	kingdoms  KingdomRepository
	buildings BuildingRepository
	resources *ResourcesService
	now       func() time.Time
}

func NewBuildingService(kingdoms KingdomRepository, buildings BuildingRepository, resources *ResourcesService) *BuildingService {
	return &BuildingService{
		kingdoms:  kingdoms,
		buildings: buildings,
		resources: resources,
		now:       time.Now,
	}
}

func (s *BuildingService) AfterKingdomCreated(ctx context.Context, kingdom domain.Kingdom) error {
	return s.EnsureForKingdom(ctx, kingdom.ID)
}

func (s *BuildingService) EnsureForKingdom(ctx context.Context, kingdomID string) error {
	return s.buildings.CreateInitial(ctx, kingdomID)
}

func (s *BuildingService) Current(ctx context.Context, userID string) ([]BuildingView, error) {
	kingdom, err := s.kingdoms.FindByUserID(ctx, userID)
	if errors.Is(err, repository.ErrKingdomNotFound) {
		return nil, ErrBuildingKingdomNotFound
	}
	if err != nil {
		return nil, err
	}

	return s.ListForKingdom(ctx, kingdom.ID)
}

func (s *BuildingService) ListForKingdom(ctx context.Context, kingdomID string) ([]BuildingView, error) {
	if err := s.EnsureForKingdom(ctx, kingdomID); err != nil {
		return nil, err
	}
	if err := s.completeFinished(ctx, kingdomID); err != nil {
		return nil, err
	}

	buildings, err := s.buildings.ListByKingdomID(ctx, kingdomID)
	if err != nil {
		return nil, err
	}

	return s.views(buildings), nil
}

func (s *BuildingService) Upgrade(ctx context.Context, userID string, buildingType string) (BuildingUpgradeResult, error) {
	if !gameconfig.IsBuildingType(buildingType) {
		return BuildingUpgradeResult{}, ErrInvalidBuildingType
	}

	kingdom, err := s.kingdoms.FindByUserID(ctx, userID)
	if errors.Is(err, repository.ErrKingdomNotFound) {
		return BuildingUpgradeResult{}, ErrBuildingKingdomNotFound
	}
	if err != nil {
		return BuildingUpgradeResult{}, err
	}

	if err := s.EnsureForKingdom(ctx, kingdom.ID); err != nil {
		return BuildingUpgradeResult{}, err
	}
	if err := s.completeFinished(ctx, kingdom.ID); err != nil {
		return BuildingUpgradeResult{}, err
	}

	building, err := s.buildings.FindByKingdomIDAndType(ctx, kingdom.ID, buildingType)
	if errors.Is(err, repository.ErrBuildingNotFound) {
		return BuildingUpgradeResult{}, ErrBuildingNotFound
	}
	if err != nil {
		return BuildingUpgradeResult{}, err
	}
	if building.IsUpgrading() {
		return BuildingUpgradeResult{}, ErrBuildingAlreadyUpgrading
	}
	if building.Level >= gameconfig.MaxBuildingLevel {
		return BuildingUpgradeResult{}, ErrBuildingMaxLevel
	}

	targetLevel := building.Level + 1
	cost := gameconfig.BuildingCost(building.Type, targetLevel)
	resources, err := s.resources.Spend(ctx, kingdom.ID, cost)
	if err != nil {
		return BuildingUpgradeResult{}, err
	}

	startedAt := s.now()
	finishesAt := startedAt.Add(time.Duration(gameconfig.BuildingUpgradeDurationSeconds(targetLevel)) * time.Second)
	started, err := s.buildings.StartUpgrade(ctx, kingdom.ID, building.Type, startedAt, finishesAt)
	if err != nil {
		return BuildingUpgradeResult{}, err
	}

	return BuildingUpgradeResult{
		Building:  s.view(started),
		Resources: resources,
	}, nil
}

func (s *BuildingService) ProductionBonus(ctx context.Context, kingdomID string) (gameconfig.ResourceValues, error) {
	if err := s.EnsureForKingdom(ctx, kingdomID); err != nil {
		return gameconfig.ResourceValues{}, err
	}
	if err := s.completeFinished(ctx, kingdomID); err != nil {
		return gameconfig.ResourceValues{}, err
	}

	buildings, err := s.buildings.ListByKingdomID(ctx, kingdomID)
	if err != nil {
		return gameconfig.ResourceValues{}, err
	}

	var bonus gameconfig.ResourceValues
	for _, building := range buildings {
		switch building.Type {
		case "farm":
			bonus.Food += int64(building.Level) * 15
		case "lumberyard":
			bonus.Wood += int64(building.Level) * 12
		case "quarry":
			bonus.Stone += int64(building.Level) * 10
		case "market":
			bonus.Gold += int64(building.Level) * 10
		}
	}

	return bonus, nil
}

func (s *BuildingService) LevelForKingdom(ctx context.Context, kingdomID string, buildingType string) (int, error) {
	if err := s.EnsureForKingdom(ctx, kingdomID); err != nil {
		return 0, err
	}
	if err := s.completeFinished(ctx, kingdomID); err != nil {
		return 0, err
	}

	building, err := s.buildings.FindByKingdomIDAndType(ctx, kingdomID, buildingType)
	if errors.Is(err, repository.ErrBuildingNotFound) {
		return 0, ErrBuildingNotFound
	}
	if err != nil {
		return 0, err
	}

	return building.Level, nil
}

func (s *BuildingService) completeFinished(ctx context.Context, kingdomID string) error {
	return s.buildings.CompleteFinished(ctx, kingdomID, s.now())
}

func (s *BuildingService) views(buildings []domain.Building) []BuildingView {
	views := make([]BuildingView, 0, len(buildings))
	for _, building := range buildings {
		views = append(views, s.view(building))
	}
	return views
}

func (s *BuildingService) view(building domain.Building) BuildingView {
	cfg := gameconfig.Buildings[building.Type]
	view := BuildingView{
		Building: building,
		Label:    cfg.Label,
		MaxLevel: gameconfig.MaxBuildingLevel,
		Effects:  cfg.Effects,
	}

	if building.IsUpgrading() {
		view.Next = nil
		return view
	}

	targetLevel := building.Level + 1
	if targetLevel > gameconfig.MaxBuildingLevel {
		reason := "max_level"
		view.Next = &BuildingNextUpgrade{
			TargetLevel:   gameconfig.MaxBuildingLevel,
			CanUpgrade:    false,
			BlockedReason: &reason,
		}
		return view
	}

	view.Next = &BuildingNextUpgrade{
		TargetLevel:     targetLevel,
		Cost:            gameconfig.BuildingCost(building.Type, targetLevel),
		DurationSeconds: gameconfig.BuildingUpgradeDurationSeconds(targetLevel),
		CanUpgrade:      true,
	}
	return view
}
