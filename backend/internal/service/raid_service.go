package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"sumerki/backend/internal/domain"
	"sumerki/backend/internal/gameconfig"
	"sumerki/backend/internal/repository"
)

var (
	ErrRaidKingdomNotFound    = errors.New("kingdom not found")
	ErrRaidTargetNotFound     = errors.New("target not found")
	ErrCannotRaidSelf         = errors.New("cannot raid self")
	ErrInvalidRaidUnitType    = errors.New("invalid unit type")
	ErrInvalidRaidUnitAmount  = errors.New("invalid unit amount")
	ErrRaidInsufficientUnits  = errors.New("insufficient units")
	ErrRaidRequirementsNotMet = errors.New("raid requirements not met")
	ErrTargetTooWeak          = errors.New("target too weak")
	ErrRaidCooldownActive     = errors.New("raid cooldown active")
	ErrTargetUnderProtection  = errors.New("target under protection")
	ErrTargetNewbieProtected  = errors.New("target newbie protected")
)

type RaidRepository interface {
	CreateRaid(ctx context.Context, attackerKingdomID string, defenderKingdomID string, startedAt time.Time, arrivesAt time.Time) (domain.Raid, error)
	CreateRaidUnit(ctx context.Context, raidID string, side string, unitType string, amount int64) (domain.RaidUnit, error)
	ListByKingdomID(ctx context.Context, kingdomID string) ([]domain.Raid, error)
	ListActiveByKingdomID(ctx context.Context, kingdomID string) ([]domain.Raid, error)
	ListUnitsByRaidID(ctx context.Context, raidID string) ([]domain.RaidUnit, error)
	UpdateRaidUnitResult(ctx context.Context, raidUnitID string, lost int64, returned int64) error
	CompleteRaid(ctx context.Context, raidID string, completedAt time.Time, result string, lootJSON []byte, attackerLossesJSON []byte, defenderLossesJSON []byte, resultJSON []byte) error
	CountRecentBetween(ctx context.Context, attackerKingdomID string, defenderKingdomID string, since time.Time) (int64, error)
	CountRecentAgainst(ctx context.Context, defenderKingdomID string, since time.Time) (int64, error)
}

type RaidKingdomRepository interface {
	FindByUserID(ctx context.Context, userID string) (domain.Kingdom, error)
	FindByID(ctx context.Context, kingdomID string) (domain.Kingdom, error)
	ListNeighbors(ctx context.Context, currentKingdomID string, limit int) ([]domain.Kingdom, error)
	AddDread(ctx context.Context, kingdomID string, amount int) error
}

type RaidReportRepository interface {
	CreateReport(ctx context.Context, kingdomID string, missionID *string, reportType string, title string, body string, result string, rewardsJSON []byte, lossesJSON []byte, phasesJSON []byte) (domain.MissionReport, error)
}

type RaidBuildingService interface {
	LevelForKingdom(ctx context.Context, kingdomID string, buildingType string) (int, error)
}

type RaidArmyService interface {
	PrepareForMission(ctx context.Context, kingdomID string) (ArmyView, error)
	SubtractForMission(ctx context.Context, kingdomID string, units map[string]int64) error
	ReturnFromMission(ctx context.Context, kingdomID string, units map[string]int64) error
	CurrentForKingdom(ctx context.Context, kingdomID string) (ArmyView, error)
}

type StartRaidUnit struct {
	UnitType string
	Amount   int64
}

type NeighborView struct {
	Kingdom       domain.Kingdom
	PowerEstimate string
	CanRaid       bool
	BlockedReason *string
}

type RaidUnitView struct {
	Unit  domain.RaidUnit
	Label string
}

type RaidView struct {
	Raid                domain.Raid
	AttackerKingdomName string
	DefenderKingdomName string
	Units               []RaidUnitView
	Loot                gameconfig.ResourceValues
}

type StartRaidResult struct {
	Raid RaidView
	Army ArmyView
}

type raidResultPayload struct {
	Result         string                    `json:"result"`
	Loot           gameconfig.ResourceValues `json:"loot"`
	AttackerLosses map[string]int64          `json:"attackerLosses"`
	DefenderLosses map[string]int64          `json:"defenderLosses"`
}

type RaidService struct {
	kingdoms  RaidKingdomRepository
	raids     RaidRepository
	reports   RaidReportRepository
	army      RaidArmyService
	resources *ResourcesService
	buildings RaidBuildingService
	now       func() time.Time
}

