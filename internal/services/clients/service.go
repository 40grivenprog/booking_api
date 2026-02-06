package clients

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	db "github.com/vention/booking_api/internal/repository"
	svcCommon "github.com/vention/booking_api/internal/services/common"
	"github.com/vention/booking_api/internal/util"
)

// Service defines the business logic operations for clients
type Service interface {
	CreateAppointment(ctx context.Context, input CreateAppointmentInput) (uuid.UUID, error)
	RegisterClient(ctx context.Context, input RegisterClientInput) (*db.Client, error)
	GetClientAppointments(ctx context.Context, clientID uuid.UUID, statusFilter string, page, pageSize int) ([]*db.GetAppointmentsByClientWithStatusRow, int, error)
	CancelAppointment(ctx context.Context, input CancelAppointmentInput) (*CancelAppointmentOutput, error)
	UpdateLocale(ctx context.Context, input UpdateLocaleInput) error
	SubscribeToProfessional(ctx context.Context, input SubscribeToProfessionalInput) error
	UnsubscribeFromProfessional(ctx context.Context, input UnsubscribeFromProfessionalInput) error
	GetSubscribedProfessionals(ctx context.Context, clientID uuid.UUID, page, pageSize int) ([]*db.GetSubscriptionsByClientIDRow, int, error)
	GetAllProfessionals(ctx context.Context, dbParams *db.GetProfessionalsParams) ([]*db.GetProfessionalsRow, int, error)
	GetClientInvites(ctx context.Context, clientID uuid.UUID) ([]*db.GetClientInvitesRow, error)
	GetClientInvite(ctx context.Context, clientID uuid.UUID, inviteID uuid.UUID) (*db.GetInviteByIDRow, error)
	DeleteClientInvite(ctx context.Context, clientID uuid.UUID, inviteID uuid.UUID) error
	AcceptClientInvite(ctx context.Context, input AcceptClientInviteInput) (*db.GetInfoForAcceptInviteNotificationRow, error)
	GetPreviosAppointments(ctx context.Context, clientID uuid.UUID, page, pageSize int) ([]*db.GetPreviousAppointmentsByClientIDRow, int, error)
}

type service struct {
	repo     ClientsRepository
	database *sql.DB
}

// NewService creates a new clients service
func NewService(repo ClientsRepository, database *sql.DB) Service {
	return &service{
		repo:     repo,
		database: database,
	}
}

