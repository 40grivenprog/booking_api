package notifications

// NotificationTranslations holds translations for notification messages
type NotificationTranslations struct {
	AppointmentRequest struct {
		Title       string
		Client      string
		Date        string
		Time        string
		Description string
		Action      string
	}
	AppointmentCancelled struct {
		Title        string
		Client       string
		Professional string
		Date         string
		Time         string
		Reason       string
	}
	AppointmentConfirmed struct {
		Title        string
		Professional string
		Date         string
		Time         string
	}
	OpenApp string
}

var translations = map[string]NotificationTranslations{
	"en": {
		AppointmentRequest: struct {
			Title       string
			Client      string
			Date        string
			Time        string
			Description string
			Action      string
		}{
			Title:       "🔔 New Appointment request!",
			Client:      "👤 Client:",
			Date:        "📅 Date:",
			Time:        "🕐 Time:",
			Description: "📝 Description:",
			Action:      "Please confirm or cancel this appointment.",
		},
		AppointmentCancelled: struct {
			Title        string
			Client       string
			Professional string
			Date         string
			Time         string
			Reason       string
		}{
			Title:        "🔔 Appointment Cancelled!",
			Client:       "👤 Client:",
			Professional: "👤 Professional:",
			Date:         "📅 Date:",
			Time:         "🕐 Time:",
			Reason:       "📝 Reason:",
		},
		AppointmentConfirmed: struct {
			Title        string
			Professional string
			Date         string
			Time         string
		}{
			Title:        "🔔 Appointment Confirmed!",
			Professional: "👤 Professional:",
			Date:         "📅 Date:",
			Time:         "🕐 Time:",
		},
		OpenApp: "📱 Open App",
	},
	"ru": {
		AppointmentRequest: struct {
			Title       string
			Client      string
			Date        string
			Time        string
			Description string
			Action      string
		}{
			Title:       "🔔 Новый запрос на запись!",
			Client:      "👤 Клиент:",
			Date:        "📅 Дата:",
			Time:        "🕐 Время:",
			Description: "📝 Описание:",
			Action:      "Пожалуйста, подтвердите или отмените эту запись.",
		},
		AppointmentCancelled: struct {
			Title        string
			Client       string
			Professional string
			Date         string
			Time         string
			Reason       string
		}{
			Title:        "🔔 Запись отменена!",
			Client:       "👤 Клиент:",
			Professional: "👤 Профессионал:",
			Date:         "📅 Дата:",
			Time:         "🕐 Время:",
			Reason:       "📝 Причина:",
		},
		AppointmentConfirmed: struct {
			Title        string
			Professional string
			Date         string
			Time         string
		}{
			Title:        "🔔 Запись подтверждена!",
			Professional: "👤 Профессионал:",
			Date:         "📅 Дата:",
			Time:         "🕐 Время:",
		},
		OpenApp: "📱 Открыть приложение",
	},
	"uk": {
		AppointmentRequest: struct {
			Title       string
			Client      string
			Date        string
			Time        string
			Description string
			Action      string
		}{
			Title:       "🔔 Новий запит на запис!",
			Client:      "👤 Клієнт:",
			Date:        "📅 Дата:",
			Time:        "🕐 Час:",
			Description: "📝 Опис:",
			Action:      "Будь ласка, підтвердіть або скасуйте цей запис.",
		},
		AppointmentCancelled: struct {
			Title        string
			Client       string
			Professional string
			Date         string
			Time         string
			Reason       string
		}{
			Title:        "🔔 Запис скасовано!",
			Client:       "👤 Клієнт:",
			Professional: "👤 Професіонал:",
			Date:         "📅 Дата:",
			Time:         "🕐 Час:",
			Reason:       "📝 Причина:",
		},
		AppointmentConfirmed: struct {
			Title        string
			Professional string
			Date         string
			Time         string
		}{
			Title:        "🔔 Запис підтверджено!",
			Professional: "👤 Професіонал:",
			Date:         "📅 Дата:",
			Time:         "🕐 Час:",
		},
		OpenApp: "📱 Відкрити додаток",
	},
	"pl": {
		AppointmentRequest: struct {
			Title       string
			Client      string
			Date        string
			Time        string
			Description string
			Action      string
		}{
			Title:       "🔔 Nowe żądanie wizyty!",
			Client:      "👤 Klient:",
			Date:        "📅 Data:",
			Time:        "🕐 Godzina:",
			Description: "📝 Opis:",
			Action:      "Proszę potwierdzić lub anulować tę wizytę.",
		},
		AppointmentCancelled: struct {
			Title        string
			Client       string
			Professional string
			Date         string
			Time         string
			Reason       string
		}{
			Title:        "🔔 Wizyta anulowana!",
			Client:       "👤 Klient:",
			Professional: "👤 Profesjonalista:",
			Date:         "📅 Data:",
			Time:         "🕐 Godzina:",
			Reason:       "📝 Powód:",
		},
		AppointmentConfirmed: struct {
			Title        string
			Professional string
			Date         string
			Time         string
		}{
			Title:        "🔔 Wizyta potwierdzona!",
			Professional: "👤 Profesjonalista:",
			Date:         "📅 Data:",
			Time:         "🕐 Godzina:",
		},
		OpenApp: "📱 Otwórz aplikację",
	},
}

// getTranslations returns translations for the given locale, defaults to "en" if locale is not supported
func getTranslations(locale string) NotificationTranslations {
	if t, ok := translations[locale]; ok {
		return t
	}
	return translations["en"]
}
