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
	ErrNoPatronSelected        = errors.New("no patron selected")
	ErrNoTributeDue            = errors.New("no tribute due")
	ErrInvalidCrisisChoice     = errors.New("invalid crisis choice")
	ErrCrisisChoiceUnavailable = errors.New("crisis choice unavailable")
)

type PatronPressureRepository interface {
	FindByKingdomID(ctx context.Context, kingdomID string) (domain.PatronPressureState, error)
	UpsertForPatron(ctx context.Context, kingdomID string, patron string, nextTributeAt time.Time, resetDebt bool) (domain.PatronPressureState, error)
	Save(ctx context.Context, state domain.PatronPressureState) (domain.PatronPressureState, error)
	ClearForKingdom(ctx context.Context, kingdomID string) error
}

type PatronPressureView struct {
	State             domain.PatronPressureState
	PatronLabel       string
	Summary           string
	AvailableActions  []string
	ProtectedMinimums gameconfig.ResourceValues
}

type PatronPressureResult struct {
	Pressure *PatronPressureView
}

type PatronPressurePaymentResult struct {
	Pressure  PatronPressureView
	Resources ResourcesResult
}

type PatronPressureBreakResult struct {
	Kingdom domain.Kingdom
}

type PatronPressureCrisisResult struct {
	Pressure *PatronPressureView
	Kingdom  *domain.Kingdom
}

type PatronPressureService struct {
	kingdoms  PatronKingdomRepository
	patrons   PatronRepository
	states    PatronPressureRepository
	resources *ResourcesService
	now       func() time.Time
}

func NewPatronPressureService(kingdoms PatronKingdomRepository, patrons PatronRepository, states PatronPressureRepository, resources *ResourcesService) *PatronPressureService {
	return &PatronPressureService{
		kingdoms:  kingdoms,
		patrons:   patrons,
		states:    states,
		resources: resources,
		now:       time.Now,
	}
}

func (s *PatronPressureService) Current(ctx context.Context, userID string) (PatronPressureResult, error) {
	kingdom, err := s.kingdomForUser(ctx, userID)
	if err != nil {
		return PatronPressureResult{}, err
	}
	if kingdom.Patron == nil {
		return PatronPressureResult{Pressure: nil}, nil
	}

	state, err := s.ensureState(ctx, kingdom)
	if err != nil {
		return PatronPressureResult{}, err
	}
	state, err = s.resolve(ctx, state)
	if err != nil {
		return PatronPressureResult{}, err
	}
	view := pressureView(state)
	return PatronPressureResult{Pressure: &view}, nil
}

func (s *PatronPressureService) PayTribute(ctx context.Context, userID string) (PatronPressurePaymentResult, error) {
	kingdom, err := s.kingdomForUser(ctx, userID)
	if err != nil {
		return PatronPressurePaymentResult{}, err
	}
	if kingdom.Patron == nil {
		return PatronPressurePaymentResult{}, ErrNoPatronSelected
	}

	state, err := s.ensureState(ctx, kingdom)
	if err != nil {
		return PatronPressurePaymentResult{}, err
	}
	state, err = s.resolve(ctx, state)
	if err != nil {
		return PatronPressurePaymentResult{}, err
	}

	due := gameconfig.ResourceValues{
		Gold: state.TributeDebtGold,
		Food: state.TributeDebtFood + state.ContributionDebtFood,
	}
	if due.Gold == 0 && due.Food == 0 {
		return PatronPressurePaymentResult{}, ErrNoTributeDue
	}

	paid, resources, err := s.resources.SpendAboveProtected(ctx, kingdom.ID, due, gameconfig.PatronPressureProtectedMinimums)
	if err != nil {
		return PatronPressurePaymentResult{}, err
	}
	paidContribution := minInt64(state.ContributionDebtFood, paid.Food)
	paidTributeFood := paid.Food - paidContribution
	state.ContributionDebtFood -= paidContribution
	state.TributeDebtFood -= paidTributeFood
	state.TributeDebtGold -= paid.Gold
	if state.TributeDebtGold < 0 {
		state.TributeDebtGold = 0
	}
	if state.TributeDebtFood < 0 {
		state.TributeDebtFood = 0
	}
	if state.ContributionDebtFood < 0 {
		state.ContributionDebtFood = 0
	}
	if state.TributeDebtGold == 0 && state.TributeDebtFood == 0 && state.ContributionDebtFood == 0 {
		state.PressureLevel = clampPressure(state.PressureLevel - 15)
	} else {
		state.PressureLevel = clampPressure(state.PressureLevel - 5)
	}
	state.CrisisStatus = crisisStatus(state.PressureLevel)
	state.LastResolvedAt = s.now()
	state, err = s.states.Save(ctx, state)
	if err != nil {
		return PatronPressurePaymentResult{}, err
	}

	return PatronPressurePaymentResult{
		Pressure:  pressureView(state),
		Resources: resources,
	}, nil
}

