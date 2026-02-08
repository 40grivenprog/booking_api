package professionals

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	apiCommon "github.com/vention/booking_api/internal/api/common"
	db "github.com/vention/booking_api/internal/repository"
	"github.com/vention/booking_api/internal/services/common"
	"github.com/vention/booking_api/internal/services/notifications"
)

// Service defines the business logic operations for professionals
type Service interface {
	SignIn(ctx context.Context, input SignInInput) (*db.Professional, error)
	ConfirmAppointment(ctx context.Context, input ConfirmAppointmentInput) (*ConfirmAppointmentOutput, error)
	GetAppointments(ctx context.Context, professionalID uuid.UUID, statusFilter, dateFilter string, page, pageSize int) ([]*db.GetAppointmentsByProfessionalWithStatusAndDateRow, int, error)
	// GetAppointmentDates(ctx context.Context, professionalID uuid.UUID, month time.Time) ([]time.Time, error)
	CancelAppointment(ctx context.Context, input CancelAppointmentInput) (*CancelAppointmentOutput, error)
	CreateUnavailableAppointment(ctx context.Context, input CreateUnavailableAppointmentInput) error
	GetAvailability(ctx context.Context, professionalID uuid.UUID, date time.Time) ([]*db.GetAppointmentsByProfessionalByDateRow, error)
	GetTimetable(ctx context.Context, professionalID uuid.UUID, date time.Time) ([]*db.GetProfessionalTimetableRow, error)
	// GetClients(ctx context.Context, professionalID uuid.UUID) ([]*db.GetProfessionalClientsRow, error)
	// GetPreviousAppointmentsByClient(ctx context.Context, professionalID uuid.UUID, clientID uuid.UUID, monthFilter *time.Time) ([]*db.GetPreviousProfessionalAppointmentsByClientRow, error)
	GenerateAvailabilitySlots(date time.Time, appointments []*db.GetAppointmentsByProfessionalByDateRow, config AvailabilityConfig) []TimeSlot
	UpdateLocale(ctx context.Context, input UpdateLocaleInput) error
	CreateGroupVisitAppointment(ctx context.Context, input CreateGroupVisitAppointmentInput) (uuid.UUID, error)
	GetSubscriptionsByProfessionalID(ctx context.Context, professionalID uuid.UUID) ([]*db.GetSubscriptionsByProfessionalIDRow, error)
	GetPreviousAppointments(ctx context.Context, professionalID uuid.UUID, page, pageSize int, dateFrom, dateTo sql.NullTime) ([]*db.GetPreviousAppointmentsByProfessionalIDRow, int, error)
	GetPreviousAppointmentsByClientID(ctx context.Context, professionalID uuid.UUID, clientID uuid.UUID, page, pageSize int, dateFrom, dateTo sql.NullTime) ([]*db.GetPreviousAppointmentsByProfessionalIDAndClientIDRow, int, error)
	GetAppointmentDetailsByProfessionalIDAndAppointmentID(ctx context.Context, professionalID uuid.UUID, appointmentID uuid.UUID) (*GetAppointmentDetailsOutput, error)
	InviteAllClients(ctx context.Context, input InviteAllClientsInput) error
	InvitePartiallySelectedClients(ctx context.Context, input InvitePartiallySelectedClientsInput) error
	UpdateAppointment(ctx context.Context, input UpdateAppointmentInput) error
	GetMissingInviteUsersForAppointment(ctx context.Context, appointmentID uuid.UUID) ([]*db.GetMissingInviteUsersForAppointmentRow, error)
	UpdatePreviousAppointment(ctx context.Context, input UpdatePreviousAppointmentInput) error
	GetMissingClientsForPreviousAppointment(ctx context.Context, appointmentID uuid.UUID, professionalID uuid.UUID) ([]*db.GetMissingClientsForPreviousAppointmentRow, error)
	GetPendingInviteUsersForAppointment(ctx context.Context, appointmentID uuid.UUID) ([]*db.GetPendingInviteUsersForAppointmentRow, error)
	CreatePackage(ctx context.Context, input CreatePackageInput) error
	GetListOfPackages(ctx context.Context, professionalID uuid.UUID) ([]*db.GetListOfPackagesByProfessionalIDRow, error)
	GetPackageDetails(ctx context.Context, packageID uuid.UUID) (*GetPackageDetailsOutput, error)
	GetClientInfoByID(ctx context.Context, clientID uuid.UUID) (*db.GetClientInfoByIDRow, error)
}

