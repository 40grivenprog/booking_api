package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	common "github.com/vention/booking_api/internal/api/common"
	db "github.com/vention/booking_api/internal/repository"
)

// RegisterClient handles POST /api/clients/register
func (h *ClientsHandler) RegisterClient(c *gin.Context) {
	var req ClientRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.HandleErrorResponse(c, http.StatusBadRequest, "validation_error", "Invalid request body", err)
		return
	}

	// Convert phone number to sql.NullString
	var phoneNumber sql.NullString
	if req.PhoneNumber != nil {
		phoneNumber = sql.NullString{String: *req.PhoneNumber, Valid: true}
	}

	// Convert chat_id to sql.NullInt64
	var chatID sql.NullInt64
	chatID = sql.NullInt64{Int64: req.ChatID, Valid: true}

	// Create new client
	user, err := h.clientsRepo.CreateClient(c.Request.Context(), &db.CreateClientParams{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: phoneNumber,
		ChatID:      chatID,
		CreatedBy:   uuid.NullUUID{}, // NULL for now
	})
	if err != nil {
		// Check if it's a unique constraint violation
		if common.IsUniqueConstraintError(err) {
			common.HandleErrorResponse(c, http.StatusConflict, "username_taken", "The provided username is already taken", nil)
			return
		}
		common.HandleErrorResponse(c, http.StatusInternalServerError, "database_error", "Failed to create client", err)
		return
	}

	// Convert to response format
	responseUser := User{
		ID:        user.ID.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      "client",
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	// Handle optional fields
	if user.ChatID.Valid {
		responseUser.ChatID = &user.ChatID.Int64
	}
	if user.PhoneNumber.Valid {
		responseUser.PhoneNumber = &user.PhoneNumber.String
	}

	response := ClientRegisterResponse{
		User: responseUser,
	}

	c.JSON(http.StatusCreated, response)
}

