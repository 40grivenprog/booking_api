package api

import (
	"context"
	db "github.com/vention/booking_api/internal/repository"
)

type AppointmentsRepository interface {
	CreateAppointmentWithDetails(ctx context.Context, arg *db.CreateAppointmentWithDetailsParams) (*db.CreateAppointmentWithDetailsRow, error)
}
