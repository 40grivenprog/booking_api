package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	common "github.com/vention/booking_api/internal/api/common"
	db "github.com/vention/booking_api/internal/repository"
	"github.com/vention/booking_api/internal/util"
	"golang.org/x/crypto/bcrypt"
)

// GetProfessionals handles GET /api/professionals
func (h *ProfessionalsHandler) GetProfessionals(c *gin.Context) {
	professionals, err := h.professionalsRepo.GetProfessionals(c.Request.Context())
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

	// Get user by username
	user, err := h.professionalsRepo.GetProfessionalByUsername(c.Request.Context(), req.Username)
	if err != nil {
		common.HandleErrorResponse(c, http.StatusUnauthorized, common.ErrorTypeValidation, common.ErrorMsgInvalidCredentials, nil)
		return
	}

	// Check if password hash exists and verify password
	if !user.PasswordHash.Valid {
		common.HandleErrorResponse(c, http.StatusUnauthorized, common.ErrorTypeValidation, common.ErrorMsgInvalidCredentials, nil)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash.String), []byte(req.Password)); err != nil {
		common.HandleErrorResponse(c, http.StatusUnauthorized, common.ErrorTypeValidation, common.ErrorMsgInvalidCredentials, nil)
		return
	}

	// Update user's chat_id after successful authentication
	updatedUser, err := h.professionalsRepo.UpdateProfessionalChatID(c.Request.Context(), &db.UpdateProfessionalChatIDParams{
		ID:     user.ID,
		ChatID: common.ToNullInt64Value(req.ChatID),
	})
	if err != nil {
		common.HandleErrorResponse(c, http.StatusInternalServerError, common.ErrorTypeDatabase, common.ErrorMsgFailedToUpdateProfessional, err)
		return
	}

	response := mapProfessionalToProfessionalSignInResponse(updatedUser)
	c.JSON(http.StatusOK, response)
}

// ConfirmAppointment handles PATCH /api/professionals/{id}/appointments/{appointment_id}/confirm
func (h *ProfessionalsHandler) ConfirmAppointment(c *gin.Context) {
	// Parse UUIDs
	professionalID, ok := common.ParseProfessionalID(c, c.Param("id"))
	if !ok {
		return
	}

	appointmentID, ok := common.ParseAppointmentID(c, c.Param("appointment_id"))
	if !ok {
		return
	}

	// Get and validate appointment
	appointment, err := h.professionalsRepo.GetAppointmentByID(c.Request.Context(), appointmentID)
	if err != nil {
		common.HandleErrorResponse(c, http.StatusInternalServerError, common.ErrorTypeDatabase, common.ErrorMsgFailedToGetAppointment, err)
		return
	}

	// Validate ownership and status
	if !common.ValidateAppointmentOwnership(c, appointment, professionalID, common.UserTypeProfessional) {
		return
	}
	if !common.ValidateAppointmentStatusIs(c, appointment, db.AppointmentStatusPending) {
		return
	}

	// Confirm appointment with details
	result, err := h.professionalsRepo.ConfirmAppointmentWithDetails(c.Request.Context(), &db.ConfirmAppointmentWithDetailsParams{
		ID:             appointmentID,
		ProfessionalID: professionalID,
	})
	if err != nil {
		common.HandleErrorResponse(c, http.StatusInternalServerError, common.ErrorTypeDatabase, common.ErrorMsgFailedToUpdateAppointment, err)
		return
	}

	response := mapAppointmentToConfirmAppointmentResponse(result)

	c.JSON(http.StatusOK, response)
}

// GetProfessionalAppointments handles GET /api/professionals/{id}/appointments
func (h *ProfessionalsHandler) GetProfessionalAppointments(c *gin.Context) {
	// Parse professional ID
	professionalID, ok := common.ParseProfessionalID(c, c.Param("id"))
	if !ok {
		return
	}

	// Validate status filter if provided
	statusFilter := c.Query("status")
	if !common.ValidateAppointmentStatus(c, statusFilter) {
		return
	}

	// Validate date filter if provided
	dateFilter := c.Query("date")
	var dateParam string
	if dateFilter != "" {
		if _, ok := common.ParseDate(c, dateFilter, common.ErrorMsgInvalidDate); !ok {
			return
		}
		dateParam = dateFilter
	}

	// Get appointments with optional status and date filters
	appointments, err := h.professionalsRepo.GetAppointmentsByProfessionalWithStatusAndDate(c.Request.Context(), &db.GetAppointmentsByProfessionalWithStatusAndDateParams{
		ProfessionalID: professionalID,
		Column2:        statusFilter,
		Column3:        dateParam,
	})
	if err != nil {
		common.HandleErrorResponse(c, http.StatusInternalServerError, common.ErrorTypeDatabase, common.ErrorMsgFailedToRetrieveAppointments, err)
		return
	}

	// Convert to response format
	response := mapAppointmentsToGetProfessionalAppointmentsResponse(appointments)

	c.JSON(http.StatusOK, response)
}