func NewRaidService(kingdoms RaidKingdomRepository, raids RaidRepository, reports RaidReportRepository, army RaidArmyService, resources *ResourcesService, buildings RaidBuildingService) *RaidService {
	return &RaidService{
		kingdoms:  kingdoms,
		raids:     raids,
		reports:   reports,
		army:      army,
		resources: resources,
		buildings: buildings,
		now:       time.Now,
	}
}

func (s *RaidService) Neighbors(ctx context.Context, userID string) ([]NeighborView, error) {
	attacker, err := s.kingdomForUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if err := s.ResolveCompleted(ctx, attacker.ID); err != nil {
		return nil, err
	}

	neighbors, err := s.kingdoms.ListNeighbors(ctx, attacker.ID, 20)
	if err != nil {
		return nil, err
	}

	attackerPower, err := s.powerScore(ctx, attacker)
	if err != nil {
		return nil, err
	}

	views := make([]NeighborView, 0, len(neighbors))
	for _, neighbor := range neighbors {
		defenderPower, err := s.powerScore(ctx, neighbor)
		if err != nil {
			return nil, err
		}
		blocked := s.targetBlockedReason(ctx, attacker, neighbor, attackerPower, defenderPower)
		views = append(views, NeighborView{
			Kingdom:       neighbor,
			PowerEstimate: powerEstimate(attackerPower, defenderPower),
			CanRaid:       blocked == nil,
			BlockedReason: blocked,
		})
	}
	return views, nil
}

func (s *RaidService) Current(ctx context.Context, userID string) ([]RaidView, error) {
	kingdom, err := s.kingdomForUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if err := s.ResolveCompleted(ctx, kingdom.ID); err != nil {
		return nil, err
	}

	raids, err := s.raids.ListByKingdomID(ctx, kingdom.ID)
	if err != nil {
		return nil, err
	}
	return s.views(ctx, raids)
}

func (s *RaidService) Start(ctx context.Context, userID string, defenderKingdomID string, units []StartRaidUnit) (StartRaidResult, error) {
	attacker, err := s.kingdomForUser(ctx, userID)
	if err != nil {
		return StartRaidResult{}, err
	}
	if defenderKingdomID == attacker.ID {
		return StartRaidResult{}, ErrCannotRaidSelf
	}
	defender, err := s.kingdoms.FindByID(ctx, defenderKingdomID)
	if errors.Is(err, repository.ErrKingdomNotFound) {
		return StartRaidResult{}, ErrRaidTargetNotFound
	}
	if err != nil {
		return StartRaidResult{}, err
	}
	if err := s.ResolveCompleted(ctx, attacker.ID); err != nil {
		return StartRaidResult{}, err
	}
	if err := s.ResolveCompleted(ctx, defender.ID); err != nil {
		return StartRaidResult{}, err
	}

	sent, err := normalizeRaidUnits(units)
	if err != nil {
		return StartRaidResult{}, err
	}
	if totalRaidUnits(sent) < gameconfig.MinimumRaidUnits {
		return StartRaidResult{}, ErrRaidRequirementsNotMet
	}
	if totalRaidUnits(sent) > gameconfig.MaximumRaidUnits {
		return StartRaidResult{}, ErrInvalidRaidUnitAmount
	}

	attackerArmy, err := s.army.PrepareForMission(ctx, attacker.ID)
	if err != nil {
		return StartRaidResult{}, err
	}
	if !raidHasUnits(attackerArmy, sent) {
		return StartRaidResult{}, ErrRaidInsufficientUnits
	}
	defenderArmy, err := s.army.CurrentForKingdom(ctx, defender.ID)
	if err != nil {
		return StartRaidResult{}, err
	}

	attackerPower, err := s.powerScore(ctx, attacker)
	if err != nil {
		return StartRaidResult{}, err
	}
	defenderPower, err := s.powerScore(ctx, defender)
	if err != nil {
		return StartRaidResult{}, err
	}
	if blocked := s.targetBlockedReason(ctx, attacker, defender, attackerPower, defenderPower); blocked != nil {
		return StartRaidResult{}, raidBlockedError(*blocked)
	}

	if err := s.army.SubtractForMission(ctx, attacker.ID, sent); err != nil {
		return StartRaidResult{}, err
	}

	startedAt := s.now()
	arrivesAt := startedAt.Add(time.Duration(gameconfig.RaidDurationSeconds) * time.Second)
	raid, err := s.raids.CreateRaid(ctx, attacker.ID, defender.ID, startedAt, arrivesAt)
	if err != nil {
		return StartRaidResult{}, err
	}
	for _, unitType := range gameconfig.UnitOrder {
		if amount := sent[unitType]; amount > 0 {
			if _, err := s.raids.CreateRaidUnit(ctx, raid.ID, "attacker", unitType, amount); err != nil {
				return StartRaidResult{}, err
			}
		}
	}
	for _, unit := range defenderArmy.Units {
		if _, err := s.raids.CreateRaidUnit(ctx, raid.ID, "defender", unit.Unit.Type, unit.Unit.Amount); err != nil {
			return StartRaidResult{}, err
		}
	}
	if err := s.kingdoms.AddDread(ctx, attacker.ID, 1); err != nil {
		return StartRaidResult{}, err
	}

	updatedArmy, err := s.army.CurrentForKingdom(ctx, attacker.ID)
	if err != nil {
		return StartRaidResult{}, err
	}
	view, err := s.view(ctx, raid)
	if err != nil {
		return StartRaidResult{}, err
	}
	return StartRaidResult{Raid: view, Army: updatedArmy}, nil
}

