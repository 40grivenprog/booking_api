package clients

import (
	"context"

	"github.com/google/uuid"
	db "github.com/vention/booking_api/internal/repository"
)

// ClientsRepository defines the database operations needed by the clients service
type ClientsRepository interface {
	CreatePersonalAppointment(ctx context.Context, arg *db.CreatePersonalAppointmentParams) (*db.CreatePersonalAppointmentRow, error)
	CreateClientAppointment(ctx context.Context, arg *db.CreateClientAppointmentParams) error
	GetAppointmentWithDetails(ctx context.Context, id uuid.UUID) (*db.GetAppointmentWithDetailsRow, error)
	CheckClientAppointmentConflict(ctx context.Context, arg *db.CheckClientAppointmentConflictParams) (bool, error)
	CreateClient(ctx context.Context, arg *db.CreateClientParams) (*db.Client, error)
	GetAppointmentsByClientWithStatus(ctx context.Context, arg *db.GetAppointmentsByClientWithStatusParams) ([]*db.GetAppointmentsByClientWithStatusRow, error)
	CountClientAppointmentsWithStatus(ctx context.Context, arg *db.CountClientAppointmentsWithStatusParams) (int64, error)
	UpdateClientLocale(ctx context.Context, arg *db.UpdateClientLocaleParams) error
	GetProfessionalInfoForNotificationByAppointmentID(ctx context.Context, id uuid.UUID) (*db.GetProfessionalInfoForNotificationByAppointmentIDRow, error)
	CreateSubscription(ctx context.Context, arg *db.CreateSubscriptionParams) error
	DeleteSubscription(ctx context.Context, arg *db.DeleteSubscriptionParams) error
	GetSubscriptionsByClientID(ctx context.Context, arg *db.GetSubscriptionsByClientIDParams) ([]*db.GetSubscriptionsByClientIDRow, error)
	CountSubscriptionsByClientID(ctx context.Context, clientID uuid.UUID) (int64, error)
	GetProfessionals(ctx context.Context, arg *db.GetProfessionalsParams) ([]*db.GetProfessionalsRow, error)
	CountProfessionals(ctx context.Context, clientID uuid.UUID) (int64, error)
	DeleteAppointmentById(ctx context.Context, id uuid.UUID) error
	GetAppointmentInfoByAppointmentID(ctx context.Context, id uuid.UUID) (*db.GetAppointmentInfoByAppointmentIDRow, error)
	GetClientInvites(ctx context.Context, clientID uuid.UUID) ([]*db.GetClientInvitesRow, error)
	GetInviteByID(ctx context.Context, arg *db.GetInviteByIDParams) (*db.GetInviteByIDRow, error)
	DeleteInviteByID(ctx context.Context, arg *db.DeleteInviteByIDParams) error
	GetInfoForAcceptInviteNotification(ctx context.Context, appointmentID uuid.UUID) (*db.GetInfoForAcceptInviteNotificationRow, error)
	GetPreviousAppointmentsByClientID(ctx context.Context, arg *db.GetPreviousAppointmentsByClientIDParams) ([]*db.GetPreviousAppointmentsByClientIDRow, error)
	CountPreviousAppointmentsByClientID(ctx context.Context, clientID uuid.UUID) (int64, error)
}
