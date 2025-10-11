package api

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/vention/booking_api/internal/services/clients"
)

// ClientsHandler handles HTTP requests for clients
type ClientsHandler struct {
	clientsService clients.Service
}

// NewClientsHandler creates a new handler with dependency injection
func NewClientsHandler(service clients.Service) *ClientsHandler {
	return &ClientsHandler{
		clientsService: service,
	}
}

// ClientsHandlerParams defines the parameters for the ClientsHandler
type ClientsHandlerParams struct {
	Router         *gin.RouterGroup
	ClientsService clients.Service
}

// ClientsRegister registers the ClientsHandler with the router
func ClientsRegister(p ClientsHandlerParams) error {
	if p.Router == nil {
		return errors.New("missing router")
	}

	if p.ClientsService == nil {
		return errors.New("missing clients service")
	}

	h := NewClientsHandler(p.ClientsService)

	clients := p.Router.Group("/clients")
	{
		clients.POST("/register", h.RegisterClient)
		clients.GET("/:id/appointments", h.GetClientAppointments)
		clients.PATCH("/:id/appointments/:appointment_id/cancel", h.CancelClientAppointment)
	}

	return nil
}
