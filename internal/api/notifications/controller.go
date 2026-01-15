package api

import (
	"fmt"
	"net/http"
	"time"

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

// formatAppointmentTime formats appointment time for display
// Input: RFC3339 format (e.g., "2024-01-15T10:00:00Z")
// Output: date (YYYY-MM-DD), startTime (HH:MM), endTime (HH:MM)
func formatAppointmentTime(startTimeStr, endTimeStr string) (string, string, string) {
	startTime, err := time.Parse(time.RFC3339, startTimeStr)
	if err != nil {
		return startTimeStr[:10], startTimeStr[11:16], endTimeStr[11:16]
	}

	endTime, err := time.Parse(time.RFC3339, endTimeStr)
	if err != nil {
		return startTimeStr[:10], startTimeStr[11:16], endTimeStr[11:16]
	}

	date := startTime.Format("2006-01-02")
	startTimeFormatted := startTime.Format("15:04")
	endTimeFormatted := endTime.Format("15:04")

	return date, startTimeFormatted, endTimeFormatted
}
