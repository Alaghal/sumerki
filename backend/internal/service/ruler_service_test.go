package service

import (
	"context"
	"errors"
	"testing"

	"sumerki/backend/internal/domain"
	"sumerki/backend/internal/repository"
)

type fakeRulerRepository struct {
	rulersByKingdomID map[string]domain.Ruler
}

func newFakeRulerRepository() *fakeRulerRepository {
	return &fakeRulerRepository{rulersByKingdomID: map[string]domain.Ruler{}}
}

func (r *fakeRulerRepository) Create(_ context.Context, ruler domain.Ruler) (domain.Ruler, error) {
	if _, exists := r.rulersByKingdomID[ruler.KingdomID]; exists {
		return domain.Ruler{}, repository.ErrRulerExists
	}

	ruler.ID = "ruler-1"
	r.rulersByKingdomID[ruler.KingdomID] = ruler
	return ruler, nil
}

func (r *fakeRulerRepository) FindByKingdomID(_ context.Context, kingdomID string) (domain.Ruler, error) {
	ruler, ok := r.rulersByKingdomID[kingdomID]
	if !ok {
		return domain.Ruler{}, repository.ErrRulerNotFound
	}

	return ruler, nil
}

func TestRulerCurrentCreatesMissingRuler(t *testing.T) {
	kingdoms := newFakeKingdomRepository()
	kingdom, err := kingdoms.Create(context.Background(), "user-1", "Blackwater", CultureNorthernPrincipality)
	if err != nil {
		t.Fatalf("create kingdom fixture failed: %v", err)
	}

	service := NewRulerService(kingdoms, newFakeRulerRepository())

	ruler, err := service.Current(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("current ruler failed: %v", err)
	}

	if ruler.KingdomID != kingdom.ID {
		t.Fatalf("expected kingdom id %q, got %q", kingdom.ID, ruler.KingdomID)
	}
	if ruler.Culture != kingdom.Culture {
		t.Fatalf("expected culture %q, got %q", kingdom.Culture, ruler.Culture)
	}
	if ruler.HealthStatus != HealthStatusHealthy {
		t.Fatalf("expected healthy ruler, got %q", ruler.HealthStatus)
	}
	assertRulerRanges(t, ruler)
}

func TestRulerCurrentReturnsKingdomNotFound(t *testing.T) {
	service := NewRulerService(newFakeKingdomRepository(), newFakeRulerRepository())

	_, err := service.Current(context.Background(), "user-1")
	if !errors.Is(err, ErrRulerKingdomNotFound) {
		t.Fatalf("expected kingdom not found, got %v", err)
	}
}

func TestGenerateForKingdomReturnsExistingRulerOnConflict(t *testing.T) {
	kingdom := domain.Kingdom{
		ID:      "kingdom-1",
		Culture: CultureFreePosad,
	}
	rulers := newFakeRulerRepository()
	service := NewRulerService(newFakeKingdomRepository(), rulers)

	first, err := service.GenerateForKingdom(context.Background(), kingdom)
	if err != nil {
		t.Fatalf("generate first ruler failed: %v", err)
	}

	second, err := service.GenerateForKingdom(context.Background(), kingdom)
	if err != nil {
		t.Fatalf("generate second ruler failed: %v", err)
	}

	if second.ID != first.ID {
		t.Fatalf("expected existing ruler %q, got %q", first.ID, second.ID)
	}
}

func assertRulerRanges(t *testing.T, ruler domain.Ruler) {
	t.Helper()

	if ruler.Age < 25 || ruler.Age > 60 {
		t.Fatalf("ruler age out of generated range: %d", ruler.Age)
	}

	stats := map[string]int{
		"authority": ruler.Authority,
		"courage":   ruler.Courage,
		"cunning":   ruler.Cunning,
		"honor":     ruler.Honor,
		"cruelty":   ruler.Cruelty,
		"ambition":  ruler.Ambition,
		"paranoia":  ruler.Paranoia,
	}

	for name, value := range stats {
		if value < 30 || value > 80 {
			t.Fatalf("%s out of generated range: %d", name, value)
		}
	}
}
