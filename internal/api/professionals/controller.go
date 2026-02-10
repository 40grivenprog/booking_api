package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	common "github.com/vention/booking_api/internal/api/common"
	"github.com/vention/booking_api/internal/services/notifications"
	"github.com/vention/booking_api/internal/services/professionals"
	"github.com/vention/booking_api/internal/util"
)

// SignInProfessional handles POST /api/professionals/sign_in
func (h *ProfessionalsHandler) SignInProfessional(c *gin.Context) {
	req, ok := common.BindAndValidate[ProfessionalSignInRequest](c)
	if !ok {
		return
	}

	professional, err := h.professionalsService.SignIn(c.Request.Context(), professionals.SignInInput{
		Username: req.Username,
		Password: req.Password,
		ChatID:   req.ChatID,
		Locale:   req.Locale,
	})
	if err != nil {
		common.HandleServiceError(c, err)
		return
	}

	token, err := h.tokenMaker.CreateToken(professional.ID, fmt.Sprintf("%s %s", professional.FirstName, professional.LastName))
	if err != nil {
		common.HandleErrorResponse(c, http.StatusInternalServerError, common.ErrorTypeInternal, common.ErrorMsgFailedToCreateToken, err)
		return
	}

	response := mapProfessionalToProfessionalSignInResponse(professional, token)
	c.JSON(http.StatusOK, response)
}

// ConfirmAppointment handles PATCH /api/professionals/appointments/{appointment_id}/confirm
func (h *ProfessionalsHandler) ConfirmAppointment(c *gin.Context) {
	professionalID, ok := common.GetUserID(c)
	if !ok {
		return
	}
	userName := common.GetUserName(c)

	appointmentID, ok := common.ParseAppointmentID(c, c.Param("appointment_id"))
	if !ok {
		return
	}

	result, err := h.professionalsService.ConfirmAppointment(c.Request.Context(), professionals.ConfirmAppointmentInput{
		ProfessionalID: professionalID,
		AppointmentID:  appointmentID,
	})
	if err != nil {
		common.HandleServiceError(c, err)
		return
	}

	for _, client := range result.ConfirmAppointmentClients {
		err = h.notificationsService.SendAppointmentConfirmationNotification(c.Request.Context(), notifications.SendAppointmentConfirmationNotificationInput{
			ChatID:           common.Int64Value(common.FromNullInt64(client.ClientChatID)),
			StartTime:        common.FormatTimeRFC3339(result.StartTime),
			EndTime:          common.FormatTimeRFC3339(result.EndTime),
			ProfessionalName: userName,
			Locale:           client.ClientLocale,
		})
	}

	c.JSON(http.StatusAccepted, nil)
}

// GetProfessionalAppointments handles GET /api/professionals/appointments
func (h *ProfessionalsHandler) GetProfessionalAppointments(c *gin.Context) {
	professionalID, ok := common.GetUserID(c)
	if !ok {
		return
	}

	statusFilter := c.Query("status")
	if !common.ValidateAppointmentStatus(c, statusFilter) {
		return
	}

	dateFilter := c.Query("date")
	if dateFilter != "" {
		if _, ok := common.ParseDate(c, dateFilter, common.ErrorMsgInvalidDate); !ok {
			return
		}
	}

	// Parse pagination parameters: page (min 1, default 1) and pageSize (min 1, default 15)
	page := common.ParseIntQuery(c, "page", 1, 1, 10000)
	pageSize := common.ParseIntQuery(c, "pageSize", 15, 1, 100)

	appointments, total, err := h.professionalsService.GetAppointments(c.Request.Context(), professionalID, statusFilter, dateFilter, page, pageSize)
	if err != nil {
		common.HandleErrorResponse(c, http.StatusInternalServerError, common.ErrorTypeDatabase, common.ErrorMsgFailedToRetrieveAppointments, err)
		return
	}

	response := mapAppointmentsToGetProfessionalAppointmentsResponse(appointments, total, page, pageSize)
	c.JSON(http.StatusOK, response)
}

