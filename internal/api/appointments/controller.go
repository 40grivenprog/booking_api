package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vention/booking_api/internal/api/common"
	db "github.com/vention/booking_api/internal/repository"
)

// CreateAppointment handles POST /api/appointments
func (h *AppointmentsHandler) CreateAppointment(c *gin.Context) {
	var req CreateAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.HandleErrorResponse(c, http.StatusBadRequest, "validation_error", "Invalid request body", err)
		return
	}

	// Parse time strings
	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		common.HandleErrorResponse(c, http.StatusBadRequest, "validation_error", "Invalid start_time format. Use RFC3339 format (e.g., 2024-01-01T10:00:00Z)", err)
		return
	}

	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		common.HandleErrorResponse(c, http.StatusBadRequest, "validation_error", "Invalid end_time format. Use RFC3339 format (e.g., 2024-01-01T11:00:00Z)", err)
		return
	}

	// Validate that end_time is after start_time
	if endTime.Before(startTime) || endTime.Equal(startTime) {
		common.HandleErrorResponse(c, http.StatusBadRequest, "validation_error", "end_time must be after start_time", nil)
		return
	}

	// Parse UUIDs
	clientID, err := uuid.Parse(req.ClientID)
	if err != nil {
		common.HandleErrorResponse(c, http.StatusBadRequest, "validation_error", "Invalid client_id format", err)
		return
	}

	professionalID, err := uuid.Parse(req.ProfessionalID)
	if err != nil {
		common.HandleErrorResponse(c, http.StatusBadRequest, "validation_error", "Invalid professional_id format", err)
		return
	}

	// Create new appointment with details
	result, err := h.appointmentsRepo.CreateAppointmentWithDetails(c.Request.Context(), &db.CreateAppointmentWithDetailsParams{
		ClientID:       uuid.NullUUID{UUID: clientID, Valid: true},
		ProfessionalID: professionalID,
		StartTime:      startTime,
		EndTime:        endTime,
	})
	if err != nil {
		common.HandleErrorResponse(c, http.StatusInternalServerError, "database_error", "Failed to create appointment", err)
		return
	}

	// Convert to response format
	response := CreateAppointmentResponse{
		Appointment: Appointment{
			ID:        result.ID.String(),
			StartTime: result.StartTime.Format(time.RFC3339),
			EndTime:   result.EndTime.Format(time.RFC3339),
			Status:    string(result.Status.AppointmentStatus),
			CreatedAt: result.CreatedAt.Format(time.RFC3339),
			UpdatedAt: result.UpdatedAt.Format(time.RFC3339),
		},
		Client: Client{
			ID:          result.ClientIDFull.String(),
			Username:    result.ClientUsername.String,
			FirstName:   result.ClientFirstName.String,
			LastName:    result.ClientLastName.String,
			PhoneNumber: result.ClientPhoneNumber.String,
			ChatID:      result.ClientChatID.Int64,
		},
		Professional: Professional{
			ID:        result.ProfessionalIDFull.String(),
			Username:  result.ProfessionalUsername.String,
			FirstName: result.ProfessionalFirstName.String,
			LastName:  result.ProfessionalLastName.String,
			ChatID:    result.ProfessionalChatID.Int64,
		},
	}

	c.JSON(http.StatusCreated, response)
}