// GetClientAppointments handles GET /api/clients/{id}/appointments
func (h *ClientsHandler) GetClientAppointments(c *gin.Context) {
	clientIDStr := c.Param("id")
	statusFilter := c.Query("status")

	// Parse client ID
	clientID, err := uuid.Parse(clientIDStr)
	if err != nil {
		common.HandleErrorResponse(c, http.StatusBadRequest, "validation_error", "Invalid client_id format", err)
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
	appointments, err := h.clientsRepo.GetAppointmentsByClientWithStatus(c.Request.Context(), &db.GetAppointmentsByClientWithStatusParams{
		ClientID: uuid.NullUUID{UUID: clientID, Valid: true},
		Status:   statusParam,
	})
	if err != nil {
		common.HandleErrorResponse(c, http.StatusInternalServerError, "database_error", "Failed to retrieve appointments", err)
		return
	}

	// Convert to response format
	var responseAppointments []ClientAppointment
	for _, appt := range appointments {
		appointment := ClientAppointment{
			ID:        appt.ID.String(),
			Type:      string(appt.Type),
			StartTime: appt.StartTime.Format(time.RFC3339),
			EndTime:   appt.EndTime.Format(time.RFC3339),
			Status:    string(appt.Status.AppointmentStatus),
			CreatedAt: appt.CreatedAt.Format(time.RFC3339),
			UpdatedAt: appt.UpdatedAt.Format(time.RFC3339),
		}
		professional := &ClientAppointmentProfessional{
			ID:        appt.ProfessionalIDFull.String(),
			Username:  appt.ProfessionalUsername.String,
			FirstName: appt.ProfessionalFirstName.String,
			LastName:  appt.ProfessionalLastName.String,
		}
		appointment.Professional = professional

		responseAppointments = append(responseAppointments, appointment)
	}

	response := GetClientAppointmentsResponse{
		Appointments: responseAppointments,
	}

	c.JSON(http.StatusOK, response)
}

// CancelClientAppointment handles PATCH /api/clients/{id}/appointments/{appointment_id}/cancel
func (h *ClientsHandler) CancelClientAppointment(c *gin.Context) {
	clientIDStr := c.Param("id")
	appointmentIDStr := c.Param("appointment_id")

	// Parse UUIDs
	clientID, err := uuid.Parse(clientIDStr)
	if err != nil {
		common.HandleErrorResponse(c, http.StatusBadRequest, "validation_error", "Invalid client_id format", err)
		return
	}

	appointmentID, err := uuid.Parse(appointmentIDStr)
	if err != nil {
		common.HandleErrorResponse(c, http.StatusBadRequest, "validation_error", "Invalid appointment_id format", err)
		return
	}

	// Parse request body
	var req CancelClientAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.HandleErrorResponse(c, http.StatusBadRequest, "validation_error", "Invalid request body", err)
		return
	}

	appointment, err := h.clientsRepo.GetAppointmentByID(c.Request.Context(), appointmentID)
	if err != nil {
		common.HandleErrorResponse(c, http.StatusInternalServerError, "database_error", "Failed to get appointment", err)
		return
	}
	if appointment.ClientID.UUID != clientID {
		common.HandleErrorResponse(c, http.StatusForbidden, "forbidden", "You are not allowed to cancel this appointment", nil)
		return
	}
	if appointment.Status.AppointmentStatus != db.AppointmentStatusPending && appointment.Status.AppointmentStatus != db.AppointmentStatusConfirmed {
		common.HandleErrorResponse(c, http.StatusBadRequest, "validation_error", "Appointment is not pending or confirmed. Please check the status of the appointment.", nil)
		return
	}

	// Cancel appointment with details
	result, err := h.clientsRepo.CancelAppointmentByClientWithDetails(c.Request.Context(), &db.CancelAppointmentByClientWithDetailsParams{
		ID:                  appointmentID,
		CancelledByClientID: uuid.NullUUID{UUID: clientID, Valid: true},
		CancellationReason:  sql.NullString{String: req.CancellationReason, Valid: true},
	})
	if err != nil {
		common.HandleErrorResponse(c, http.StatusInternalServerError, "database_error", "Failed to cancel appointment", err)
		return
	}

	// Convert to response format
	response := CancelClientAppointmentResponse{
		Appointment: CancelledAppointment{
			ID:                 result.ID.String(),
			Type:               string(result.Type),
			StartTime:          result.StartTime.Format(time.RFC3339),
			EndTime:            result.EndTime.Format(time.RFC3339),
			Status:             string(result.Status.AppointmentStatus),
			CancellationReason: result.CancellationReason.String,
			CancelledBy:        "client",
			CreatedAt:          result.CreatedAt.Format(time.RFC3339),
			UpdatedAt:          result.UpdatedAt.Format(time.RFC3339),
		},
		Client: ClientAppointmentClient{
			ID:        result.ClientIDFull.String(),
			FirstName: result.ClientFirstName.String,
			LastName:  result.ClientLastName.String,
		},
		Professional: ClientAppointmentProfessional{
			ID:        result.ProfessionalIDFull.String(),
			Username:  result.ProfessionalUsername.String,
			FirstName: result.ProfessionalFirstName.String,
			LastName:  result.ProfessionalLastName.String,
		},
	}

	// Handle optional fields
	if result.ClientPhoneNumber.Valid {
		response.Client.PhoneNumber = &result.ClientPhoneNumber.String
	}
	if result.ClientChatID.Valid {
		response.Client.ChatID = &result.ClientChatID.Int64
	}
	if result.ProfessionalPhoneNumber.Valid {
		response.Professional.PhoneNumber = &result.ProfessionalPhoneNumber.String
	}
	if result.ProfessionalChatID.Valid {
		response.Professional.ChatID = &result.ProfessionalChatID.Int64
	}

	c.JSON(http.StatusOK, response)
}
