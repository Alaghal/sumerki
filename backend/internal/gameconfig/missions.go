package gameconfig

type MissionRequirements struct {
	TotalUnits int64
	Scouts     int64
}

type MissionConfig struct {
	Key                   string
	Label                 string
	Type                  string
	Description           string
	DurationSeconds       int
	MinimumRequirements   MissionRequirements
	RecommendedTotalUnits int64
	BaseRewards           ResourceValues
	Risk                  string
	BaseLossPercent       int64
	MaxLossPercent        int64
}

var MissionOrder = []string{
	"black_forest_expedition",
	"old_kurgan_expedition",
	"dry_ford_scouting",
}

var Missions = map[string]MissionConfig{
	"black_forest_expedition": {
		Key:             "black_forest_expedition",
		Label:           "Чёрный Лес",
		Type:            "expedition",
		Description:     "Охотники шепчут, что в Чёрном Лесу пропадают тропы, но старые склады всё ещё можно найти.",
		DurationSeconds: 120,
		MinimumRequirements: MissionRequirements{
			TotalUnits: 5,
		},
		RecommendedTotalUnits: 7,
		BaseRewards: ResourceValues{
			Gold: 40,
			Food: 80,
			Wood: 120,
		},
		Risk:            "medium",
		BaseLossPercent: 8,
		MaxLossPercent:  12,
	},
	"old_kurgan_expedition": {
		Key:             "old_kurgan_expedition",
		Label:           "Старый Курган",
		Type:            "expedition",
		Description:     "В кургане лежат кости старых дружин. Иногда они не любят гостей.",
		DurationSeconds: 180,
		MinimumRequirements: MissionRequirements{
			TotalUnits: 8,
		},
		RecommendedTotalUnits: 10,
		BaseRewards: ResourceValues{
			Gold:  150,
			Wood:  30,
			Stone: 100,
		},
		Risk:            "medium",
		BaseLossPercent: 14,
		MaxLossPercent:  18,
	},
	"dry_ford_scouting": {
		Key:             "dry_ford_scouting",
		Label:           "Сухой Брод",
		Type:            "scouting",
		Description:     "Сухой Брод тихий только на карте. На деле там считают чужие следы.",
		DurationSeconds: 75,
		MinimumRequirements: MissionRequirements{
			Scouts: 1,
		},
		RecommendedTotalUnits: 3,
		BaseRewards: ResourceValues{
			Gold:  50,
			Food:  40,
			Wood:  20,
			Stone: 20,
		},
		Risk:            "low",
		BaseLossPercent: 3,
		MaxLossPercent:  5,
	},
}

func IsMissionKey(key string) bool {
	_, ok := Missions[key]
	return ok
}
