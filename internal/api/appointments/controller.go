package api

import (
	"database/sql"
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
		common.HandleErrorResponse(c, http.StatusBadRequest, common.ErrorTypeValidation, common.ErrorMsgInvalidRequestBody, err)
		return
	}

	// Parse time strings and convert to application timezone
	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		common.HandleErrorResponse(c, http.StatusBadRequest, common.ErrorTypeValidation, common.ErrorMsgInvalidTime, err)
		return
	}

	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		common.HandleErrorResponse(c, http.StatusBadRequest, common.ErrorTypeValidation, common.ErrorMsgInvalidTime, err)
		return
	}

	// Store times in application timezone (Europe/Berlin)
	startTime = util.ConvertToAppTimezone(startTime)
	endTime = util.ConvertToAppTimezone(endTime)

	// Validate that start_time is in the future
	if startTime.Before(util.NowInAppTimezone()) {
		common.HandleErrorResponse(c, http.StatusBadRequest, common.ErrorTypeValidation, common.ErrorMsgFutureTimeRequired, nil)
		return
	}

	// Validate that end_time is after start_time
	if endTime.Before(startTime) || endTime.Equal(startTime) {
		common.HandleErrorResponse(c, http.StatusBadRequest, common.ErrorTypeValidation, common.ErrorMsgInvalidTime, nil)
		return
	}

	// Parse UUIDs
	clientID, err := uuid.Parse(req.ClientID)
	if err != nil {
		common.HandleErrorResponse(c, http.StatusBadRequest, common.ErrorTypeValidation, common.ErrorMsgInvalidClientID, err)
		return
	}

	professionalID, err := uuid.Parse(req.ProfessionalID)
	if err != nil {
		common.HandleErrorResponse(c, http.StatusBadRequest, common.ErrorTypeValidation, common.ErrorMsgInvalidProfessionalID, err)
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
		common.HandleErrorResponse(c, http.StatusInternalServerError, common.ErrorTypeDatabase, common.ErrorMsgFailedToCreateAppointment, err)
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
