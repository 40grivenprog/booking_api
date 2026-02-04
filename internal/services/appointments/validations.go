package appointments

// import (
// 	"context"
// 	"time"

// 	"github.com/google/uuid"
// 	db "github.com/vention/booking_api/internal/repository"
// 	svcCommon "github.com/vention/booking_api/internal/services/common"
// )

// // validateAppointmentTime validates the appointment time range
// func (s *service) validateAppointmentTime(startTime, endTime time.Time) error {
// 	now := time.Now()

// 	// Check if start time is in the future
// 	if startTime.Before(now) {
// 		return svcCommon.ErrPastTime
// 	}

// 	// Check if end time is after start time
// 	if endTime.Before(startTime) || endTime.Equal(startTime) {
// 		return svcCommon.ErrInvalidTimeRange
// 	}

// 	return nil
// }

// // validateAppointmentConflict checks if the client already has an appointment at the same time with the same professional
// func (s *service) validateAppointmentConflict(ctx context.Context, clientID, professionalID uuid.UUID, startTime time.Time) error {
// 	hasConflict, err := s.repo.CheckClientAppointmentConflict(ctx, &db.CheckClientAppointmentConflictParams{
// 		ClientID:       uuid.NullUUID{UUID: clientID, Valid: true},
// 		ProfessionalID: professionalID,
// 		StartTime:      startTime,
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	if hasConflict {
// 		return svcCommon.ErrAppointmentTimeConflict
// 	}

// 	return nil
// }
