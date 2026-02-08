package notifications

import (
	"context"
	"fmt"
	"net/url"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Service interface {
	SendAppointmentRequestNotification(ctx context.Context, input SendAppointmentRequestNotificationInput) error
	SendAppointmentCancellationNotification(ctx context.Context, input SendAppointmentCancellationNotificationInput) error
	SendAppointmentConfirmationNotification(ctx context.Context, input SendAppointmentConfirmationNotificationInput) error
	SendSubscriptionNotification(ctx context.Context, input SendSubscriptionNotificationInput) error
	SendGroupVisitAppointmentNotification(ctx context.Context, input SendGroupVisitAppointmentNotificationInput) error
	SendAcceptInviteNotification(ctx context.Context, input SendAcceptInviteNotificationInput) error
	SendPackageCreatedNotification(ctx context.Context, input SendPackageCreatedNotificationInput) error
}

type service struct {
	bot    *tgbotapi.BotAPI
	appURL string
}

func NewService(botToken string, appURL string) (Service, error) {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return nil, err
	}

	return &service{
		bot:    bot,
		appURL: appURL,
	}, nil
}

func (s *service) SendAppointmentRequestNotification(ctx context.Context, input SendAppointmentRequestNotificationInput) error {
	// Format time for display
	date, startTime, endTime := formatAppointmentTime(input.StartTime, input.EndTime)

	// Get locale (default to "en" if not provided)
	locale := input.Locale
	if locale == "" {
		locale = "en"
	}
	t := getTranslations(locale)

	// Format localized message
	messageText := fmt.Sprintf(
		"%s\n\n%s %s\n%s %s\n%s %s - %s\n%s %s\n\n%s",
		t.AppointmentRequest.Title,
		t.AppointmentRequest.Client,
		input.ClientName,
		t.AppointmentRequest.Date,
		date,
		t.AppointmentRequest.Time,
		startTime,
		endTime,
		t.AppointmentRequest.Description,
		input.Description,
		t.AppointmentRequest.Action,
	)

	// Create inline keyboard with Web App button (using URL button)
	payload := "appointment_" + input.AppointmentID.String() // или "appointment_"+id
	link := s.appURL + "?startapp=" + url.QueryEscape(payload)

	webAppButton := tgbotapi.NewInlineKeyboardButtonURL(t.CheckNotification, link)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(webAppButton),
	)

	// Send message via Telegram Bot API
	msg := tgbotapi.NewMessage(input.ChatID, messageText)
	msg.ReplyMarkup = keyboard

	_, err := s.bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) SendAppointmentCancellationNotification(ctx context.Context, input SendAppointmentCancellationNotificationInput) error {
	// Format time for display
	date, startTime, endTime := formatAppointmentTime(input.StartTime, input.EndTime)

	// Get locale (default to "en" if not provided)
	locale := input.Locale
	if locale == "" {
		locale = "en"
	}
	t := getTranslations(locale)

	// Determine recipient label based on type
	recipientLabel := t.AppointmentCancelled.Client
	if input.Type == "professional" {
		recipientLabel = t.AppointmentCancelled.Professional
	}

	// Format localized message
	messageText := fmt.Sprintf(
		"%s\n\n%s %s\n%s %s\n%s %s - %s\n%s %s",
		t.AppointmentCancelled.Title,
		recipientLabel,
		input.RespondentName,
		t.AppointmentCancelled.Date,
		date,
		t.AppointmentCancelled.Time,
		startTime,
		endTime,
		t.AppointmentCancelled.Reason,
		input.CancellationReason,
	)

	// Send message via Telegram Bot API
	msg := tgbotapi.NewMessage(input.ChatID, messageText)

	_, err := s.bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) SendAppointmentConfirmationNotification(ctx context.Context, input SendAppointmentConfirmationNotificationInput) error {
	// Format time for display
	date, startTime, endTime := formatAppointmentTime(input.StartTime, input.EndTime)

	locale := input.Locale
	if locale == "" {
		locale = "en"
	}
	t := getTranslations(locale)

	// Format localized message
	messageText := fmt.Sprintf(
		"%s\n\n%s %s\n%s %s\n%s %s - %s",
		t.AppointmentConfirmed.Title,
		t.AppointmentConfirmed.Professional,
		input.ProfessionalName,
		t.AppointmentConfirmed.Date,
		date,
		t.AppointmentConfirmed.Time,
		startTime,
		endTime,
	)

	// Send message via Telegram Bot API
	msg := tgbotapi.NewMessage(input.ChatID, messageText)
	_, err := s.bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) SendSubscriptionNotification(ctx context.Context, input SendSubscriptionNotificationInput) error {
	// Get locale (default to "en" if not provided)
	locale := input.Locale
	if locale == "" {
		locale = "en"
	}
	t := getTranslations(locale)

	// Format localized message
	messageText := fmt.Sprintf(
		"%s\n\n%s %s",
		t.Subscription.Title,
		t.Subscription.Client,
		input.ClientName,
	)

	// Send message via Telegram Bot API
	msg := tgbotapi.NewMessage(input.ChatID, messageText)

	_, err := s.bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) SendGroupVisitAppointmentNotification(ctx context.Context, input SendGroupVisitAppointmentNotificationInput) error {
	date, startTime, endTime := formatAppointmentTime(input.StartTime, input.EndTime)

	locale := input.Locale
	if locale == "" {
		locale = "en"
	}
	t := getTranslations(locale)

	// Format localized message
	messageText := fmt.Sprintf(
		"%s %s\n%s\n%s %s\n%s %s - %s\n%s %s",
		t.GroupVisitAppointment.Professional,
		input.ProfessionalName,
		t.GroupVisitAppointment.Message,
		t.GroupVisitAppointment.Date,
		date,
		t.GroupVisitAppointment.Time,
		startTime,
		endTime,
		t.GroupVisitAppointment.Description,
		input.Description,
	)
	// Create inline keyboard with Web App button (using URL button)
	payload := "invite_" + input.InviteID.String()
	webAppButton := tgbotapi.NewInlineKeyboardButtonURL(t.CheckNotification, s.appURL+"?startapp="+url.QueryEscape(payload))
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(webAppButton),
	)

	// Send message via Telegram Bot API
	msg := tgbotapi.NewMessage(input.ChatID, messageText)
	msg.ReplyMarkup = keyboard
	_, err := s.bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) SendAcceptInviteNotification(ctx context.Context, input SendAcceptInviteNotificationInput) error {
	// Format time for display
	date, startTime, endTime := formatAppointmentTime(input.StartTime, input.EndTime)

	// Get locale (default to "en" if not provided)
	locale := input.Locale
	if locale == "" {
		locale = "en"
	}
	t := getTranslations(locale)

	// Format type for display
	typeDisplay := input.Type
	switch typeDisplay {
	case "split":
		typeDisplay = "Split"
	case "group":
		typeDisplay = "Group"
	default:
		typeDisplay = "Unknown"
	}

	// Format localized message
	messageText := fmt.Sprintf(
		"%s\n\n%s\n%s %s\n%s %s\n%s %s\n%s %s",
		fmt.Sprintf(t.AcceptInvite.JoinedMessage, input.ClientName),
		t.AcceptInvite.DetailsTitle,
		t.AcceptInvite.Date,
		date,
		t.AcceptInvite.StartTime,
		startTime,
		t.AcceptInvite.EndTime,
		endTime,
		t.AcceptInvite.Type,
		typeDisplay,
	)

	// Send message via Telegram Bot API
	msg := tgbotapi.NewMessage(input.ChatID, messageText)
	_, err := s.bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) SendPackageCreatedNotification(ctx context.Context, input SendPackageCreatedNotificationInput) error {
	// Get locale (default to "en" if not provided)
	locale := input.Locale
	if locale == "" {
		locale = "en"
	}
	t := getTranslations(locale)

	// Format dates for display
	issuedAtFormatted := formatDate(input.IssuedAt)
	expiresAtFormatted := formatDate(input.ExpiresAt)

	// Format localized message
	messageText := fmt.Sprintf(
		"%s\n\n%s\n%s %d\n%s %s\n%s %s",
		fmt.Sprintf(t.PackageCreated.Title, input.ProfessionalName),
		t.PackageCreated.DetailsTitle,
		t.PackageCreated.AppointmentsNumber,
		input.ApppointmentsNumber,
		t.PackageCreated.IssuedAt,
		issuedAtFormatted,
		t.PackageCreated.ExpiresAt,
		expiresAtFormatted,
	)

	// Send message via Telegram Bot API
	msg := tgbotapi.NewMessage(input.ChatID, messageText)
	_, err := s.bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}
