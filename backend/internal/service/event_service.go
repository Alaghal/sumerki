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

const maxActiveEvents = 3

var (
	ErrEventKingdomNotFound    = errors.New("kingdom not found")
	ErrEventNotFound           = errors.New("event not found")
	ErrEventExpired            = errors.New("event expired")
	ErrEventAlreadyResolved    = errors.New("event already resolved")
	ErrInvalidEventChoice      = errors.New("invalid event choice")
	ErrEventChoiceNotAvailable = errors.New("event choice not available")
)

type EventRepository interface {
	ListActiveGameEvents(ctx context.Context) ([]domain.GameEvent, error)
	ListActiveByKingdomID(ctx context.Context, kingdomID string) ([]domain.KingdomEvent, error)
	ListRecentResolvedByKingdomID(ctx context.Context, kingdomID string, limit int) ([]domain.KingdomEvent, error)
	CreateKingdomEvent(ctx context.Context, kingdomID string, gameEventID string, generatedAt time.Time, expiresAt time.Time) (domain.KingdomEvent, error)
	ExpireStale(ctx context.Context, kingdomID string, now time.Time) error
	FindByIDAndKingdomID(ctx context.Context, eventID string, kingdomID string) (domain.KingdomEvent, error)
	MarkResolved(ctx context.Context, eventID string, kingdomID string, choiceKey string, resolvedAt time.Time, resultJSON []byte) (domain.KingdomEvent, error)
	HasRecentByEventKey(ctx context.Context, kingdomID string, eventKey string, since time.Time) (bool, error)
}

type EventKingdomRepository interface {
	FindByUserID(ctx context.Context, userID string) (domain.Kingdom, error)
	ApplyReputationDelta(ctx context.Context, kingdomID string, dreadDelta int, honorDelta int) (domain.Kingdom, error)
}

type EventArmyService interface {
	ApplyUnitDeltaForKingdom(ctx context.Context, kingdomID string, units map[string]int64) (ArmyView, error)
}

type EventPatronRepository interface {
	FindByKingdomID(ctx context.Context, kingdomID string) (domain.PatronRelation, error)
	AdjustFavor(ctx context.Context, kingdomID string, delta int) (domain.PatronRelation, error)
}

type EventReportRepository interface {
	CreateReport(ctx context.Context, kingdomID string, missionID *string, reportType string, title string, body string, result string, rewardsJSON []byte, lossesJSON []byte, phasesJSON []byte) (domain.MissionReport, error)
}

type EventListResult struct {
	Events      []KingdomEventView
	ActiveCount int
}

type EventChoiceResult struct {
	Event     KingdomEventView
	Resources *ResourcesResult
	Army      *ArmyView
	Kingdom   domain.Kingdom
}

type KingdomEventView struct {
	Event   domain.KingdomEvent
	Choices []EventChoiceView
	Result  *EventResult
}

type EventChoiceView struct {
	Choice domain.EventChoice
}

type EventResult struct {
	Title          string              `json:"title"`
	Body           string              `json:"body"`
	AppliedEffects EventAppliedEffects `json:"appliedEffects"`
}

type EventAppliedEffects struct {
	ResourceDelta    *EventResourceDelta `json:"resourceDelta,omitempty"`
	UnitDelta        map[string]int64    `json:"unitDelta,omitempty"`
	KingdomDelta     *EventKingdomDelta  `json:"kingdomDelta,omitempty"`
	PatronFavorDelta int                 `json:"patronFavorDelta,omitempty"`
}

type EventResourceDelta struct {
	Gold       int64 `json:"gold,omitempty"`
	Food       int64 `json:"food,omitempty"`
	Wood       int64 `json:"wood,omitempty"`
	Stone      int64 `json:"stone,omitempty"`
	Population int64 `json:"population,omitempty"`
}

type EventKingdomDelta struct {
	Dread int `json:"dread,omitempty"`
	Honor int `json:"honor,omitempty"`
}

