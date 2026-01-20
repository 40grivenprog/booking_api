package notifications

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Service interface {
	SendAppointmentRequestNotification(ctx context.Context, input SendAppointmentRequestNotificationInput) error
	SendAppointmentCancellationNotification(ctx context.Context, input SendAppointmentCancellationNotificationInput) error
	SendAppointmentConfirmationNotification(ctx context.Context, input SendAppointmentConfirmationNotificationInput) error
}

type service struct {
	bot *tgbotapi.BotAPI
}

func NewService(botToken string) (Service, error) {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return nil, err
	}

	return &service{
		bot: bot,
	}, nil
}

func (s *service) SendAppointmentRequestNotification(ctx context.Context, input SendAppointmentRequestNotificationInput) error {
	// Format time for display
	date, startTime, endTime := formatAppointmentTime(input.StartTime, input.EndTime)

	// Format message (same as in booking_client)
	messageText := fmt.Sprintf(
		"🔔 New Appointment request!\n\n👤 Client: %s\n📅 Date: %s\n🕐 Time: %s - %s\n📝 Description: %s\n\nPlease confirm or cancel this appointment.",
		input.ClientName,
		date,
		startTime,
		endTime,
		input.Description,
	)

	// Create inline keyboard with Web App button (using URL button)
	webAppButton := tgbotapi.NewInlineKeyboardButtonURL("📱 Open App", "https://t.me/testMfiAppBot/someRandomTestApp777")
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

	recepient := "Client"
	if input.Type == "professional" {
		recepient = "Professional"
	}

	// Format message
	messageText := fmt.Sprintf(
		"🔔 Appointment Cancelled!\n\n👤 %s: %s\n📅 Date: %s\n🕐 Time: %s - %s\n📝 Reason: %s",
		recepient,
		input.RespondentName,
		date,
		startTime,
		endTime,
		input.CancellationReason,
	)
	// Create inline keyboard with Web App button (using URL button)
	webAppButton := tgbotapi.NewInlineKeyboardButtonURL("📱 Open App", "https://t.me/testMfiAppBot/someRandomTestApp777")
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

func (s *service) SendAppointmentConfirmationNotification(ctx context.Context, input SendAppointmentConfirmationNotificationInput) error {
	// Format time for display
	date, startTime, endTime := formatAppointmentTime(input.StartTime, input.EndTime)

	// Format message (same as in booking_client)
	messageText := fmt.Sprintf(
		"🔔 Appointment Confirmed!\n\n👤 Professional: %s\n📅 Date: %s\n🕐 Time: %s - %s",
		input.ProfessionalName,
		date,
		startTime,
		endTime,
	)
	// Create inline keyboard with Web App button (using URL button)
	webAppButton := tgbotapi.NewInlineKeyboardButtonURL("📱 Open App", "https://t.me/testMfiAppBot/someRandomTestApp777")
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
