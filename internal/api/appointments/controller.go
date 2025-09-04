package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	common "github.com/vention/booking_api/internal/api/common"
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
			Username:    "", // Clients don't have username
			FirstName:   result.ClientFirstName.String,
			LastName:    result.ClientLastName.String,
			PhoneNumber: result.ClientPhoneNumber.String,
		},
		Professional: Professional{
			ID:        result.ProfessionalIDFull.String(),
			Username:  result.ProfessionalUsername.String,
			FirstName: result.ProfessionalFirstName.String,
			LastName:  result.ProfessionalLastName.String,
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

// // GetAppointmentsByUser handles GET /api/users/{id}/appointments?status={status}
// func (h *AppointmentsHandler) GetAppointmentsByUser(c *gin.Context) {
// 	userIDStr := c.Param("id")
// 	status := c.Query("status")

// 	// Validate status parameter
// 	if status == "" {
// 		common.HandleErrorResponse(c, http.StatusBadRequest, "validation_error", "status parameter is required", nil)
// 		return
// 	}

// 	// Validate status value
// 	validStatuses := []string{"pending", "confirmed", "cancelled", "completed"}
// 	isValidStatus := false
// 	for _, validStatus := range validStatuses {
// 		if status == validStatus {
// 			isValidStatus = true
// 			break
// 		}
// 	}
// 	if !isValidStatus {
// 		common.HandleErrorResponse(c, http.StatusBadRequest, "validation_error", "Invalid status. Must be one of: pending, confirmed, cancelled, completed", nil)
// 		return
// 	}

// 	// Parse UUID
// 	userID, err := uuid.Parse(userIDStr)
// 	if err != nil {
// 		common.HandleErrorResponse(c, http.StatusBadRequest, "validation_error", "Invalid user_id format", err)
// 		return
// 	}

// 	// Get appointments by user and status
// 	appointments, err := h.appointmentsRepo.GetAppointmentsByUserAndStatus(c.Request.Context(), &db.GetAppointmentsByUserAndStatusParams{
// 		ID:     userID,
// 		Status: db.AppointmentStatus(status),
// 	})
// 	if err != nil {
// 		common.HandleErrorResponse(c, http.StatusInternalServerError, "database_error", "Failed to get appointments", err)
// 		return
// 	}

// 	// Convert to response format
// 	response := GetAppointmentsResponse{
// 		Appointments: make([]AppointmentWithDetails, len(appointments)),
// 	}

// 	for i, apt := range appointments {
// 		appointment := AppointmentWithDetails{
// 			ID:        apt.ID.String(),
// 			Type:      string(apt.Type.AppointmentType),
// 			StartTime: apt.StartTime.Format(time.RFC3339),
// 			EndTime:   apt.EndTime.Format(time.RFC3339),
// 			Status:    string(apt.Status.AppointmentStatus),
// 			CreatedAt: apt.CreatedAt.Format(time.RFC3339),
// 			UpdatedAt: apt.UpdatedAt.Format(time.RFC3339),
// 			Professional: ProfessionalDetails{
// 				ID:        apt.ProfessionalIDFull.String(),
// 				Username:  apt.ProfessionalUsername.String,
// 				FirstName: apt.ProfessionalFirstName.String,
// 				LastName:  apt.ProfessionalLastName.String,
// 			},
// 		}

// 		// Handle optional professional ChatID
// 		if apt.ProfessionalChatID.Valid {
// 			appointment.Professional.ChatID = apt.ProfessionalChatID.Int64
// 		}

// 		// Handle optional client details
// 		if apt.ClientIDFull.Valid {
// 			appointment.Client = &ClientDetails{
// 				ID:        apt.ClientIDFull.UUID.String(),
// 				FirstName: apt.ClientFirstName.String,
// 				LastName:  apt.ClientLastName.String,
// 			}
// 			if apt.ClientChatID.Valid {
// 				appointment.Client.ChatID = apt.ClientChatID.Int64
// 			}
// 		}

// 		response.Appointments[i] = appointment
// 	}

// 	c.JSON(http.StatusOK, response)
// }
