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
		"menu.play":          "Играть",
		"menu.language":      "Язык",
		"menu.mode":          "Режим игры",
		"menu.skins":         "Скины",
		"menu.player":        "Игрок",
		"menu.leaderboard":   "Лидерборд",
		"game.press_space":   "Нажмите Пробел, чтобы начать. Esc - выйти в меню",
		"game.pause":         "Пауза (Пробел — продолжить)",
		"game.game_over":     "Игра окончена! Счёт: %d. Пробел — начать заново",
		"game.score":         "Счёт: %d",
		"mode.normal":        "Обычный",
		"mode.defaltyk":      "Сергей Андреевич Дефалтук",
		"skins.normal":       "Обычный",
		"skins.halfup":       "Карим",
		"skins.rantlol":      "Илья",
		"name.title":         "Введите ваше имя",
		"name.placeholder":   "Имя",
		"name.hint_required": "Enter - сохранить",
		"name.hint_optional": "Enter - сохранить, Esc - отменить",
		"name.edit_hint":     "N - изменить имя",
		"name.not_set":       "не задано",
		"name.error_empty":   "Имя не может быть пустым",
		"name.error_long":    "Имя не может быть длиннее 20 символов",
		"name.error_invalid": "Имя содержит недопустимые символы",
		"name.error_save":    "Не удалось сохранить имя",

		"leaderboard.title":         "Лидерборд",
		"leaderboard.loading":       "Загрузка...",
		"leaderboard.error":         "Не удалось загрузить лидерборд",
		"leaderboard.empty":         "Пока нет результатов",
		"leaderboard.refresh_hint":  "R/Пробел - обновить, Esc - меню",
		"leaderboard.rank_header":   "Место",
		"leaderboard.name_header":   "Игрок",
		"leaderboard.score_header":  "Счёт",
		"game.score_submitting":     "Отправляем счёт...",
		"game.score_submitted":      "Счёт отправлен! Место: #%d",
		"game.score_best_unchanged": "Лучший счёт уже в таблице: #%d",
		"game.score_submit_failed":  "Не удалось отправить счёт",
		"code.title":                "Введите секретный код",
		"code.placeholder":          "Код",
		"code.hint":                 "Enter - подтвердить, Esc - отменить",
		"code.error_wrong":          "Неверный код",
		"code.success":              "Разблокировано! Скины и моды разблокированы",
		"code.menu_hint":            "C - ввести код",
		"lock.locked":               "Заблокировано",
	},
	LangEN: {
		"menu.play":          "Play",
		"menu.language":      "Language",
		"menu.skins":         "Game skins",
		"menu.player":        "Player",
		"menu.mode":          "Game mode",
		"menu.leaderboard":   "Leaderboard",
		"game.press_space":   "Press Space to Start. Esc - exit to menu",
		"game.pause":         "Pause (Space to Continue)",
		"game.game_over":     "Game Over! Score: %d. Press Space to Restart",
		"game.score":         "Score: %d",
		"mode.normal":        "Normal",
		"mode.defaltyk":      "Defaltyk",
		"skins.normal":       "Normal",
		"skins.halfup":       "Halfup",
		"skins.rantlol":      "Rantlol",
		"name.title":         "Enter your name",
		"name.placeholder":   "Name",
		"name.hint_required": "Enter - save",
		"name.hint_optional": "Enter - save, Esc - cancel",
		"name.edit_hint":     "N - change name",
		"name.not_set":       "not set",
		"name.error_empty":   "Name cannot be empty",
		"name.error_long":    "Name cannot be longer than 20 characters",
		"name.error_invalid": "Name contains unsupported characters",
		"name.error_save":    "Failed to save name",

		"leaderboard.title":         "Leaderboard",
		"leaderboard.loading":       "Loading...",
		"leaderboard.error":         "Could not load leaderboard",
		"leaderboard.empty":         "No scores yet",
		"leaderboard.refresh_hint":  "R/Space - refresh, Esc - menu",
		"leaderboard.rank_header":   "Rank",
		"leaderboard.name_header":   "Player",
		"leaderboard.score_header":  "Score",
		"game.score_submitting":     "Submitting score...",
		"game.score_submitted":      "Score submitted! Rank: #%d",
		"game.score_best_unchanged": "Best score remains ranked #%d",
		"game.score_submit_failed":  "Could not submit score",
		"code.title":                "Enter secret code",
		"code.placeholder":          "Code",
		"code.hint":                 "Enter - confirm, Esc - cancel",
		"code.error_wrong":          "Wrong code",
		"code.success":              "Unlocked! Skins and modes unlocked",
		"code.menu_hint":            "C - enter code",
		"lock.locked":               "Locked",
	},
	LangUA: {
		"menu.play":          "Грати",
		"menu.language":      "Мова",
		"menu.player":        "Гравець",
		"menu.skins":         "Скiни гри",
		"menu.mode":          "Режим гри",
		"menu.leaderboard":   "Лідерборд",
		"game.press_space":   "Натисніть Пробiл, щоб почати. Esc - вийти в меню",
		"game.pause":         "Пауза (Пробiл — продовжити)",
		"game.game_over":     "Гра закінчена! Рахунок: %d. Пробiл — почати заново",
		"game.score":         "Рахунок: %d",
		"mode.normal":        "Normal",
		"mode.defaltyk":      "Сергій Андрійович Дефалтук",
		"skins.normal":       "Normal",
		"skins.halfup":       "Карім",
		"skins.rantlol":      "Ілля",
		"name.title":         "Введіть ваше ім'я",
		"name.placeholder":   "Ім'я",
		"name.hint_required": "Enter - зберегти",
		"name.hint_optional": "Enter - зберегти, Esc - скасувати",
		"name.edit_hint":     "N - змінити ім'я",
		"name.not_set":       "не задано",
		"name.error_empty":   "Ім'я не може бути порожнім",
		"name.error_long":    "Ім'я не може бути довшим за 20 символів",
		"name.error_invalid": "Ім'я містить недопустимі символи",
		"name.error_save":    "Не вдалося зберегти ім'я",

		"leaderboard.title":         "Лідерборд",
		"leaderboard.loading":       "Завантаження...",
		"leaderboard.error":         "Не вдалося завантажити лідерборд",
		"leaderboard.empty":         "Поки немає результатів",
		"leaderboard.refresh_hint":  "R/Пробіл - оновити, Esc - меню",
		"leaderboard.rank_header":   "Місце",
		"leaderboard.name_header":   "Гравець",
		"leaderboard.score_header":  "Рахунок",
		"game.score_submitting":     "Надсилаємо рахунок...",
		"game.score_submitted":      "Рахунок надіслано! Місце: #%d",
		"game.score_best_unchanged": "Найкращий рахунок уже в таблиці: #%d",
		"game.score_submit_failed":  "Не вдалося надіслати рахунок",
		"code.title":                "Введіть секретний код",
		"code.placeholder":          "Код",
		"code.hint":                 "Enter - підтвердити, Esc - скасувати",
		"code.error_wrong":          "Неправильний код",
		"code.success":              "Розблоковано! Скіни і режими розблоковані",
		"code.menu_hint":            "C - ввести код",
		"lock.locked":               "Заблоковано",
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
