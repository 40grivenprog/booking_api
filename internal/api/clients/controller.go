package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	common "github.com/vention/booking_api/internal/api/common"
	db "github.com/vention/booking_api/internal/repository"
	"github.com/vention/booking_api/internal/services/clients"
	"github.com/vention/booking_api/internal/services/notifications"
)

// RegisterClient handles POST /api/clients/register
func (h *ClientsHandler) RegisterClient(c *gin.Context) {
	req, ok := common.BindAndValidate[ClientRegisterRequest](c)
	if !ok {
		return
	}

	client, err := h.clientsService.RegisterClient(c.Request.Context(), clients.RegisterClientInput{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		ChatID:    req.ChatID,
		Locale:    req.Locale,
	})
	if err != nil {
		if common.IsUniqueConstraintError(err) {
			common.HandleErrorResponse(c, http.StatusConflict, common.ErrorTypeConflict, common.ErrorMsgUsernameAlreadyExists, nil)
			return
		}
		common.HandleErrorResponse(c, http.StatusInternalServerError, common.ErrorTypeDatabase, common.ErrorMsgFailedToCreateClient, err)
		return
	}

	token, err := h.tokenMaker.CreateToken(client.ID, fmt.Sprintf("%s %s", client.FirstName, client.LastName))
	if err != nil {
		common.HandleErrorResponse(c, http.StatusInternalServerError, common.ErrorTypeInternal, common.ErrorMsgFailedToCreateToken, err)
		return
	}

	response := mapClientToClientRegisterResponse(client, token)
	c.JSON(http.StatusCreated, response)
}

// GetAllProfessionals handles GET /api/clients/professionals/all
func (h *ClientsHandler) GetAllProfessionals(c *gin.Context) {
	clientID, ok := common.GetUserID(c)
	if !ok {
		return
	}
	page := common.ParseIntQuery(c, "page", 1, 1, 10000)
	pageSize := common.ParseIntQuery(c, "pageSize", 15, 1, 100)

	offset := (page - 1) * pageSize

	dbParams := &db.GetProfessionalsParams{
		ClientID: clientID,
		Limit:    int32(pageSize),
		Offset:   int32(offset),
	}

	professionals, total, err := h.clientsService.GetAllProfessionals(c.Request.Context(), dbParams)
	if err != nil {
		common.HandleServiceError(c, err)
		return
	}

	response := mapProfessionalsToGetProfessionalsResponse(professionals, total, page, pageSize)

	c.JSON(http.StatusOK, response)
}

// GetSubscribedProfessionals handles GET /api/clients/professionals
func (h *ClientsHandler) GetSubscribedProfessionals(c *gin.Context) {
	clientID, ok := common.GetUserID(c)
	if !ok {
		return
	}

	page := common.ParseIntQuery(c, "page", 1, 1, 10000)
	pageSize := common.ParseIntQuery(c, "pageSize", 15, 1, 100)

	professionals, total, err := h.clientsService.GetSubscribedProfessionals(c.Request.Context(), clientID, page, pageSize)
	if err != nil {
		common.HandleServiceError(c, err)
		return
	}
	response := mapProfessionalsToGetSubscribedProfessionalsResponse(professionals, total, page, pageSize)

	c.JSON(http.StatusOK, response)
}

// SubscribeToProfessional handles POST /api/clients/subscribe
func (h *ClientsHandler) SubscribeToProfessional(c *gin.Context) {
	clientID, ok := common.GetUserID(c)
	if !ok {
		return
	}
	clientName := common.GetUserName(c)

	req, ok := common.BindAndValidate[SubscribedProfessionalRequest](c)
	if !ok {
		return
	}

	professionalID, ok := common.ParseProfessionalID(c, req.ProfessionalID)
	if !ok {
		return
	}

	err := h.clientsService.SubscribeToProfessional(c.Request.Context(), clients.SubscribeToProfessionalInput{
		ClientID:       clientID,
		ProfessionalID: professionalID,
	})
	if err != nil {
		common.HandleServiceError(c, err)
		return
	}
	err = h.notificationsService.SendSubscriptionNotification(c.Request.Context(), notifications.SendSubscriptionNotificationInput{
		ChatID:     req.ChatID,
		ClientName: clientName,
		Locale:     req.Locale,
	})
	if err != nil {
		common.HandleNotificationError(c, err)
		return
	}

	c.JSON(http.StatusAccepted, nil)
}

