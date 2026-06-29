package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"sumerki/backend/internal/domain"
	"sumerki/backend/internal/gameconfig"
	"sumerki/backend/internal/repository"
)

var (
	ErrMissionKingdomNotFound    = errors.New("kingdom not found")
	ErrInvalidMissionKey         = errors.New("invalid mission key")
	ErrInvalidMissionUnitType    = errors.New("invalid unit type")
	ErrInvalidMissionUnitAmount  = errors.New("invalid unit amount")
	ErrInsufficientUnits         = errors.New("insufficient units")
	ErrMissionRequirementsNotMet = errors.New("mission requirements not met")
)

type MissionRepository interface {
	CreateMission(ctx context.Context, kingdomID string, missionKey string, missionType string, startedAt time.Time, finishesAt time.Time) (domain.Mission, error)
	CreateMissionUnit(ctx context.Context, missionID string, unitType string, amount int64) (domain.MissionUnit, error)
	ListByKingdomID(ctx context.Context, kingdomID string) ([]domain.Mission, error)
	ListActiveByKingdomID(ctx context.Context, kingdomID string) ([]domain.Mission, error)
	ListUnitsByMissionID(ctx context.Context, missionID string) ([]domain.MissionUnit, error)
	UpdateMissionUnitResult(ctx context.Context, missionUnitID string, lost int64, returned int64) error
	CompleteMission(ctx context.Context, missionID string, completedAt time.Time, resultJSON []byte) error
}

type MissionReportRepository interface {
	CreateMissionReport(ctx context.Context, kingdomID string, missionID string, title string, body string, result string, rewardsJSON []byte, lossesJSON []byte) (domain.MissionReport, error)
	ListByKingdomID(ctx context.Context, kingdomID string) ([]domain.MissionReport, error)
}

type MissionArmyService interface {
	PrepareForMission(ctx context.Context, kingdomID string) (ArmyView, error)
	SubtractForMission(ctx context.Context, kingdomID string, units map[string]int64) error
	ReturnFromMission(ctx context.Context, kingdomID string, units map[string]int64) error
	CurrentForKingdom(ctx context.Context, kingdomID string) (ArmyView, error)
}

type MissionView struct {
	Mission domain.Mission
	Label   string
	Units   []MissionUnitView
	Result  *MissionResult
}

type MissionUnitView struct {
	Unit  domain.MissionUnit
	Label string
}

type MissionResult struct {
	Result  string                    `json:"result"`
	Rewards gameconfig.ResourceValues `json:"rewards"`
	Losses  map[string]int64          `json:"losses"`
}

type MissionReportView struct {
	Report  domain.MissionReport
	Rewards gameconfig.ResourceValues
	Losses  map[string]int64
}

type StartMissionUnit struct {
	UnitType string
	Amount   int64
}

type StartMissionResult struct {
	Mission MissionView
	Army    ArmyView
}

type MissionService struct {
	kingdoms  KingdomRepository
	missions  MissionRepository
	reports   MissionReportRepository
	army      MissionArmyService
	resources *ResourcesService
	now       func() time.Time
}

func NewMissionService(kingdoms KingdomRepository, missions MissionRepository, reports MissionReportRepository, army MissionArmyService, resources *ResourcesService) *MissionService {
	return &MissionService{
		kingdoms:  kingdoms,
		missions:  missions,
		reports:   reports,
		army:      army,
		resources: resources,
		now:       time.Now,
	}
}

func (s *MissionService) Available(ctx context.Context, userID string) ([]gameconfig.MissionConfig, error) {
	if _, err := s.kingdomForUser(ctx, userID); err != nil {
		return nil, err
	}

	missions := make([]gameconfig.MissionConfig, 0, len(gameconfig.MissionOrder))
	for _, key := range gameconfig.MissionOrder {
		missions = append(missions, gameconfig.Missions[key])
	}
	return missions, nil
}

