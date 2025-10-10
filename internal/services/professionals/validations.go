package professionals

import (
	"time"

	"github.com/google/uuid"
	db "github.com/vention/booking_api/internal/repository"
	svcCommon "github.com/vention/booking_api/internal/services/common"
	"golang.org/x/crypto/bcrypt"
)

// validatePassword validates the professional's password
func (s *service) validatePassword(professional *db.Professional, password string) error {
	// Check if password hash exists
	if !professional.PasswordHash.Valid {
		return svcCommon.ErrInvalidCredentials
	}

	// Compare password with hash
	if err := bcrypt.CompareHashAndPassword([]byte(professional.PasswordHash.String), []byte(password)); err != nil {
		return svcCommon.ErrInvalidCredentials
	}

	return nil
}

// validateAppointmentOwnership validates that the appointment belongs to the professional
func (s *service) validateAppointmentOwnership(appointment *db.Appointment, professionalID uuid.UUID) error {
	if appointment.ProfessionalID != professionalID {
		return svcCommon.ErrForbidden
	}
	return nil
}

// validateAppointmentPending validates that the appointment is in pending status
func (s *service) validateAppointmentPending(appointment *db.Appointment) error {
	if appointment.Status.AppointmentStatus != db.AppointmentStatusPending {
		return svcCommon.ErrAppointmentNotPending
	}
	return nil
}

// validateAppointmentCancellable validates that the appointment can be cancelled
func (s *service) validateAppointmentCancellable(appointment *db.Appointment) error {
	if appointment.Status.AppointmentStatus != db.AppointmentStatusPending &&
		appointment.Status.AppointmentStatus != db.AppointmentStatusConfirmed {
		return svcCommon.ErrAppointmentNotPendingOrConfirmed
	}
	return nil
}

// validateTimeRange validates time range for appointments
func (s *service) validateTimeRange(startTime, endTime time.Time) error {
	// Check if start time is in the future
	if startTime.Before(time.Now()) {
		return svcCommon.ErrPastTime
	}

	// Check if end time is after start time
	if endTime.Before(startTime) || endTime.Equal(startTime) {
		return svcCommon.ErrInvalidTimeRange
	}

	return nil
}
