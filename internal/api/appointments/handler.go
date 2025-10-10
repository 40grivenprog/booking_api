package api

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/vention/booking_api/internal/services/appointments"
)

// AppointmentsHandler handles HTTP requests for appointments
type AppointmentsHandler struct {
	appointmentsService appointments.Service
}

// NewAppointmentsHandler creates a new handler with dependency injection
func NewAppointmentsHandler(service appointments.Service) *AppointmentsHandler {
	return &AppointmentsHandler{
		appointmentsService: service,
	}
}

type AppointmentsHandlerParams struct {
	Router              *gin.Engine
	AppointmentsService appointments.Service
}

func AppointmentsRegister(p AppointmentsHandlerParams) error {
	if p.Router == nil {
		return errors.New("missing router")
	}

	if p.AppointmentsService == nil {
		return errors.New("missing appointments service")
	}

	h := NewAppointmentsHandler(p.AppointmentsService)

	api := p.Router.Group("/api")
	{
		appointments := api.Group("/appointments")
		{
			appointments.POST("/", h.CreateAppointment)
		}
	}
	return nil
}