func (s *MissionService) Current(ctx context.Context, userID string) ([]MissionView, error) {
	kingdom, err := s.kingdomForUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if err := s.ResolveCompleted(ctx, kingdom.ID); err != nil {
		return nil, err
	}

	missions, err := s.missions.ListByKingdomID(ctx, kingdom.ID)
	if err != nil {
		return nil, err
	}

	return s.views(ctx, missions)
}

func (s *MissionService) Reports(ctx context.Context, userID string) ([]MissionReportView, error) {
	kingdom, err := s.kingdomForUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if err := s.ResolveCompleted(ctx, kingdom.ID); err != nil {
		return nil, err
	}

	reports, err := s.reports.ListByKingdomID(ctx, kingdom.ID)
	if err != nil {
		return nil, err
	}

	views := make([]MissionReportView, 0, len(reports))
	for _, report := range reports {
		reward, _ := decodeResourceValues(report.RewardsJSON)
		losses, _ := decodeLosses(report.LossesJSON)
		views = append(views, MissionReportView{
			Report:  report,
			Rewards: reward,
			Losses:  losses,
		})
	}
	return views, nil
}

func (s *MissionService) Start(ctx context.Context, userID string, missionKey string, units []StartMissionUnit) (StartMissionResult, error) {
	cfg, ok := gameconfig.Missions[missionKey]
	if !ok {
		return StartMissionResult{}, ErrInvalidMissionKey
	}

	kingdom, err := s.kingdomForUser(ctx, userID)
	if err != nil {
		return StartMissionResult{}, err
	}
	if err := s.ResolveCompleted(ctx, kingdom.ID); err != nil {
		return StartMissionResult{}, err
	}

	sent, err := normalizeMissionUnits(units)
	if err != nil {
		return StartMissionResult{}, err
	}

	army, err := s.army.PrepareForMission(ctx, kingdom.ID)
	if err != nil {
		return StartMissionResult{}, err
	}
	if !hasUnits(army, sent) {
		return StartMissionResult{}, ErrInsufficientUnits
	}
	if !missionRequirementsMet(cfg, sent) {
		return StartMissionResult{}, ErrMissionRequirementsNotMet
	}

	if err := s.army.SubtractForMission(ctx, kingdom.ID, sent); err != nil {
		return StartMissionResult{}, err
	}

	startedAt := s.now()
	finishesAt := startedAt.Add(time.Duration(cfg.DurationSeconds) * time.Second)
	mission, err := s.missions.CreateMission(ctx, kingdom.ID, cfg.Key, cfg.Type, startedAt, finishesAt)
	if err != nil {
		return StartMissionResult{}, err
	}
	for _, unitType := range gameconfig.UnitOrder {
		amount := sent[unitType]
		if amount <= 0 {
			continue
		}
		if _, err := s.missions.CreateMissionUnit(ctx, mission.ID, unitType, amount); err != nil {
			return StartMissionResult{}, err
		}
	}

	updatedArmy, err := s.army.CurrentForKingdom(ctx, kingdom.ID)
	if err != nil {
		return StartMissionResult{}, err
	}
	view, err := s.view(ctx, mission)
	if err != nil {
		return StartMissionResult{}, err
	}

	return StartMissionResult{
		Mission: view,
		Army:    updatedArmy,
	}, nil
}

func (s *MissionService) ResolveCompleted(ctx context.Context, kingdomID string) error {
	active, err := s.missions.ListActiveByKingdomID(ctx, kingdomID)
	if err != nil {
		return err
	}

	now := s.now()
	for _, mission := range active {
		if mission.FinishesAt.After(now) {
			continue
		}
		if err := s.resolveMission(ctx, mission, now); err != nil && !errors.Is(err, repository.ErrMissionNotFound) {
			return err
		}
	}

	return nil
}

