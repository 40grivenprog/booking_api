package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	common "github.com/vention/booking_api/internal/api/common"
	db "github.com/vention/booking_api/internal/repository"
	"github.com/vention/booking_api/internal/util"
	"golang.org/x/crypto/bcrypt"
)

// GetProfessionals handles GET /api/professionals
func (h *ProfessionalsHandler) GetProfessionals(c *gin.Context) {
	professionals, err := h.professionalsRepo.GetProfessionals(c.Request.Context())
	if err != nil {
		common.HandleErrorResponse(c, http.StatusInternalServerError, "database_error", "Failed to retrieve professionals", err)
		return
	}

	// Convert to response format
	var responseUsers []User
	for _, prof := range professionals {
		user := User{
			ID:        prof.ID.String(),
			Username:  prof.Username,
			FirstName: prof.FirstName,
			LastName:  prof.LastName,
			UserType:  "professional",
			CreatedAt: prof.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: prof.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}

		// Handle optional fields
		if prof.ChatID.Valid {
			user.ChatID = &prof.ChatID.Int64
		}
		if prof.PhoneNumber.Valid {
			user.PhoneNumber = &prof.PhoneNumber.String
		}

		responseUsers = append(responseUsers, user)
	}

	response := GetProfessionalsResponse{
		Professionals: responseUsers,
	}

	c.JSON(http.StatusOK, response)
}

// SignInProfessional handles POST /api/professionals/sign_in
func (h *ProfessionalsHandler) SignInProfessional(c *gin.Context) {
	var req ProfessionalSignInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.HandleErrorResponse(c, http.StatusBadRequest, "validation_error", "Invalid request body", err)
		return
	}

	// Get user by username
	user, err := h.professionalsRepo.GetProfessionalByUsername(c.Request.Context(), req.Username)
	if err != nil {
		common.HandleErrorResponse(c, http.StatusUnauthorized, "invalid_credentials", "Username or password is incorrect", nil)
		return
	}

	// Check if password hash exists
	if !user.PasswordHash.Valid {
		common.HandleErrorResponse(c, http.StatusUnauthorized, "invalid_credentials", "Username or password is incorrect", nil)
		return
	}

	// Verify password using bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash.String), []byte(req.Password))
	if err != nil {
		common.HandleErrorResponse(c, http.StatusUnauthorized, "invalid_credentials", "Username or password is incorrect", nil)
		return
	}

	// Update user's chat_id after successful authentication
	updatedUser, err := h.professionalsRepo.UpdateProfessionalChatID(c.Request.Context(), &db.UpdateProfessionalChatIDParams{
		ID:     user.ID,
		ChatID: sql.NullInt64{Int64: req.ChatID, Valid: true},
	})
	if err != nil {
		common.HandleErrorResponse(c, http.StatusInternalServerError, "database_error", "Failed to update user chat_id", err)
		return
	}

	// Convert to response format
	responseUser := User{
		ID:        updatedUser.ID.String(),
		Username:  updatedUser.Username,
		FirstName: updatedUser.FirstName,
		LastName:  updatedUser.LastName,
		Role:      "professional",
		CreatedAt: updatedUser.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: updatedUser.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	// Handle optional fields
	if updatedUser.ChatID.Valid {
		responseUser.ChatID = &updatedUser.ChatID.Int64
	}
	if updatedUser.PhoneNumber.Valid {
		responseUser.PhoneNumber = &updatedUser.PhoneNumber.String
	}

	response := ProfessionalSignInResponse{
		User: responseUser,
	}

	c.JSON(http.StatusOK, response)
}

