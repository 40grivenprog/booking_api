package clients

import (
	"github.com/google/uuid"
	db "github.com/vention/booking_api/internal/repository"
	svcCommon "github.com/vention/booking_api/internal/services/common"
)

// validateAppointmentOwnership validates that the appointment belongs to the client
func (s *service) validateAppointmentOwnership(appointment *db.Appointment, clientID uuid.UUID) error {
	if appointment.ClientID.UUID != clientID {
		return svcCommon.ErrForbidden
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
