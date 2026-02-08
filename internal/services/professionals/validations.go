package professionals

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	db "github.com/vention/booking_api/internal/repository"
	svcCommon "github.com/vention/booking_api/internal/services/common"
	"golang.org/x/crypto/bcrypt"
)

// validatePassword validates the professional's password
func (s *service) validatePassword(professional *db.Professional, password string) error {
	// Compare password with hash
	if err := bcrypt.CompareHashAndPassword([]byte(professional.PasswordHash), []byte(password)); err != nil {
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
	if appointment.Status != "pending" {
		return svcCommon.ErrAppointmentNotPending
	}
	return nil
}

// validateAppointmentCancellable validates that the appointment can be cancelled
func (s *service) validateAppointmentCancellable(appointment *db.Appointment) error {
	if appointment.Status != "pending" &&
		appointment.Status != "confirmed" {
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

// validateAppointmentNotInPast validates that the appointment is not in the past
func (s *service) validateAppointmentNotInPast(appointment *db.Appointment) error {
	if appointment.EndTime.Before(time.Now()) {
		return svcCommon.ErrPastAppointment
	}
	return nil
}

// validateTypeForAppointmentUpdate validates the type for an appointment update
func (s *service) validateTypeForAppointmentUpdate(ctx context.Context, appointmentID uuid.UUID, apptType string) error {
	clientsCount, err := s.repo.GetAppintmentClientsCountByAppointmentID(ctx, appointmentID)
	if err != nil {
		return err
	}

	switch apptType {
	case "personal":
		if clientsCount != 1 {
			return errors.New("personal appointment must have exactly one client")
		}
	case "split":
		if clientsCount != 2 {
			return errors.New("split appointment must have exactly 2 clients")
		}
	case "group":
		if clientsCount < 3 {
			return errors.New("group appointment must have at least 3 clients")
		}
	default:
		return errors.New("invalid appointment type")
	}

	return nil
}

// validateTypeForPreviousAppointmentUpdate validates the type based on final client count after add/remove
func (s *service) validateTypeForPreviousAppointmentUpdate(currentClientsCount int, addedCount int, removedCount int, apptType string) error {
	finalClientsCount := currentClientsCount + addedCount - removedCount

	switch apptType {
	case "personal":
		if finalClientsCount != 1 {
			return errors.New("personal appointment must have exactly one client")
		}
	case "split":
		if finalClientsCount != 2 {
			return errors.New("split appointment must have exactly 2 clients")
		}
	case "group":
		if finalClientsCount < 3 {
			return errors.New("group appointment must have at least 3 clients")
		}
	default:
		return errors.New("invalid appointment type")
	}

	return nil
}

// validatePreviousPackage validates the previous package
func (s *service) validatePreviousPackage(ctx context.Context, clientID uuid.UUID, professionalID uuid.UUID) (uuid.UUID, error) {
	currentPackage, err := s.repo.GetPackageByClientIDAndProfessionalID(ctx, &db.GetPackageByClientIDAndProfessionalIDParams{
		ClientID:       clientID,
		ProfessionalID: professionalID,
	})
	if errors.Is(err, sql.ErrNoRows) {
		return uuid.Nil, nil
	}

	if err != nil {
		return uuid.Nil, err
	}

	if currentPackage.DeactivatedAt.Valid {
		return uuid.Nil, nil
	}

	currentAppointmentsCount, err := s.repo.GetAppintmentNumberByClientIDAndProfessionalID(ctx, &db.GetAppintmentNumberByClientIDAndProfessionalIDParams{
		ProfessionalID: professionalID,
		ClientID:       clientID,
		StartTime:      currentPackage.IssuedAt,
	})
	if err != nil {
		return uuid.Nil, err
	}

	if int(currentAppointmentsCount) < int(currentPackage.ApppointmentsNumber) && currentPackage.ExpiresAt.After(time.Now()) {
		return uuid.Nil, svcCommon.ErrClientsCurrentPackageIsActive
	}

	return currentPackage.ID, nil
}
