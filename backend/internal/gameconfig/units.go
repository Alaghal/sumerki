package gameconfig

const (
	MaxTrainingAmount = 50
)

type UnitStats struct {
	Attack  int
	Defense int
	Speed   int
	Supply  int
}

type UnitConfig struct {
	Type                  string
	Label                 string
	Role                  string
	Stats                 UnitStats
	Cost                  ResourceValues
	SecondsPerUnit        int
	RequiredBarracksLevel int
}

var UnitOrder = []string{
	"militia",
	"spearmen",
	"archers",
	"cavalry",
	"scouts",
}

var InitialUnitAmounts = map[string]int64{
	"militia":  10,
	"spearmen": 0,
	"archers":  0,
	"cavalry":  0,
	"scouts":   2,
}

var Units = map[string]UnitConfig{
	"militia": {
		Type:  "militia",
		Label: "Ополчение",
		Role:  "cheap defensive bodies",
		Stats: UnitStats{
			Attack:  2,
			Defense: 3,
			Speed:   2,
			Supply:  1,
		},
		Cost: ResourceValues{
			Gold:       15,
			Food:       10,
			Population: 1,
		},
		SecondsPerUnit: 5,
	},
	"spearmen": {
		Type:  "spearmen",
		Label: "Копейщики",
		Role:  "stable infantry, good defense",
		Stats: UnitStats{
			Attack:  4,
			Defense: 6,
			Speed:   2,
			Supply:  2,
		},
		Cost: ResourceValues{
			Gold:       30,
			Food:       15,
			Wood:       10,
			Population: 1,
		},
		SecondsPerUnit:        8,
		RequiredBarracksLevel: 1,
	},
	"archers": {
		Type:  "archers",
		Label: "Лучники",
		Role:  "ranged support",
		Stats: UnitStats{
			Attack:  6,
			Defense: 3,
			Speed:   2,
			Supply:  2,
		},
		Cost: ResourceValues{
			Gold:       35,
			Food:       15,
			Wood:       20,
			Population: 1,
		},
		SecondsPerUnit:        8,
		RequiredBarracksLevel: 1,
	},
	"cavalry": {
		Type:  "cavalry",
		Label: "Конница",
		Role:  "fast attack and future raids",
		Stats: UnitStats{
			Attack:  10,
			Defense: 5,
			Speed:   5,
			Supply:  4,
		},
		Cost: ResourceValues{
			Gold:       80,
			Food:       35,
			Wood:       20,
			Population: 1,
		},
		SecondsPerUnit:        12,
		RequiredBarracksLevel: 2,
	},
	"scouts": {
		Type:  "scouts",
		Label: "Разведчики",
		Role:  "scouting and future mission safety",
		Stats: UnitStats{
			Attack:  2,
			Defense: 2,
			Speed:   6,
			Supply:  1,
		},
		Cost: ResourceValues{
			Gold:       25,
			Food:       15,
			Wood:       10,
			Population: 1,
		},
		SecondsPerUnit: 6,
	},
}

func IsUnitType(unitType string) bool {
	_, ok := Units[unitType]
	return ok
}

func UnitTrainingCost(unitType string, amount int64) ResourceValues {
	cfg := Units[unitType]
	return ResourceValues{
		Gold:       cfg.Cost.Gold * amount,
		Food:       cfg.Cost.Food * amount,
		Wood:       cfg.Cost.Wood * amount,
		Stone:      cfg.Cost.Stone * amount,
		Population: cfg.Cost.Population * amount,
	}
}

func UnitTrainingDurationSeconds(unitType string, amount int64) int {
	return Units[unitType].SecondsPerUnit * int(amount)
}
