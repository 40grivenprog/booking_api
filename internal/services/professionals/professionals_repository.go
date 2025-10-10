package professionals

import (
	"context"
	"time"

	"github.com/google/uuid"
	db "github.com/vention/booking_api/internal/repository"
)

// ProfessionalsRepository defines the database operations needed by the professionals service
type ProfessionalsRepository interface {
	GetProfessionals(ctx context.Context) ([]*db.Professional, error)
	GetProfessionalByUsername(ctx context.Context, username string) (*db.Professional, error)
	UpdateProfessionalChatID(ctx context.Context, arg *db.UpdateProfessionalChatIDParams) (*db.Professional, error)
	GetAppointmentByID(ctx context.Context, id uuid.UUID) (*db.Appointment, error)
	ConfirmAppointmentWithDetails(ctx context.Context, arg *db.ConfirmAppointmentWithDetailsParams) (*db.ConfirmAppointmentWithDetailsRow, error)
	CancelAppointmentByProfessionalWithDetails(ctx context.Context, arg *db.CancelAppointmentByProfessionalWithDetailsParams) (*db.CancelAppointmentByProfessionalWithDetailsRow, error)
	CreateUnavailableAppointment(ctx context.Context, arg *db.CreateUnavailableAppointmentParams) (*db.Appointment, error)
	GetAppointmentsByProfessionalWithStatusAndDate(ctx context.Context, arg *db.GetAppointmentsByProfessionalWithStatusAndDateParams) ([]*db.GetAppointmentsByProfessionalWithStatusAndDateRow, error)
	GetProfessionalAppointmentDates(ctx context.Context, arg *db.GetProfessionalAppointmentDatesParams) ([]time.Time, error)
	GetAppointmentsByProfessionalAndDateWithClient(ctx context.Context, arg *db.GetAppointmentsByProfessionalAndDateWithClientParams) ([]*db.GetAppointmentsByProfessionalAndDateWithClientRow, error)
	GetProfessionalTimetable(ctx context.Context, arg *db.GetProfessionalTimetableParams) ([]*db.GetProfessionalTimetableRow, error)
}
