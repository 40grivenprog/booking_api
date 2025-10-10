package appointments

import (
	"time"

	svcCommon "github.com/vention/booking_api/internal/services/common"
)

// validateAppointmentTime validates the appointment time range
func (s *service) validateAppointmentTime(startTime, endTime time.Time) error {
	now := time.Now()

	// Check if start time is in the future
	if startTime.Before(now) {
		return svcCommon.ErrPastTime
	}

	// Check if end time is after start time
	if endTime.Before(startTime) || endTime.Equal(startTime) {
		return svcCommon.ErrInvalidTimeRange
	}

	return nil
}