// CancelAppointment handles PATCH /api/professionals/appointments/{appointment_id}/cancel
func (h *ProfessionalsHandler) CancelAppointment(c *gin.Context) {
	professionalID, ok := common.GetUserID(c)
	if !ok {
		return
	}
	userName := common.GetUserName(c)

	appointmentID, ok := common.ParseAppointmentID(c, c.Param("appointment_id"))
	if !ok {
		return
	}

	req, ok := common.BindAndValidate[CancelAppointmentRequest](c)
	if !ok {
		return
	}

	result, err := h.professionalsService.CancelAppointment(c.Request.Context(), professionals.CancelAppointmentInput{
		ProfessionalID:     professionalID,
		AppointmentID:      appointmentID,
		CancellationReason: req.CancellationReason,
	})
	if err != nil {
		common.HandleServiceError(c, err)
		return
	}

	if result.StartTime.After(time.Now()) {
		for _, client := range result.CancelAppointmentClients {
			err = h.notificationsService.SendAppointmentCancellationNotification(c.Request.Context(), notifications.SendAppointmentCancellationNotificationInput{
				ChatID:             common.Int64Value(common.FromNullInt64(client.ClientChatID)),
				StartTime:          common.FormatTimeRFC3339(result.StartTime),
				EndTime:            common.FormatTimeRFC3339(result.EndTime),
				RespondentName:     userName,
				Locale:             client.ClientLocale,
				CancellationReason: req.CancellationReason,
			})
		}
	}

	c.JSON(http.StatusAccepted, nil)
}

// // CreateUnavailableAppointment handles POST /api/professionals/unavailable_appointments
func (h *ProfessionalsHandler) CreateUnavailableAppointment(c *gin.Context) {
	professionalID, ok := common.GetUserID(c)
	if !ok {
		return
	}

	req, ok := common.BindAndValidate[CreateUnavailableAppointmentRequest](c)
	if !ok {
		return
	}

	startTime, ok := common.ParseTime(c, req.StartAt, common.ErrorMsgInvalidTime)
	if !ok {
		return
	}

	endTime, ok := common.ParseTime(c, req.EndAt, common.ErrorMsgInvalidTime)
	if !ok {
		return
	}

	err := h.professionalsService.CreateUnavailableAppointment(c.Request.Context(), professionals.CreateUnavailableAppointmentInput{
		ProfessionalID: professionalID,
		StartTime:      startTime,
		EndTime:        endTime,
		Description:    req.Description,
	})
	if err != nil {
		common.HandleServiceError(c, err)
		return
	}

	c.JSON(http.StatusAccepted, nil)
}

// GetProfessionalAvailability handles GET /api/professionals/:id/availability
func (h *ProfessionalsHandler) GetProfessionalAvailability(c *gin.Context) {
	professionalID, ok := common.ParseProfessionalID(c, c.Param("id"))
	if !ok {
		return
	}

	dateStr, ok := common.RequireQueryParam(c, "date")
	if !ok {
		return
	}

	date, ok := common.ParseDate(c, dateStr, common.ErrorMsgInvalidDate)
	if !ok {
		return
	}

	dateApp := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, util.GetAppTimezone())

	appointments, err := h.professionalsService.GetAvailability(c.Request.Context(), professionalID, dateApp)
	if err != nil {
		common.HandleErrorResponse(c, http.StatusInternalServerError, common.ErrorTypeDatabase, common.ErrorMsgFailedToRetrieveAppointments, err)
		return
	}

	// Generate availability slots using service
	slots := h.professionalsService.GenerateAvailabilitySlots(date, appointments, professionals.AvailabilityConfig{
		WorkingHoursStart: common.WorkingHoursStart,
		WorkingHoursEnd:   common.WorkingHoursEnd,
		AppTimezone:       util.GetAppTimezone(),
	})

	// Map service slots to response slots
	responseSlots := make([]TimeSlot, len(slots))
	for i, slot := range slots {
		responseSlots[i] = TimeSlot{
			StartTime: slot.StartTime,
			EndTime:   slot.EndTime,
			Available: slot.Available,
		}
	}

	response := GetProfessionalAvailabilityResponse{
		Date:  dateStr,
		Slots: responseSlots,
	}

	c.JSON(http.StatusOK, response)
}

