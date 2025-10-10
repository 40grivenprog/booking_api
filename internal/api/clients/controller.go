package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	common "github.com/vention/booking_api/internal/api/common"
	db "github.com/vention/booking_api/internal/repository"
)

// RegisterClient handles POST /api/clients/register
func (h *ClientsHandler) RegisterClient(c *gin.Context) {
	req, ok := common.BindAndValidate[ClientRegisterRequest](c)
	if !ok {
		return
	}

	// Create new client
	user, err := h.clientsRepo.CreateClient(c.Request.Context(), &db.CreateClientParams{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: common.ToNullString(req.PhoneNumber),
		ChatID:      common.ToNullInt64Value(req.ChatID),
		CreatedBy:   uuid.NullUUID{}, // NULL for now
	})
	if err != nil {
		// Check if it's a unique constraint violation
		if common.IsUniqueConstraintError(err) {
			common.HandleErrorResponse(c, http.StatusConflict, common.ErrorTypeConflict, common.ErrorMsgUsernameAlreadyExists, nil)
			return
		}
		common.HandleErrorResponse(c, http.StatusInternalServerError, common.ErrorTypeDatabase, common.ErrorMsgFailedToCreateClient, err)
		return
	}

	response := mapClientToClientRegisterResponse(user)
	c.JSON(http.StatusCreated, response)
}

// GetClientAppointments handles GET /api/clients/{id}/appointments
func (h *ClientsHandler) GetClientAppointments(c *gin.Context) {
	// Parse and validate client ID
	clientID, ok := common.ParseClientID(c, c.Param("id"))
	if !ok {
		return
	}

	// Validate status filter if provided
	statusFilter := c.Query("status")
	if !common.ValidateAppointmentStatus(c, statusFilter) {
		return
	}

	// Prepare parameters for the query
	var statusParam db.NullAppointmentStatus
	if statusFilter != "" {
		statusParam = db.NullAppointmentStatus{AppointmentStatus: db.AppointmentStatus(statusFilter), Valid: true}
	}

	// Get appointments with optional status filter
	appointments, err := h.clientsRepo.GetAppointmentsByClientWithStatus(c.Request.Context(), &db.GetAppointmentsByClientWithStatusParams{
		ClientID: common.ToNullUUID(clientID),
		Status:   statusParam,
	})
	if err != nil {
		common.HandleErrorResponse(c, http.StatusInternalServerError, common.ErrorTypeDatabase, common.ErrorMsgFailedToRetrieveAppointments, err)
		return
	}

	// Convert to response format
	response := mapAppointmentToGetClientAppointmentsResponse(appointments)

	c.JSON(http.StatusOK, response)
}

// CancelClientAppointment handles PATCH /api/clients/{id}/appointments/{appointment_id}/cancel
func (h *ClientsHandler) CancelClientAppointment(c *gin.Context) {
	// Parse UUIDs
	clientID, ok := common.ParseClientID(c, c.Param("id"))
	if !ok {
		return
	}

	appointmentID, ok := common.ParseAppointmentID(c, c.Param("appointment_id"))
	if !ok {
		return
	}

	// Parse request body
	req, ok := common.BindAndValidate[CancelClientAppointmentRequest](c)
	if !ok {
		return
	}

	// Get and validate appointment
	appointment, err := h.clientsRepo.GetAppointmentByID(c.Request.Context(), appointmentID)
	if err != nil {
		common.HandleErrorResponse(c, http.StatusInternalServerError, common.ErrorTypeDatabase, common.ErrorMsgFailedToGetAppointment, err)
		return
	}

	// Validate ownership and status
	if !common.ValidateAppointmentOwnership(c, appointment, clientID, common.UserTypeClient) {
		return
	}
	if !common.ValidateAppointmentStatusIs(c, appointment, db.AppointmentStatusPending, db.AppointmentStatusConfirmed) {
		return
	}

	// Cancel appointment with details
	result, err := h.clientsRepo.CancelAppointmentByClientWithDetails(c.Request.Context(), &db.CancelAppointmentByClientWithDetailsParams{
		ID:                  appointmentID,
		CancelledByClientID: common.ToNullUUID(clientID),
		CancellationReason:  common.ToNullStringValue(req.CancellationReason),
	})
	if err != nil {
		common.HandleErrorResponse(c, http.StatusInternalServerError, common.ErrorTypeDatabase, common.ErrorMsgFailedToUpdateAppointment, err)
		return
	}

	response := mapAppointmentToCancelClientAppointmentResponse(result)
	c.JSON(http.StatusOK, response)
}
