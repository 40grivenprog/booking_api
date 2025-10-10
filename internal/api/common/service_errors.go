package common

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	svcCommon "github.com/vention/booking_api/internal/services/common"
)

// HandleServiceError maps service-layer errors to HTTP responses
func HandleServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, svcCommon.ErrInvalidTimeRange):
		HandleErrorResponse(c, http.StatusBadRequest, ErrorTypeValidation, ErrorMsgInvalidTime, err)

	case errors.Is(err, svcCommon.ErrPastTime):
		HandleErrorResponse(c, http.StatusBadRequest, ErrorTypeValidation, ErrorMsgFutureTimeRequired, err)

	case errors.Is(err, svcCommon.ErrInvalidCredentials):
		HandleErrorResponse(c, http.StatusUnauthorized, ErrorTypeValidation, ErrorMsgInvalidCredentials, err)

	case errors.Is(err, svcCommon.ErrForbidden):
		HandleErrorResponse(c, http.StatusForbidden, ErrorTypeForbidden, ErrorMsgNotAllowedToAccessResource, err)

	case errors.Is(err, svcCommon.ErrAppointmentNotPending):
		HandleErrorResponse(c, http.StatusBadRequest, ErrorTypeValidation, ErrorMsgAppointmentNotPending, err)

	case errors.Is(err, svcCommon.ErrAppointmentNotPendingOrConfirmed):
		HandleErrorResponse(c, http.StatusBadRequest, ErrorTypeValidation, ErrorMsgAppointmentNotPendingOrConfirmed, err)

	default:
		// For unknown errors, return internal server error
		HandleErrorResponse(c, http.StatusInternalServerError, ErrorTypeInternal, ErrorMsgInternalServerError, err)
	}
}