func (s *RaidService) ResolveCompleted(ctx context.Context, kingdomID string) error {
	active, err := s.raids.ListActiveByKingdomID(ctx, kingdomID)
	if err != nil {
		return err
	}
	now := s.now()
	for _, raid := range active {
		if raid.ArrivesAt.After(now) {
			continue
		}
		if err := s.resolveRaid(ctx, raid, now); err != nil && !errors.Is(err, repository.ErrRaidNotFound) {
			return err
		}
	}
	return nil
}

func (s *RaidService) resolveRaid(ctx context.Context, raid domain.Raid, now time.Time) error {
	units, err := s.raids.ListUnitsByRaidID(ctx, raid.ID)
	if err != nil {
		return err
	}
	attackerUnits, defenderUnits := splitRaidUnits(units)
	attackerScore := raidAttackerScore(attackerUnits)
	defenderScore := raidDefenderScore(defenderUnits)
	wallsLevel, _ := s.buildings.LevelForKingdom(ctx, raid.DefenderKingdomID, "walls")
	defenderScore = defenderScore * int64(100+wallsLevel*10) / 100
	defender, err := s.kingdoms.FindByID(ctx, raid.DefenderKingdomID)
	if err != nil {
		return err
	}
	if defender.Patron != nil && (*defender.Patron == "empire_of_dusk" || *defender.Patron == "old_pact") {
		defenderScore = defenderScore * 105 / 100
	}

	result := raidResult(attackerScore, defenderScore)
	attackerLosses := raidLosses(attackerUnits, attackerLossPercent(result))
	defenderLosses := raidLosses(defenderUnits, defenderLossPercent(result))
	returned := raidReturned(attackerUnits, attackerLosses)
	for _, unit := range attackerUnits {
		if err := s.raids.UpdateRaidUnitResult(ctx, unit.ID, attackerLosses[unit.UnitType], returned[unit.UnitType]); err != nil {
			return err
		}
	}
	for _, unit := range defenderUnits {
		if err := s.raids.UpdateRaidUnitResult(ctx, unit.ID, defenderLosses[unit.UnitType], unit.AmountSent-defenderLosses[unit.UnitType]); err != nil {
			return err
		}
	}
	if err := s.army.ReturnFromMission(ctx, raid.AttackerKingdomID, returned); err != nil {
		return err
	}

	loot := gameconfig.ResourceValues{}
	switch result {
	case "attacker_success":
		loot, err = s.resources.TransferRaidLoot(ctx, raid.AttackerKingdomID, raid.DefenderKingdomID, gameconfig.MaxRaidLootPercent)
	case "bloody_stalemate":
		loot, err = s.resources.TransferRaidStalemateLoot(ctx, raid.AttackerKingdomID, raid.DefenderKingdomID, gameconfig.StalemateLootPercent)
	}
	if err != nil {
		return err
	}

	if result == "attacker_success" {
		if err := s.kingdoms.AddDread(ctx, raid.AttackerKingdomID, 2); err != nil {
			return err
		}
	} else {
		if err := s.kingdoms.AddDread(ctx, raid.AttackerKingdomID, 1); err != nil {
			return err
		}
	}

	lootJSON, _ := json.Marshal(loot)
	attackerLossesJSON, _ := json.Marshal(attackerLosses)
	defenderLossesJSON, _ := json.Marshal(defenderLosses)
	resultJSON, _ := json.Marshal(raidResultPayload{Result: result, Loot: loot, AttackerLosses: attackerLosses, DefenderLosses: defenderLosses})
	if err := s.raids.CompleteRaid(ctx, raid.ID, now, result, lootJSON, attackerLossesJSON, defenderLossesJSON, resultJSON); err != nil {
		return err
	}
	return s.createRaidReports(ctx, raid, result, loot, attackerLosses, defenderLosses)
}