// GetProfessionalTimetable handles GET /api/professionals/:id/timetable
func (h *ProfessionalsHandler) GetProfessionalTimetable(c *gin.Context) {
	professionalID, ok := common.GetUserID(c)
	if !ok {
		return
	}

	dateStr, ok := common.RequireQueryParam(c, "date")
	if !ok {
		return
	}

	date, ok := common.ParseDate(c, dateStr, common.ErrorMsgInvalidDate)
	if !ok {
		return
	}

	appointments, err := h.professionalsService.GetTimetable(c.Request.Context(), professionalID, date)
	if err != nil {
		common.HandleErrorResponse(c, http.StatusInternalServerError, common.ErrorTypeDatabase, common.ErrorMsgFailedToGetTimetable, err)
		return
	}

	response := mapTimetableAppointmentsToGetProfessionalTimetableResponse(appointments, dateStr)
	c.JSON(http.StatusOK, response)
}

// CreateGroupVisitAppointment handles POST /api/professionals/group_visit_appointments
func (h *ProfessionalsHandler) CreateGroupVisitAppointment(c *gin.Context) {
	professionalID, ok := common.GetUserID(c)
	if !ok {
		return
	}
	userName := common.GetUserName(c)

	req, ok := common.BindAndValidate[CreateGroupVisitAppointmentRequest](c)
	if !ok {
		return
	}

	startTime, ok := common.ParseTime(c, req.StartAt, common.ErrorMsgInvalidTime)
	if !ok {
		return
	}

	endTime, ok := common.ParseTime(c, req.EndAt, common.ErrorMsgInvalidTime)
	if !ok {
		return
	}

	appointmentID, err := h.professionalsService.CreateGroupVisitAppointment(c.Request.Context(), professionals.CreateGroupVisitAppointmentInput{
		ProfessionalID: professionalID,
		StartTime:      startTime,
		EndTime:        endTime,
		Description:    req.Description,
		Type:           req.Type,
	})
	if err != nil {
		common.HandleServiceError(c, err)
		return
	}

	if req.ClientsSelected == "all" {
		h.professionalsService.InviteAllClients(c.Request.Context(), professionals.InviteAllClientsInput{
			ProfessionalID:       professionalID,
			AppointmentID:        appointmentID,
			StartTime:            startTime,
			EndTime:              endTime,
			Description:          req.Description,
			Type:                 req.Type,
			ProfessionalName:     userName,
			NotificationsService: h.notificationsService,
		})
	} else {
		clientsInput := convertCreateGroupVisitAppointmentClientToInput(req.Clients)
		h.professionalsService.InvitePartiallySelectedClients(c.Request.Context(), professionals.InvitePartiallySelectedClientsInput{
			AppointmentID:        appointmentID,
			StartTime:            startTime,
			EndTime:              endTime,
			Description:          req.Description,
			Type:                 req.Type,
			Clients:              clientsInput,
			ProfessionalName:     userName,
			NotificationsService: h.notificationsService,
		})
	}

	c.JSON(http.StatusAccepted, nil)
}

// UpdateLocale handles PATCH /api/professionals/update_locale
func (h *ProfessionalsHandler) UpdateLocale(c *gin.Context) {
	professionalID, ok := common.GetUserID(c)
	if !ok {
		return
	}

	req, ok := common.BindAndValidate[UpdateLocaleRequest](c)
	if !ok {
		return
	}

	err := h.professionalsService.UpdateLocale(c.Request.Context(), professionals.UpdateLocaleInput{
		ProfessionalID: professionalID,
		Locale:         req.Locale,
	})
	if err != nil {
		common.HandleServiceError(c, err)
		return
	}

	c.JSON(http.StatusAccepted, nil)
}

// GetProfessionalSubscriptions handles GET /api/professionals/subscriptions
func (h *ProfessionalsHandler) GetProfessionalSubscriptions(c *gin.Context) {
	professionalID, ok := common.GetUserID(c)
	if !ok {
		return
	}
	subscriptions, err := h.professionalsService.GetSubscriptionsByProfessionalID(c.Request.Context(), professionalID)
	if err != nil {
		common.HandleServiceError(c, err)
		return
	}

	response := mapSubscriptionsToGetProfessionalSubscriptionsResponse(subscriptions)
	c.JSON(http.StatusOK, response)
}

