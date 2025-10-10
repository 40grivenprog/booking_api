package appointments

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	db "github.com/vention/booking_api/internal/repository"
)

// Service defines the business logic operations for appointments
type Service interface {
	CreateAppointment(ctx context.Context, input CreateAppointmentInput) (*db.CreateAppointmentWithDetailsRow, error)
}

type service struct {
	repo AppointmentsRepository
}

// NewService creates a new appointments service
func NewService(repo AppointmentsRepository) Service {
	return &service{
		repo: repo,
	}
}

// CreateAppointment creates a new appointment with business logic validation
func (s *service) CreateAppointment(ctx context.Context, input CreateAppointmentInput) (*db.CreateAppointmentWithDetailsRow, error) {
	if err := s.validateAppointmentTime(input.StartTime, input.EndTime); err != nil {
		return nil, err
	}

	// Create appointment in database
	result, err := s.repo.CreateAppointmentWithDetails(ctx, &db.CreateAppointmentWithDetailsParams{
		ClientID:       uuid.NullUUID{UUID: input.ClientID, Valid: true},
		ProfessionalID: input.ProfessionalID,
		StartTime:      input.StartTime,
		EndTime:        input.EndTime,
		Description:    sql.NullString{String: input.Description, Valid: input.Description != ""},
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}