type service struct {
	repo     ProfessionalsRepository
	database *sql.DB
}

// NewService creates a new professionals service
func NewService(repo ProfessionalsRepository, database *sql.DB) Service {
	return &service{
		repo:     repo,
		database: database,
	}
}

// SignIn authenticates a professional and updates their chat ID
func (s *service) SignIn(ctx context.Context, input SignInInput) (*db.Professional, error) {
	// Get professional by username
	professional, err := s.repo.GetProfessionalByUsername(ctx, input.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, common.ErrProfessionalNotFound
		}
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
		Locale: input.Locale,
	})
	if err != nil {
		return nil, err
	}

	return updatedProfessional, nil
}

// ConfirmAppointment confirms an appointment with validation
func (s *service) ConfirmAppointment(ctx context.Context, input ConfirmAppointmentInput) (*ConfirmAppointmentOutput, error) {
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
		fmt.Println("appointment not pending", appointment.Status)
		return nil, err
	}

	// Validate appointment is not in the past
	if err := s.validateAppointmentNotInPast(appointment); err != nil {
		return nil, err
	}

	// Confirm appointment
	err = s.repo.ConfirmAppointmentById(ctx, input.AppointmentID)
	if err != nil {
		return nil, err
	}

	// Get clients by appointment ID
	clients, err := s.repo.GetAppointmentClientsByAppointmentID(ctx, input.AppointmentID)
	if err != nil {
		return nil, err
	}

	result := &ConfirmAppointmentOutput{
		StartTime:                 appointment.StartTime,
		EndTime:                   appointment.EndTime,
		ConfirmAppointmentClients: clients,
	}

	return result, nil
}

