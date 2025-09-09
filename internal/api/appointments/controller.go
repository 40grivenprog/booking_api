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
)

// CreateAppointment handles POST /api/appointments
func (h *AppointmentsHandler) CreateAppointment(c *gin.Context) {
	var req CreateAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.HandleErrorResponse(c, http.StatusBadRequest, "validation_error", "Invalid request body", err)
		return
	}
	fmt.Println("req.StartTime", req.StartTime)
	fmt.Println("req.EndTime", req.EndTime)

	// Parse time strings and convert to application timezone
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

	// Store times in application timezone (Europe/Berlin)
	startTime = util.ConvertToAppTimezone(startTime)
	endTime = util.ConvertToAppTimezone(endTime)
	fmt.Println("startTime", startTime)
	fmt.Println("endTime", endTime)

	// Validate that start_time is in the future
	if startTime.Before(util.NowInAppTimezone()) {
		common.HandleErrorResponse(c, http.StatusBadRequest, "validation_error", "start_time must be in the future", nil)
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
		Description:    sql.NullString{String: "Personal training", Valid: true},
	})
	if err != nil {
		common.HandleErrorResponse(c, http.StatusInternalServerError, "database_error", "Failed to create appointment", err)
		return
	}

	// Convert to response format
	response := CreateAppointmentResponse{
		Appointment: Appointment{
			ID:          result.ID.String(),
			StartTime:   result.StartTime.Format(time.RFC3339),
			EndTime:     result.EndTime.Format(time.RFC3339),
			Status:      string(result.Status.AppointmentStatus),
			Description: result.Description.String,
			CreatedAt:   result.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   result.UpdatedAt.Format(time.RFC3339),
		},
		Client: Client{
			ID:          result.ClientIDFull.String(),
			FirstName:   result.ClientFirstName.String,
			LastName:    result.ClientLastName.String,
			PhoneNumber: result.ClientPhoneNumber.String,
		},
		Professional: Professional{
			ID:          result.ProfessionalIDFull.String(),
			Username:    result.ProfessionalUsername.String,
			FirstName:   result.ProfessionalFirstName.String,
			LastName:    result.ProfessionalLastName.String,
			PhoneNumber: result.ProfessionalPhoneNumber.String,
		},
	}

	// Handle optional fields
	if result.ClientChatID.Valid {
		response.Client.ChatID = result.ClientChatID.Int64
	}
	if result.ProfessionalChatID.Valid {
		response.Professional.ChatID = result.ProfessionalChatID.Int64
	}

	c.JSON(http.StatusCreated, response)
}
