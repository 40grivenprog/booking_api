package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	common "github.com/vention/booking_api/internal/api/common"
	"github.com/vention/booking_api/internal/services/clients"
)

// RegisterClient handles POST /api/clients/register
func (h *ClientsHandler) RegisterClient(c *gin.Context) {
	req, ok := common.BindAndValidate[ClientRegisterRequest](c)
	if !ok {
		return
	}

	phoneNumber := ""
	if req.PhoneNumber != nil {
		phoneNumber = *req.PhoneNumber
	}

	client, err := h.clientsService.RegisterClient(c.Request.Context(), clients.RegisterClientInput{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: phoneNumber,
		ChatID:      req.ChatID,
	})
	if err != nil {
		if common.IsUniqueConstraintError(err) {
			common.HandleErrorResponse(c, http.StatusConflict, common.ErrorTypeConflict, common.ErrorMsgUsernameAlreadyExists, nil)
			return
		}
		common.HandleErrorResponse(c, http.StatusInternalServerError, common.ErrorTypeDatabase, common.ErrorMsgFailedToCreateClient, err)
		return
	}

	response := mapClientToClientRegisterResponse(client)
	c.JSON(http.StatusCreated, response)
}

// GetClientAppointments handles GET /api/clients/{id}/appointments
func (h *ClientsHandler) GetClientAppointments(c *gin.Context) {
	clientID, ok := common.ParseClientID(c, c.Param("id"))
	if !ok {
		return
	}

	statusFilter := c.Query("status")
	if !common.ValidateAppointmentStatus(c, statusFilter) {
		return
	}

	appointments, err := h.clientsService.GetClientAppointments(c.Request.Context(), clientID, statusFilter)
	if err != nil {
		common.HandleErrorResponse(c, http.StatusInternalServerError, common.ErrorTypeDatabase, common.ErrorMsgFailedToRetrieveAppointments, err)
		return
	}

	response := mapAppointmentToGetClientAppointmentsResponse(appointments)
	c.JSON(http.StatusOK, response)
}

// CancelClientAppointment handles PATCH /api/clients/{id}/appointments/{appointment_id}/cancel
func (h *ClientsHandler) CancelClientAppointment(c *gin.Context) {
	clientID, ok := common.ParseClientID(c, c.Param("id"))
	if !ok {
		return
	}

	appointmentID, ok := common.ParseAppointmentID(c, c.Param("appointment_id"))
	if !ok {
		return
	}

	req, ok := common.BindAndValidate[CancelClientAppointmentRequest](c)
	if !ok {
		return
	}

	result, err := h.clientsService.CancelAppointment(c.Request.Context(), clients.CancelAppointmentInput{
		ClientID:           clientID,
		AppointmentID:      appointmentID,
		CancellationReason: req.CancellationReason,
	})
	if err != nil {
		common.HandleServiceError(c, err)
		return
	}

	response := mapAppointmentToCancelClientAppointmentResponse(result)
	c.JSON(http.StatusOK, response)
}
