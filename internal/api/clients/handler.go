package api

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/vention/booking_api/internal/services/clients"
	"github.com/vention/booking_api/internal/token"
)

// ClientsHandler handles HTTP requests for clients
type ClientsHandler struct {
	clientsService clients.Service
	tokenMaker     token.Maker
}

// NewClientsHandler creates a new handler with dependency injection
func NewClientsHandler(service clients.Service, tokenMaker token.Maker) *ClientsHandler {
	return &ClientsHandler{
		clientsService: service,
		tokenMaker:     tokenMaker,
	}
}

// ClientsHandlerParams defines the parameters for the ClientsHandler
type ClientsHandlerParams struct {
	Router         *gin.RouterGroup
	ClientsService clients.Service
	TokenMaker     token.Maker
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

	h := NewClientsHandler(p.ClientsService, p.TokenMaker)

	clients := p.Router.Group("/clients")
	{
		clients.POST("/book_appointment", h.BookAppointment)
		clients.POST("/register", h.RegisterClient)
		clients.GET("/appointments", h.GetClientAppointments)
		clients.PATCH("/appointments/:appointment_id/cancel", h.CancelClientAppointment)
	}

	return nil
}
