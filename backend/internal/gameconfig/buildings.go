package gameconfig

const MaxBuildingLevel = 5

type BuildingConfig struct {
	Type     string
	Label    string
	Purpose  string
	BaseCost ResourceValues
	Effects  []string
}

var BuildingOrder = []string{
	"town_hall",
	"farm",
	"lumberyard",
	"quarry",
	"market",
	"barracks",
	"walls",
	"shrine",
}

var InitialBuildingLevels = map[string]int{
	"town_hall":  1,
	"farm":       1,
	"lumberyard": 1,
	"quarry":     1,
	"market":     1,
	"barracks":   0,
	"walls":      0,
	"shrine":     0,
}

var Buildings = map[string]BuildingConfig{
	"town_hall": {
		Type:    "town_hall",
		Label:   "Ратуша",
		Purpose: "unlocks future development",
		BaseCost: ResourceValues{
			Gold:  150,
			Wood:  120,
			Stone: 100,
		},
		Effects: []string{"Нет эффекта в текущей версии"},
	},
	"farm": {
		Type:    "farm",
		Label:   "Ферма",
		Purpose: "increases food production",
		BaseCost: ResourceValues{
			Gold:  80,
			Wood:  80,
			Stone: 20,
		},
		Effects: []string{"+15 еды/час за уровень"},
	},
	"lumberyard": {
		Type:    "lumberyard",
		Label:   "Лесопилка",
		Purpose: "increases wood production",
		BaseCost: ResourceValues{
			Gold:  80,
			Wood:  40,
			Stone: 30,
		},
		Effects: []string{"+12 дерева/час за уровень"},
	},
	"quarry": {
		Type:    "quarry",
		Label:   "Каменоломня",
		Purpose: "increases stone production",
		BaseCost: ResourceValues{
			Gold:  90,
			Wood:  60,
			Stone: 40,
		},
		Effects: []string{"+10 камня/час за уровень"},
	},
	"barracks": {
		Type:    "barracks",
		Label:   "Казармы",
		Purpose: "future unit training",
		BaseCost: ResourceValues{
			Gold:  120,
			Wood:  90,
			Stone: 60,
		},
		Effects: []string{"Нет эффекта в текущей версии"},
	},
	"market": {
		Type:    "market",
		Label:   "Рынок",
		Purpose: "increases gold production",
		BaseCost: ResourceValues{
			Gold:  100,
			Wood:  80,
			Stone: 40,
		},
		Effects: []string{"+10 золота/час за уровень"},
	},
	"walls": {
		Type:    "walls",
		Label:   "Стены",
		Purpose: "future raid defense",
		BaseCost: ResourceValues{
			Gold:  130,
			Wood:  100,
			Stone: 120,
		},
		Effects: []string{"Нет эффекта в текущей версии"},
	},
	"shrine": {
		Type:    "shrine",
		Label:   "Святилище",
		Purpose: "future events and patrons",
		BaseCost: ResourceValues{
			Gold:  110,
			Wood:  70,
			Stone: 90,
		},
		Effects: []string{"Нет эффекта в текущей версии"},
	},
}

func BuildingCost(buildingType string, targetLevel int) ResourceValues {
	base := Buildings[buildingType].BaseCost
	multiplier := int64(targetLevel)
	return ResourceValues{
		Gold:  base.Gold * multiplier,
		Food:  base.Food * multiplier,
		Wood:  base.Wood * multiplier,
		Stone: base.Stone * multiplier,
	}
}

func BuildingUpgradeDurationSeconds(targetLevel int) int {
	return 60 * targetLevel
}

func IsBuildingType(buildingType string) bool {
	_, ok := Buildings[buildingType]
	return ok
}
