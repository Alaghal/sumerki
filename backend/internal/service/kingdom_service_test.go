package service

import (
	"context"
	"errors"
	"strings"
	"testing"

	"sumerki/backend/internal/domain"
	"sumerki/backend/internal/repository"
)

type fakeKingdomRepository struct {
	kingdomByUserID map[string]domain.Kingdom
}

type fakeRulerGenerator struct {
	createdForKingdomID string
}

func (g *fakeRulerGenerator) GenerateForKingdom(_ context.Context, kingdom domain.Kingdom) (domain.Ruler, error) {
	g.createdForKingdomID = kingdom.ID
	return domain.Ruler{KingdomID: kingdom.ID}, nil
}

func newFakeKingdomRepository() *fakeKingdomRepository {
	return &fakeKingdomRepository{kingdomByUserID: map[string]domain.Kingdom{}}
}

func (r *fakeKingdomRepository) Create(_ context.Context, userID string, name string, culture string) (domain.Kingdom, error) {
	if _, exists := r.kingdomByUserID[userID]; exists {
		return domain.Kingdom{}, repository.ErrKingdomExists
	}

	kingdom := domain.Kingdom{
		ID:      "kingdom-1",
		UserID:  userID,
		Name:    name,
		Culture: culture,
	}
	r.kingdomByUserID[userID] = kingdom

	return kingdom, nil
}

func (r *fakeKingdomRepository) FindByUserID(_ context.Context, userID string) (domain.Kingdom, error) {
	kingdom, ok := r.kingdomByUserID[userID]
	if !ok {
		return domain.Kingdom{}, repository.ErrKingdomNotFound
	}

	return kingdom, nil
}

func TestCreateKingdomTrimsName(t *testing.T) {
	service := NewKingdomService(newFakeKingdomRepository())

	kingdom, err := service.Create(context.Background(), "user-1", "  Blackwater  ", CultureNorthernPrincipality)
	if err != nil {
		t.Fatalf("create kingdom failed: %v", err)
	}

	if kingdom.Name != "Blackwater" {
		t.Fatalf("expected trimmed name, got %q", kingdom.Name)
	}
}

func TestCreateKingdomRejectsInvalidCulture(t *testing.T) {
	service := NewKingdomService(newFakeKingdomRepository())

	_, err := service.Create(context.Background(), "user-1", "Blackwater", "bad_culture")
	if !errors.Is(err, ErrInvalidCulture) {
		t.Fatalf("expected invalid culture, got %v", err)
	}
}

func TestCreateKingdomRejectsShortName(t *testing.T) {
	service := NewKingdomService(newFakeKingdomRepository())

	_, err := service.Create(context.Background(), "user-1", "ab", CultureNorthernPrincipality)
	if !errors.Is(err, ErrKingdomNameTooShort) {
		t.Fatalf("expected short name error, got %v", err)
	}
}

func TestCreateKingdomRejectsLongName(t *testing.T) {
	service := NewKingdomService(newFakeKingdomRepository())

	_, err := service.Create(context.Background(), "user-1", strings.Repeat("a", 33), CultureNorthernPrincipality)
	if !errors.Is(err, ErrKingdomNameTooLong) {
		t.Fatalf("expected long name error, got %v", err)
	}
}

func TestCreateKingdomRejectsSecondKingdom(t *testing.T) {
	service := NewKingdomService(newFakeKingdomRepository())

	if _, err := service.Create(context.Background(), "user-1", "Blackwater", CultureNorthernPrincipality); err != nil {
		t.Fatalf("create kingdom failed: %v", err)
	}

	_, err := service.Create(context.Background(), "user-1", "Whitewater", CultureFreePosad)
	if !errors.Is(err, ErrKingdomAlreadyExists) {
		t.Fatalf("expected kingdom already exists, got %v", err)
	}
}

func TestCreateKingdomCreatesRuler(t *testing.T) {
	rulers := &fakeRulerGenerator{}
	service := NewKingdomService(newFakeKingdomRepository(), rulers)

	kingdom, err := service.Create(context.Background(), "user-1", "Blackwater", CultureNorthernPrincipality)
	if err != nil {
		t.Fatalf("create kingdom failed: %v", err)
	}

	if rulers.createdForKingdomID != kingdom.ID {
		t.Fatalf("expected ruler for kingdom %q, got %q", kingdom.ID, rulers.createdForKingdomID)
	}
}

func TestCurrentKingdomReturnsNilBeforeCreation(t *testing.T) {
	service := NewKingdomService(newFakeKingdomRepository())

	kingdom, err := service.Current(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("current kingdom failed: %v", err)
	}
	if kingdom != nil {
		t.Fatalf("expected nil kingdom, got %#v", kingdom)
	}
}