// ConfirmAppointment handles PATCH /api/professionals/{id}/appointments/{appointment_id}/confirm
func (h *ProfessionalsHandler) ConfirmAppointment(c *gin.Context) {
	professionalIDStr := c.Param("id")
	appointmentIDStr := c.Param("appointment_id")

	// Parse UUIDs
	professionalID, err := uuid.Parse(professionalIDStr)
	if err != nil {
		common.HandleErrorResponse(c, http.StatusBadRequest, "validation_error", "Invalid professional_id format", err)
		return
	}

	appointmentID, err := uuid.Parse(appointmentIDStr)
	if err != nil {
		common.HandleErrorResponse(c, http.StatusBadRequest, "validation_error", "Invalid appointment_id format", err)
		return
	}

	appointment, err := h.professionalsRepo.GetAppointmentByID(c.Request.Context(), appointmentID)
	if err != nil {
		common.HandleErrorResponse(c, http.StatusInternalServerError, "database_error", "Failed to get appointment", err)
		return
	}
	if appointment.ProfessionalID != professionalID {
		common.HandleErrorResponse(c, http.StatusForbidden, "forbidden", "You are not allowed to confirm this appointment", nil)
		return
	}
	if appointment.Status.AppointmentStatus != db.AppointmentStatusPending {
		common.HandleErrorResponse(c, http.StatusBadRequest, "validation_error", "Appointment is not pending", nil)
		return
	}

	// Confirm appointment with details
	result, err := h.professionalsRepo.ConfirmAppointmentWithDetails(c.Request.Context(), &db.ConfirmAppointmentWithDetailsParams{
		ID:             appointmentID,
		ProfessionalID: professionalID,
	})
	if err != nil {
		common.HandleErrorResponse(c, http.StatusInternalServerError, "database_error", "Failed to confirm appointment", err)
		return
	}

	// Convert to response format
	response := ConfirmAppointmentResponse{
		Appointment: AppointmentConfirm{
			ID:        result.ID.String(),
			Status:    string(result.Status.AppointmentStatus),
			StartTime: result.StartTime.Format(time.RFC3339),
			EndTime:   result.EndTime.Format(time.RFC3339),
			CreatedAt: result.CreatedAt.Format(time.RFC3339),
			UpdatedAt: result.UpdatedAt.Format(time.RFC3339),
		},
		Client: ClientConfirm{
			ID:        result.ClientID.UUID.String(),
			FirstName: result.ClientFirstName.String,
			LastName:  result.ClientLastName.String,
			ChatID:    result.ClientChatID.Int64,
		},
		Professional: ProfessionalConfirm{
			ID:        result.ProfessionalIDFull.String(),
			Username:  result.ProfessionalUsername.String,
			FirstName: result.ProfessionalFirstName.String,
			LastName:  result.ProfessionalLastName.String,
		},
	}

	c.JSON(http.StatusOK, response)
}

// GetProfessionalAppointments handles GET /api/professionals/{id}/appointments
func (h *ProfessionalsHandler) GetProfessionalAppointments(c *gin.Context) {
	professionalIDStr := c.Param("id")
	statusFilter := c.Query("status")

	// Parse professional ID
	professionalID, err := uuid.Parse(professionalIDStr)
	if err != nil {
		common.HandleErrorResponse(c, http.StatusBadRequest, "validation_error", "Invalid professional_id format", err)
		return
	}

	// Validate status filter if provided
	if statusFilter != "" && statusFilter != "pending" && statusFilter != "confirmed" && statusFilter != "cancelled" && statusFilter != "completed" {
		common.HandleErrorResponse(c, http.StatusBadRequest, "validation_error", "Invalid status. Must be one of: pending, confirmed, cancelled, completed", nil)
		return
	}

	// Prepare parameters for the query
	var statusParam db.NullAppointmentStatus
	if statusFilter != "" {
		statusParam = db.NullAppointmentStatus{AppointmentStatus: db.AppointmentStatus(statusFilter), Valid: true}
	}

	// Get appointments with optional status filter
	appointments, err := h.professionalsRepo.GetAppointmentsByProfessionalWithStatus(c.Request.Context(), &db.GetAppointmentsByProfessionalWithStatusParams{
		ProfessionalID: professionalID,
		Status:         statusParam,
	})
	if err != nil {
		common.HandleErrorResponse(c, http.StatusInternalServerError, "database_error", "Failed to retrieve appointments", err)
		return
	}

	// Convert to response format
	var responseAppointments []ProfessionalAppointment
	for _, appt := range appointments {
		appointment := ProfessionalAppointment{
			ID:        appt.ID.String(),
			Type:      string(appt.Type),
			StartTime: appt.StartTime.Format(time.RFC3339),
			EndTime:   appt.EndTime.Format(time.RFC3339),
			Status:    string(appt.Status.AppointmentStatus),
			CreatedAt: appt.CreatedAt.Format(time.RFC3339),
			UpdatedAt: appt.UpdatedAt.Format(time.RFC3339),
		}
		appointment.Client = &ProfessionalAppointmentClient{
			ID:          appt.ClientID.UUID.String(),
			FirstName:   appt.ClientFirstName.String,
			LastName:    appt.ClientLastName.String,
			PhoneNumber: &appt.ClientPhoneNumber.String,
		}

		responseAppointments = append(responseAppointments, appointment)
	}

	response := GetProfessionalAppointmentsResponse{
		Appointments: responseAppointments,
	}

	c.JSON(http.StatusOK, response)
}

