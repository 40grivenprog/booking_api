package clients

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	db "github.com/vention/booking_api/internal/repository"
)

// Service defines the business logic operations for clients
type Service interface {
	RegisterClient(ctx context.Context, input RegisterClientInput) (*db.Client, error)
	GetClientAppointments(ctx context.Context, clientID uuid.UUID, statusFilter string) ([]*db.GetAppointmentsByClientWithStatusRow, error)
	CancelAppointment(ctx context.Context, input CancelAppointmentInput) (*db.CancelAppointmentByClientWithDetailsRow, error)
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

// RegisterClient registers a new client
func (s *service) RegisterClient(ctx context.Context, input RegisterClientInput) (*db.Client, error) {
	params := &db.CreateClientParams{
		FirstName: input.FirstName,
		LastName:  input.LastName,
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

	client, err := s.repo.CreateClient(ctx, params)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// GetClientAppointments retrieves appointments for a client with optional status filter
func (s *service) GetClientAppointments(ctx context.Context, clientID uuid.UUID, statusFilter string) ([]*db.GetAppointmentsByClientWithStatusRow, error) {
	params := &db.GetAppointmentsByClientWithStatusParams{
		ClientID: uuid.NullUUID{UUID: clientID, Valid: true},
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
		return nil, err
	}

	return appointments, nil
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
