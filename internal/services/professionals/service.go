package professionals

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	db "github.com/vention/booking_api/internal/repository"
	"github.com/vention/booking_api/internal/util"
)

// Service defines the business logic operations for professionals
type Service interface {
	GetProfessionals(ctx context.Context) ([]*db.Professional, error)
	SignIn(ctx context.Context, input SignInInput) (*db.Professional, error)
	ConfirmAppointment(ctx context.Context, input ConfirmAppointmentInput) (*db.ConfirmAppointmentWithDetailsRow, error)
	GetAppointments(ctx context.Context, professionalID uuid.UUID, statusFilter, dateFilter string) ([]*db.GetAppointmentsByProfessionalWithStatusAndDateRow, error)
	GetAppointmentDates(ctx context.Context, professionalID uuid.UUID, month time.Time) ([]time.Time, error)
	CancelAppointment(ctx context.Context, input CancelAppointmentInput) (*db.CancelAppointmentByProfessionalWithDetailsRow, error)
	CreateUnavailableAppointment(ctx context.Context, input CreateUnavailableAppointmentInput) (*db.Appointment, error)
	GetAvailability(ctx context.Context, professionalID uuid.UUID, date time.Time) ([]*db.GetAppointmentsByProfessionalAndDateWithClientRow, error)
	GetTimetable(ctx context.Context, professionalID uuid.UUID, date time.Time) ([]*db.GetProfessionalTimetableRow, error)
	GenerateAvailabilitySlots(date time.Time, appointments []*db.GetAppointmentsByProfessionalAndDateWithClientRow, config AvailabilityConfig) []TimeSlot
}

type service struct {
	repo ProfessionalsRepository
}

// NewService creates a new professionals service
func NewService(repo ProfessionalsRepository) Service {
	return &service{
		repo: repo,
	}
}

// GetProfessionals retrieves all professionals
func (s *service) GetProfessionals(ctx context.Context) ([]*db.Professional, error) {
	return s.repo.GetProfessionals(ctx)
}

// SignIn authenticates a professional and updates their chat ID
func (s *service) SignIn(ctx context.Context, input SignInInput) (*db.Professional, error) {
	// Get professional by username
	professional, err := s.repo.GetProfessionalByUsername(ctx, input.Username)
	if err != nil {
		return nil, err
	}

	// Validate password
	if err := s.validatePassword(professional, input.Password); err != nil {
		return nil, err
	}

	// Update chat ID
	updatedProfessional, err := s.repo.UpdateProfessionalChatID(ctx, &db.UpdateProfessionalChatIDParams{
		ID: professional.ID,
		ChatID: sql.NullInt64{
			Int64: input.ChatID,
			Valid: input.ChatID != 0,
		},
	})
	if err != nil {
		return nil, err
	}

	return updatedProfessional, nil
}

// ConfirmAppointment confirms an appointment with validation
func (s *service) ConfirmAppointment(ctx context.Context, input ConfirmAppointmentInput) (*db.ConfirmAppointmentWithDetailsRow, error) {
	// Get appointment
	appointment, err := s.repo.GetAppointmentByID(ctx, input.AppointmentID)
	if err != nil {
		return nil, err
	}

	// Validate ownership
	if err := s.validateAppointmentOwnership(appointment, input.ProfessionalID); err != nil {
		return nil, err
	}

	// Validate status
	if err := s.validateAppointmentPending(appointment); err != nil {
		return nil, err
	}

	// Confirm appointment
	result, err := s.repo.ConfirmAppointmentWithDetails(ctx, &db.ConfirmAppointmentWithDetailsParams{
		ID:             input.AppointmentID,
		ProfessionalID: input.ProfessionalID,
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetAppointments retrieves appointments with optional filters
func (s *service) GetAppointments(ctx context.Context, professionalID uuid.UUID, statusFilter, dateFilter string) ([]*db.GetAppointmentsByProfessionalWithStatusAndDateRow, error) {
	return s.repo.GetAppointmentsByProfessionalWithStatusAndDate(ctx, &db.GetAppointmentsByProfessionalWithStatusAndDateParams{
		ProfessionalID: professionalID,
		Column2:        statusFilter,
		Column3:        dateFilter,
	})
}

// GetAppointmentDates retrieves distinct dates with appointments for a month
func (s *service) GetAppointmentDates(ctx context.Context, professionalID uuid.UUID, month time.Time) ([]time.Time, error) {
	// Normalize to start of month in application timezone (business rule)
	startOfMonth := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, util.GetAppTimezone())
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	return s.repo.GetProfessionalAppointmentDates(ctx, &db.GetProfessionalAppointmentDatesParams{
		ProfessionalID: professionalID,
		StartTime:      startOfMonth,
		StartTime_2:    endOfMonth,
	})
}

// CancelAppointment cancels an appointment with validation
func (s *service) CancelAppointment(ctx context.Context, input CancelAppointmentInput) (*db.CancelAppointmentByProfessionalWithDetailsRow, error) {
	// Get appointment
	appointment, err := s.repo.GetAppointmentByID(ctx, input.AppointmentID)
	if err != nil {
		return nil, err
	}

	// Validate ownership
	if err := s.validateAppointmentOwnership(appointment, input.ProfessionalID); err != nil {
		return nil, err
	}

	// Validate status
	if err := s.validateAppointmentCancellable(appointment); err != nil {
		return nil, err
	}

	// Cancel appointment
	result, err := s.repo.CancelAppointmentByProfessionalWithDetails(ctx, &db.CancelAppointmentByProfessionalWithDetailsParams{
		ID: input.AppointmentID,
		CancelledByProfessionalID: uuid.NullUUID{
			UUID:  input.ProfessionalID,
			Valid: true,
		},
		CancellationReason: sql.NullString{
			String: input.CancellationReason,
			Valid:  input.CancellationReason != "",
		},
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

// CreateUnavailableAppointment creates an unavailable time slot with validation
func (s *service) CreateUnavailableAppointment(ctx context.Context, input CreateUnavailableAppointmentInput) (*db.Appointment, error) {
	// Validate time range
	if err := s.validateTimeRange(input.StartTime, input.EndTime); err != nil {
		return nil, err
	}

	// Create unavailable appointment
	appointment, err := s.repo.CreateUnavailableAppointment(ctx, &db.CreateUnavailableAppointmentParams{
		ProfessionalID: input.ProfessionalID,
		StartTime:      input.StartTime,
		EndTime:        input.EndTime,
		Description: sql.NullString{
			String: input.Description,
			Valid:  input.Description != "",
		},
	})
	if err != nil {
		return nil, err
	}

	return appointment, nil
}

// GetAvailability retrieves appointments for availability calculation
func (s *service) GetAvailability(ctx context.Context, professionalID uuid.UUID, date time.Time) ([]*db.GetAppointmentsByProfessionalAndDateWithClientRow, error) {
	return s.repo.GetAppointmentsByProfessionalAndDateWithClient(ctx, &db.GetAppointmentsByProfessionalAndDateWithClientParams{
		ProfessionalID: professionalID,
		StartTime:      date,
	})
}

// GetTimetable retrieves timetable for a specific date
func (s *service) GetTimetable(ctx context.Context, professionalID uuid.UUID, date time.Time) ([]*db.GetProfessionalTimetableRow, error) {
	return s.repo.GetProfessionalTimetable(ctx, &db.GetProfessionalTimetableParams{
		ProfessionalID: professionalID,
		StartTime:      date,
	})
}