type eventConditions struct {
	RequiresPatron string `json:"requiresPatron"`
}

type EventService struct {
	kingdoms  EventKingdomRepository
	events    EventRepository
	resources *ResourcesService
	army      EventArmyService
	patrons   EventPatronRepository
	reports   EventReportRepository
	now       func() time.Time
}

func NewEventService(kingdoms EventKingdomRepository, events EventRepository, resources *ResourcesService, army EventArmyService, patrons EventPatronRepository, reports EventReportRepository) *EventService {
	return &EventService{
		kingdoms:  kingdoms,
		events:    events,
		resources: resources,
		army:      army,
		patrons:   patrons,
		reports:   reports,
		now:       time.Now,
	}
}

func (s *EventService) Current(ctx context.Context, userID string, includeResolved bool, limit int) (EventListResult, error) {
	kingdom, err := s.kingdomForUser(ctx, userID)
	if err != nil {
		return EventListResult{}, err
	}
	if err := s.refresh(ctx, kingdom); err != nil {
		return EventListResult{}, err
	}

	active, err := s.events.ListActiveByKingdomID(ctx, kingdom.ID)
	if err != nil {
		return EventListResult{}, err
	}
	all := make([]domain.KingdomEvent, 0, limit)
	all = append(all, active...)
	if includeResolved && len(all) < limit {
		recent, err := s.events.ListRecentResolvedByKingdomID(ctx, kingdom.ID, limit-len(all))
		if err != nil {
			return EventListResult{}, err
		}
		all = append(all, recent...)
	}

	views, err := s.views(all)
	if err != nil {
		return EventListResult{}, err
	}
	return EventListResult{Events: views, ActiveCount: len(active)}, nil
}

func (s *EventService) Choose(ctx context.Context, userID string, eventID string, choiceKey string) (EventChoiceResult, error) {
	kingdom, err := s.kingdomForUser(ctx, userID)
	if err != nil {
		return EventChoiceResult{}, err
	}
	if err := s.events.ExpireStale(ctx, kingdom.ID, s.now()); err != nil {
		return EventChoiceResult{}, err
	}

	event, err := s.events.FindByIDAndKingdomID(ctx, eventID, kingdom.ID)
	if errors.Is(err, repository.ErrEventNotFound) {
		return EventChoiceResult{}, ErrEventNotFound
	}
	if err != nil {
		return EventChoiceResult{}, err
	}
	if event.Status == "resolved" {
		return EventChoiceResult{}, ErrEventAlreadyResolved
	}
	if event.Status == "expired" || !event.ExpiresAt.After(s.now()) {
		return EventChoiceResult{}, ErrEventExpired
	}

	choice, ok := findEventChoice(event.Choices, choiceKey)
	if !ok {
		if choiceKey == "" {
			return EventChoiceResult{}, ErrInvalidEventChoice
		}
		return EventChoiceResult{}, ErrEventChoiceNotAvailable
	}

	effects, err := decodeEventEffects(choice.EffectsJSON)
	if err != nil {
		return EventChoiceResult{}, err
	}

	var resources *ResourcesResult
	if effects.ResourceDelta != nil {
		delta := resourceDeltaValues(*effects.ResourceDelta)
		result, err := s.resources.ApplyDeltaForKingdom(ctx, kingdom.ID, delta)
		if err != nil {
			return EventChoiceResult{}, err
		}
		resources = &result
	}

	var army *ArmyView
	if len(effects.UnitDelta) > 0 {
		result, err := s.army.ApplyUnitDeltaForKingdom(ctx, kingdom.ID, effects.UnitDelta)
		if err != nil {
			return EventChoiceResult{}, err
		}
		army = &result
	}

	updatedKingdom := kingdom
	if effects.KingdomDelta != nil {
		updated, err := s.kingdoms.ApplyReputationDelta(ctx, kingdom.ID, effects.KingdomDelta.Dread, effects.KingdomDelta.Honor)
		if err != nil {
			return EventChoiceResult{}, err
		}
		updatedKingdom = updated
	}

	if effects.PatronFavorDelta != 0 {
		if _, err := s.patrons.AdjustFavor(ctx, kingdom.ID, effects.PatronFavorDelta); err != nil && !errors.Is(err, repository.ErrPatronRelationNotFound) {
			return EventChoiceResult{}, err
		}
	}

	eventResult := EventResult{
		Title:          choice.ResultTitle,
		Body:           choice.ResultBody,
		AppliedEffects: effects,
	}
	resultJSON, err := json.Marshal(eventResult)
	if err != nil {
		return EventChoiceResult{}, err
	}
	resolved, err := s.events.MarkResolved(ctx, event.ID, kingdom.ID, choice.Key, s.now(), resultJSON)
	if errors.Is(err, repository.ErrEventNotFound) {
		return EventChoiceResult{}, ErrEventAlreadyResolved
	}
	if err != nil {
		return EventChoiceResult{}, err
	}
	if err := s.createReport(ctx, kingdom.ID, event, choice, effects); err != nil {
		return EventChoiceResult{}, err
	}
	view, err := s.view(resolved)
	if err != nil {
		return EventChoiceResult{}, err
	}
	return EventChoiceResult{
		Event:     view,
		Resources: resources,
		Army:      army,
		Kingdom:   updatedKingdom,
	}, nil
}

