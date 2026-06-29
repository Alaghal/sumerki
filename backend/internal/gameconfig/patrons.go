package gameconfig

type PatronConfig struct {
	Key              string
	Label            string
	ShortDescription string
	Flavor           string
	CurrentEffects   []string
	FutureEffects    []string
}

var PatronOrder = []string{
	"independent",
	"empire_of_dusk",
	"old_pact",
}

var Patrons = map[string]PatronConfig{
	"independent": {
		Key:              "independent",
		Label:            "Независимость",
		ShortDescription: "Ты никому не служишь и никому не платишь. Свобода полная, защита только своя.",
		Flavor:           "Свободные владетели живут без печатей и клятв. Но когда приходит беда, никто не обязан идти им на помощь.",
		CurrentEffects: []string{
			"Нет дани",
			"Нет защиты",
			"Полная свобода решений",
		},
		FutureEffects: []string{
			"Безопасность зависит только от собственных сил",
		},
	},
	"empire_of_dusk": {
		Key:              "empire_of_dusk",
		Label:            "Империя Заката",
		ShortDescription: "Империя обещает порядок, дороги и защиту. Цена будет позже.",
		Flavor:           "Белые печати, латунные списки и ровные дороги. Закат не сжигает город сразу. Он сначала считает его.",
		CurrentEffects: []string{
			"Дань ещё не активна в этой версии",
			"Защита ещё не активна в этой версии",
		},
		FutureEffects: []string{
			"Позже Империя сможет требовать дань",
			"Позже Империя сможет давать защиту",
			"Позже появятся гарнизоны и давление",
		},
	},
	"old_pact": {
		Key:              "old_pact",
		Label:            "Старый Договор",
		ShortDescription: "Старые клятвы, волхвы, княжеские круги и помощь против внешней угрозы. Долги будут позже.",
		Flavor:           "Старый Договор не делает всех друзьями. Он только напоминает, что некоторые враги должны стоять на одной стене.",
		CurrentEffects: []string{
			"Обязательства ещё не активны в этой версии",
			"Помощь ещё не активна в этой версии",
		},
		FutureEffects: []string{
			"Позже Старый Договор сможет требовать вклад",
			"Позже Старый Договор сможет помогать в обороне",
			"Позже появятся клятвенные долги",
		},
	},
}

func IsPatronKey(key string) bool {
	_, ok := Patrons[key]
	return ok
}
