package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	common "github.com/vention/booking_api/internal/api/common"
)

// SendAppointmentRequest handles POST /api/notifications/send_appointment_request
func (h *NotificationsHandler) SendAppointmentRequest(c *gin.Context) {
	req, ok := common.BindAndValidate[SendAppointmentRequestRequest](c)
	if !ok {
		return
	}

	// Format time for display
	date, startTime, endTime := formatAppointmentTime(req.StartTime, req.EndTime)

	// Format message (same as in booking_client)
	messageText := fmt.Sprintf(
		"🔔 New Appointment Request!\n\n👤 Client: %s\n📅 Date: %s\n🕐 Time: %s - %s\n📝 Description: %s\n\nPlease confirm or cancel this appointment.",
		req.ClientName,
		date,
		startTime,
		endTime,
		req.Description,
	)

	// Create inline keyboard with Web App button (using URL button)
	webAppButton := tgbotapi.NewInlineKeyboardButtonURL("📱 Open App", "https://t.me/testMfiAppBot/someRandomTestApp777")
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(webAppButton),
	)

	// Send message via Telegram Bot API
	msg := tgbotapi.NewMessage(req.ChatID, messageText)
	msg.ReplyMarkup = keyboard

	_, err := h.bot.Send(msg)
	if err != nil {
		common.HandleErrorResponse(c, http.StatusInternalServerError, common.ErrorTypeInternal, "Failed to send notification", err)
		return
	}

	response := SendAppointmentRequestResponse{
		Success: true,
		Message: "Notification sent successfully",
	}
	c.JSON(http.StatusOK, response)
}

// SendAppointmentCancellationNotification handles POST /api/notifications/send_appointment_cancellation_notification
func (h *NotificationsHandler) SendAppointmentCancellationNotification(c *gin.Context) {
	req, ok := common.BindAndValidate[SendAppointmentCancellationNotificationRequest](c)
	if !ok {
		return
	}

	// Format time for display
	date, startTime, endTime := formatAppointmentTime(req.StartTime, req.EndTime)

	// Format message
	messageText := fmt.Sprintf(
		"❌ Appointment Cancelled\n\n👤 Client: %s %s\n📅 Date: %s\n🕐 Time: %s - %s\n📝 Reason: %s",
		req.FirstName,
		req.LastName,
		date,
		startTime,
		endTime,
		req.CancellationReason,
	)

	// Create inline keyboard with Web App button
	webAppButton := tgbotapi.NewInlineKeyboardButtonURL("📱 Open App", "https://t.me/testMfiAppBot/someRandomTestApp777")
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(webAppButton),
	)

	// Send message via Telegram Bot API
	msg := tgbotapi.NewMessage(req.ChatID, messageText)
	msg.ReplyMarkup = keyboard

	_, err := h.bot.Send(msg)
	if err != nil {
		common.HandleErrorResponse(c, http.StatusInternalServerError, common.ErrorTypeInternal, "Failed to send notification", err)
		return
	}

	response := SendAppointmentRequestResponse{
		Success: true,
		Message: "Cancellation notification sent successfully",
	}
	c.JSON(http.StatusOK, response)
}