// GetAppointments retrieves appointments with optional filters and pagination
func (s *service) GetAppointments(ctx context.Context, professionalID uuid.UUID, statusFilter, dateFilter string, page, pageSize int) ([]*db.GetAppointmentsByProfessionalWithStatusAndDateRow, int, error) {
	offset := (page - 1) * pageSize

	params := &db.GetAppointmentsByProfessionalWithStatusAndDateParams{
		ProfessionalID: professionalID,
		Column2:        statusFilter,
		Limit:          int32(pageSize),
		Offset:         int32(offset),
	}

	appointments, err := s.repo.GetAppointmentsByProfessionalWithStatusAndDate(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	// Get total count for pagination
	countParams := &db.CountProfessionalAppointmentsWithStatusAndDateParams{
		ProfessionalID: professionalID,
		Column2:        statusFilter,
	}
	total, err := s.repo.CountProfessionalAppointmentsWithStatusAndDate(ctx, countParams)
	if err != nil {
		return nil, 0, err
	}

	return appointments, int(total), nil
}

// // GetAppointmentDates retrieves distinct dates with appointments for a month
// func (s *service) GetAppointmentDates(ctx context.Context, professionalID uuid.UUID, month time.Time) ([]time.Time, error) {
// 	appTimezone := util.GetAppTimezone()
// 	now := time.Now().In(appTimezone)

// 	// Normalize to start of month in application timezone
// 	startOfMonth := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, appTimezone)
// 	endOfMonth := startOfMonth.AddDate(0, 1, 0)

// 	// If target month is current month, start from today (don't show past dates)
// 	var startTime time.Time
// 	if month.Year() == now.Year() && month.Month() == now.Month() {
// 		// Start from today at 00:00:00
// 		startTime = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, appTimezone)
// 	} else {
// 		// Start from beginning of the month
// 		startTime = startOfMonth
// 	}

// 	return s.repo.GetProfessionalAppointmentDates(ctx, &db.GetProfessionalAppointmentDatesParams{
// 		ProfessionalID: professionalID,
// 		StartTime:      startTime,
// 		StartTime_2:    endOfMonth,
// 	})
// }

// CancelAppointment cancels an appointment with validation
func (s *service) CancelAppointment(ctx context.Context, input CancelAppointmentInput) (*CancelAppointmentOutput, error) {
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

	// Get clients by appointment ID
	clients, err := s.repo.GetAppointmentClientsByAppointmentID(ctx, input.AppointmentID)
	if err != nil {
		return nil, err
	}

	// Delete appointment
	err = s.repo.DeleteAppointmentById(ctx, input.AppointmentID)
	if err != nil {
		return nil, err
	}

	result := &CancelAppointmentOutput{
		StartTime:                appointment.StartTime,
		EndTime:                  appointment.EndTime,
		CancelAppointmentClients: clients,
	}

	return result, nil
}

// CreateUnavailableAppointment creates an unavailable time slot with validation
func (s *service) CreateUnavailableAppointment(ctx context.Context, input CreateUnavailableAppointmentInput) error {
	// Validate time range
	if err := s.validateTimeRange(input.StartTime, input.EndTime); err != nil {
		return err
	}

	// Create unavailable appointment
	err := s.repo.CreateUnavailableAppointment(ctx, &db.CreateUnavailableAppointmentParams{
		ProfessionalID: input.ProfessionalID,
		StartTime:      input.StartTime,
		EndTime:        input.EndTime,
		Description: sql.NullString{
			String: input.Description,
			Valid:  input.Description != "",
		},
	})
	if err != nil {
		return err
	}

	return nil
}

// GetAvailability retrieves appointments for availability calculation
func (s *service) GetAvailability(ctx context.Context, professionalID uuid.UUID, date time.Time) ([]*db.GetAppointmentsByProfessionalByDateRow, error) {
	return s.repo.GetAppointmentsByProfessionalByDate(ctx, &db.GetAppointmentsByProfessionalByDateParams{
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

// CreateGroupVisitAppointment creates a group visit appointment
func (s *service) CreateGroupVisitAppointment(ctx context.Context, input CreateGroupVisitAppointmentInput) (uuid.UUID, error) {
	return s.repo.CreateGroupVisitAppointment(ctx, &db.CreateGroupVisitAppointmentParams{
		ProfessionalID: input.ProfessionalID,
		StartTime:      input.StartTime,
		EndTime:        input.EndTime,
		Type:           input.Type,
	})
}

// UpdateLocale updates the locale of a professional
func (s *service) UpdateLocale(ctx context.Context, input UpdateLocaleInput) error {
	return s.repo.UpdateProfessionalLocale(ctx, &db.UpdateProfessionalLocaleParams{
		ID:     input.ProfessionalID,
		Locale: input.Locale,
	})
}

// GetSubscriptionsByProfessionalID retrieves all subscriptions by professional ID
func (s *service) GetSubscriptionsByProfessionalID(ctx context.Context, professionalID uuid.UUID) ([]*db.GetSubscriptionsByProfessionalIDRow, error) {
	return s.repo.GetSubscriptionsByProfessionalID(ctx, professionalID)
}

// GetPreviousAppointments retrieves previous appointments by professional ID
func (s *service) GetPreviousAppointments(ctx context.Context, professionalID uuid.UUID, page, pageSize int, dateFrom, dateTo sql.NullTime) ([]*db.GetPreviousAppointmentsByProfessionalIDRow, int, error) {
	offset := (page - 1) * pageSize

	// Convert sql.NullTime to time.Time (zero value if not valid)
	var dateFromTime, dateToTime time.Time
	if dateFrom.Valid {
		dateFromTime = dateFrom.Time
	}
	if dateTo.Valid {
		dateToTime = dateTo.Time
	}

	appointments, err := s.repo.GetPreviousAppointmentsByProfessionalID(ctx, &db.GetPreviousAppointmentsByProfessionalIDParams{
		ProfessionalID: professionalID,
		Column2:        dateFromTime,
		Column3:        dateToTime,
		Limit:          int32(pageSize),
		Offset:         int32(offset),
	})
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.CountPreviousAppointmentsByProfessionalID(ctx, &db.CountPreviousAppointmentsByProfessionalIDParams{
		ProfessionalID: professionalID,
		Column2:        dateFromTime,
		Column3:        dateToTime,
	})
	if err != nil {
		return nil, 0, err
	}

	return appointments, int(total), nil
}

// GetPreviousAppointmentsByClientID retrieves previous appointments by professional ID and client ID
func (s *service) GetPreviousAppointmentsByClientID(ctx context.Context, professionalID uuid.UUID, clientID uuid.UUID, page, pageSize int, dateFrom, dateTo sql.NullTime) ([]*db.GetPreviousAppointmentsByProfessionalIDAndClientIDRow, int, error) {
	offset := (page - 1) * pageSize

	// Convert sql.NullTime to time.Time (zero value if not valid)
	var dateFromTime, dateToTime time.Time
	if dateFrom.Valid {
		dateFromTime = dateFrom.Time
	}
	if dateTo.Valid {
		dateToTime = dateTo.Time
	}

	appointments, err := s.repo.GetPreviousAppointmentsByProfessionalIDAndClientID(ctx, &db.GetPreviousAppointmentsByProfessionalIDAndClientIDParams{
		ProfessionalID: professionalID,
		ClientID:       clientID,
		Column3:        dateFromTime,
		Column4:        dateToTime,
		Limit:          int32(pageSize),
		Offset:         int32(offset),
	})
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.CountPreviousAppointmentsByProfessionalIDAndClientID(ctx, &db.CountPreviousAppointmentsByProfessionalIDAndClientIDParams{
		ProfessionalID: professionalID,
		ClientID:       clientID,
		Column3:        dateFromTime,
		Column4:        dateToTime,
	})
	if err != nil {
		return nil, 0, err
	}

	return appointments, int(total), nil
}

// GetAppointmentDetailsByProfessionalIDAndAppointmentID retrieves appointment details by professional ID and appointment ID
func (s *service) GetAppointmentDetailsByProfessionalIDAndAppointmentID(ctx context.Context, professionalID uuid.UUID, appointmentID uuid.UUID) (*GetAppointmentDetailsOutput, error) {
	appointmentDetails, err := s.repo.GetAppointmentDetailsByProfessionalIDAndAppointmentID(ctx, &db.GetAppointmentDetailsByProfessionalIDAndAppointmentIDParams{
		ProfessionalID: professionalID,
		ID:             appointmentID,
	})

	if err != nil {
		return nil, err
	}

	clients, err := s.repo.GetAppointmentClientsByAppointmentID(ctx, appointmentID)
	if err != nil {
		return nil, err
	}

	result := &GetAppointmentDetailsOutput{
		ID:          appointmentDetails.ID,
		StartTime:   appointmentDetails.StartTime,
		EndTime:     appointmentDetails.EndTime,
		Description: appointmentDetails.Description,
		Type:        appointmentDetails.Type,
		Clients:     clients,
	}

	return result, nil
}

// CreateInvite creates an invite
func (s *service) InviteAllClients(ctx context.Context, input InviteAllClientsInput) error {
	clients, err := s.GetSubscriptionsByProfessionalID(ctx, input.ProfessionalID)
	if err != nil {
		return err
	}

	appointmentIDs := make([]uuid.UUID, len(clients))
	startTime := make([]time.Time, len(clients))
	endTime := make([]time.Time, len(clients))
	professionalNames := make([]string, len(clients))
	types := make([]string, len(clients))
	descriptions := make([]string, len(clients))
	clientIDs := make([]uuid.UUID, len(clients))

	for i, client := range clients {
		appointmentIDs[i] = input.AppointmentID
		startTime[i] = input.StartTime
		endTime[i] = input.EndTime
		professionalNames[i] = input.ProfessionalName
		types[i] = input.Type
		descriptions[i] = input.Description
		clientIDs[i] = client.ID
	}
	err = s.repo.CreateInvites(ctx, &db.CreateInvitesParams{
		Column1: appointmentIDs,
		Column2: startTime,
		Column3: endTime,
		Column4: clientIDs,
		Column5: descriptions,
		Column6: types,
		Column7: professionalNames,
	})
	if err != nil {
		return err
	}

	invites, err := s.repo.GetInvitesByAppointmentIDAndClientIs(ctx, &db.GetInvitesByAppointmentIDAndClientIsParams{
		AppointmentID: input.AppointmentID,
		Column2:       clientIDs,
	})
	if err != nil {
		return err
	}

	mapClientIdInviteID := buildMapClientIdInviteID(invites)

	for _, client := range clients {
		err = input.NotificationsService.SendGroupVisitAppointmentNotification(ctx, notifications.SendGroupVisitAppointmentNotificationInput{
			ChatID:           apiCommon.Int64Value(apiCommon.FromNullInt64(client.ChatID)),
			StartTime:        apiCommon.FormatTimeRFC3339(input.StartTime),
			EndTime:          apiCommon.FormatTimeRFC3339(input.EndTime),
			ProfessionalName: input.ProfessionalName,
			Locale:           client.Locale,
			Description:      input.Description,
			InviteID:         mapClientIdInviteID[client.ID],
		})
	}

	return nil
}

// InvitePartiallySelectedClients invites partially selected clients
func (s *service) InvitePartiallySelectedClients(ctx context.Context, input InvitePartiallySelectedClientsInput) error {
	appointmentIDs := make([]uuid.UUID, len(input.Clients))
	startTime := make([]time.Time, len(input.Clients))
	endTime := make([]time.Time, len(input.Clients))
	professionalNames := make([]string, len(input.Clients))
	types := make([]string, len(input.Clients))
	descriptions := make([]string, len(input.Clients))
	clientIDs := make([]uuid.UUID, len(input.Clients))

	for i, client := range input.Clients {
		appointmentIDs[i] = input.AppointmentID
		startTime[i] = input.StartTime
		endTime[i] = input.EndTime
		professionalNames[i] = input.ProfessionalName
		types[i] = input.Type
		descriptions[i] = input.Description
		clientIDs[i] = client.ID
	}

	err := s.repo.CreateInvites(ctx, &db.CreateInvitesParams{
		Column1: appointmentIDs,
		Column2: startTime,
		Column3: endTime,
		Column4: clientIDs,
		Column5: descriptions,
		Column6: types,
		Column7: professionalNames,
	})
	if err != nil {
		return err
	}
	invites, err := s.repo.GetInvitesByAppointmentIDAndClientIs(ctx, &db.GetInvitesByAppointmentIDAndClientIsParams{
		AppointmentID: input.AppointmentID,
		Column2:       clientIDs,
	})
	if err != nil {
		return err
	}
	mapClientIdInviteID := buildMapClientIdInviteID(invites)

	for _, client := range input.Clients {
		input.NotificationsService.SendGroupVisitAppointmentNotification(ctx, notifications.SendGroupVisitAppointmentNotificationInput{
			ChatID:           client.ChatID,
			StartTime:        apiCommon.FormatTimeRFC3339(input.StartTime),
			EndTime:          apiCommon.FormatTimeRFC3339(input.EndTime),
			ProfessionalName: input.ProfessionalName,
			Locale:           client.Locale,
			Description:      input.Description,
			InviteID:         mapClientIdInviteID[client.ID],
		})
	}
	return nil
}

// UpdateAppointment updates an appointment
func (s *service) UpdateAppointment(ctx context.Context, input UpdateAppointmentInput) error {
	// Get appointment to validate ownership
	appointment, err := s.repo.GetAppointmentByID(ctx, input.AppointmentID)
	if err != nil {
		return err
	}

	// Validate ownership
	if err := s.validateAppointmentOwnership(appointment, input.ProfessionalID); err != nil {
		return err
	}

	// Validate type if provided
	if input.Type != nil {
		err := s.validateTypeForAppointmentUpdate(ctx, input.AppointmentID, *input.Type)
		if err != nil {
			return err
		}
	}

	// Prepare update parameters
	var updateType sql.NullString
	var updateDescription sql.NullString

	if input.Type != nil {
		updateType = sql.NullString{
			String: *input.Type,
			Valid:  true,
		}
	}

	if input.Description != nil {
		updateDescription = sql.NullString{
			String: *input.Description,
			Valid:  true,
		}
	}

	// Update appointment type and description
	err = s.repo.UpdateAppointment(ctx, &db.UpdateAppointmentParams{
		ID:          input.AppointmentID,
		Type:        updateType.String,
		Description: updateDescription,
	})
	if err != nil {
		return err
	}

	return nil
}

// GetMissingInviteUsersForAppointment retrieves missing invite users for an appointment
func (s *service) GetMissingInviteUsersForAppointment(ctx context.Context, appointmentID uuid.UUID) ([]*db.GetMissingInviteUsersForAppointmentRow, error) {
	missingInviteUsers, err := s.repo.GetMissingInviteUsersForAppointment(ctx, appointmentID)
	if err != nil {
		return nil, err
	}
	return missingInviteUsers, nil
}

// UpdatePreviousAppointment updates a previous appointment's type and clients in a transaction
func (s *service) UpdatePreviousAppointment(ctx context.Context, input UpdatePreviousAppointmentInput) error {
	appointment, err := s.repo.GetAppointmentByID(ctx, input.AppointmentID)
	if err != nil {
		return err
	}

	if err := s.validateAppointmentOwnership(appointment, input.ProfessionalID); err != nil {
		return err
	}

	currentClientsCount, err := s.repo.GetAppintmentClientsCountByAppointmentID(ctx, input.AppointmentID)
	if err != nil {
		return err
	}

	if err := s.validateTypeForPreviousAppointmentUpdate(int(currentClientsCount), len(input.ClientsAdded), len(input.ClientsRemoved), input.Type); err != nil {
		return err
	}

	return db.WithTransaction(ctx, s.database, func(q *db.Queries) error {
		// Update appointment type
		err := q.UpdateAppointment(ctx, &db.UpdateAppointmentParams{
			ID:   input.AppointmentID,
			Type: input.Type,
		})
		if err != nil {
			return err
		}

		// Remove clients
		if len(input.ClientsRemoved) > 0 {
			err = q.DeleteClientAppointmentsByAppointmentIDAndClientIDs(ctx, &db.DeleteClientAppointmentsByAppointmentIDAndClientIDsParams{
				AppointmentID: input.AppointmentID,
				Column2:       input.ClientsRemoved,
			})
			if err != nil {
				return err
			}
		}

		// Add clients
		for _, clientID := range input.ClientsAdded {
			err = q.CreateClientAppointment(ctx, &db.CreateClientAppointmentParams{
				ClientID:      clientID,
				AppointmentID: input.AppointmentID,
			})
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// GetMissingClientsForPreviousAppointment retrieves missing clients for a previous appointment
func (s *service) GetMissingClientsForPreviousAppointment(ctx context.Context, professionalID, appointmentID uuid.UUID) ([]*db.GetMissingClientsForPreviousAppointmentRow, error) {
	missingClients, err := s.repo.GetMissingClientsForPreviousAppointment(ctx, &db.GetMissingClientsForPreviousAppointmentParams{
		AppointmentID:  appointmentID,
		ProfessionalID: professionalID,
	})
	if err != nil {
		return nil, err
	}
	return missingClients, nil
}

// GetPendingInviteUsersForAppointment retrieves pending invite users for an appointment
func (s *service) GetPendingInviteUsersForAppointment(ctx context.Context, appointmentID uuid.UUID) ([]*db.GetPendingInviteUsersForAppointmentRow, error) {
	pendingInviteUsers, err := s.repo.GetPendingInviteUsersForAppointment(ctx, appointmentID)
	if err != nil {
		return nil, err
	}
	return pendingInviteUsers, nil
}

// CreatePackage creates a package for a client
// If there's an active package, it deactivates it and creates a new one in a transaction
func (s *service) CreatePackage(ctx context.Context, input CreatePackageInput) error {
	packageID, err := s.validatePreviousPackage(ctx, input.ClientID, input.ProfessionalID)
	if err != nil {
		return err
	}

	return db.WithTransaction(ctx, s.database, func(q *db.Queries) error {
		// If there's an active package, deactivate it first
		if packageID != uuid.Nil {
			if err := q.DeactivatePackage(ctx, packageID); err != nil {
				return err
			}
		}

		// Create new package
		if err := q.CreatePackage(ctx, &db.CreatePackageParams{
			ClientID:            input.ClientID,
			ProfessionalID:      input.ProfessionalID,
			IssuedAt:            input.IssuedAt,
			ExpiresAt:           input.ExpiresAt,
			ApppointmentsNumber: int32(input.ApppointmentsNumber),
		}); err != nil {
			return err
		}

		return nil
	})
}

// GetListOfPackages retrieves a list of packages for a professional
func (s *service) GetListOfPackages(ctx context.Context, professionalID uuid.UUID) ([]*db.GetListOfPackagesByProfessionalIDRow, error) {
	packages, err := s.repo.GetListOfPackagesByProfessionalID(ctx, professionalID)
	if err != nil {
		return nil, err
	}
	return packages, nil
}

// GetPackageDetails retrieves the details of a package
func (s *service) GetPackageDetails(ctx context.Context, packageID uuid.UUID) (*GetPackageDetailsOutput, error) {
	fmt.Println("packageID", packageID)
	packageDetails, err := s.repo.GetPackageById(ctx, packageID)
	if err != nil {
		return nil, err
	}
	fmt.Println("clientID", packageDetails.ClientID)
	fmt.Println("professionalID", packageDetails.ProfessionalID)
	fmt.Println("issuedAt", packageDetails.IssuedAt)
	appointments, err := s.repo.GetAppointmentsForThePackage(ctx, &db.GetAppointmentsForThePackageParams{
		ClientID:       packageDetails.ClientID,
		ProfessionalID: packageDetails.ProfessionalID,
		StartTime:      packageDetails.IssuedAt,
	})
	if err != nil {
		return nil, err
	}

	result := &GetPackageDetailsOutput{
		ID:                  packageDetails.ID,
		ClientID:            packageDetails.ClientID,
		ProfessionalID:      packageDetails.ProfessionalID,
		IssuedAt:            packageDetails.IssuedAt,
		ExpiresAt:           packageDetails.ExpiresAt,
		ApppointmentsNumber: packageDetails.ApppointmentsNumber,
		Appointments:        appointments,
	}
	return result, nil
}

// GetClientInfoByID retrieves the client info by ID
func (s *service) GetClientInfoByID(ctx context.Context, clientID uuid.UUID) (*db.GetClientInfoByIDRow, error) {
	clientInfo, err := s.repo.GetClientInfoByID(ctx, clientID)
	if err != nil {
		return nil, err
	}
	return clientInfo, nil
}
