package api

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/vention/booking_api/internal/services/professionals"
)

type ProfessionalsHandler struct {
	professionalsService professionals.Service
}

func NewProfessionalsHandler(service professionals.Service) *ProfessionalsHandler {
	return &ProfessionalsHandler{
		professionalsService: service,
	}
}

type ProfessionalsHandlerParams struct {
	Router               *gin.Engine
	ProfessionalsService professionals.Service
}

func ProfessionalsRegister(p ProfessionalsHandlerParams) error {
	if p.Router == nil {
		return errors.New("missing router")
	}

	if p.ProfessionalsService == nil {
		return errors.New("missing professionals service")
	}

	h := NewProfessionalsHandler(p.ProfessionalsService)

	api := p.Router.Group("/api")
	{
		professionals := api.Group("/professionals")
		{
			professionals.GET("", h.GetProfessionals)
			professionals.POST("/sign_in", h.SignInProfessional)
			professionals.GET("/:id/appointments", h.GetProfessionalAppointments)
			professionals.GET("/:id/appointment_dates", h.GetProfessionalAppointmentDates)
			professionals.PATCH("/:id/appointments/:appointment_id/confirm", h.ConfirmAppointment)
			professionals.PATCH("/:id/appointments/:appointment_id/cancel", h.CancelAppointment)
			professionals.POST("/:id/unavailable_appointments", h.CreateUnavailableAppointment)
			professionals.GET("/:id/availability", h.GetProfessionalAvailability)
			professionals.GET("/:id/timetable", h.GetProfessionalTimetable)
		}
	}

	return nil
}