func (s *RaidService) createRaidReports(ctx context.Context, raid domain.Raid, result string, loot gameconfig.ResourceValues, attackerLosses map[string]int64, defenderLosses map[string]int64) error {
	attacker, err := s.kingdoms.FindByID(ctx, raid.AttackerKingdomID)
	if err != nil {
		return err
	}
	defender, err := s.kingdoms.FindByID(ctx, raid.DefenderKingdomID)
	if err != nil {
		return err
	}
	attackerTemplate := attackerRaidReport(defender.Name, result)
	defenderTemplate := defenderRaidReport(attacker.Name, result)
	lootJSON, _ := json.Marshal(loot)
	attackerLossesJSON, _ := json.Marshal(attackerLosses)
	defenderLossesJSON, _ := json.Marshal(defenderLosses)
	attackerPhasesJSON, _ := json.Marshal(attackerTemplate.Phases)
	defenderPhasesJSON, _ := json.Marshal(defenderTemplate.Phases)
	if _, err := s.reports.CreateReport(ctx, attacker.ID, nil, "pvp_raid_attacker", attackerTemplate.Title, attackerTemplate.Body, result, lootJSON, attackerLossesJSON, attackerPhasesJSON); err != nil {
		return err
	}
	if _, err := s.reports.CreateReport(ctx, defender.ID, nil, "pvp_raid_defender", defenderTemplate.Title, defenderTemplate.Body, result, lootJSON, defenderLossesJSON, defenderPhasesJSON); err != nil {
		return err
	}
	return nil
}

func (s *RaidService) kingdomForUser(ctx context.Context, userID string) (domain.Kingdom, error) {
	kingdom, err := s.kingdoms.FindByUserID(ctx, userID)
	if errors.Is(err, repository.ErrKingdomNotFound) {
		return domain.Kingdom{}, ErrRaidKingdomNotFound
	}
	if err != nil {
		return domain.Kingdom{}, err
	}
	return kingdom, nil
}

func (s *RaidService) targetBlockedReason(ctx context.Context, attacker domain.Kingdom, defender domain.Kingdom, attackerPower int64, defenderPower int64) *string {
	now := s.now()
	if now.Sub(defender.CreatedAt) < gameconfig.NewKingdomProtectionAge {
		return stringPtr("target_newbie_protected")
	}
	if defenderPower > 0 && attackerPower > defenderPower*3 {
		return stringPtr("target_too_weak")
	}
	if count, err := s.raids.CountRecentBetween(ctx, attacker.ID, defender.ID, now.Add(-gameconfig.SameTargetRaidCooldown)); err == nil && count > 0 {
		return stringPtr("raid_cooldown_active")
	}
	if count, err := s.raids.CountRecentAgainst(ctx, defender.ID, now.Add(-gameconfig.DefenderProtectionWindow)); err == nil && count >= gameconfig.DefenderProtectionMaxHits {
		return stringPtr("target_under_protection")
	}
	return nil
}

func (s *RaidService) powerScore(ctx context.Context, kingdom domain.Kingdom) (int64, error) {
	army, err := s.army.CurrentForKingdom(ctx, kingdom.ID)
	if err != nil {
		return 0, err
	}
	var score int64
	for _, unit := range army.Units {
		stats := gameconfig.Units[unit.Unit.Type].Stats
		score += unit.Unit.Amount * int64(stats.Attack+stats.Defense+stats.Speed)
	}
	buildings := []string{"town_hall", "farm", "lumberyard", "quarry", "market", "barracks", "walls", "shrine"}
	for _, buildingType := range buildings {
		level, err := s.buildings.LevelForKingdom(ctx, kingdom.ID, buildingType)
		if err != nil {
			continue
		}
		score += int64(level) * 20
		if buildingType == "town_hall" {
			score += int64(level) * 50
		}
	}
	return score, nil
}

