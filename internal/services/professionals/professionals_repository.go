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
	CreateGroupVisitAppointment(ctx context.Context, arg *db.CreateGroupVisitAppointmentParams) error
	GetSubscriptionsByProfessionalID(ctx context.Context, professionalID uuid.UUID) ([]*db.GetSubscriptionsByProfessionalIDRow, error)
}
