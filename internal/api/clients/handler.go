package api

import (
	"errors"

	"github.com/gin-gonic/gin"
)

type ClientsHandler struct {
	clientsRepo ClientsRepository
}

type ClientsHandlerParams struct {
	Router      *gin.Engine
	ClientsRepo ClientsRepository
}

func ClientsRegister(p ClientsHandlerParams) error {
	if p.Router == nil {
		return errors.New("missing router")
	}

	if p.ClientsRepo == nil {
		return errors.New("missing clients repository")
	}

	h := &ClientsHandler{
		clientsRepo: p.ClientsRepo,
	}

	api := p.Router.Group("/api")
	{
		clients := api.Group("/clients")
		{
			clients.POST("/register", h.RegisterClient)
			clients.GET("/:id/appointments", h.GetClientAppointments)
			clients.PATCH("/:id/appointments/:appointment_id/cancel", h.CancelClientAppointment)
		}
	}

	return nil
}
