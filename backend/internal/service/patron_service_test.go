package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"sumerki/backend/internal/domain"
	"sumerki/backend/internal/repository"
)

func TestPatronOptionsReturnsAllPatrons(t *testing.T) {
	service := NewPatronService(newFakePatronKingdomRepository(), newFakePatronRepository())

	options := service.Options()
	if len(options) != 3 {
		t.Fatalf("len(options) = %d, want 3", len(options))
	}
	if options[0].Key != "independent" || options[1].Key != "empire_of_dusk" || options[2].Key != "old_pact" {
		t.Fatalf("unexpected patron order: %#v", options)
	}
}

func TestPatronCurrentWithoutKingdom(t *testing.T) {
	service := NewPatronService(newFakePatronKingdomRepository(), newFakePatronRepository())

	_, err := service.Current(context.Background(), "user-1")
	if !errors.Is(err, ErrPatronKingdomNotFound) {
		t.Fatalf("Current error = %v, want ErrPatronKingdomNotFound", err)
	}
}

func TestPatronJoinInvalid(t *testing.T) {
	kingdoms := newFakePatronKingdomRepository()
	kingdoms.kingdomByUserID["user-1"] = domain.Kingdom{ID: "kingdom-1", UserID: "user-1"}
	service := NewPatronService(kingdoms, newFakePatronRepository())

	_, err := service.Join(context.Background(), "user-1", "bad_patron")
	if !errors.Is(err, ErrInvalidPatron) {
		t.Fatalf("Join error = %v, want ErrInvalidPatron", err)
	}
}

func TestPatronJoinCreatesAndSwitchesRelation(t *testing.T) {
	ctx := context.Background()
	kingdoms := newFakePatronKingdomRepository()
	kingdoms.kingdomByUserID["user-1"] = domain.Kingdom{ID: "kingdom-1", UserID: "user-1"}
	patrons := newFakePatronRepository()
	service := NewPatronService(kingdoms, patrons)

	first, err := service.Join(ctx, "user-1", "old_pact")
	if err != nil {
		t.Fatalf("Join old_pact returned error: %v", err)
	}
	if first.Patron.Relation.Patron != "old_pact" || first.Kingdom.Patron == nil || *first.Kingdom.Patron != "old_pact" {
		t.Fatalf("Join old_pact did not update relation and kingdom: %#v", first)
	}

	second, err := service.Join(ctx, "user-1", "old_pact")
	if err != nil {
		t.Fatalf("idempotent Join returned error: %v", err)
	}
	if second.Patron.Relation.Patron != "old_pact" {
		t.Fatalf("idempotent Join patron = %q, want old_pact", second.Patron.Relation.Patron)
	}

	switched, err := service.Join(ctx, "user-1", "empire_of_dusk")
	if err != nil {
		t.Fatalf("switch Join returned error: %v", err)
	}
	if switched.Patron.Relation.Patron != "empire_of_dusk" || switched.Kingdom.Patron == nil || *switched.Kingdom.Patron != "empire_of_dusk" {
		t.Fatalf("Join switch did not update relation and kingdom: %#v", switched)
	}
	if len(patrons.relationByKingdomID) != 1 {
		t.Fatalf("len(relationByKingdomID) = %d, want 1", len(patrons.relationByKingdomID))
	}
}

func TestPatronBreakClearsKingdomAndIsIdempotent(t *testing.T) {
	ctx := context.Background()
	kingdoms := newFakePatronKingdomRepository()
	patron := "old_pact"
	kingdoms.kingdomByUserID["user-1"] = domain.Kingdom{ID: "kingdom-1", UserID: "user-1", Patron: &patron}
	service := NewPatronService(kingdoms, newFakePatronRepository())

	if _, err := service.Join(ctx, "user-1", "old_pact"); err != nil {
		t.Fatalf("Join returned error: %v", err)
	}
	first, err := service.Break(ctx, "user-1")
	if err != nil {
		t.Fatalf("Break returned error: %v", err)
	}
	second, err := service.Break(ctx, "user-1")
	if err != nil {
		t.Fatalf("idempotent Break returned error: %v", err)
	}
	if first.Kingdom.Patron != nil || second.Kingdom.Patron != nil {
		t.Fatalf("Break did not clear kingdom patron: first=%#v second=%#v", first.Kingdom.Patron, second.Kingdom.Patron)
	}
}

