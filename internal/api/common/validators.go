package common

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "github.com/vention/booking_api/internal/repository"
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

// ValidateTimeRange validates that end time is after start time and both are valid
func ValidateTimeRange(c *gin.Context, startTime, endTime time.Time) bool {
	if endTime.Before(startTime) || endTime.Equal(startTime) {
		HandleErrorResponse(c, http.StatusBadRequest, ErrorTypeValidation, ErrorMsgInvalidTime, nil)
		return false
	}
	return true
}

// ValidateFutureTime validates that time is in the future
func ValidateFutureTime(c *gin.Context, t time.Time, now time.Time) bool {
	if t.Before(now) {
		HandleErrorResponse(c, http.StatusBadRequest, ErrorTypeValidation, ErrorMsgFutureTimeRequired, nil)
		return false
	}
	return true
}

// ValidateAppointmentOwnership validates that appointment belongs to the given user
func ValidateAppointmentOwnership(c *gin.Context, appointment *db.Appointment, userID uuid.UUID, userType string) bool {
	switch userType {
	case UserTypeClient:
		if appointment.ClientID.UUID != userID {
			HandleErrorResponse(c, http.StatusForbidden, ErrorTypeForbidden, ErrorMsgNotAllowedToCancelAppointment, nil)
			return false
		}
	case UserTypeProfessional:
		if appointment.ProfessionalID != userID {
			HandleErrorResponse(c, http.StatusForbidden, ErrorTypeForbidden, ErrorMsgNotAllowedToConfirmAppointment, nil)
			return false
		}
	default:
		HandleErrorResponse(c, http.StatusForbidden, ErrorTypeForbidden, ErrorMsgNotAllowedToAccessResource, nil)
		return false
	}
	return true
}

// ValidateAppointmentStatus validates that appointment has one of the allowed statuses
func ValidateAppointmentStatusIs(c *gin.Context, appointment *db.Appointment, allowedStatuses ...db.AppointmentStatus) bool {
	currentStatus := appointment.Status.AppointmentStatus

	for _, status := range allowedStatuses {
		if currentStatus == status {
			return true
		}
	}

	// Build error message based on allowed statuses
	if len(allowedStatuses) == 1 && allowedStatuses[0] == db.AppointmentStatusPending {
		HandleErrorResponse(c, http.StatusBadRequest, ErrorTypeValidation, ErrorMsgAppointmentNotPending, nil)
	} else {
		HandleErrorResponse(c, http.StatusBadRequest, ErrorTypeValidation, ErrorMsgAppointmentNotPendingOrConfirmed, nil)
	}
	return false
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