func (s *MissionService) resolveMission(ctx context.Context, mission domain.Mission, now time.Time) error {
	cfg := gameconfig.Missions[mission.Key]
	units, err := s.missions.ListUnitsByMissionID(ctx, mission.ID)
	if err != nil {
		return err
	}

	result := calculateMissionResult(cfg, units)
	returned := map[string]int64{}
	for _, unit := range units {
		lost := result.Losses[unit.UnitType]
		returnedAmount := unit.AmountSent - lost
		returned[unit.UnitType] = returnedAmount
		if err := s.missions.UpdateMissionUnitResult(ctx, unit.ID, lost, returnedAmount); err != nil {
			return err
		}
	}
	if err := s.army.ReturnFromMission(ctx, mission.KingdomID, returned); err != nil {
		return err
	}
	if _, err := s.resources.Grant(ctx, mission.KingdomID, result.Rewards); err != nil {
		return err
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		return err
	}
	rewardsJSON, err := json.Marshal(result.Rewards)
	if err != nil {
		return err
	}
	lossesJSON, err := json.Marshal(result.Losses)
	if err != nil {
		return err
	}
	if _, err := s.reports.CreateMissionReport(ctx, mission.KingdomID, mission.ID, reportTitle(cfg), reportBody(cfg, result), result.Result, rewardsJSON, lossesJSON); err != nil {
		return err
	}

	return s.missions.CompleteMission(ctx, mission.ID, now, resultJSON)
}

func (s *MissionService) kingdomForUser(ctx context.Context, userID string) (domain.Kingdom, error) {
	kingdom, err := s.kingdoms.FindByUserID(ctx, userID)
	if errors.Is(err, repository.ErrKingdomNotFound) {
		return domain.Kingdom{}, ErrMissionKingdomNotFound
	}
	if err != nil {
		return domain.Kingdom{}, err
	}
	return kingdom, nil
}

func (s *MissionService) views(ctx context.Context, missions []domain.Mission) ([]MissionView, error) {
	views := make([]MissionView, 0, len(missions))
	for _, mission := range missions {
		view, err := s.view(ctx, mission)
		if err != nil {
			return nil, err
		}
		views = append(views, view)
	}
	return views, nil
}

func (s *MissionService) view(ctx context.Context, mission domain.Mission) (MissionView, error) {
	units, err := s.missions.ListUnitsByMissionID(ctx, mission.ID)
	if err != nil {
		return MissionView{}, err
	}

	unitViews := make([]MissionUnitView, 0, len(units))
	for _, unit := range units {
		unitViews = append(unitViews, MissionUnitView{
			Unit:  unit,
			Label: gameconfig.Units[unit.UnitType].Label,
		})
	}

	var result *MissionResult
	if len(mission.ResultJSON) > 0 {
		var decoded MissionResult
		if err := json.Unmarshal(mission.ResultJSON, &decoded); err == nil {
			result = &decoded
		}
	}

	return MissionView{
		Mission: mission,
		Label:   gameconfig.Missions[mission.Key].Label,
		Units:   unitViews,
		Result:  result,
	}, nil
}

func normalizeMissionUnits(units []StartMissionUnit) (map[string]int64, error) {
	normalized := map[string]int64{}
	for _, unit := range units {
		if !gameconfig.IsUnitType(unit.UnitType) {
			return nil, ErrInvalidMissionUnitType
		}
		if unit.Amount < 0 {
			return nil, ErrInvalidMissionUnitAmount
		}
		if unit.Amount == 0 {
			continue
		}
		normalized[unit.UnitType] += unit.Amount
	}
	if totalSent(normalized) == 0 {
		return nil, ErrInvalidMissionUnitAmount
	}
	return normalized, nil
}

func hasUnits(army ArmyView, sent map[string]int64) bool {
	available := map[string]int64{}
	for _, unit := range army.Units {
		available[unit.Unit.Type] = unit.Unit.Amount
	}
	for unitType, amount := range sent {
		if available[unitType] < amount {
			return false
		}
	}
	return true
}

func missionRequirementsMet(cfg gameconfig.MissionConfig, sent map[string]int64) bool {
	if totalSent(sent) < cfg.MinimumRequirements.TotalUnits {
		return false
	}
	if sent["scouts"] < cfg.MinimumRequirements.Scouts {
		return false
	}
	return true
}