// GetProfessionalAppointmentDates handles GET /api/professionals/{id}/appointment-dates
func (h *ProfessionalsHandler) GetProfessionalAppointmentDates(c *gin.Context) {
	// Parse professional ID
	professionalID, ok := common.ParseProfessionalID(c, c.Param("id"))
	if !ok {
		return
	}

	// Parse month parameter
	monthStr := c.Query("month") // Format: "2025-09"
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

	// Get start and end of month
	startOfMonth := time.Date(targetMonth.Year(), targetMonth.Month(), 1, 0, 0, 0, 0, util.GetAppTimezone())
	endOfMonth := startOfMonth.AddDate(0, 1, 0) // First day of next month

	// Get distinct dates with confirmed appointments for the month
	appointmentDates, err := h.professionalsRepo.GetProfessionalAppointmentDates(c.Request.Context(), &db.GetProfessionalAppointmentDatesParams{
		ProfessionalID: professionalID,
		StartTime:      startOfMonth,
		StartTime_2:    endOfMonth,
	})
	if err != nil {
		common.HandleErrorResponse(c, http.StatusInternalServerError, common.ErrorTypeDatabase, common.ErrorMsgFailedToRetrieveAppointments, err)
		return
	}

	// Convert to string slice
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
	// Parse and validate IDs
	professionalID, ok := common.ParseProfessionalID(c, c.Param("id"))
	if !ok {
		return
	}

	appointmentID, ok := common.ParseAppointmentID(c, c.Param("appointment_id"))
	if !ok {
		return
	}

	// Parse request body
	req, ok := common.BindAndValidate[CancelAppointmentRequest](c)
	if !ok {
		return
	}

	appointment, err := h.professionalsRepo.GetAppointmentByID(c.Request.Context(), appointmentID)
	if err != nil {
		common.HandleErrorResponse(c, http.StatusInternalServerError, common.ErrorTypeDatabase, common.ErrorMsgFailedToGetAppointment, err)
		return
	}
	if appointment.ProfessionalID != professionalID {
		common.HandleErrorResponse(c, http.StatusForbidden, common.ErrorTypeForbidden, common.ErrorMsgNotAllowedToConfirmAppointment, nil)
		return
	}
	if appointment.Status.AppointmentStatus != db.AppointmentStatusPending && appointment.Status.AppointmentStatus != db.AppointmentStatusConfirmed {
		common.HandleErrorResponse(c, http.StatusBadRequest, common.ErrorTypeValidation, common.ErrorMsgAppointmentNotPendingOrConfirmed, nil)
		return
	}

	// Cancel appointment with details
	result, err := h.professionalsRepo.CancelAppointmentByProfessionalWithDetails(c.Request.Context(), &db.CancelAppointmentByProfessionalWithDetailsParams{
		ID:                        appointmentID,
		CancelledByProfessionalID: common.ToNullUUID(professionalID),
		CancellationReason:        common.ToNullStringValue(req.CancellationReason),
	})
	if err != nil {
		common.HandleErrorResponse(c, http.StatusInternalServerError, common.ErrorTypeDatabase, common.ErrorMsgFailedToUpdateAppointment, err)
		return
	}

	// Convert to response format
	response := mapAppointmentToCancelAppointmentResponse(result)

	c.JSON(http.StatusOK, response)
}

// CreateUnavailableAppointment handles POST /api/professionals/{id}/unavailable_appointments
func (h *ProfessionalsHandler) CreateUnavailableAppointment(c *gin.Context) {
	// Parse and validate professional ID
	professionalID, ok := common.ParseProfessionalID(c, c.Param("id"))
	if !ok {
		return
	}

	// Parse request body
	req, ok := common.BindAndValidate[CreateUnavailableAppointmentRequest](c)
	if !ok {
		return
	}

	// Parse start and end times
	startTime, err := time.Parse(time.RFC3339, req.StartAt)
	if err != nil {
		common.HandleErrorResponse(c, http.StatusBadRequest, common.ErrorTypeValidation, common.ErrorMsgInvalidTime, err)
		return
	}

	endTime, err := time.Parse(time.RFC3339, req.EndAt)
	if err != nil {
		common.HandleErrorResponse(c, http.StatusBadRequest, common.ErrorTypeValidation, common.ErrorMsgInvalidTime, err)
		return
	}

	// Validate that start time is in the future
	if startTime.Before(time.Now()) {
		common.HandleErrorResponse(c, http.StatusBadRequest, common.ErrorTypeValidation, common.ErrorMsgFutureTimeRequired, nil)
		return
	}

	// Validate that end time is after start time
	if endTime.Before(startTime) || endTime.Equal(startTime) {
		common.HandleErrorResponse(c, http.StatusBadRequest, common.ErrorTypeValidation, common.ErrorMsgInvalidTime, nil)
		return
	}

	// Create unavailable appointment
	appointment, err := h.professionalsRepo.CreateUnavailableAppointment(c.Request.Context(), &db.CreateUnavailableAppointmentParams{
		ProfessionalID: professionalID,
		Description:    common.ToNullStringValue(req.Description),
		StartTime:      startTime,
		EndTime:        endTime,
	})
	if err != nil {
		common.HandleErrorResponse(c, http.StatusInternalServerError, common.ErrorTypeDatabase, common.ErrorMsgFailedToCreateAppointment, err)
		return
	}

	// Convert to response format
	response := mapAppointmentToCreateUnavailableAppointmentResponse(appointment)

	c.JSON(http.StatusCreated, response)
}