// GetPreviousAppointments handles GET /api/professionals/previous_appointments
func (h *ProfessionalsHandler) GetPreviousAppointments(c *gin.Context) {
	professionalID, ok := common.GetUserID(c)
	if !ok {
		return
	}

	// Parse pagination parameters: page (min 1, default 1) and pageSize (min 1, default 15)
	page := common.ParseIntQuery(c, "page", 1, 1, 10000)
	pageSize := common.ParseIntQuery(c, "pageSize", 15, 1, 100)

	// Parse optional date range filters (YYYY-MM-DD format)
	// Only set dates if query parameters are provided
	var dateFrom, dateTo sql.NullTime
	if dateFromStr := c.Query("date_from"); dateFromStr != "" {
		t, err := time.Parse("2006-01-02", dateFromStr)
		if err != nil {
			common.HandleErrorResponse(c, http.StatusBadRequest, common.ErrorTypeValidation, "Invalid date_from format, expected YYYY-MM-DD", err)
			return
		}
		dateFrom = sql.NullTime{Time: t, Valid: true}
	}
	if dateToStr := c.Query("date_to"); dateToStr != "" {
		t, err := time.Parse("2006-01-02", dateToStr)
		if err != nil {
			common.HandleErrorResponse(c, http.StatusBadRequest, common.ErrorTypeValidation, "Invalid date_to format, expected YYYY-MM-DD", err)
			return
		}
		// Add one day so date_to is inclusive (end of the day)
		dateTo = sql.NullTime{Time: t.AddDate(0, 0, 1), Valid: true}
	}

	clientIDStr := c.Query("client_id")
	var response GetPreviousAppointmentsResponse

	if clientIDStr != "" {
		clientID, ok := common.ParseClientID(c, clientIDStr)
		if !ok {
			return
		}

		previousAppointments, total, err := h.professionalsService.GetPreviousAppointmentsByClientID(c.Request.Context(), professionalID, clientID, page, pageSize, dateFrom, dateTo)
		if err != nil {
			common.HandleServiceError(c, err)
			return
		}
		response = mapPreviousAppointmentsByClientIDToGetPreviousAppointmentsResponse(previousAppointments, total, page, pageSize)
	} else {
		previousAppointments, total, err := h.professionalsService.GetPreviousAppointments(c.Request.Context(), professionalID, page, pageSize, dateFrom, dateTo)
		if err != nil {
			common.HandleServiceError(c, err)
			return
		}

		response = mapPreviousAppointmentsToGetPreviousAppointmentsResponse(previousAppointments, total, page, pageSize)
	}

	c.JSON(http.StatusOK, response)
}

// GetAppointmentDetails handles GET /api/professionals/appointments/{appointment_id}
func (h *ProfessionalsHandler) GetAppointmentDetails(c *gin.Context) {
	professionalID, ok := common.GetUserID(c)
	if !ok {
		return
	}

	appointmentID, ok := common.ParseAppointmentID(c, c.Param("appointment_id"))
	if !ok {
		return
	}

	appointmentDetails, err := h.professionalsService.GetAppointmentDetailsByProfessionalIDAndAppointmentID(c.Request.Context(), professionalID, appointmentID)
	if err != nil {
		common.HandleServiceError(c, err)
		return
	}

	response := mapAppointmentDetailsToGetAppointmentDetailsResponse(appointmentDetails)
	c.JSON(http.StatusOK, response)
}

// UpdateAppointment handles PATCH /api/professionals/appointments/{appointment_id}/update
func (h *ProfessionalsHandler) UpdateAppointment(c *gin.Context) {
	professionalID, ok := common.GetUserID(c)
	if !ok {
		return
	}
	professionalName := common.GetUserName(c)

	req, ok := common.BindAndValidate[UpdateAppointmentRequest](c)
	if !ok {
		return
	}
	appointmentID, ok := common.ParseAppointmentID(c, c.Param("appointment_id"))
	if !ok {
		return
	}

	err := h.professionalsService.UpdateAppointment(c.Request.Context(), professionals.UpdateAppointmentInput{
		ProfessionalID:       professionalID,
		ProfessionalName:     professionalName,
		AppointmentID:        appointmentID,
		Description:          req.Description,
		Type:                 req.Type,
		NotificationsService: h.notificationsService,
	})
	if err != nil {
		common.HandleServiceError(c, err)
		return
	}

	c.JSON(http.StatusAccepted, nil)
}

