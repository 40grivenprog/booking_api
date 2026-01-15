package clients

import (
	"context"

	"github.com/google/uuid"
	db "github.com/vention/booking_api/internal/repository"
)

// ClientsRepository defines the database operations needed by the clients service
type ClientsRepository interface {
	CreateAppointmentWithDetails(ctx context.Context, arg *db.CreateAppointmentWithDetailsParams) (*db.CreateAppointmentWithDetailsRow, error)
	CheckClientAppointmentConflict(ctx context.Context, arg *db.CheckClientAppointmentConflictParams) (bool, error)
	CreateClient(ctx context.Context, arg *db.CreateClientParams) (*db.Client, error)
	GetAppointmentsByClientWithStatus(ctx context.Context, arg *db.GetAppointmentsByClientWithStatusParams) ([]*db.GetAppointmentsByClientWithStatusRow, error)
	CountClientAppointmentsWithStatus(ctx context.Context, arg *db.CountClientAppointmentsWithStatusParams) (int64, error)
	GetAppointmentByID(ctx context.Context, id uuid.UUID) (*db.Appointment, error)
	CancelAppointmentByClientWithDetails(ctx context.Context, arg *db.CancelAppointmentByClientWithDetailsParams) (*db.CancelAppointmentByClientWithDetailsRow, error)
}
