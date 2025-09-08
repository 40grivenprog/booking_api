package api

import (
	"context"

	"github.com/google/uuid"
	db "github.com/vention/booking_api/internal/repository"
)

type ClientsRepository interface {
	CreateClient(ctx context.Context, arg *db.CreateClientParams) (*db.Client, error)
	GetAppointmentsByClientWithStatus(ctx context.Context, arg *db.GetAppointmentsByClientWithStatusParams) ([]*db.GetAppointmentsByClientWithStatusRow, error)
	CancelAppointmentByClientWithDetails(ctx context.Context, arg *db.CancelAppointmentByClientWithDetailsParams) (*db.CancelAppointmentByClientWithDetailsRow, error)
	GetAppointmentByID(ctx context.Context, id uuid.UUID) (*db.Appointment, error)
}
