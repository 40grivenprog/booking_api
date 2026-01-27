package clients

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	db "github.com/vention/booking_api/internal/repository"
	"github.com/vention/booking_api/internal/util"
)

// Service defines the business logic operations for clients
type Service interface {
	CreateAppointment(ctx context.Context, input CreateAppointmentInput) (*db.CreateAppointmentWithDetailsRow, error)
	RegisterClient(ctx context.Context, input RegisterClientInput) (*db.Client, error)
	GetClientAppointments(ctx context.Context, clientID uuid.UUID, statusFilter string, page, pageSize int) ([]*db.GetAppointmentsByClientWithStatusRow, int, error)
	CancelAppointment(ctx context.Context, input CancelAppointmentInput) (*db.CancelAppointmentByClientWithDetailsRow, error)
	UpdateLocale(ctx context.Context, input UpdateLocaleInput) error
}

type service struct {
	repo ClientsRepository
}

// NewService creates a new clients service
func NewService(repo ClientsRepository) Service {
	return &service{
		repo: repo,
	}
}

// CreateAppointment creates a new appointment with business logic validation
func (s *service) CreateAppointment(ctx context.Context, input CreateAppointmentInput) (*db.CreateAppointmentWithDetailsRow, error) {
	// Convert times to application timezone (business rule)
	startTime := util.ConvertToAppTimezone(input.StartTime)
	endTime := util.ConvertToAppTimezone(input.EndTime)

	// Validate appointment time
	if err := s.validateAppointmentTime(startTime, endTime); err != nil {
		return nil, err
	}

	// Check for appointment conflicts
	if err := s.validateAppointmentConflict(ctx, input.ClientID, input.ProfessionalID, startTime); err != nil {
		return nil, err
	}

	// Create appointment in database
	result, err := s.repo.CreateAppointmentWithDetails(ctx, &db.CreateAppointmentWithDetailsParams{
		ClientID:       uuid.NullUUID{UUID: input.ClientID, Valid: true},
		ProfessionalID: input.ProfessionalID,
		StartTime:      startTime,
		EndTime:        endTime,
		Description:    sql.NullString{String: input.Description, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

// RegisterClient registers a new client
func (s *service) RegisterClient(ctx context.Context, input RegisterClientInput) (*db.Client, error) {
	params := &db.CreateClientParams{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Locale:    sql.NullString{String: input.Locale, Valid: true},
	}

	// Set optional phone number
	if input.PhoneNumber != "" {
		params.PhoneNumber.String = input.PhoneNumber
		params.PhoneNumber.Valid = true
	}

	// Set optional chat ID
	if input.ChatID != 0 {
		params.ChatID.Int64 = input.ChatID
		params.ChatID.Valid = true
	}

	// CreatedBy is NULL for self-registration
	params.CreatedBy = uuid.NullUUID{}
	fmt.Println("params.Locale", params.Locale)

	client, err := s.repo.CreateClient(ctx, params)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// GetClientAppointments retrieves appointments for a client with optional status filter and pagination
func (s *service) GetClientAppointments(ctx context.Context, clientID uuid.UUID, statusFilter string, page, pageSize int) ([]*db.GetAppointmentsByClientWithStatusRow, int, error) {
	offset := (page - 1) * pageSize

	params := &db.GetAppointmentsByClientWithStatusParams{
		ClientID: uuid.NullUUID{UUID: clientID, Valid: true},
		Limit:    int32(pageSize),
		Offset:   int32(offset),
	}

	// Set optional status filter
	if statusFilter != "" {
		params.Status = db.NullAppointmentStatus{
			AppointmentStatus: db.AppointmentStatus(statusFilter),
			Valid:             true,
		}
	}

	appointments, err := s.repo.GetAppointmentsByClientWithStatus(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	// Get total count for pagination
	countParams := &db.CountClientAppointmentsWithStatusParams{
		ClientID: uuid.NullUUID{UUID: clientID, Valid: true},
	}
	if statusFilter != "" {
		countParams.Status = db.NullAppointmentStatus{
			AppointmentStatus: db.AppointmentStatus(statusFilter),
			Valid:             true,
		}
	}

	total, err := s.repo.CountClientAppointmentsWithStatus(ctx, countParams)
	if err != nil {
		return nil, 0, err
	}

	return appointments, int(total), nil
}

// CancelAppointment cancels an appointment with business logic validation
func (s *service) CancelAppointment(ctx context.Context, input CancelAppointmentInput) (*db.CancelAppointmentByClientWithDetailsRow, error) {
	// Get appointment for validation
	appointment, err := s.repo.GetAppointmentByID(ctx, input.AppointmentID)
	if err != nil {
		return nil, err
	}

	// Validate ownership
	if err := s.validateAppointmentOwnership(appointment, input.ClientID); err != nil {
		return nil, err
	}

	// Validate status
	if err := s.validateAppointmentCancellable(appointment); err != nil {
		return nil, err
	}

	// Cancel appointment
	result, err := s.repo.CancelAppointmentByClientWithDetails(ctx, &db.CancelAppointmentByClientWithDetailsParams{
		ID: input.AppointmentID,
		CancelledByClientID: uuid.NullUUID{
			UUID:  input.ClientID,
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

func (s *service) UpdateLocale(ctx context.Context, input UpdateLocaleInput) error {
	return s.repo.UpdateClientLocale(ctx, &db.UpdateClientLocaleParams{
		ID:     input.ClientID,
		Locale: sql.NullString{String: input.Locale, Valid: true},
	})
}
