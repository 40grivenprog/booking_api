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
	Subscription struct {
		Title  string
		Client string
	}
	GroupVisitAppointment struct {
		Professional string
		Message      string
		Date         string
		Time         string
		Description  string
	}
	AcceptInvite struct {
		JoinedMessage string
		DetailsTitle  string
		Date          string
		StartTime     string
		EndTime       string
		Type          string
	}
	PackageCreated struct {
		Title              string
		DetailsTitle       string
		AppointmentsNumber string
		IssuedAt           string
		ExpiresAt          string
	}
	OpenApp           string
	CheckNotification string
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
		OpenApp:           "📱 Open App",
		CheckNotification: "🔔 Check Notification",
		Subscription: struct {
			Title  string
			Client string
		}{
			Title:  "🔔 New client subscription!",
			Client: "👤 Client:",
		},
		GroupVisitAppointment: struct {
			Professional string
			Message      string
			Date         string
			Time         string
			Description  string
		}{
			Professional: "👤 Trainer:",
			Message:      "📋 Posted a slot for a group visit",
			Date:         "📅 Date:",
			Time:         "🕐 Time:",
			Description:  "📝 Description:",
		},
		AcceptInvite: struct {
			JoinedMessage string
			DetailsTitle  string
			Date          string
			StartTime     string
			EndTime       string
			Type          string
		}{
			JoinedMessage: "%s joined the training",
			DetailsTitle:  "Training Details",
			Date:          "📅 Date:",
			StartTime:     "🕐 Start:",
			EndTime:       "🕐 End:",
			Type:          "📋 Type:",
		},
		PackageCreated: struct {
			Title              string
			DetailsTitle       string
			AppointmentsNumber string
			IssuedAt           string
			ExpiresAt          string
		}{
			Title:              "%s add new package for you!",
			DetailsTitle:       "Details:",
			AppointmentsNumber: "📦 Number of appointments:",
			IssuedAt:           "📅 Issued at:",
			ExpiresAt:          "⏰ Expires at:",
		},
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
		OpenApp:           "📱 Открыть приложение",
		CheckNotification: "🔔 Проверить уведомления",
		Subscription: struct {
			Title  string
			Client string
		}{
			Title:  "🔔 Новая подписка клиента!",
			Client: "👤 Клиент:",
		},
		GroupVisitAppointment: struct {
			Professional string
			Message      string
			Date         string
			Time         string
			Description  string
		}{
			Professional: "👤 Тренер:",
			Message:      "📋 Выложил слот для группового визита",
			Date:         "📅 Дата:",
			Time:         "🕐 Время:",
			Description:  "📝 Описание:",
		},
		AcceptInvite: struct {
			JoinedMessage string
			DetailsTitle  string
			Date          string
			StartTime     string
			EndTime       string
			Type          string
		}{
			JoinedMessage: "%s присоединился к тренировке",
			DetailsTitle:  "Детали тренировки",
			Date:          "📅 Дата:",
			StartTime:     "🕐 Начало:",
			EndTime:       "🕐 Конец:",
			Type:          "📋 Тип:",
		},
		PackageCreated: struct {
			Title              string
			DetailsTitle       string
			AppointmentsNumber string
			IssuedAt           string
			ExpiresAt          string
		}{
			Title:              "%s добавил новый пакет для вас!",
			DetailsTitle:       "Детали:",
			AppointmentsNumber: "📦 Количество записей:",
			IssuedAt:           "📅 Выдан:",
			ExpiresAt:          "⏰ Истекает:",
		},
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
		OpenApp:           "📱 Відкрити додаток",
		CheckNotification: "🔔 Перевірити сповіщення",
		Subscription: struct {
			Title  string
			Client string
		}{
			Title:  "🔔 Нова підписка клієнта!",
			Client: "👤 Клієнт:",
		},
		GroupVisitAppointment: struct {
			Professional string
			Message      string
			Date         string
			Time         string
			Description  string
		}{
			Professional: "👤 Тренер:",
			Message:      "📋 Виклав слот для групового візиту",
			Date:         "📅 Дата:",
			Time:         "🕐 Час:",
			Description:  "📝 Опис:",
		},
		AcceptInvite: struct {
			JoinedMessage string
			DetailsTitle  string
			Date          string
			StartTime     string
			EndTime       string
			Type          string
		}{
			JoinedMessage: "%s приєднався до тренування",
			DetailsTitle:  "Деталі тренування",
			Date:          "📅 Дата:",
			StartTime:     "🕐 Початок:",
			EndTime:       "🕐 Кінець:",
			Type:          "📋 Тип:",
		},
		PackageCreated: struct {
			Title              string
			DetailsTitle       string
			AppointmentsNumber string
			IssuedAt           string
			ExpiresAt          string
		}{
			Title:              "%s додав новий пакет для вас!",
			DetailsTitle:       "Деталі:",
			AppointmentsNumber: "📦 Кількість записів:",
			IssuedAt:           "📅 Видано:",
			ExpiresAt:          "⏰ Закінчується:",
		},
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
		OpenApp:           "📱 Otwórz aplikację",
		CheckNotification: "🔔 Sprawdź powiadomienia",
		Subscription: struct {
			Title  string
			Client string
		}{
			Title:  "🔔 Nowa subskrypcja klienta!",
			Client: "👤 Klient:",
		},
		GroupVisitAppointment: struct {
			Professional string
			Message      string
			Date         string
			Time         string
			Description  string
		}{
			Professional: "👤 Trener:",
			Message:      "📋 Opublikował slot na wizytę grupową",
			Date:         "📅 Data:",
			Time:         "🕐 Godzina:",
			Description:  "📝 Opis:",
		},
		AcceptInvite: struct {
			JoinedMessage string
			DetailsTitle  string
			Date          string
			StartTime     string
			EndTime       string
			Type          string
		}{
			JoinedMessage: "%s dołączył do treningu",
			DetailsTitle:  "Szczegóły treningu",
			Date:          "📅 Data:",
			StartTime:     "🕐 Początek:",
			EndTime:       "🕐 Koniec:",
			Type:          "📋 Typ:",
		},
		PackageCreated: struct {
			Title              string
			DetailsTitle       string
			AppointmentsNumber string
			IssuedAt           string
			ExpiresAt          string
		}{
			Title:              "%s dodał nowy pakiet dla Ciebie!",
			DetailsTitle:       "Szczegóły:",
			AppointmentsNumber: "📦 Liczba wizyt:",
			IssuedAt:           "📅 Wydano:",
			ExpiresAt:          "⏰ Wygasa:",
		},
	},
}

// getTranslations returns translations for the given locale, defaults to "en" if locale is not supported
func getTranslations(locale string) NotificationTranslations {
	if t, ok := translations[locale]; ok {
		return t
	}
	return translations["en"]
}
