package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	common "github.com/vention/booking_api/internal/api/common"
	db "github.com/vention/booking_api/internal/repository"
	"github.com/vention/booking_api/internal/util"
)

// CreateAppointment handles POST /api/appointments
func (h *AppointmentsHandler) CreateAppointment(c *gin.Context) {
	req, ok := common.BindAndValidate[CreateAppointmentRequest](c)
	if !ok {
		return
	}

	// Parse and validate time strings
	startTime, ok := common.ParseTime(c, req.StartTime, common.ErrorMsgInvalidTime)
	if !ok {
		return
	}

	endTime, ok := common.ParseTime(c, req.EndTime, common.ErrorMsgInvalidTime)
	if !ok {
		return
	}

	// Store times in application timezone
	startTime = util.ConvertToAppTimezone(startTime)
	endTime = util.ConvertToAppTimezone(endTime)

	// Validate time range
	if !common.ValidateFutureTime(c, startTime, util.NowInAppTimezone()) {
		return
	}
	if !common.ValidateTimeRange(c, startTime, endTime) {
		return
	}

	// Parse UUIDs
	clientID, ok := common.ParseClientID(c, req.ClientID)
	if !ok {
		return
	}

	professionalID, ok := common.ParseProfessionalID(c, req.ProfessionalID)
	if !ok {
		return
	}

	// Create new appointment with details
	result, err := h.appointmentsRepo.CreateAppointmentWithDetails(c.Request.Context(), &db.CreateAppointmentWithDetailsParams{
		ClientID:       common.ToNullUUID(clientID),
		ProfessionalID: professionalID,
		StartTime:      startTime,
		EndTime:        endTime,
		Description:    common.ToNullStringValue("Personal training"),
	})
	if err != nil {
		common.HandleErrorResponse(c, http.StatusInternalServerError, common.ErrorTypeDatabase, common.ErrorMsgFailedToCreateAppointment, err)
		return
	}

	response := mapAppointmentToCreateAppointmentResponse(result)

	c.JSON(http.StatusCreated, response)
}
