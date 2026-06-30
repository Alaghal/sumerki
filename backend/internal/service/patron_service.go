package service

import (
	"context"
	"errors"

	"sumerki/backend/internal/domain"
	"sumerki/backend/internal/gameconfig"
	"sumerki/backend/internal/repository"
)

var (
	ErrPatronKingdomNotFound = errors.New("kingdom not found")
	ErrInvalidPatron         = errors.New("invalid patron")
)

type PatronRepository interface {
	FindByKingdomID(ctx context.Context, kingdomID string) (domain.PatronRelation, error)
	UpsertForKingdom(ctx context.Context, kingdomID string, patron string) (domain.PatronRelation, error)
	BreakForKingdom(ctx context.Context, kingdomID string) error
	BackfillForKingdom(ctx context.Context, kingdomID string, patron string) (domain.PatronRelation, error)
}

type PatronKingdomRepository interface {
	FindByUserID(ctx context.Context, userID string) (domain.Kingdom, error)
	UpdatePatronByID(ctx context.Context, kingdomID string, patron *string) (domain.Kingdom, error)
}

type PatronRelationView struct {
	Relation       domain.PatronRelation
	Label          string
	CurrentEffects []string
	FutureEffects  []string
}

type PatronStatus struct {
	Patron           *PatronRelationView
	AvailablePatrons []string
}

type PatronJoinResult struct {
	Patron  PatronRelationView
	Kingdom domain.Kingdom
}

type PatronBreakResult struct {
	Kingdom domain.Kingdom
}

type PatronService struct {
	kingdoms PatronKingdomRepository
	patrons  PatronRepository
	pressure PatronPressureLifecycle
}

func NewPatronService(kingdoms PatronKingdomRepository, patrons PatronRepository) *PatronService {
	return &PatronService{kingdoms: kingdoms, patrons: patrons}
}

type PatronPressureLifecycle interface {
	EnsureForJoin(ctx context.Context, kingdomID string, patron string, samePatron bool) error
	ClearForBreak(ctx context.Context, kingdomID string) error
	ResolveForKingdom(ctx context.Context, kingdom domain.Kingdom) error
}

func (s *PatronService) SetPressureLifecycle(pressure PatronPressureLifecycle) {
	s.pressure = pressure
}

func (s *PatronService) Options() []gameconfig.PatronConfig {
	options := make([]gameconfig.PatronConfig, 0, len(gameconfig.PatronOrder))
	for _, key := range gameconfig.PatronOrder {
		options = append(options, gameconfig.Patrons[key])
	}
	return options
}

func (s *PatronService) Current(ctx context.Context, userID string) (PatronStatus, error) {
	kingdom, err := s.kingdomForUser(ctx, userID)
	if err != nil {
		return PatronStatus{}, err
	}

	relation, err := s.patrons.FindByKingdomID(ctx, kingdom.ID)
	if errors.Is(err, repository.ErrPatronRelationNotFound) {
		if kingdom.Patron == nil {
			return PatronStatus{AvailablePatrons: availablePatronKeys()}, nil
		}
		relation, err = s.patrons.BackfillForKingdom(ctx, kingdom.ID, *kingdom.Patron)
	}
	if err != nil {
		return PatronStatus{}, err
	}
	if s.pressure != nil {
		if err := s.pressure.ResolveForKingdom(ctx, kingdom); err != nil {
			return PatronStatus{}, err
		}
	}

	view := patronRelationView(relation)
	return PatronStatus{
		Patron:           &view,
		AvailablePatrons: availablePatronKeys(),
	}, nil
}

func (s *PatronService) Join(ctx context.Context, userID string, patron string) (PatronJoinResult, error) {
	if !gameconfig.IsPatronKey(patron) {
		return PatronJoinResult{}, ErrInvalidPatron
	}

	kingdom, err := s.kingdomForUser(ctx, userID)
	if err != nil {
		return PatronJoinResult{}, err
	}

	relation, err := s.patrons.UpsertForKingdom(ctx, kingdom.ID, patron)
	if err != nil {
		return PatronJoinResult{}, err
	}
	samePatron := kingdom.Patron != nil && *kingdom.Patron == patron
	updatedKingdom, err := s.kingdoms.UpdatePatronByID(ctx, kingdom.ID, &patron)
	if err != nil {
		return PatronJoinResult{}, err
	}
	if s.pressure != nil {
		if err := s.pressure.EnsureForJoin(ctx, kingdom.ID, patron, samePatron); err != nil {
			return PatronJoinResult{}, err
		}
	}

	return PatronJoinResult{
		Patron:  patronRelationView(relation),
		Kingdom: updatedKingdom,
	}, nil
}

func (s *PatronService) Break(ctx context.Context, userID string) (PatronBreakResult, error) {
	kingdom, err := s.kingdomForUser(ctx, userID)
	if err != nil {
		return PatronBreakResult{}, err
	}
	if err := s.patrons.BreakForKingdom(ctx, kingdom.ID); err != nil {
		return PatronBreakResult{}, err
	}
	if s.pressure != nil {
		if err := s.pressure.ClearForBreak(ctx, kingdom.ID); err != nil {
			return PatronBreakResult{}, err
		}
	}

	updatedKingdom, err := s.kingdoms.UpdatePatronByID(ctx, kingdom.ID, nil)
	if err != nil {
		return PatronBreakResult{}, err
	}
	return PatronBreakResult{Kingdom: updatedKingdom}, nil
}

func (s *PatronService) kingdomForUser(ctx context.Context, userID string) (domain.Kingdom, error) {
	kingdom, err := s.kingdoms.FindByUserID(ctx, userID)
	if errors.Is(err, repository.ErrKingdomNotFound) {
		return domain.Kingdom{}, ErrPatronKingdomNotFound
	}
	if err != nil {
		return domain.Kingdom{}, err
	}
	return kingdom, nil
}

func availablePatronKeys() []string {
	keys := make([]string, 0, len(gameconfig.PatronOrder))
	keys = append(keys, gameconfig.PatronOrder...)
	return keys
}

func patronRelationView(relation domain.PatronRelation) PatronRelationView {
	config := gameconfig.Patrons[relation.Patron]
	return PatronRelationView{
		Relation:       relation,
		Label:          config.Label,
		CurrentEffects: config.CurrentEffects,
		FutureEffects:  config.FutureEffects,
	}
}