func (s *EventService) refresh(ctx context.Context, kingdom domain.Kingdom) error {
	now := s.now()
	if err := s.events.ExpireStale(ctx, kingdom.ID, now); err != nil {
		return err
	}
	active, err := s.events.ListActiveByKingdomID(ctx, kingdom.ID)
	if err != nil {
		return err
	}
	if len(active) >= maxActiveEvents {
		return nil
	}

	activeKeys := map[string]bool{}
	for _, event := range active {
		activeKeys[event.GameEvent.Key] = true
	}

	catalog, err := s.events.ListActiveGameEvents(ctx)
	if err != nil {
		return err
	}
	for _, event := range catalog {
		if len(active) >= maxActiveEvents {
			break
		}
		if activeKeys[event.Key] || !eligibleEvent(event, kingdom) {
			continue
		}
		since := now.Add(-time.Duration(event.CooldownSeconds) * time.Second)
		recent, err := s.events.HasRecentByEventKey(ctx, kingdom.ID, event.Key, since)
		if err != nil {
			return err
		}
		if recent {
			continue
		}
		instance, err := s.events.CreateKingdomEvent(ctx, kingdom.ID, event.ID, now, now.Add(time.Duration(event.ExpiresAfterSeconds)*time.Second))
		if err != nil {
			return err
		}
		active = append(active, instance)
		activeKeys[event.Key] = true
	}
	return nil
}

func (s *EventService) kingdomForUser(ctx context.Context, userID string) (domain.Kingdom, error) {
	kingdom, err := s.kingdoms.FindByUserID(ctx, userID)
	if errors.Is(err, repository.ErrKingdomNotFound) {
		return domain.Kingdom{}, ErrEventKingdomNotFound
	}
	if err != nil {
		return domain.Kingdom{}, err
	}
	return kingdom, nil
}

func (s *EventService) views(events []domain.KingdomEvent) ([]KingdomEventView, error) {
	views := make([]KingdomEventView, 0, len(events))
	for _, event := range events {
		view, err := s.view(event)
		if err != nil {
			return nil, err
		}
		views = append(views, view)
	}
	return views, nil
}

func (s *EventService) view(event domain.KingdomEvent) (KingdomEventView, error) {
	var result *EventResult
	if len(event.ResultJSON) > 0 {
		decoded := EventResult{}
		if err := json.Unmarshal(event.ResultJSON, &decoded); err != nil {
			return KingdomEventView{}, err
		}
		result = &decoded
	}
	choices := []EventChoiceView{}
	if event.Status == "active" {
		choices = make([]EventChoiceView, 0, len(event.Choices))
		for _, choice := range event.Choices {
			choices = append(choices, EventChoiceView{Choice: choice})
		}
	}
	return KingdomEventView{
		Event:   event,
		Choices: choices,
		Result:  result,
	}, nil
}