func (s *PatronPressureService) ChooseCrisis(ctx context.Context, userID string, choice string) (PatronPressureCrisisResult, error) {
	kingdom, err := s.kingdomForUser(ctx, userID)
	if err != nil {
		return PatronPressureCrisisResult{}, err
	}
	if kingdom.Patron == nil {
		return PatronPressureCrisisResult{}, ErrNoPatronSelected
	}

	switch choice {
	case "ask_delay":
		state, err := s.ensureState(ctx, kingdom)
		if err != nil {
			return PatronPressureCrisisResult{}, err
		}
		if state.Patron == "independent" {
			return PatronPressureCrisisResult{}, ErrCrisisChoiceUnavailable
		}
		now := s.now()
		delayUntil := now.Add(gameconfig.PatronDelayDuration)
		state.DelayUntil = &delayUntil
		state.CrisisStatus = "delayed"
		state.PressureLevel = clampPressure(state.PressureLevel + 5)
		state.LastResolvedAt = now
		state, err = s.states.Save(ctx, state)
		if err != nil {
			return PatronPressureCrisisResult{}, err
		}
		view := pressureView(state)
		return PatronPressureCrisisResult{Pressure: &view}, nil
	case "break_patron":
		if err := s.patrons.BreakForKingdom(ctx, kingdom.ID); err != nil {
			return PatronPressureCrisisResult{}, err
		}
		if err := s.ClearForBreak(ctx, kingdom.ID); err != nil {
			return PatronPressureCrisisResult{}, err
		}
		updated, err := s.kingdoms.UpdatePatronByID(ctx, kingdom.ID, nil)
		if err != nil {
			return PatronPressureCrisisResult{}, err
		}
		return PatronPressureCrisisResult{Pressure: nil, Kingdom: &updated}, nil
	default:
		return PatronPressureCrisisResult{}, ErrInvalidCrisisChoice
	}
}

func (s *PatronPressureService) EnsureForJoin(ctx context.Context, kingdomID string, patron string, samePatron bool) error {
	next := s.now()
	switch patron {
	case "empire_of_dusk":
		next = next.Add(gameconfig.EmpireTributeInterval)
	case "old_pact":
		next = next.Add(gameconfig.OldPactInterval)
	}
	resetDebt := !samePatron || patron == "independent" || patron == "old_pact"
	_, err := s.states.UpsertForPatron(ctx, kingdomID, patron, next, resetDebt)
	return err
}

func (s *PatronPressureService) ClearForBreak(ctx context.Context, kingdomID string) error {
	return s.states.ClearForKingdom(ctx, kingdomID)
}

func (s *PatronPressureService) ResolveForKingdom(ctx context.Context, kingdom domain.Kingdom) error {
	if kingdom.Patron == nil {
		return nil
	}
	state, err := s.ensureState(ctx, kingdom)
	if err != nil {
		return err
	}
	_, err = s.resolve(ctx, state)
	return err
}

