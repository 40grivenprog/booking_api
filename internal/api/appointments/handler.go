package api

import (
	"errors"

	"github.com/gin-gonic/gin"
)

type AppointmentsHandler struct {
	appointmentsRepo AppointmentsRepository
}

func NewAppointmentsHandler(appointmentsRepo AppointmentsRepository) *AppointmentsHandler {
	return &AppointmentsHandler{appointmentsRepo: appointmentsRepo}
}

type AppointmentsHandlerParams struct {
	Router           *gin.Engine
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
		appointments := api.Group("/appointments")
		{
			appointments.POST("/", h.CreateAppointment)
		}
	}
	return nil
}