// CreateAppointment creates a new personal appointment with business logic validation
// It creates both the appointment and client_appointments record in a transaction
func (s *service) CreateAppointment(ctx context.Context, input CreateAppointmentInput) (uuid.UUID, error) {
	// Convert times to application timezone (business rule)
	startTime := util.ConvertToAppTimezone(input.StartTime)
	endTime := util.ConvertToAppTimezone(input.EndTime)

	// Validate appointment time
	if err := s.validateAppointmentTime(startTime, endTime); err != nil {
		return uuid.Nil, err
	}

	// Check for appointment conflicts
	if err := s.validateAppointmentConflict(ctx, input.ClientID, input.ProfessionalID, startTime); err != nil {
		return uuid.Nil, err
	}

	var appointmentID uuid.UUID

	// Use transaction to create appointment and client_appointments atomically
	err := db.WithTransaction(ctx, s.database, func(q *db.Queries) error {
		// Create personal appointment
		createdAppointment, err := q.CreatePersonalAppointment(ctx, &db.CreatePersonalAppointmentParams{
			ProfessionalID: input.ProfessionalID,
			StartTime:      startTime,
			EndTime:        endTime,
		})
		if err != nil {
			return err
		}

		// Create client_appointments link
		if err := q.CreateClientAppointment(ctx, &db.CreateClientAppointmentParams{
			ClientID:      input.ClientID,
			AppointmentID: createdAppointment.ID,
		}); err != nil {
			return err
		}

		appointmentID = createdAppointment.ID

		return nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	return appointmentID, nil
}

// RegisterClient registers a new client
func (s *service) RegisterClient(ctx context.Context, input RegisterClientInput) (*db.Client, error) {
	params := &db.CreateClientParams{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Locale:    input.Locale,
	}

	// Set optional chat ID
	if input.ChatID != 0 {
		params.ChatID.Int64 = input.ChatID
		params.ChatID.Valid = true
	}

	client, err := s.repo.CreateClient(ctx, params)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// // GetClientAppointments retrieves appointments for a client with optional status filter and pagination
func (s *service) GetClientAppointments(ctx context.Context, clientID uuid.UUID, statusFilter string, page, pageSize int) ([]*db.GetAppointmentsByClientWithStatusRow, int, error) {
	offset := (page - 1) * pageSize

	params := &db.GetAppointmentsByClientWithStatusParams{
		ClientID: clientID,
		Limit:    int32(pageSize),
		Offset:   int32(offset),
		Status:   statusFilter,
	}

	appointments, err := s.repo.GetAppointmentsByClientWithStatus(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	// Get total count for pagination
	countParams := &db.CountClientAppointmentsWithStatusParams{
		ClientID: clientID,
		Status:   statusFilter,
	}

	total, err := s.repo.CountClientAppointmentsWithStatus(ctx, countParams)
	if err != nil {
		return nil, 0, err
	}

	return appointments, int(total), nil
}

// CancelAppointment cancels an appointment with business logic validation
func (s *service) CancelAppointment(ctx context.Context, input CancelAppointmentInput) (*CancelAppointmentOutput, error) {
	professionalInfo, err := s.repo.GetProfessionalInfoForNotificationByAppointmentID(ctx, input.AppointmentID)
	if err != nil {
		return nil, err
	}

	appointmentInfo, err := s.repo.GetAppointmentInfoByAppointmentID(ctx, input.AppointmentID)
	if err != nil {
		return nil, err
	}

	err = s.repo.DeleteAppointmentById(ctx, input.AppointmentID)
	if err != nil {
		return nil, err
	}
	result := &CancelAppointmentOutput{
		StartTime:          appointmentInfo.StartTime,
		EndTime:            appointmentInfo.EndTime,
		ProfessionalChatID: professionalInfo.ChatID,
		ProfessionalLocale: professionalInfo.Locale,
	}

	return result, nil
}

// UpdateLocale updates the locale for a client
func (s *service) UpdateLocale(ctx context.Context, input UpdateLocaleInput) error {
	return s.repo.UpdateClientLocale(ctx, &db.UpdateClientLocaleParams{
		ID:     input.ClientID,
		Locale: input.Locale,
	})
}

// SubscribeToProfessional subscribes a client to a professional
func (s *service) SubscribeToProfessional(ctx context.Context, input SubscribeToProfessionalInput) error {
	return s.repo.CreateSubscription(ctx, &db.CreateSubscriptionParams{
		ClientID:       input.ClientID,
		ProfessionalID: input.ProfessionalID,
	})
}

// UnsubscribeFromProfessional unsubscribes a client from a professional
func (s *service) UnsubscribeFromProfessional(ctx context.Context, input UnsubscribeFromProfessionalInput) error {
	return s.repo.DeleteSubscription(ctx, &db.DeleteSubscriptionParams{
		ClientID:       input.ClientID,
		ProfessionalID: input.ProfessionalID,
	})
}

// GetSubscribedProfessionals gets the professionals a client is subscribed to
func (s *service) GetSubscribedProfessionals(ctx context.Context, clientID uuid.UUID, page, pageSize int) ([]*db.GetSubscriptionsByClientIDRow, int, error) {
	offset := (page - 1) * pageSize

	professionals, err := s.repo.GetSubscriptionsByClientID(ctx, &db.GetSubscriptionsByClientIDParams{
		ClientID: clientID,
		Limit:    int32(pageSize),
		Offset:   int32(offset),
	})
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.CountSubscriptionsByClientID(ctx, clientID)
	if err != nil {
		return nil, 0, err
	}

	return professionals, int(total), nil
}

// GetAllProfessionals gets all professionals
func (s *service) GetAllProfessionals(ctx context.Context, dbParams *db.GetProfessionalsParams) ([]*db.GetProfessionalsRow, int, error) {
	professionals, err := s.repo.GetProfessionals(ctx, dbParams)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.CountProfessionals(ctx, dbParams.ClientID)
	if err != nil {
		return nil, 0, err
	}

	return professionals, int(total), nil
}

// GetClientInvites gets the invites for a client
func (s *service) GetClientInvites(ctx context.Context, clientID uuid.UUID) ([]*db.GetClientInvitesRow, error) {
	invites, err := s.repo.GetClientInvites(ctx, clientID)
	if err != nil {
		return nil, err
	}
	return invites, nil
}

// GetClientInvite gets an invite for a client
func (s *service) GetClientInvite(ctx context.Context, clientID uuid.UUID, inviteID uuid.UUID) (*db.GetInviteByIDRow, error) {
	invite, err := s.repo.GetInviteByID(ctx, &db.GetInviteByIDParams{
		ID:       inviteID,
		ClientID: clientID,
	})
	if err != nil {
		return nil, err
	}

	return invite, nil
}

// DeleteClientInvite deletes an invite for a client
func (s *service) DeleteClientInvite(ctx context.Context, clientID uuid.UUID, inviteID uuid.UUID) error {
	return s.repo.DeleteInviteByID(ctx, &db.DeleteInviteByIDParams{
		ID:       inviteID,
		ClientID: clientID,
	})
}

// AcceptClientInvite accepts a client invite
func (s *service) AcceptClientInvite(ctx context.Context, input AcceptClientInviteInput) (*db.GetInfoForAcceptInviteNotificationRow, error) {
	var resultErr error
	err := db.WithTransaction(ctx, s.database, func(q *db.Queries) error {
		// If type is "split", check if appointment is already fully booked (>= 2 clients)
		if input.Type == "split" {
			// Get all clients for this appointment
			clients, err := q.GetAppointmentClientsByAppointmentID(ctx, input.AppointmentID)
			if err != nil {
				return err
			}

			// If already 2 or more clients, appointment is fully booked
			if len(clients) >= 2 {
				// Delete invite
				if err := q.DeleteInviteByID(ctx, &db.DeleteInviteByIDParams{
					ID:       input.InviteID,
					ClientID: input.ClientID,
				}); err != nil {
					return err
				}
				resultErr = svcCommon.ErrAppointmentFullyBooked
				return resultErr
			}
		}

		// Create client_appointments link
		if err := q.CreateClientAppointment(ctx, &db.CreateClientAppointmentParams{
			ClientID:      input.ClientID,
			AppointmentID: input.AppointmentID,
		}); err != nil {
			return err
		}

		// Delete invite
		if err := q.DeleteInviteByID(ctx, &db.DeleteInviteByIDParams{
			ID:       input.InviteID,
			ClientID: input.ClientID,
		}); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	// Get appointment info first to check type
	notificationInfo, err := s.repo.GetInfoForAcceptInviteNotification(ctx, input.AppointmentID)
	if err != nil {
		return nil, err
	}

	return notificationInfo, nil
}

// GetPreviosAppointments gets the previous appointments for a client
func (s *service) GetPreviosAppointments(ctx context.Context, clientID uuid.UUID, page, pageSize int) ([]*db.GetPreviousAppointmentsByClientIDRow, int, error) {
	offset := (page - 1) * pageSize

	appointments, err := s.repo.GetPreviousAppointmentsByClientID(ctx, &db.GetPreviousAppointmentsByClientIDParams{
		ClientID: clientID,
		Limit:    int32(pageSize),
		Offset:   int32(offset),
	})
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.CountPreviousAppointmentsByClientID(ctx, clientID)
	if err != nil {
		return nil, 0, err
	}

	return appointments, int(total), nil
}