func TestPatronCurrentBackfillsKingdomPatron(t *testing.T) {
	ctx := context.Background()
	kingdoms := newFakePatronKingdomRepository()
	patron := "independent"
	kingdoms.kingdomByUserID["user-1"] = domain.Kingdom{ID: "kingdom-1", UserID: "user-1", Patron: &patron}
	service := NewPatronService(kingdoms, newFakePatronRepository())

	status, err := service.Current(ctx, "user-1")
	if err != nil {
		t.Fatalf("Current returned error: %v", err)
	}
	if status.Patron == nil || status.Patron.Relation.Patron != "independent" {
		t.Fatalf("Current did not backfill patron relation: %#v", status.Patron)
	}
}

func TestPatronUserCannotReadAnotherUserRelation(t *testing.T) {
	ctx := context.Background()
	kingdoms := newFakePatronKingdomRepository()
	kingdoms.kingdomByUserID["user-1"] = domain.Kingdom{ID: "kingdom-1", UserID: "user-1"}
	kingdoms.kingdomByUserID["user-2"] = domain.Kingdom{ID: "kingdom-2", UserID: "user-2"}
	patrons := newFakePatronRepository()
	service := NewPatronService(kingdoms, patrons)

	if _, err := service.Join(ctx, "user-2", "old_pact"); err != nil {
		t.Fatalf("Join user-2 returned error: %v", err)
	}
	status, err := service.Current(ctx, "user-1")
	if err != nil {
		t.Fatalf("Current user-1 returned error: %v", err)
	}
	if status.Patron != nil {
		t.Fatalf("user-1 saw another user's patron relation: %#v", status.Patron)
	}
}

type fakePatronKingdomRepository struct {
	kingdomByUserID map[string]domain.Kingdom
}

func newFakePatronKingdomRepository() *fakePatronKingdomRepository {
	return &fakePatronKingdomRepository{kingdomByUserID: map[string]domain.Kingdom{}}
}

func (r *fakePatronKingdomRepository) FindByUserID(ctx context.Context, userID string) (domain.Kingdom, error) {
	kingdom, ok := r.kingdomByUserID[userID]
	if !ok {
		return domain.Kingdom{}, repository.ErrKingdomNotFound
	}
	return kingdom, nil
}

func (r *fakePatronKingdomRepository) UpdatePatronByID(ctx context.Context, kingdomID string, patron *string) (domain.Kingdom, error) {
	for userID, kingdom := range r.kingdomByUserID {
		if kingdom.ID == kingdomID {
			kingdom.Patron = patron
			r.kingdomByUserID[userID] = kingdom
			return kingdom, nil
		}
	}
	return domain.Kingdom{}, repository.ErrKingdomNotFound
}

type fakePatronRepository struct {
	relationByKingdomID map[string]domain.PatronRelation
}

func newFakePatronRepository() *fakePatronRepository {
	return &fakePatronRepository{relationByKingdomID: map[string]domain.PatronRelation{}}
}

func (r *fakePatronRepository) FindByKingdomID(ctx context.Context, kingdomID string) (domain.PatronRelation, error) {
	relation, ok := r.relationByKingdomID[kingdomID]
	if !ok || relation.LeftAt != nil {
		return domain.PatronRelation{}, repository.ErrPatronRelationNotFound
	}
	return relation, nil
}

func (r *fakePatronRepository) UpsertForKingdom(ctx context.Context, kingdomID string, patron string) (domain.PatronRelation, error) {
	relation := r.relationByKingdomID[kingdomID]
	if relation.ID == "" {
		relation.ID = kingdomID + "-patron"
		relation.KingdomID = kingdomID
		relation.CreatedAt = time.Now()
	}
	relation.Patron = patron
	relation.Favor = 0
	relation.Standing = "neutral"
	relation.JoinedAt = time.Now()
	relation.LeftAt = nil
	relation.UpdatedAt = time.Now()
	r.relationByKingdomID[kingdomID] = relation
	return relation, nil
}

func (r *fakePatronRepository) BreakForKingdom(ctx context.Context, kingdomID string) error {
	relation, ok := r.relationByKingdomID[kingdomID]
	if !ok {
		return nil
	}
	now := time.Now()
	relation.LeftAt = &now
	relation.UpdatedAt = now
	r.relationByKingdomID[kingdomID] = relation
	return nil
}

func (r *fakePatronRepository) BackfillForKingdom(ctx context.Context, kingdomID string, patron string) (domain.PatronRelation, error) {
	return r.UpsertForKingdom(ctx, kingdomID, patron)
}
