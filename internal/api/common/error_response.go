package common

import (
	"github.com/gin-gonic/gin"
	. "github.com/vention/booking_api/internal/op"
)

// HandleErrorResponse creates a standardized error response
func HandleErrorResponse(c *gin.Context, statusCode int, errorType, message string, err error) {
	logger := GetLogger(c)
	requestID := GetRequestID(c)

	// Log the error with stack trace if it's a system error
	if err != nil {
		logger.Error().
			Err(err).
			Str("stack_trace", Stack()).
			Msg(message)
	} else {
		logger.Warn().
			Str("error_type", errorType).
			Msg(message)
	}

	errorResp := ErrorResponse{
		Error:     errorType,
		Message:   message,
		RequestID: requestID,
	}

	c.JSON(statusCode, errorResp)
}

// ErrorResponse represents error responses
type ErrorResponse struct {
	Error     string `json:"error"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
}