func (s *PatronPressureService) ensureState(ctx context.Context, kingdom domain.Kingdom) (domain.PatronPressureState, error) {
	state, err := s.states.FindByKingdomID(ctx, kingdom.ID)
	if errors.Is(err, repository.ErrPatronPressureNotFound) {
		next := s.now()
		if kingdom.Patron != nil {
			switch *kingdom.Patron {
			case "empire_of_dusk":
				next = next.Add(gameconfig.EmpireTributeInterval)
			case "old_pact":
				next = next.Add(gameconfig.OldPactInterval)
			}
			return s.states.UpsertForPatron(ctx, kingdom.ID, *kingdom.Patron, next, true)
		}
	}
	if err != nil {
		return domain.PatronPressureState{}, err
	}
	if kingdom.Patron != nil && state.Patron != *kingdom.Patron {
		return s.states.UpsertForPatron(ctx, kingdom.ID, *kingdom.Patron, s.now().Add(patronInterval(*kingdom.Patron)), true)
	}
	return state, nil
}

func (s *PatronPressureService) resolve(ctx context.Context, state domain.PatronPressureState) (domain.PatronPressureState, error) {
	now := s.now()
	if state.Patron == "independent" {
		state.TributeDebtGold = 0
		state.TributeDebtFood = 0
		state.ContributionDebtFood = 0
		state.PressureLevel = 0
		state.CrisisStatus = "none"
		state.DelayUntil = nil
		state.LastResolvedAt = now
		return s.states.Save(ctx, state)
	}
	if state.DelayUntil != nil && now.Before(*state.DelayUntil) {
		state.CrisisStatus = "delayed"
		state.LastResolvedAt = now
		return s.states.Save(ctx, state)
	}
	if state.DelayUntil != nil && !now.Before(*state.DelayUntil) {
		state.DelayUntil = nil
	}

	interval := patronInterval(state.Patron)
	processed := 0
	for !state.NextTributeAt.After(now) && processed < gameconfig.PatronPressureMaxIntervals {
		var err error
		if state.Patron == "empire_of_dusk" {
			state, err = s.resolveEmpireInterval(ctx, state)
		} else if state.Patron == "old_pact" {
			state, err = s.resolveOldPactInterval(ctx, state)
		}
		if err != nil {
			return domain.PatronPressureState{}, err
		}
		state.NextTributeAt = state.NextTributeAt.Add(interval)
		processed++
	}
	state.LastResolvedAt = now
	state.CrisisStatus = crisisStatus(state.PressureLevel)
	return s.states.Save(ctx, state)
}

func (s *PatronPressureService) resolveEmpireInterval(ctx context.Context, state domain.PatronPressureState) (domain.PatronPressureState, error) {
	result, err := s.resources.CurrentForKingdom(ctx, state.KingdomID)
	if err != nil {
		return domain.PatronPressureState{}, err
	}
	due := empireDue(result.Resources)
	if due.Gold == 0 && due.Food == 0 {
		due.Gold = gameconfig.EmpireMinimumGoldDue
		due.Food = gameconfig.EmpireMinimumFoodDue
	}
	paid, _, err := s.resources.SpendAboveProtected(ctx, state.KingdomID, due, gameconfig.PatronPressureProtectedMinimums)
	if err != nil {
		return domain.PatronPressureState{}, err
	}
	unpaidGold := due.Gold - paid.Gold
	unpaidFood := due.Food - paid.Food
	state.TributeDebtGold += unpaidGold
	state.TributeDebtFood += unpaidFood
	if unpaidGold == 0 && unpaidFood == 0 {
		state.PressureLevel = clampPressure(state.PressureLevel - 5)
	} else if paid.Gold > 0 || paid.Food > 0 {
		state.PressureLevel = clampPressure(state.PressureLevel + 5)
	} else {
		state.PressureLevel = clampPressure(state.PressureLevel + 10)
	}
	return state, nil
}

