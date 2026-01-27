package api

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/vention/booking_api/internal/services/clients"
	"github.com/vention/booking_api/internal/services/notifications"
	"github.com/vention/booking_api/internal/token"
)

// ClientsHandler handles HTTP requests for clients
type ClientsHandler struct {
	clientsService       clients.Service
	tokenMaker           token.Maker
	notificationsService notifications.Service
}

// NewClientsHandler creates a new handler with dependency injection
func NewClientsHandler(service clients.Service, tokenMaker token.Maker, notificationsService notifications.Service) *ClientsHandler {
	return &ClientsHandler{
		clientsService:       service,
		tokenMaker:           tokenMaker,
		notificationsService: notificationsService,
	}
}

// ClientsHandlerParams defines the parameters for the ClientsHandler
type ClientsHandlerParams struct {
	Router               *gin.RouterGroup
	ClientsService       clients.Service
	TokenMaker           token.Maker
	NotificationsService notifications.Service
}

// ClientsRegister registers the ClientsHandler with the router
func ClientsRegister(p ClientsHandlerParams) error {
	if p.Router == nil {
		return errors.New("missing router")
	}

	if p.ClientsService == nil {
		return errors.New("missing clients service")
	}

	if p.TokenMaker == nil {
		return errors.New("missing token maker")
	}

	if p.NotificationsService == nil {
		return errors.New("missing notifications service")
	}

	h := NewClientsHandler(p.ClientsService, p.TokenMaker, p.NotificationsService)

	clients := p.Router.Group("/clients")
	{
		clients.POST("/book_appointment", h.BookAppointment)
		clients.POST("/register", h.RegisterClient)
		clients.GET("/appointments", h.GetClientAppointments)
		clients.PATCH("/appointments/:appointment_id/cancel", h.CancelClientAppointment)
		clients.PATCH("/update_locale", h.UpdateLocale)
	}

	return nil
}
