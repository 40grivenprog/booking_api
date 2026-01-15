package api

import (
	"errors"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// NotificationsHandler handles HTTP requests for notifications
type NotificationsHandler struct {
	bot *tgbotapi.BotAPI
}

// NewNotificationsHandler creates a new handler with dependency injection
func NewNotificationsHandler(botToken string) (*NotificationsHandler, error) {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return nil, err
	}

	return &NotificationsHandler{
		bot: bot,
	}, nil
}

// NotificationsHandlerParams defines the parameters for the NotificationsHandler
type NotificationsHandlerParams struct {
	Router   *gin.RouterGroup
	BotToken string
}

// NotificationsRegister registers the NotificationsHandler with the router
func NotificationsRegister(p NotificationsHandlerParams) error {
	if p.Router == nil {
		return errors.New("missing router")
	}

	if p.BotToken == "" {
		return errors.New("missing bot token")
	}

	h, err := NewNotificationsHandler(p.BotToken)
	if err != nil {
		return err
	}

	notifications := p.Router.Group("/notifications")
	{
		notifications.POST("/send_appointment_request", h.SendAppointmentRequest)
		notifications.POST("/send_appointment_cancellation_notification", h.SendAppointmentCancellationNotification)
	}

	return nil
}
