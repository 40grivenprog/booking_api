package api

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/vention/booking_api/internal/services/notifications"
	"github.com/vention/booking_api/internal/services/professionals"
	"github.com/vention/booking_api/internal/token"
)

type ProfessionalsHandler struct {
	professionalsService professionals.Service
	tokenMaker           token.Maker
	notificationsService notifications.Service
}

func NewProfessionalsHandler(service professionals.Service, tokenMaker token.Maker, notificationsService notifications.Service) *ProfessionalsHandler {
	return &ProfessionalsHandler{
		professionalsService: service,
		tokenMaker:           tokenMaker,
		notificationsService: notificationsService,
	}
}

type ProfessionalsHandlerParams struct {
	Router               *gin.RouterGroup
	ProfessionalsService professionals.Service
	TokenMaker           token.Maker
	NotificationsService notifications.Service
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

	if p.NotificationsService == nil {
		return errors.New("missing notifications service")
	}

	h := NewProfessionalsHandler(p.ProfessionalsService, p.TokenMaker, p.NotificationsService)

	professionals := p.Router.Group("/professionals")
	{
		professionals.GET("", h.GetProfessionals)
		professionals.POST("/sign_in", h.SignInProfessional)
		professionals.GET("/appointments", h.GetProfessionalAppointments)
		professionals.GET("/:id/appointment_dates", h.GetProfessionalAppointmentDates)
		professionals.PATCH("/appointments/:appointment_id/confirm", h.ConfirmAppointment)
		professionals.PATCH("/appointments/:appointment_id/cancel", h.CancelAppointment)
		professionals.POST("/unavailable_appointments", h.CreateUnavailableAppointment)
		professionals.GET("/:id/availability", h.GetProfessionalAvailability)
		professionals.GET("/timetable", h.GetProfessionalTimetable)
		professionals.GET("/:id/clients", h.GetProfessionalClients)
		professionals.GET("/:id/previous_appointments", h.GetPreviousAppointmentsByClient)
	}

	return nil
}
