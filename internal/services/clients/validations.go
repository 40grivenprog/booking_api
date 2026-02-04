package clients

import (
	"context"
	"time"

	"github.com/google/uuid"
	db "github.com/vention/booking_api/internal/repository"
	svcCommon "github.com/vention/booking_api/internal/services/common"
)

// import (
// 	"context"
// 	"time"

// 	"github.com/google/uuid"
// 	db "github.com/vention/booking_api/internal/repository"
// 	svcCommon "github.com/vention/booking_api/internal/services/common"
// )

// // validateAppointmentOwnership validates that the appointment belongs to the client
// func (s *service) validateAppointmentOwnership(appointment *db.Appointment, clientID uuid.UUID) error {
// 	if appointment.ClientID.UUID != clientID {
// 		return svcCommon.ErrForbidden
// 	}
// 	return nil
// }

// // validateAppointmentCancellable validates that the appointment can be cancelled
// func (s *service) validateAppointmentCancellable(appointment *db.Appointment) error {
// 	if appointment.Status.AppointmentStatus != db.AppointmentStatusPending &&
// 		appointment.Status.AppointmentStatus != db.AppointmentStatusConfirmed {
// 		return svcCommon.ErrAppointmentNotPendingOrConfirmed
// 	}
// 	return nil
// }

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

// validateAppointmentConflict checks if the client already has an appointment at the same time with the same professional
func (s *service) validateAppointmentConflict(ctx context.Context, clientID, professionalID uuid.UUID, startTime time.Time) error {
	hasConflict, err := s.repo.CheckClientAppointmentConflict(ctx, &db.CheckClientAppointmentConflictParams{
		ClientID:       clientID,
		ProfessionalID: professionalID,
		StartTime:      startTime,
	})
	if err != nil {
		return err
	}

	if hasConflict {
		return svcCommon.ErrAppointmentTimeConflict
	}

	return nil
}