func (s *RaidService) views(ctx context.Context, raids []domain.Raid) ([]RaidView, error) {
	views := make([]RaidView, 0, len(raids))
	for _, raid := range raids {
		view, err := s.view(ctx, raid)
		if err != nil {
			return nil, err
		}
		views = append(views, view)
	}
	return views, nil
}

func (s *RaidService) view(ctx context.Context, raid domain.Raid) (RaidView, error) {
	attacker, err := s.kingdoms.FindByID(ctx, raid.AttackerKingdomID)
	if err != nil {
		return RaidView{}, err
	}
	defender, err := s.kingdoms.FindByID(ctx, raid.DefenderKingdomID)
	if err != nil {
		return RaidView{}, err
	}
	units, err := s.raids.ListUnitsByRaidID(ctx, raid.ID)
	if err != nil {
		return RaidView{}, err
	}
	unitViews := make([]RaidUnitView, 0, len(units))
	for _, unit := range units {
		if unit.Side != "attacker" {
			continue
		}
		unitViews = append(unitViews, RaidUnitView{Unit: unit, Label: gameconfig.Units[unit.UnitType].Label})
	}
	loot, _ := decodeResourceValues(raid.LootJSON)
	return RaidView{Raid: raid, AttackerKingdomName: attacker.Name, DefenderKingdomName: defender.Name, Units: unitViews, Loot: loot}, nil
}

func normalizeRaidUnits(units []StartRaidUnit) (map[string]int64, error) {
	normalized := map[string]int64{}
	for _, unit := range units {
		if !gameconfig.IsUnitType(unit.UnitType) {
			return nil, ErrInvalidRaidUnitType
		}
		if unit.Amount < 0 {
			return nil, ErrInvalidRaidUnitAmount
		}
		if unit.Amount == 0 {
			continue
		}
		normalized[unit.UnitType] += unit.Amount
	}
	if totalRaidUnits(normalized) == 0 {
		return nil, ErrInvalidRaidUnitAmount
	}
	return normalized, nil
}