// GetMissingInviteUsersForAppointment handles GET /api/professionals/appointments/{appointment_id}/missing_invite_users
func (h *ProfessionalsHandler) GetMissingInviteUsersForAppointment(c *gin.Context) {
	appointmentID, ok := common.ParseAppointmentID(c, c.Param("appointment_id"))
	if !ok {
		return
	}

	missingInviteUsers, err := h.professionalsService.GetMissingInviteUsersForAppointment(c.Request.Context(), appointmentID)
	if err != nil {
		common.HandleServiceError(c, err)
		return
	}
	response := mapMissingInviteUsersToGetMissingInviteUsersForAppointmentResponse(missingInviteUsers)

	c.JSON(http.StatusOK, response)
}

// GetPendingInviteUsersForAppointment handles GET /api/professionals/appointments/{appointment_id}/pending_invite_users
func (h *ProfessionalsHandler) GetPendingInviteUsersForAppointment(c *gin.Context) {
	appointmentID, ok := common.ParseAppointmentID(c, c.Param("appointment_id"))
	if !ok {
		return
	}

	pendingInviteUsers, err := h.professionalsService.GetPendingInviteUsersForAppointment(c.Request.Context(), appointmentID)
	if err != nil {
		common.HandleServiceError(c, err)
		return
	}
	response := mapPendingInviteUsersToGetPendingInviteUsersForAppointmentResponse(pendingInviteUsers)
	c.JSON(http.StatusOK, response)
}

// UpdatePreviousAppointment handles PATCH /api/professionals/previous_appointments/{appointment_id}/update
func (h *ProfessionalsHandler) UpdatePreviousAppointment(c *gin.Context) {
	professionalID, ok := common.GetUserID(c)
	if !ok {
		return
	}

	appointmentID, ok := common.ParseAppointmentID(c, c.Param("appointment_id"))
	if !ok {
		return
	}

	req, ok := common.BindAndValidate[UpdatePreviousAppointmentRequest](c)
	if !ok {
		return
	}

	clientsAdded := make([]uuid.UUID, len(req.ClientsAdded))
	for i, id := range req.ClientsAdded {
		parsed, err := uuid.Parse(id)
		if err != nil {
			common.HandleErrorResponse(c, http.StatusBadRequest, common.ErrorTypeValidation, "Invalid client ID in clients_added", err)
			return
		}
		clientsAdded[i] = parsed
	}

	clientsRemoved := make([]uuid.UUID, len(req.ClientsRemoved))
	for i, id := range req.ClientsRemoved {
		parsed, err := uuid.Parse(id)
		if err != nil {
			common.HandleErrorResponse(c, http.StatusBadRequest, common.ErrorTypeValidation, "Invalid client ID in clients_removed", err)
			return
		}
		clientsRemoved[i] = parsed
	}

	err := h.professionalsService.UpdatePreviousAppointment(c.Request.Context(), professionals.UpdatePreviousAppointmentInput{
		ProfessionalID: professionalID,
		AppointmentID:  appointmentID,
		Type:           req.Type,
		ClientsAdded:   clientsAdded,
		ClientsRemoved: clientsRemoved,
	})
	if err != nil {
		common.HandleServiceError(c, err)
		return
	}

	c.JSON(http.StatusAccepted, nil)
}

// GetMissingClientsForPreviousAppointment handles GET /api/professionals/previous_appointments/{appointment_id}/missing_clients
func (h *ProfessionalsHandler) GetMissingClientsForPreviousAppointment(c *gin.Context) {
	professionalID, ok := common.GetUserID(c)
	if !ok {
		return
	}

	appointmentID, ok := common.ParseAppointmentID(c, c.Param("appointment_id"))
	if !ok {
		return
	}

	missingClients, err := h.professionalsService.GetMissingClientsForPreviousAppointment(c.Request.Context(), professionalID, appointmentID)
	if err != nil {
		common.HandleServiceError(c, err)
		return
	}

	response := mapMissingClientsToGetMissingClientsForPreviousAppointmentResponse(missingClients)
	c.JSON(http.StatusOK, response)
}

