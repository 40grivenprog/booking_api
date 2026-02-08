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
		professionals.POST("/sign_in", h.SignInProfessional)
		professionals.GET("/appointments", h.GetProfessionalAppointments)
		professionals.PATCH("/appointments/:appointment_id/confirm", h.ConfirmAppointment)
		professionals.PATCH("/appointments/:appointment_id/cancel", h.CancelAppointment)
		professionals.POST("/unavailable_appointments", h.CreateUnavailableAppointment)
		professionals.GET("/:id/availability", h.GetProfessionalAvailability)
		professionals.GET("/timetable", h.GetProfessionalTimetable)
		professionals.POST("/group_visit_appointments", h.CreateGroupVisitAppointment)
		professionals.GET("/subscriptions", h.GetProfessionalSubscriptions)
		professionals.GET("/previous_appointments", h.GetPreviousAppointments)
		professionals.PATCH("/update_locale", h.UpdateLocale)
		professionals.GET("/appointments/:appointment_id", h.GetAppointmentDetails)
		professionals.PATCH("/appointments/:appointment_id/update", h.UpdateAppointment)
		professionals.PATCH("/previous_appointments/:appointment_id/update", h.UpdatePreviousAppointment)
		professionals.GET("/appointments/:appointment_id/missing_invite_users", h.GetMissingInviteUsersForAppointment)
		professionals.GET("/previous_appointments/:appointment_id/missing_clients", h.GetMissingClientsForPreviousAppointment)
		professionals.GET("/appointments/:appointment_id/pending_invite_users", h.GetPendingInviteUsersForAppointment)
		professionals.POST("/appointments/:appointment_id/invite", h.CreateInviteForAppointment)
		professionals.POST("/packages", h.CreatePackage)
		professionals.GET("/packages", h.GetListOfPackages)
		professionals.GET("/packages/:package_id", h.GetPackageDetails)
	}

	return nil
}