func (s *EventService) createReport(ctx context.Context, kingdomID string, event domain.KingdomEvent, choice domain.EventChoice, effects EventAppliedEffects) error {
	rewardsJSON, _ := json.Marshal(positiveResourceEffects(effects))
	lossesJSON, _ := json.Marshal(negativeEffects(effects))
	phases := []gameconfig.ReportPhase{
		{Title: "Событие", Body: event.GameEvent.Body},
		{Title: "Выбор", Body: choice.Label},
		{Title: "Последствия", Body: choice.ResultBody},
	}
	phasesJSON, _ := json.Marshal(phases)
	_, err := s.reports.CreateReport(ctx, kingdomID, nil, "event", choice.ResultTitle, choice.ResultBody, "success", rewardsJSON, lossesJSON, phasesJSON)
	return err
}

func eligibleEvent(event domain.GameEvent, kingdom domain.Kingdom) bool {
	conditions := eventConditions{}
	if len(event.ConditionsJSON) > 0 {
		if err := json.Unmarshal(event.ConditionsJSON, &conditions); err != nil {
			return false
		}
	}
	switch conditions.RequiresPatron {
	case "":
		return true
	case "any":
		return kingdom.Patron != nil
	case "none":
		return kingdom.Patron == nil
	case "independent", "empire_of_dusk", "old_pact":
		return kingdom.Patron != nil && *kingdom.Patron == conditions.RequiresPatron
	default:
		return false
	}
}

func findEventChoice(choices []domain.EventChoice, key string) (domain.EventChoice, bool) {
	for _, choice := range choices {
		if choice.Key == key {
			return choice, true
		}
	}
	return domain.EventChoice{}, false
}

func decodeEventEffects(data []byte) (EventAppliedEffects, error) {
	if len(data) == 0 {
		return EventAppliedEffects{}, nil
	}
	var effects EventAppliedEffects
	if err := json.Unmarshal(data, &effects); err != nil {
		return EventAppliedEffects{}, err
	}
	return effects, nil
}

func resourceDeltaValues(delta EventResourceDelta) gameconfig.ResourceValues {
	return gameconfig.ResourceValues{
		Gold:       delta.Gold,
		Food:       delta.Food,
		Wood:       delta.Wood,
		Stone:      delta.Stone,
		Population: delta.Population,
	}
}

func positiveResourceEffects(effects EventAppliedEffects) gameconfig.ResourceValues {
	if effects.ResourceDelta == nil {
		return gameconfig.ResourceValues{}
	}
	delta := effects.ResourceDelta
	return gameconfig.ResourceValues{
		Gold:       maxInt64(0, delta.Gold),
		Food:       maxInt64(0, delta.Food),
		Wood:       maxInt64(0, delta.Wood),
		Stone:      maxInt64(0, delta.Stone),
		Population: maxInt64(0, delta.Population),
	}
}

func negativeEffects(effects EventAppliedEffects) map[string]int64 {
	losses := map[string]int64{}
	if effects.ResourceDelta != nil {
		addNegativeEffect(losses, "gold", effects.ResourceDelta.Gold)
		addNegativeEffect(losses, "food", effects.ResourceDelta.Food)
		addNegativeEffect(losses, "wood", effects.ResourceDelta.Wood)
		addNegativeEffect(losses, "stone", effects.ResourceDelta.Stone)
		addNegativeEffect(losses, "population", effects.ResourceDelta.Population)
	}
	for unitType, delta := range effects.UnitDelta {
		addNegativeEffect(losses, unitType, delta)
	}
	return losses
}

func addNegativeEffect(losses map[string]int64, key string, delta int64) {
	if delta < 0 {
		losses[key] = -delta
	}
}