// CancelAppointment handles PATCH /api/professionals/{id}/appointments/{appointment_id}/cancel
func (h *ProfessionalsHandler) CancelAppointment(c *gin.Context) {
	professionalIDStr := c.Param("id")
	appointmentIDStr := c.Param("appointment_id")

	// Parse UUIDs
	professionalID, err := uuid.Parse(professionalIDStr)
	if err != nil {
		common.HandleErrorResponse(c, http.StatusBadRequest, "validation_error", "Invalid professional_id format", err)
		return
	}

	appointmentID, err := uuid.Parse(appointmentIDStr)
	if err != nil {
		common.HandleErrorResponse(c, http.StatusBadRequest, "validation_error", "Invalid appointment_id format", err)
		return
	}

	// Parse request body
	var req CancelAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.HandleErrorResponse(c, http.StatusBadRequest, "validation_error", "Invalid request body", err)
		return
	}

	appointment, err := h.professionalsRepo.GetAppointmentByID(c.Request.Context(), appointmentID)
	if err != nil {
		common.HandleErrorResponse(c, http.StatusInternalServerError, "database_error", "Failed to get appointment", err)
		return
	}
	if appointment.ProfessionalID != professionalID {
		common.HandleErrorResponse(c, http.StatusForbidden, "forbidden", "You are not allowed to confirm this appointment", nil)
		return
	}
	if appointment.Status.AppointmentStatus != db.AppointmentStatusPending && appointment.Status.AppointmentStatus != db.AppointmentStatusConfirmed {
		common.HandleErrorResponse(c, http.StatusBadRequest, "validation_error", "Appointment is not pending or confirmed. Please check the status of the appointment.", nil)
		return
	}

	// Cancel appointment with details
	result, err := h.professionalsRepo.CancelAppointmentByProfessionalWithDetails(c.Request.Context(), &db.CancelAppointmentByProfessionalWithDetailsParams{
		ID:                        appointmentID,
		CancelledByProfessionalID: uuid.NullUUID{UUID: professionalID, Valid: true},
		CancellationReason:        sql.NullString{String: req.CancellationReason, Valid: true},
	})
	if err != nil {
		common.HandleErrorResponse(c, http.StatusInternalServerError, "database_error", "Failed to cancel appointment", err)
		return
	}

	// Convert to response format
	response := CancelAppointmentResponse{
		Appointment: CancelledAppointment{
			ID:                 result.ID.String(),
			Type:               string(result.Type),
			StartTime:          result.StartTime.Format(time.RFC3339),
			EndTime:            result.EndTime.Format(time.RFC3339),
			Status:             string(result.Status.AppointmentStatus),
			CancellationReason: result.CancellationReason.String,
			CancelledBy:        "professional",
			CreatedAt:          result.CreatedAt.Format(time.RFC3339),
			UpdatedAt:          result.UpdatedAt.Format(time.RFC3339),
		},
		Client: ProfessionalAppointmentClient{
			ID:        result.ClientIDFull.String(),
			FirstName: result.ClientFirstName.String,
			LastName:  result.ClientLastName.String,
			ChatID:    &result.ClientChatID.Int64,
		},
		Professional: ProfessionalInfo{
			ID:        result.ProfessionalIDFull.String(),
			Username:  result.ProfessionalUsername.String,
			FirstName: result.ProfessionalFirstName.String,
			LastName:  result.ProfessionalLastName.String,
		},
	}

	c.JSON(http.StatusOK, response)
}