// GetProfessionalAvailability handles GET /api/professionals/{id}/availability
func (h *ProfessionalsHandler) GetProfessionalAvailability(c *gin.Context) {
	// Parse and validate professional ID
	professionalID, ok := common.ParseProfessionalID(c, c.Param("id"))
	if !ok {
		return
	}

	// Get and validate date parameter
	dateStr, ok := common.RequireQueryParam(c, "date")
	if !ok {
		return
	}

	// Parse date
	date, ok := common.ParseDate(c, dateStr, common.ErrorMsgInvalidDate)
	if !ok {
		return
	}

	// Get appointments for the specific date with client information
	// Convert the date to application timezone
	dateApp := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, util.GetAppTimezone())

	appointments, err := h.professionalsRepo.GetAppointmentsByProfessionalAndDateWithClient(c.Request.Context(), &db.GetAppointmentsByProfessionalAndDateWithClientParams{
		ProfessionalID: professionalID,
		StartTime:      dateApp,
	})
	if err != nil {
		common.HandleErrorResponse(c, http.StatusInternalServerError, common.ErrorTypeDatabase, common.ErrorMsgFailedToRetrieveAppointments, err)
		return
	}

	// Generate time slots from current time to 23:00
	slots := make([]TimeSlot, 0, 18)

	// Use centralized timezone
	localNow := util.NowInAppTimezone()

	// Create base date in application timezone
	baseDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, util.GetAppTimezone())

	// Determine the starting hour based on current time
	startHour := common.WorkingHoursStart // Default start hour
	if date.Year() == localNow.Year() && date.Month() == localNow.Month() && date.Day() == localNow.Day() {
		// If it's today, start from current hour (rounded up to next hour)
		currentHour := localNow.Hour()
		if localNow.Minute() > 0 {
			currentHour++ // Round up to next hour if we're past the hour mark
		}
		if currentHour > common.WorkingHoursStart {
			startHour = currentHour
		}
	}

	for hour := startHour; hour < common.WorkingHoursEnd; hour++ {
		startTime := baseDate.Add(time.Duration(hour) * time.Hour)
		endTime := startTime.Add(time.Hour)

		// Skip if the slot is in the past (additional safety check)
		if startTime.Before(localNow) {
			continue
		}

		slot := TimeSlot{
			StartTime: common.FormatTimeRFC3339(startTime),
			EndTime:   common.FormatTimeRFC3339(endTime),
			Available: true,
		}

		// Check if this slot conflicts with any existing appointment
		for _, appointment := range appointments {
			apptStart := appointment.StartTime
			apptEnd := appointment.EndTime

			// Convert appointment times to application timezone for comparison
			apptStartLocal := apptStart.In(util.GetAppTimezone())
			apptEndLocal := apptEnd.In(util.GetAppTimezone())

			// Check if the slot overlaps with the appointment (both in application timezone)
			if startTime.Before(apptEndLocal) && endTime.After(apptStartLocal) {
				slot.Available = false
				slot.Type = string(appointment.Type)

				// Generate description based on appointment type and client info
				if appointment.Description.Valid {
					if appointment.ClientID.Valid && appointment.ClientFirstName.Valid && appointment.ClientLastName.Valid {
						// Appointment with client - show client info + description
						slot.Description = fmt.Sprintf("%s %s - %s",
							appointment.ClientFirstName.String,
							appointment.ClientLastName.String,
							appointment.Description.String)
					} else {
						// Unavailable period - show just description
						slot.Description = appointment.Description.String
					}
				}
				break
			}
		}

		slots = append(slots, slot)
	}

	response := GetProfessionalAvailabilityResponse{
		Date:  dateStr,
		Slots: slots,
	}

	c.JSON(http.StatusOK, response)
}

// GetProfessionalTimetable handles GET /api/professionals/:id/timetable
func (h *ProfessionalsHandler) GetProfessionalTimetable(c *gin.Context) {
	// Parse and validate professional ID
	professionalID, ok := common.ParseProfessionalID(c, c.Param("id"))
	if !ok {
		return
	}

	// Get and validate date parameter
	dateStr, ok := common.RequireQueryParam(c, "date")
	if !ok {
		return
	}

	// Parse date
	date, ok := common.ParseDate(c, dateStr, common.ErrorMsgInvalidDate)
	if !ok {
		return
	}

	// Get appointments for the day
	appointments, err := h.professionalsRepo.GetProfessionalTimetable(c.Request.Context(), &db.GetProfessionalTimetableParams{
		ProfessionalID: professionalID,
		StartTime:      date,
	})
	if err != nil {
		common.HandleErrorResponse(c, http.StatusInternalServerError, common.ErrorTypeDatabase, common.ErrorMsgFailedToGetTimetable, err)
		return
	}

	// Convert to response format
	response := mapTimetableAppointmentsToGetProfessionalTimetableResponse(appointments, dateStr)

	c.JSON(http.StatusOK, response)
}
