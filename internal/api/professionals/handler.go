package api

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/vention/booking_api/internal/services/professionals"
	"github.com/vention/booking_api/internal/token"
)

type ProfessionalsHandler struct {
	professionalsService professionals.Service
	tokenMaker           token.Maker
}

func NewProfessionalsHandler(service professionals.Service, tokenMaker token.Maker) *ProfessionalsHandler {
	return &ProfessionalsHandler{
		professionalsService: service,
		tokenMaker:           tokenMaker,
	}
}

type ProfessionalsHandlerParams struct {
	Router               *gin.RouterGroup
	ProfessionalsService professionals.Service
	TokenMaker           token.Maker
}

func ProfessionalsRegister(p ProfessionalsHandlerParams) error {
	if p.Router == nil {
		return errors.New("missing router")
	}

	if p.ProfessionalsService == nil {
		return errors.New("missing professionals service")
	}

	if p.TokenMaker == nil {
		return errors.New("missing token maker")
	}

	h := NewProfessionalsHandler(p.ProfessionalsService, p.TokenMaker)

	professionals := p.Router.Group("/professionals")
	{
		professionals.GET("", h.GetProfessionals)
		professionals.POST("/sign_in", h.SignInProfessional)
		professionals.GET("/appointments", h.GetProfessionalAppointments)
		professionals.GET("/:id/appointment_dates", h.GetProfessionalAppointmentDates)
		professionals.PATCH("/:id/appointments/:appointment_id/confirm", h.ConfirmAppointment)
		professionals.PATCH("/:id/appointments/:appointment_id/cancel", h.CancelAppointment)
		professionals.POST("/:id/unavailable_appointments", h.CreateUnavailableAppointment)
		professionals.GET("/:id/availability", h.GetProfessionalAvailability)
		professionals.GET("/:id/timetable", h.GetProfessionalTimetable)
		professionals.GET("/:id/clients", h.GetProfessionalClients)
		professionals.GET("/:id/previous_appointments", h.GetPreviousAppointmentsByClient)
	}

	return nil
}