// CreateUnavailableAppointment handles POST /api/professionals/{id}/unavailable_appointments
func (h *ProfessionalsHandler) CreateUnavailableAppointment(c *gin.Context) {
	professionalIDStr := c.Param("id")

	// Parse professional ID
	professionalID, err := uuid.Parse(professionalIDStr)
	if err != nil {
		common.HandleErrorResponse(c, http.StatusBadRequest, "validation_error", "Invalid professional_id format", err)
		return
	}

	// Parse request body
	var req CreateUnavailableAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.HandleErrorResponse(c, http.StatusBadRequest, "validation_error", "Invalid request body", err)
		return
	}

	// Parse start and end times
	startTime, err := time.Parse(time.RFC3339, req.StartAt)
	if err != nil {
		common.HandleErrorResponse(c, http.StatusBadRequest, "validation_error", "Invalid start_at format. Use RFC3339 format (e.g., 2024-01-15T10:00:00Z)", err)
		return
	}

	endTime, err := time.Parse(time.RFC3339, req.EndAt)
	if err != nil {
		common.HandleErrorResponse(c, http.StatusBadRequest, "validation_error", "Invalid end_at format. Use RFC3339 format (e.g., 2024-01-15T11:00:00Z)", err)
		return
	}

	// Validate that start time is in the future
	if startTime.Before(time.Now()) {
		common.HandleErrorResponse(c, http.StatusBadRequest, "validation_error", "start_at must be in the future", nil)
		return
	}

	// Validate that end time is after start time
	if endTime.Before(startTime) || endTime.Equal(startTime) {
		common.HandleErrorResponse(c, http.StatusBadRequest, "validation_error", "end_at must be after start_at", nil)
		return
	}

	// Create unavailable appointment
	appointment, err := h.professionalsRepo.CreateUnavailableAppointment(c.Request.Context(), &db.CreateUnavailableAppointmentParams{
		ProfessionalID: professionalID,
		StartTime:      startTime,
		EndTime:        endTime,
	})
	if err != nil {
		common.HandleErrorResponse(c, http.StatusInternalServerError, "database_error", "Failed to create unavailable appointment", err)
		return
	}

	// Convert to response format
	response := CreateUnavailableAppointmentResponse{
		Appointment: UnavailableAppointment{
			ID:        appointment.ID.String(),
			Type:      string(appointment.Type),
			StartTime: appointment.StartTime.Format(time.RFC3339),
			EndTime:   appointment.EndTime.Format(time.RFC3339),
			Status:    string(appointment.Status.AppointmentStatus),
			CreatedAt: appointment.CreatedAt.Format(time.RFC3339),
			UpdatedAt: appointment.UpdatedAt.Format(time.RFC3339),
		},
	}

	c.JSON(http.StatusCreated, response)
}

// GetProfessionalAvailability handles GET /api/professionals/{id}/availability
func (h *ProfessionalsHandler) GetProfessionalAvailability(c *gin.Context) {
	professionalIDStr := c.Param("id")
	dateStr := c.Query("date")

	// Parse professional ID
	professionalID, err := uuid.Parse(professionalIDStr)
	if err != nil {
		common.HandleErrorResponse(c, http.StatusBadRequest, "validation_error", "Invalid professional_id format", err)
		return
	}

	// Validate and parse date
	if dateStr == "" {
		common.HandleErrorResponse(c, http.StatusBadRequest, "validation_error", "date query parameter is required", nil)
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		common.HandleErrorResponse(c, http.StatusBadRequest, "validation_error", "Invalid date format. Use YYYY-MM-DD format", err)
		return
	}

	// Get appointments for the specific date
	appointments, err := h.professionalsRepo.GetAppointmentsByProfessionalAndDate(c.Request.Context(), &db.GetAppointmentsByProfessionalAndDateParams{
		ProfessionalID: professionalID,
		StartTime:      date,
	})
	if err != nil {
		common.HandleErrorResponse(c, http.StatusInternalServerError, "database_error", "Failed to retrieve appointments", err)
		return
	}

	// Generate time slots from current time to 23:00
	slots := make([]TimeSlot, 0, 18)
	now := time.Now()
	fmt.Println("now UTC", now)

	// Use centralized timezone
	localNow := util.NowInAppTimezone()
	fmt.Println("now local", localNow)

	// Create base date in application timezone
	baseDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, util.GetAppTimezone())

	// Determine the starting hour based on current time
	startHour := 5 // Default start hour
	if date.Year() == localNow.Year() && date.Month() == localNow.Month() && date.Day() == localNow.Day() {
		// If it's today, start from current hour (rounded up to next hour)
		currentHour := localNow.Hour()
		if localNow.Minute() > 0 {
			currentHour++ // Round up to next hour if we're past the hour mark
		}
		if currentHour > 5 {
			startHour = currentHour
		}
	}

	for hour := startHour; hour < 23; hour++ {
		startTime := baseDate.Add(time.Duration(hour) * time.Hour)
		endTime := startTime.Add(time.Hour)

		// Skip if the slot is in the past (additional safety check)
		if startTime.Before(localNow) {
			continue
		}

		slot := TimeSlot{
			StartTime: startTime.Format(time.RFC3339),
			EndTime:   endTime.Format(time.RFC3339),
			Available: true,
		}

		// Check if this slot conflicts with any existing appointment
		for _, appointment := range appointments {
			apptStart := appointment.StartTime
			apptEnd := appointment.EndTime

			// Check if the slot overlaps with the appointment
			if startTime.Before(apptEnd) && endTime.After(apptStart) {
				slot.Available = false
				slot.Type = string(appointment.Type)
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
