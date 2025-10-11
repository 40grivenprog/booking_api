package common

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ParseUUID parses UUID from string and handles error response automatically
// Returns parsed UUID and boolean indicating success
func ParseUUID(c *gin.Context, idStr string, errorMsg string) (uuid.UUID, bool) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		HandleErrorResponse(c, http.StatusBadRequest, ErrorTypeValidation, errorMsg, err)
		return uuid.UUID{}, false
	}
	return id, true
}

// ParseClientID is a convenience wrapper for parsing client IDs
func ParseClientID(c *gin.Context, idStr string) (uuid.UUID, bool) {
	return ParseUUID(c, idStr, ErrorMsgInvalidClientID)
}

// ParseProfessionalID is a convenience wrapper for parsing professional IDs
func ParseProfessionalID(c *gin.Context, idStr string) (uuid.UUID, bool) {
	return ParseUUID(c, idStr, ErrorMsgInvalidProfessionalID)
}

// ParseAppointmentID is a convenience wrapper for parsing appointment IDs
func ParseAppointmentID(c *gin.Context, idStr string) (uuid.UUID, bool) {
	return ParseUUID(c, idStr, ErrorMsgInvalidAppointmentID)
}

// ParseTime parses RFC3339 time string and handles error response automatically
func ParseTime(c *gin.Context, timeStr string, errorMsg string) (time.Time, bool) {
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		HandleErrorResponse(c, http.StatusBadRequest, ErrorTypeValidation, errorMsg, err)
		return time.Time{}, false
	}
	return t, true
}

// ParseDate parses date string (YYYY-MM-DD format) and handles error response automatically
func ParseDate(c *gin.Context, dateStr string, errorMsg string) (time.Time, bool) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		HandleErrorResponse(c, http.StatusBadRequest, ErrorTypeValidation, errorMsg, err)
		return time.Time{}, false
	}
	return date, true
}

// ParseMonth parses month string (YYYY-MM format) and handles error response automatically
func ParseMonth(c *gin.Context, monthStr string) (time.Time, bool) {
	month, err := time.Parse("2006-01", monthStr)
	if err != nil {
		HandleErrorResponse(c, http.StatusBadRequest, ErrorTypeValidation, "Invalid month format. Use YYYY-MM", err)
		return time.Time{}, false
	}
	return month, true
}

// ValidAppointmentStatuses contains all valid appointment statuses
var ValidAppointmentStatuses = map[string]bool{
	"pending":   true,
	"confirmed": true,
	"cancelled": true,
	"completed": true,
}

// ValidateAppointmentStatus validates appointment status and handles error response automatically
func ValidateAppointmentStatus(c *gin.Context, status string) bool {
	if status == "" {
		return true // Empty status is allowed for optional filters
	}

	if !ValidAppointmentStatuses[status] {
		HandleErrorResponse(c, http.StatusBadRequest, ErrorTypeValidation, ErrorMsgInvalidStatus, nil)
		return false
	}
	return true
}

// RequireQueryParam validates that a required query parameter is present
func RequireQueryParam(c *gin.Context, paramName string) (string, bool) {
	value := c.Query(paramName)
	if value == "" {
		HandleErrorResponse(c, http.StatusBadRequest, ErrorTypeValidation, ErrorMsgMissingRequiredField, nil)
		return "", false
	}
	return value, true
}
