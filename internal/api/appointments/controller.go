package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	common "github.com/vention/booking_api/internal/api/common"
	"github.com/vention/booking_api/internal/services/appointments"
	"github.com/vention/booking_api/internal/util"
)

// CreateAppointment handles POST /api/appointments
func (h *AppointmentsHandler) CreateAppointment(c *gin.Context) {
	req, ok := common.BindAndValidate[CreateAppointmentRequest](c)
	if !ok {
		return
	}

	// Parse time strings
	startTime, ok := common.ParseTime(c, req.StartTime, common.ErrorMsgInvalidTime)
	if !ok {
		return
	}

	endTime, ok := common.ParseTime(c, req.EndTime, common.ErrorMsgInvalidTime)
	if !ok {
		return
	}

	// Convert to application timezone
	startTime = util.ConvertToAppTimezone(startTime)
	endTime = util.ConvertToAppTimezone(endTime)

	// Parse UUIDs
	clientID, ok := common.ParseClientID(c, req.ClientID)
	if !ok {
		return
	}

	professionalID, ok := common.ParseProfessionalID(c, req.ProfessionalID)
	if !ok {
		return
	}

	result, err := h.appointmentsService.CreateAppointment(c.Request.Context(), appointments.CreateAppointmentInput{
		ClientID:       clientID,
		ProfessionalID: professionalID,
		StartTime:      startTime,
		EndTime:        endTime,
		Description:    "Personal training",
	})
	if err != nil {
		common.HandleServiceError(c, err)
		return
	}

	response := mapAppointmentToCreateAppointmentResponse(result)

	c.JSON(http.StatusCreated, response)
}
