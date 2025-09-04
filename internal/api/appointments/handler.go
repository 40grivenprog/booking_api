package api

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/vention/booking_api/internal/config"
)

type AppointmentsHandler struct {
	appointmentsRepo AppointmentsRepository
}

func NewAppointmentsHandler(appointmentsRepo AppointmentsRepository) *AppointmentsHandler {
	return &AppointmentsHandler{appointmentsRepo: appointmentsRepo}
}

type AppointmentsHandlerParams struct {
	Router           *gin.Engine
	Cfg              *config.Config
	AppointmentsRepo AppointmentsRepository
}

func AppointmentsRegister(p AppointmentsHandlerParams) error {
	if p.Router == nil {
		return errors.New("missing router")
	}

	if p.AppointmentsRepo == nil {
		return errors.New("missing clients repository")
	}

	h := &AppointmentsHandler{
		appointmentsRepo: p.AppointmentsRepo,
	}

	api := p.Router.Group("/api")
	{
		clients := api.Group("/appointments")
		{
			clients.POST("/", h.CreateAppointment)
		}
	}

	return nil
}
