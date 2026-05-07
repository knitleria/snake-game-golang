package i18n

type Lang int

const (
	LangRU Lang = iota
	LangEN
	LangUA
)

var langOrder = []Lang{LangRU, LangEN, LangUA}

var dict = map[Lang]map[string]string{
	LangRU: {
		"menu.play":        "Играть",
		"menu.language":    "Язык",
		"menu.mode":        "Режим игры (TBD)",
		"menu.skins":       "Скины",
		"game.press_space": "Нажмите Пробел, чтобы начать. Esc - выйти в меню",
		"game.pause":       "Пауза (Пробел — продолжить)",
		"game.game_over":   "Игра окончена! Счёт: %d. Пробел — начать заново",
		"game.score":       "Счёт: %d",
		"mode.normal":      "Обычный",
		"mode.defaltyk":    "Сергей Андреевич Дефалтук",
		"skins.normal":     "Обычный",
		"skins.halfup":     "Карим",
		"skins.rantlol":    "Илья",
	},
	LangEN: {
		"menu.play":        "Play",
		"menu.language":    "Language",
		"menu.skins":       "Game skins",
		"menu.mode":        "Game mode (TBD)",
		"game.press_space": "Press Space to Start. Esc - exit to menu",
		"game.pause":       "Pause (Space to Continue)",
		"game.game_over":   "Game Over! Score: %d. Press Space to Restart",
		"game.score":       "Score: %d",
		"mode.normal":      "Normal",
		"mode.defaltyk":    "Defaltyk",
		"skins.normal":     "Normal",
		"skins.halfup":     "Halfup",
		"skins.rantlol":    "Rantlol",
	},
	LangUA: {
		"menu.play":        "Грати",
		"menu.language":    "Мова",
		"menu.skins":       "Скiни гри",
		"menu.mode":        "Режим гри (TBD)",
		"game.press_space": "Натисніть Пробiл, щоб почати. Esc - вийти в меню",
		"game.pause":       "Пауза (Пробiл — продовжити)",
		"game.game_over":   "Гра закінчена! Рахунок: %d. Пробiл — почати заново",
		"game.score":       "Рахунок: %d",
		"mode.normal":      "Normal",
		"mode.defaltyk":    "Сергій Андрійович Дефалтук",
		"skins.normal":     "Normal",
		"skins.halfup":     "Карім",
		"skins.rantlol":    "Ілля",
	},
}

var currentLang = LangRU

func LangOrder() []Lang {
	out := make([]Lang, len(langOrder))
	copy(out, langOrder)
	return out
}

func SetLang(l Lang) {
	currentLang = l
}

func CurrentLang() Lang {
	return currentLang
}

func NextLang() Lang {
	for i, l := range langOrder {
		if l == currentLang {
			return langOrder[(i+1)%len(langOrder)]
		}
	}
	return langOrder[0]
}

func T(key string) string {
	if s, ok := dict[currentLang][key]; ok {
		return s
	}
	if s, ok := dict[LangEN][key]; ok {
		return s
	}
	return key
}