// CreatePackage handles POST /api/professionals/packages
func (h *ProfessionalsHandler) CreatePackage(c *gin.Context) {
	professionalID, ok := common.GetUserID(c)
	if !ok {
		return
	}
	professionalName := common.GetUserName(c)

	req, ok := common.BindAndValidate[CreatePackageRequest](c)
	if !ok {
		return
	}

	issuedAt, ok := common.ParseTime(c, req.IssuedAt, common.ErrorMsgInvalidTime)
	if !ok {
		return
	}

	expiresAt, ok := common.ParseTime(c, req.ExpiresAt, common.ErrorMsgInvalidTime)
	if !ok {
		return
	}

	clientID, ok := common.ParseClientID(c, req.Client.ID)
	if !ok {
		return
	}

	err := h.professionalsService.CreatePackage(c.Request.Context(), professionals.CreatePackageInput{
		ProfessionalID:      professionalID,
		ClientID:            clientID,
		IssuedAt:            issuedAt,
		ExpiresAt:           expiresAt,
		ApppointmentsNumber: req.ApppointmentsNumber,
	})
	if err != nil {
		common.HandleServiceError(c, err)
		return
	}

	h.notificationsService.SendPackageCreatedNotification(c.Request.Context(), notifications.SendPackageCreatedNotificationInput{
		ChatID:              req.Client.ChatID,
		Locale:              req.Client.Locale,
		ProfessionalName:    professionalName,
		ApppointmentsNumber: req.ApppointmentsNumber,
		IssuedAt:            common.FormatTimeRFC3339(issuedAt),
		ExpiresAt:           common.FormatTimeRFC3339(expiresAt),
	})

	c.JSON(http.StatusAccepted, nil)
}

// GetListOfPackages handles GET /api/professionals/packages
func (h *ProfessionalsHandler) GetListOfPackages(c *gin.Context) {
	professionalID, ok := common.GetUserID(c)
	if !ok {
		return
	}

	packages, err := h.professionalsService.GetListOfPackages(c.Request.Context(), professionalID)
	if err != nil {
		common.HandleServiceError(c, err)
		return
	}

	response := mapPackagesToGetListOfPackagesResponse(packages)
	c.JSON(http.StatusOK, response)
}

// GetPackageDetails handles GET /api/professionals/packages/{package_id}
func (h *ProfessionalsHandler) GetPackageDetails(c *gin.Context) {
	packageID, ok := common.ParsePackageID(c, c.Param("package_id"))
	if !ok {
		return
	}

	packageDetails, err := h.professionalsService.GetPackageDetails(c.Request.Context(), packageID)
	if err != nil {
		common.HandleServiceError(c, err)
		return
	}

	response := mapPackageDetailsToGetPackageDetailsResponse(packageDetails)
	c.JSON(http.StatusOK, response)

}

// CreateInviteForAppointment handles POST /api/professionals/appointments/{appointment_id}/invite
func (h *ProfessionalsHandler) CreateInviteForAppointment(c *gin.Context) {
	professionalName := common.GetUserName(c)
	professionalID, ok := common.GetUserID(c)
	if !ok {
		return
	}

	appointmentID, ok := common.ParseAppointmentID(c, c.Param("appointment_id"))
	if !ok {
		return
	}

	req, ok := common.BindAndValidate[CreateInvitesForAppointmentRequest](c)
	if !ok {
		return
	}
	appointmentDetails, err := h.professionalsService.GetAppointmentDetailsByProfessionalIDAndAppointmentID(c.Request.Context(), professionalID, appointmentID)
	if err != nil {
		common.HandleServiceError(c, err)
		return
	}

	clientsInput := convertCreateInviteForAppointmentClientToInput(req.Clients)

	description := ""
	if appointmentDetails.Description.Valid {
		description = appointmentDetails.Description.String
	}

	h.professionalsService.InvitePartiallySelectedClients(c.Request.Context(), professionals.InvitePartiallySelectedClientsInput{
		AppointmentID:        appointmentID,
		StartTime:            appointmentDetails.StartTime,
		EndTime:              appointmentDetails.EndTime,
		Description:          description,
		Type:                 appointmentDetails.Type,
		Clients:              clientsInput,
		ProfessionalName:     professionalName,
		NotificationsService: h.notificationsService,
	})

	c.JSON(http.StatusAccepted, nil)
}
