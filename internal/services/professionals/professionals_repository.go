package professionals

import (
	"context"

	"github.com/google/uuid"
	db "github.com/vention/booking_api/internal/repository"
)

// ProfessionalsRepository defines the database operations needed by the professionals service
type ProfessionalsRepository interface {
	GetProfessionalByUsername(ctx context.Context, username string) (*db.Professional, error)
	UpdateProfessionalChatID(ctx context.Context, arg *db.UpdateProfessionalChatIDParams) (*db.Professional, error)
	GetAppointmentByID(ctx context.Context, id uuid.UUID) (*db.Appointment, error)
	// CancelAppointmentByProfessionalWithDetails(ctx context.Context, arg *db.CancelAppointmentByProfessionalWithDetailsParams) (*db.CancelAppointmentByProfessionalWithDetailsRow, error)
	CreateUnavailableAppointment(ctx context.Context, arg *db.CreateUnavailableAppointmentParams) error
	GetAppointmentsByProfessionalWithStatusAndDate(ctx context.Context, arg *db.GetAppointmentsByProfessionalWithStatusAndDateParams) ([]*db.GetAppointmentsByProfessionalWithStatusAndDateRow, error)
	CountProfessionalAppointmentsWithStatusAndDate(ctx context.Context, arg *db.CountProfessionalAppointmentsWithStatusAndDateParams) (int64, error)
	// GetProfessionalAppointmentDates(ctx context.Context, arg *db.GetProfessionalAppointmentDatesParams) ([]time.Time, error)
	GetProfessionalTimetable(ctx context.Context, arg *db.GetProfessionalTimetableParams) ([]*db.GetProfessionalTimetableRow, error)
	// GetProfessionalClients(ctx context.Context, professionalID uuid.UUID) ([]*db.GetProfessionalClientsRow, error)
	// GetPreviousProfessionalAppointmentsByClient(ctx context.Context, arg *db.GetPreviousProfessionalAppointmentsByClientParams) ([]*db.GetPreviousProfessionalAppointmentsByClientRow, error)
	// GetPreviousAppointmentsByClientForMonth(ctx context.Context, arg *db.GetPreviousAppointmentsByClientForMonthParams) ([]*db.GetPreviousAppointmentsByClientForMonthRow, error)
	UpdateProfessionalLocale(ctx context.Context, arg *db.UpdateProfessionalLocaleParams) error
	GetAppointmentsByProfessionalByDate(ctx context.Context, arg *db.GetAppointmentsByProfessionalByDateParams) ([]*db.GetAppointmentsByProfessionalByDateRow, error)
	ConfirmAppointmentById(ctx context.Context, id uuid.UUID) error
	GetAppointmentClientsByAppointmentID(ctx context.Context, appointmentID uuid.UUID) ([]*db.GetAppointmentClientsByAppointmentIDRow, error)
	DeleteAppointmentById(ctx context.Context, id uuid.UUID) error
	CreateGroupVisitAppointment(ctx context.Context, arg *db.CreateGroupVisitAppointmentParams) (uuid.UUID, error)
	GetSubscriptionsByProfessionalID(ctx context.Context, professionalID uuid.UUID) ([]*db.GetSubscriptionsByProfessionalIDRow, error)
	GetPreviousAppointmentsByProfessionalID(ctx context.Context, arg *db.GetPreviousAppointmentsByProfessionalIDParams) ([]*db.GetPreviousAppointmentsByProfessionalIDRow, error)
	CountPreviousAppointmentsByProfessionalID(ctx context.Context, arg *db.CountPreviousAppointmentsByProfessionalIDParams) (int64, error)
	GetPreviousAppointmentsByProfessionalIDAndClientID(ctx context.Context, arg *db.GetPreviousAppointmentsByProfessionalIDAndClientIDParams) ([]*db.GetPreviousAppointmentsByProfessionalIDAndClientIDRow, error)
	CountPreviousAppointmentsByProfessionalIDAndClientID(ctx context.Context, arg *db.CountPreviousAppointmentsByProfessionalIDAndClientIDParams) (int64, error)
	GetAppointmentDetailsByProfessionalIDAndAppointmentID(ctx context.Context, arg *db.GetAppointmentDetailsByProfessionalIDAndAppointmentIDParams) (*db.GetAppointmentDetailsByProfessionalIDAndAppointmentIDRow, error)
	CreateInvites(ctx context.Context, arg *db.CreateInvitesParams) error
	GetInvitesByAppointmentIDAndClientIs(ctx context.Context, arg *db.GetInvitesByAppointmentIDAndClientIsParams) ([]*db.GetInvitesByAppointmentIDAndClientIsRow, error)
	GetAppintmentClientsCountByAppointmentID(ctx context.Context, appointmentID uuid.UUID) (int64, error)
	UpdateAppointment(ctx context.Context, arg *db.UpdateAppointmentParams) error
	CreateClientAppointment(ctx context.Context, arg *db.CreateClientAppointmentParams) error
	GetMissingInviteUsersForAppointment(ctx context.Context, appointmentID uuid.UUID) ([]*db.GetMissingInviteUsersForAppointmentRow, error)
	DeleteClientAppointmentsByAppointmentIDAndClientIDs(ctx context.Context, arg *db.DeleteClientAppointmentsByAppointmentIDAndClientIDsParams) error
	GetMissingClientsForPreviousAppointment(ctx context.Context, arg *db.GetMissingClientsForPreviousAppointmentParams) ([]*db.GetMissingClientsForPreviousAppointmentRow, error)
	GetPendingInviteUsersForAppointment(ctx context.Context, appointmentID uuid.UUID) ([]*db.GetPendingInviteUsersForAppointmentRow, error)
	GetPackageByClientIDAndProfessionalID(ctx context.Context, arg *db.GetPackageByClientIDAndProfessionalIDParams) (*db.Package, error)
	CreatePackage(ctx context.Context, arg *db.CreatePackageParams) error
	GetAppintmentNumberByClientIDAndProfessionalID(ctx context.Context, arg *db.GetAppintmentNumberByClientIDAndProfessionalIDParams) (int64, error)
	GetListOfPackagesByProfessionalID(ctx context.Context, professionalID uuid.UUID) ([]*db.GetListOfPackagesByProfessionalIDRow, error)
	DeactivatePackage(ctx context.Context, id uuid.UUID) error
	GetAppointmentsForThePackage(ctx context.Context, arg *db.GetAppointmentsForThePackageParams) ([]*db.GetAppointmentsForThePackageRow, error)
	GetPackageById(ctx context.Context, id uuid.UUID) (*db.Package, error)
	GetClientInfoByID(ctx context.Context, id uuid.UUID) (*db.GetClientInfoByIDRow, error)
}