func raidHasUnits(army ArmyView, sent map[string]int64) bool {
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

func totalRaidUnits(units map[string]int64) int64 {
	var total int64
	for _, amount := range units {
		total += amount
	}
	return total
}

func splitRaidUnits(units []domain.RaidUnit) ([]domain.RaidUnit, []domain.RaidUnit) {
	attacker := []domain.RaidUnit{}
	defender := []domain.RaidUnit{}
	for _, unit := range units {
		if unit.Side == "attacker" {
			attacker = append(attacker, unit)
		} else {
			defender = append(defender, unit)
		}
	}
	return attacker, defender
}

func raidAttackerScore(units []domain.RaidUnit) int64 {
	var score int64
	for _, unit := range units {
		stats := gameconfig.Units[unit.UnitType].Stats
		score += unit.AmountSent * int64(stats.Attack)
		score += unit.AmountSent * int64(stats.Speed) / 2
	}
	if score <= 0 {
		return 1
	}
	return score
}

func raidDefenderScore(units []domain.RaidUnit) int64 {
	var score int64
	for _, unit := range units {
		stats := gameconfig.Units[unit.UnitType].Stats
		score += unit.AmountSent * int64(stats.Defense)
	}
	if score <= 0 {
		return 1
	}
	return score
}

func raidResult(attackerScore int64, defenderScore int64) string {
	if attackerScore*100 >= defenderScore*115 {
		return "attacker_success"
	}
	if defenderScore*100 >= attackerScore*115 {
		return "defender_success"
	}
	return "bloody_stalemate"
}

func attackerLossPercent(result string) int64 {
	switch result {
	case "attacker_success":
		return 5
	case "bloody_stalemate":
		return 12
	default:
		return 20
	}
}

func defenderLossPercent(result string) int64 {
	switch result {
	case "attacker_success":
		return 3
	case "bloody_stalemate":
		return 2
	default:
		return 0
	}
}

func raidLosses(units []domain.RaidUnit, percent int64) map[string]int64 {
	losses := map[string]int64{}
	for _, unit := range units {
		losses[unit.UnitType] = unit.AmountSent * percent / 100
	}
	return losses
}

func raidReturned(units []domain.RaidUnit, losses map[string]int64) map[string]int64 {
	returned := map[string]int64{}
	for _, unit := range units {
		returned[unit.UnitType] = unit.AmountSent - losses[unit.UnitType]
	}
	return returned
}

func powerEstimate(attackerPower int64, defenderPower int64) string {
	if defenderPower*2 <= attackerPower {
		return "much_weaker"
	}
	if defenderPower*4 <= attackerPower*3 {
		return "weaker"
	}
	if defenderPower >= attackerPower*2 {
		return "much_stronger"
	}
	if defenderPower*3 >= attackerPower*4 {
		return "stronger"
	}
	return "similar"
}

func raidBlockedError(reason string) error {
	switch reason {
	case "target_newbie_protected":
		return ErrTargetNewbieProtected
	case "target_too_weak":
		return ErrTargetTooWeak
	case "raid_cooldown_active":
		return ErrRaidCooldownActive
	case "target_under_protection":
		return ErrTargetUnderProtection
	default:
		return nil
	}
}

type raidReportTemplate struct {
	Title  string
	Body   string
	Phases []gameconfig.ReportPhase
}

func attackerRaidReport(defenderName string, result string) raidReportTemplate {
	switch result {
	case "attacker_success":
		return raidReportTemplate{
			Title: "Набег на " + defenderName,
			Body:  "Отряд вернулся с добычей. Стены врага не рухнули, но складские двери всё же открылись.",
			Phases: []gameconfig.ReportPhase{
				{Title: "Сбор отряда", Body: "Люди вышли до рассвета, пряча железо под мокрыми плащами."},
				{Title: "Дорога к цели", Body: "Разведчики вели короткой тропой, где не слышно колоколов."},
				{Title: "Стычка у стен", Body: "Стража дрогнула, и складские двери раскрылись под ударами."},
				{Title: "Возвращение", Body: "К воротам принесли добычу и новые дурные слухи."},
			},
		}
	case "defender_success":
		return raidReportTemplate{
			Title: "Набег на " + defenderName + " отброшен",
			Body:  "Отряд вернулся ни с чем. Дорога назад была длиннее, чем обещали разведчики.",
			Phases: []gameconfig.ReportPhase{
				{Title: "Сбор отряда", Body: "Копья считали быстро, будто удача уже ждала за воротами."},
				{Title: "Дорога к цели", Body: "У цели оказалось больше огней и часовых, чем говорили слухи."},
				{Title: "Стычка у стен", Body: "Защитники встретили набег плотным строем и не пустили к складам."},
				{Title: "Возвращение", Body: "Люди вернулись усталыми, неся только потери и злость."},
			},
		}
	default:
		return raidReportTemplate{
			Title: "Кровавый набег на " + defenderName,
			Body:  "Набег не стал победой, но и защитники запомнят эту ночь.",
			Phases: []gameconfig.ReportPhase{
				{Title: "Сбор отряда", Body: "Отряд ушёл тихо, без песен и лишнего огня."},
				{Title: "Дорога к цели", Body: "Ночь тянулась долго, пока впереди не показались стены."},
				{Title: "Стычка у стен", Body: "У ворот смешались крики, стрелы и тяжёлое дыхание."},
				{Title: "Возвращение", Body: "Добыча мала, но набег уже стал разговором на дорогах."},
			},
		}
	}
}

func defenderRaidReport(attackerName string, result string) raidReportTemplate {
	body := "Чужой отряд подошёл к владению, но город устоял. Потери ограничены, но слухи уже пошли по рынку."
	if result == "attacker_success" {
		body = "Чужой отряд добрался до складов и ушёл с частью припасов. Город устоял, но рынки говорят шёпотом."
	}
	return raidReportTemplate{
		Title: "Ночной набег",
		Body:  body,
		Phases: []gameconfig.ReportPhase{
			{Title: "Тревога", Body: "Дозорные подняли крик, когда у дороги заметили людей " + attackerName + "."},
			{Title: "Стража у ворот", Body: "Ворота удержали, и стража не дала бою перейти к домам."},
			{Title: "Исход набега", Body: "Когда шум стих, писари пошли считать урон и пропавшие мешки."},
			{Title: "Последствия", Body: "Город цел, но память о ночном набеге останется у стен."},
		},
	}
}

func stringPtr(value string) *string {
	return &value
}