func (s *PatronPressureService) resolveOldPactInterval(ctx context.Context, state domain.PatronPressureState) (domain.PatronPressureState, error) {
	result, err := s.resources.CurrentForKingdom(ctx, state.KingdomID)
	if err != nil {
		return domain.PatronPressureState{}, err
	}
	surplusFood := result.Resources.Food - gameconfig.PatronPressureProtectedMinimums.Food
	if surplusFood > 0 {
		state.ContributionDebtFood += surplusFood * gameconfig.OldPactFoodPercent / 100
		if state.PressureLevel < gameconfig.OldPactPressureCap {
			state.PressureLevel += 2
			if state.PressureLevel > gameconfig.OldPactPressureCap {
				state.PressureLevel = gameconfig.OldPactPressureCap
			}
		}
	}
	return state, nil
}

func (s *PatronPressureService) kingdomForUser(ctx context.Context, userID string) (domain.Kingdom, error) {
	kingdom, err := s.kingdoms.FindByUserID(ctx, userID)
	if errors.Is(err, repository.ErrKingdomNotFound) {
		return domain.Kingdom{}, ErrPatronKingdomNotFound
	}
	if err != nil {
		return domain.Kingdom{}, err
	}
	return kingdom, nil
}

func pressureView(state domain.PatronPressureState) PatronPressureView {
	return PatronPressureView{
		State:             state,
		PatronLabel:       gameconfig.Patrons[state.Patron].Label,
		Summary:           pressureSummary(state),
		AvailableActions:  pressureActions(state),
		ProtectedMinimums: pressureProtectedMinimums(state.Patron),
	}
}

func pressureSummary(state domain.PatronPressureState) string {
	switch state.Patron {
	case "independent":
		return "Независимые владения не платят дань и не имеют покровительского давления."
	case "old_pact":
		return "Старый Договор ждёт вклад, но строгие обязательства появятся позже."
	default:
		if state.CrisisStatus == "delayed" {
			return "Имперский сборщик получил отсрочку, но запомнил слабость двора."
		}
		return "Империя ждёт дань, но не может забрать неприкосновенный запас."
	}
}

func pressureActions(state domain.PatronPressureState) []string {
	if state.Patron == "independent" {
		return []string{}
	}
	actions := []string{"break_patron"}
	if state.TributeDebtGold > 0 || state.TributeDebtFood > 0 || state.ContributionDebtFood > 0 {
		actions = append([]string{"pay_tribute"}, actions...)
	}
	if state.CrisisStatus == "warning" || state.CrisisStatus == "active" {
		actions = append(actions, "ask_delay")
	}
	return actions
}

func pressureProtectedMinimums(patron string) gameconfig.ResourceValues {
	if patron == "old_pact" {
		return gameconfig.ResourceValues{Food: gameconfig.PatronPressureProtectedMinimums.Food}
	}
	if patron == "independent" {
		return gameconfig.ResourceValues{}
	}
	return gameconfig.PatronPressureProtectedMinimums
}

func empireDue(resources domain.Resources) gameconfig.ResourceValues {
	goldSurplus := resources.Gold - gameconfig.PatronPressureProtectedMinimums.Gold
	foodSurplus := resources.Food - gameconfig.PatronPressureProtectedMinimums.Food
	var due gameconfig.ResourceValues
	if goldSurplus > 0 {
		due.Gold = maxInt64(gameconfig.EmpireMinimumGoldDue, goldSurplus*gameconfig.EmpireTributeGoldPercent/100)
	}
	if foodSurplus > 0 {
		due.Food = maxInt64(gameconfig.EmpireMinimumFoodDue, foodSurplus*gameconfig.EmpireTributeFoodPercent/100)
	}
	return due
}

func patronInterval(patron string) time.Duration {
	if patron == "old_pact" {
		return gameconfig.OldPactInterval
	}
	if patron == "empire_of_dusk" {
		return gameconfig.EmpireTributeInterval
	}
	return 0
}

func crisisStatus(pressure int) string {
	if pressure >= 60 {
		return "active"
	}
	if pressure >= 30 {
		return "warning"
	}
	return "none"
}

func clampPressure(value int) int {
	if value < 0 {
		return 0
	}
	if value > 100 {
		return 100
	}
	return value
}

func minInt64(a int64, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func maxInt64(a int64, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
