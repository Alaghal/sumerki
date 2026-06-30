package service

import (
	"testing"

	"sumerki/backend/internal/domain"
	"sumerki/backend/internal/gameconfig"
)

func TestRaidResultThresholds(t *testing.T) {
	tests := []struct {
		name          string
		attackerScore int64
		defenderScore int64
		want          string
	}{
		{name: "attacker success", attackerScore: 115, defenderScore: 100, want: "attacker_success"},
		{name: "defender success", attackerScore: 100, defenderScore: 115, want: "defender_success"},
		{name: "stalemate", attackerScore: 110, defenderScore: 100, want: "bloody_stalemate"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := raidResult(tt.attackerScore, tt.defenderScore)
			if got != tt.want {
				t.Fatalf("raidResult() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRaidLootNeverDropsBelowProtectedMinimum(t *testing.T) {
	resources := domain.Resources{
		Gold:       105,
		Food:       300,
		Wood:       100,
		Stone:      80,
		Population: 999,
	}

	loot := raidLoot(resources, gameconfig.ProtectedRaidResources, 10)
	if loot.Gold != 5 {
		t.Fatalf("Gold loot = %d, want 5", loot.Gold)
	}
	if loot.Food != 30 {
		t.Fatalf("Food loot = %d, want 30", loot.Food)
	}
	if loot.Wood != 0 {
		t.Fatalf("Wood loot = %d, want 0", loot.Wood)
	}
	if loot.Stone != 5 {
		t.Fatalf("Stone loot = %d, want 5", loot.Stone)
	}
	if loot.Population != 0 {
		t.Fatalf("Population loot = %d, want 0", loot.Population)
	}
}

func TestNormalizeRaidUnitsRejectsBadInput(t *testing.T) {
	if _, err := normalizeRaidUnits([]StartRaidUnit{{UnitType: "bad", Amount: 1}}); err != ErrInvalidRaidUnitType {
		t.Fatalf("invalid type error = %v, want ErrInvalidRaidUnitType", err)
	}
	if _, err := normalizeRaidUnits([]StartRaidUnit{{UnitType: "militia", Amount: -1}}); err != ErrInvalidRaidUnitAmount {
		t.Fatalf("negative amount error = %v, want ErrInvalidRaidUnitAmount", err)
	}
	if _, err := normalizeRaidUnits([]StartRaidUnit{{UnitType: "militia", Amount: 0}}); err != ErrInvalidRaidUnitAmount {
		t.Fatalf("empty amount error = %v, want ErrInvalidRaidUnitAmount", err)
	}
}
