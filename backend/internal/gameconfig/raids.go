package gameconfig

import "time"

const (
	RaidDurationSeconds       = 120
	MinimumRaidUnits          = int64(3)
	MaximumRaidUnits          = int64(100)
	MaxRaidLootPercent        = int64(10)
	StalemateLootPercent      = int64(3)
	SameTargetRaidCooldown    = 12 * time.Hour
	DefenderProtectionWindow  = 24 * time.Hour
	DefenderProtectionMaxHits = 3
	NewKingdomProtectionAge   = 24 * time.Hour
)

var ProtectedRaidResources = ResourceValues{
	Gold:  100,
	Food:  100,
	Wood:  100,
	Stone: 75,
}
