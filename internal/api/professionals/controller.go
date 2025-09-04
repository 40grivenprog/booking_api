package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	common "github.com/vention/booking_api/internal/api/common"
	db "github.com/vention/booking_api/internal/repository"
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
		UserType:  "professional",
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
			CreatedAt: result.CreatedAt.Format(time.RFC3339),
			UpdatedAt: result.UpdatedAt.Format(time.RFC3339),
		},
		Client: ClientConfirm{
			ID:        result.ClientID.UUID.String(),
			FirstName: result.ClientFirstName.String,
			LastName:  result.ClientLastName.String,
		},
	}

	// Handle optional ChatID field
	if result.ClientChatID.Valid {
		response.Client.ChatID = result.ClientChatID.Int64
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
