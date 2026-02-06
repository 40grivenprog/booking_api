package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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
				ChatID:         common.Int64Value(common.FromNullInt64(client.ClientChatID)),
				StartTime:      common.FormatTimeRFC3339(result.StartTime),
				EndTime:        common.FormatTimeRFC3339(result.EndTime),
				RespondentName: userName,
				Locale:         client.ClientLocale,
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
		clientsInput := ConvertCreateGroupVisitAppointmentClientToInput(req.Clients)
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

	clientIDStr := c.Query("client_id")
	var response GetPreviousAppointmentsResponse

	if clientIDStr != "" {
		clientID, ok := common.ParseClientID(c, clientIDStr)
		if !ok {
			return
		}

		previousAppointments, total, err := h.professionalsService.GetPreviousAppointmentsByClientID(c.Request.Context(), professionalID, clientID, page, pageSize)
		if err != nil {
			common.HandleServiceError(c, err)
			return
		}
		response = mapPreviousAppointmentsByClientIDToGetPreviousAppointmentsResponse(previousAppointments, total, page, pageSize)
	} else {
		previousAppointments, total, err := h.professionalsService.GetPreviousAppointments(c.Request.Context(), professionalID, page, pageSize)
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

// // GetProfessionalAppointmentDates handles GET /api/professionals/{id}/appointment-dates
// func (h *ProfessionalsHandler) GetProfessionalAppointmentDates(c *gin.Context) {
// 	professionalID, ok := common.ParseProfessionalID(c, c.Param("id"))
// 	if !ok {
// 		return
// 	}

// 	monthStr := c.Query("month")
// 	var targetMonth time.Time
// 	if monthStr == "" {
// 		targetMonth = time.Now()
// 	} else {
// 		var ok bool
// 		targetMonth, ok = common.ParseMonth(c, monthStr)
// 		if !ok {
// 			return
// 		}
// 	}

// 	appointmentDates, err := h.professionalsService.GetAppointmentDates(c.Request.Context(), professionalID, targetMonth)
// 	if err != nil {
// 		common.HandleErrorResponse(c, http.StatusInternalServerError, common.ErrorTypeDatabase, common.ErrorMsgFailedToRetrieveAppointments, err)
// 		return
// 	}

// 	var dates []string
// 	for _, appointmentDate := range appointmentDates {
// 		dates = append(dates, common.FormatDate(appointmentDate))
// 	}

// 	response := GetProfessionalAppointmentDatesResponse{
// 		Month: monthStr,
// 		Dates: dates,
// 	}
// 	c.JSON(http.StatusOK, response)
// }