func calculateMissionResult(cfg gameconfig.MissionConfig, units []domain.MissionUnit) MissionResult {
	sent := map[string]int64{}
	for _, unit := range units {
		sent[unit.UnitType] += unit.AmountSent
	}

	outcome := "partial_success"
	if totalSent(sent) >= cfg.RecommendedTotalUnits {
		outcome = "success"
	}

	lossPercent := adjustedLossPercent(cfg, sent)
	losses := distributeLosses(sent, lossPercent, cfg.MaxLossPercent)
	return MissionResult{
		Result:  outcome,
		Rewards: scaledRewards(cfg.BaseRewards, outcome),
		Losses:  losses,
	}
}

func adjustedLossPercent(cfg gameconfig.MissionConfig, sent map[string]int64) int64 {
	lossPercent := cfg.BaseLossPercent
	switch cfg.Key {
	case "black_forest_expedition":
		if sent["scouts"] >= 2 {
			lossPercent -= 3
		}
	case "old_kurgan_expedition":
		if sent["spearmen"] >= 3 {
			lossPercent -= 4
		}
		if sent["archers"] >= 3 {
			lossPercent -= 3
		}
	case "dry_ford_scouting":
		if sent["scouts"] >= 3 {
			lossPercent -= 3
		}
	}
	if lossPercent < 0 {
		return 0
	}
	return lossPercent
}

func distributeLosses(sent map[string]int64, lossPercent int64, maxLossPercent int64) map[string]int64 {
	total := totalSent(sent)
	totalLoss := total * lossPercent / 100
	maxLoss := total * maxLossPercent / 100
	if totalLoss > maxLoss {
		totalLoss = maxLoss
	}

	losses := map[string]int64{}
	remaining := totalLoss
	for _, unitType := range gameconfig.UnitOrder {
		if remaining <= 0 {
			break
		}
		available := sent[unitType]
		lost := available
		if lost > remaining {
			lost = remaining
		}
		losses[unitType] = lost
		remaining -= lost
	}
	return losses
}

func scaledRewards(base gameconfig.ResourceValues, outcome string) gameconfig.ResourceValues {
	scale := int64(100)
	switch outcome {
	case "partial_success":
		scale = 60
	case "failure":
		scale = 20
	}
	return gameconfig.ResourceValues{
		Gold:       base.Gold * scale / 100,
		Food:       base.Food * scale / 100,
		Wood:       base.Wood * scale / 100,
		Stone:      base.Stone * scale / 100,
		Population: base.Population * scale / 100,
	}
}

func totalSent(sent map[string]int64) int64 {
	var total int64
	for _, amount := range sent {
		total += amount
	}
	return total
}

func reportTitle(cfg gameconfig.MissionConfig) string {
	if cfg.Type == "scouting" {
		return "Разведка: " + cfg.Label
	}
	return "Экспедиция в " + cfg.Label
}

func reportBody(cfg gameconfig.MissionConfig, result MissionResult) string {
	switch result.Result {
	case "success":
		return fmt.Sprintf("Отряд вернулся из %s с добычей и вестями.", cfg.Label)
	case "partial_success":
		return fmt.Sprintf("Отряд вернулся из %s не с пустыми руками, но путь был тяжёлым.", cfg.Label)
	default:
		return fmt.Sprintf("Поход к %s едва не сорвался.", cfg.Label)
	}
}

func decodeResourceValues(data []byte) (gameconfig.ResourceValues, error) {
	var values gameconfig.ResourceValues
	if len(data) == 0 {
		return values, nil
	}
	err := json.Unmarshal(data, &values)
	return values, err
}

func decodeLosses(data []byte) (map[string]int64, error) {
	losses := map[string]int64{}
	if len(data) == 0 {
		return losses, nil
	}
	err := json.Unmarshal(data, &losses)
	return losses, err
}
