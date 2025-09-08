package api

import (
	"context"

	"github.com/google/uuid"
	db "github.com/vention/booking_api/internal/repository"
)

type ProfessionalsRepository interface {
	GetProfessionals(ctx context.Context) ([]*db.Professional, error)
	GetProfessionalByUsername(ctx context.Context, username string) (*db.Professional, error)
	UpdateProfessionalChatID(ctx context.Context, arg *db.UpdateProfessionalChatIDParams) (*db.Professional, error)
	ConfirmAppointmentWithDetails(ctx context.Context, arg *db.ConfirmAppointmentWithDetailsParams) (*db.ConfirmAppointmentWithDetailsRow, error)
	GetAppointmentsByProfessionalWithStatus(ctx context.Context, arg *db.GetAppointmentsByProfessionalWithStatusParams) ([]*db.GetAppointmentsByProfessionalWithStatusRow, error)
	CancelAppointmentByProfessionalWithDetails(ctx context.Context, arg *db.CancelAppointmentByProfessionalWithDetailsParams) (*db.CancelAppointmentByProfessionalWithDetailsRow, error)
	CreateUnavailableAppointment(ctx context.Context, arg *db.CreateUnavailableAppointmentParams) (*db.Appointment, error)
	GetAppointmentsByProfessionalAndDate(ctx context.Context, arg *db.GetAppointmentsByProfessionalAndDateParams) ([]*db.Appointment, error)
	GetAppointmentByID(ctx context.Context, id uuid.UUID) (*db.Appointment, error)
}
