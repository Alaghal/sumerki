package gameconfig

import "time"

const (
	PatronPressureMaxIntervals = 6
	EmpireTributeInterval      = time.Hour
	EmpireTributeGoldPercent   = int64(10)
	EmpireTributeFoodPercent   = int64(5)
	EmpireMinimumGoldDue       = int64(20)
	EmpireMinimumFoodDue       = int64(10)
	OldPactInterval            = 2 * time.Hour
	OldPactFoodPercent         = int64(5)
	OldPactPressureCap         = 20
	PatronDelayDuration        = 2 * time.Hour
)

var PatronPressureProtectedMinimums = ResourceValues{
	Gold:  150,
	Food:  150,
	Wood:  100,
	Stone: 75,
}
