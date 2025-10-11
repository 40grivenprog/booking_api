package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	common "github.com/vention/booking_api/internal/api/common"
	"github.com/vention/booking_api/internal/services/professionals"
	"github.com/vention/booking_api/internal/util"
)

// GetProfessionals handles GET /api/professionals
func (h *ProfessionalsHandler) GetProfessionals(c *gin.Context) {
	professionals, err := h.professionalsService.GetProfessionals(c.Request.Context())
	if err != nil {
		common.HandleErrorResponse(c, http.StatusInternalServerError, common.ErrorTypeDatabase, common.ErrorMsgFailedToRetrieveProfessionals, err)
		return
	}

	response := mapProfessionalsToGetProfessionalsResponse(professionals)
	c.JSON(http.StatusOK, response)
}

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
	})
	if err != nil {
		common.HandleServiceError(c, err)
		return
	}

	response := mapProfessionalToProfessionalSignInResponse(professional)
	c.JSON(http.StatusOK, response)
}

// ConfirmAppointment handles PATCH /api/professionals/{id}/appointments/{appointment_id}/confirm
func (h *ProfessionalsHandler) ConfirmAppointment(c *gin.Context) {
	professionalID, ok := common.ParseProfessionalID(c, c.Param("id"))
	if !ok {
		return
	}

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

	response := mapAppointmentToConfirmAppointmentResponse(result)
	c.JSON(http.StatusOK, response)
}

// GetProfessionalAppointments handles GET /api/professionals/{id}/appointments
func (h *ProfessionalsHandler) GetProfessionalAppointments(c *gin.Context) {
	professionalID, ok := common.ParseProfessionalID(c, c.Param("id"))
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

	appointments, err := h.professionalsService.GetAppointments(c.Request.Context(), professionalID, statusFilter, dateFilter)
	if err != nil {
		common.HandleErrorResponse(c, http.StatusInternalServerError, common.ErrorTypeDatabase, common.ErrorMsgFailedToRetrieveAppointments, err)
		return
	}

	response := mapAppointmentsToGetProfessionalAppointmentsResponse(appointments)
	c.JSON(http.StatusOK, response)
}

// GetProfessionalAppointmentDates handles GET /api/professionals/{id}/appointment-dates
func (h *ProfessionalsHandler) GetProfessionalAppointmentDates(c *gin.Context) {
	professionalID, ok := common.ParseProfessionalID(c, c.Param("id"))
	if !ok {
		return
	}

	monthStr := c.Query("month")
	var targetMonth time.Time
	if monthStr == "" {
		targetMonth = time.Now()
	} else {
		var ok bool
		targetMonth, ok = common.ParseMonth(c, monthStr)
		if !ok {
			return
		}
	}

	appointmentDates, err := h.professionalsService.GetAppointmentDates(c.Request.Context(), professionalID, targetMonth)
	if err != nil {
		common.HandleErrorResponse(c, http.StatusInternalServerError, common.ErrorTypeDatabase, common.ErrorMsgFailedToRetrieveAppointments, err)
		return
	}

	var dates []string
	for _, appointmentDate := range appointmentDates {
		dates = append(dates, common.FormatDate(appointmentDate))
	}

	response := GetProfessionalAppointmentDatesResponse{
		Month: monthStr,
		Dates: dates,
	}
	c.JSON(http.StatusOK, response)
}

// CancelAppointment handles PATCH /api/professionals/{id}/appointments/{appointment_id}/cancel
func (h *ProfessionalsHandler) CancelAppointment(c *gin.Context) {
	professionalID, ok := common.ParseProfessionalID(c, c.Param("id"))
	if !ok {
		return
	}

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

	response := mapAppointmentToCancelAppointmentResponse(result)
	c.JSON(http.StatusOK, response)
}

// CreateUnavailableAppointment handles POST /api/professionals/{id}/unavailable_appointments
func (h *ProfessionalsHandler) CreateUnavailableAppointment(c *gin.Context) {
	professionalID, ok := common.ParseProfessionalID(c, c.Param("id"))
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

	appointment, err := h.professionalsService.CreateUnavailableAppointment(c.Request.Context(), professionals.CreateUnavailableAppointmentInput{
		ProfessionalID: professionalID,
		StartTime:      startTime,
		EndTime:        endTime,
		Description:    req.Description,
	})
	if err != nil {
		common.HandleServiceError(c, err)
		return
	}

	response := mapAppointmentToCreateUnavailableAppointmentResponse(appointment)
	c.JSON(http.StatusCreated, response)
}

// GetProfessionalAvailability handles GET /api/professionals/{id}/availability
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
			StartTime:   slot.StartTime,
			EndTime:     slot.EndTime,
			Available:   slot.Available,
			Type:        slot.Type,
			Description: slot.Description,
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

	appointments, err := h.professionalsService.GetTimetable(c.Request.Context(), professionalID, date)
	if err != nil {
		common.HandleErrorResponse(c, http.StatusInternalServerError, common.ErrorTypeDatabase, common.ErrorMsgFailedToGetTimetable, err)
		return
	}

	response := mapTimetableAppointmentsToGetProfessionalTimetableResponse(appointments, dateStr)
	c.JSON(http.StatusOK, response)
}