// BookAppointment handles POST /api/clients/book_appointment
func (h *ClientsHandler) BookAppointment(c *gin.Context) {
	clientID, ok := common.GetUserID(c)
	if !ok {
		return
	}

	clientName := common.GetUserName(c)
	req, ok := common.BindAndValidate[CreateAppointmentRequest](c)
	if !ok {
		return
	}

	startTime, ok := common.ParseTime(c, req.StartTime, common.ErrorMsgInvalidTime)
	if !ok {
		return
	}

	endTime, ok := common.ParseTime(c, req.EndTime, common.ErrorMsgInvalidTime)
	if !ok {
		return
	}

	professionalID, ok := common.ParseProfessionalID(c, req.ProfessionalID)
	if !ok {
		return
	}

	appointmentID, err := h.clientsService.CreateAppointment(c.Request.Context(), clients.CreateAppointmentInput{
		ClientID:       clientID,
		ProfessionalID: professionalID,
		StartTime:      startTime,
		EndTime:        endTime,
	})

	if err != nil {
		common.HandleServiceError(c, err)
		return
	}

	err = h.notificationsService.SendAppointmentRequestNotification(c.Request.Context(), notifications.SendAppointmentRequestNotificationInput{
		ChatID:        req.ProfessionalChatID,
		ClientName:    clientName,
		StartTime:     common.FormatTimeRFC3339(startTime),
		EndTime:       common.FormatTimeRFC3339(endTime),
		Locale:        req.ProfessionalLocale,
		AppointmentID: appointmentID,
	})
	if err != nil {
		common.HandleNotificationError(c, err)
		return
	}

	c.JSON(http.StatusCreated, nil)
}

// GetClientAppointments handles GET /api/clients/appointments
func (h *ClientsHandler) GetClientAppointments(c *gin.Context) {
	clientID, ok := common.GetUserID(c)
	if !ok {
		return
	}

	statusFilter := c.Query("status")
	if !common.ValidateAppointmentStatus(c, statusFilter) {
		return
	}

	// Parse pagination parameters: page (min 1, default 1) and pageSize (min 1, default 15)
	page := common.ParseIntQuery(c, "page", 1, 1, 10000)
	pageSize := common.ParseIntQuery(c, "pageSize", 15, 1, 100)

	appointments, total, err := h.clientsService.GetClientAppointments(c.Request.Context(), clientID, statusFilter, page, pageSize)
	if err != nil {
		common.HandleErrorResponse(c, http.StatusInternalServerError, common.ErrorTypeDatabase, common.ErrorMsgFailedToRetrieveAppointments, err)
		return
	}

	response := mapAppointmentToGetClientAppointmentsResponse(appointments, total, page, pageSize)
	c.JSON(http.StatusOK, response)
}

// CancelClientAppointment handles DELETE /api/clients/appointments/{appointment_id}/cancel
func (h *ClientsHandler) CancelClientAppointment(c *gin.Context) {
	clientID, ok := common.GetUserID(c)
	if !ok {
		return
	}
	clientName := common.GetUserName(c)

	appointmentID, ok := common.ParseAppointmentID(c, c.Param("appointment_id"))
	if !ok {
		return
	}

	req, ok := common.BindAndValidate[CancelClientAppointmentRequest](c)
	if !ok {
		return
	}

	result, err := h.clientsService.CancelAppointment(c.Request.Context(), clients.CancelAppointmentInput{
		ClientID:           clientID,
		AppointmentID:      appointmentID,
		CancellationReason: req.CancellationReason,
	})
	if err != nil {
		common.HandleServiceError(c, err)
		return
	}

	err = h.notificationsService.SendAppointmentCancellationNotification(c.Request.Context(), notifications.SendAppointmentCancellationNotificationInput{
		ChatID:             common.Int64Value(common.FromNullInt64(result.ProfessionalChatID)),
		StartTime:          common.FormatTimeRFC3339(result.StartTime),
		EndTime:            common.FormatTimeRFC3339(result.EndTime),
		RespondentName:     clientName,
		CancellationReason: req.CancellationReason,
		Type:               "client",
		Locale:             result.ProfessionalLocale,
	})

	c.JSON(http.StatusAccepted, nil)
}

// UpdateLocale handles PATCH /api/clients/update_locale
func (h *ClientsHandler) UpdateLocale(c *gin.Context) {
	clientID, ok := common.GetUserID(c)
	if !ok {
		return
	}

	req, ok := common.BindAndValidate[UpdateLocaleRequest](c)
	if !ok {
		return
	}

	err := h.clientsService.UpdateLocale(c.Request.Context(), clients.UpdateLocaleInput{
		ClientID: clientID,
		Locale:   req.Locale,
	})
	if err != nil {
		common.HandleServiceError(c, err)
		return
	}

	c.JSON(http.StatusAccepted, nil)
}

// UnsubscribeFromProfessional handles DELETE /api/clients/unsubscribe/{professional_id}
func (h *ClientsHandler) UnsubscribeFromProfessional(c *gin.Context) {
	clientID, ok := common.GetUserID(c)
	if !ok {
		return
	}

	professionalID, ok := common.ParseProfessionalID(c, c.Param("professional_id"))
	if !ok {
		return
	}

	err := h.clientsService.UnsubscribeFromProfessional(c.Request.Context(), clients.UnsubscribeFromProfessionalInput{
		ClientID:       clientID,
		ProfessionalID: professionalID,
	})

	if err != nil {
		common.HandleServiceError(c, err)
		return
	}

	c.JSON(http.StatusAccepted, nil)
}

// GetClientInvites handles GET /api/clients/invites
func (h *ClientsHandler) GetClientInvites(c *gin.Context) {
	clientID, ok := common.GetUserID(c)
	if !ok {
		return
	}

	invites, err := h.clientsService.GetClientInvites(c.Request.Context(), clientID)
	if err != nil {
		common.HandleServiceError(c, err)
		return
	}
}


