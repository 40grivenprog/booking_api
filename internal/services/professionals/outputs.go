package professionals

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	db "github.com/vention/booking_api/internal/repository"
)

type ConfirmAppointmentOutput struct {
	StartTime                 time.Time
	EndTime                   time.Time
	ConfirmAppointmentClients []*db.GetAppointmentClientsByAppointmentIDRow
}

type CancelAppointmentOutput struct {
	StartTime                time.Time
	EndTime                  time.Time
	CancelAppointmentClients []*db.GetAppointmentClientsByAppointmentIDRow
}

type GetAppointmentDetailsOutput struct {
	ID          uuid.UUID
	StartTime   time.Time
	EndTime     time.Time
	Description sql.NullString
	Type        string
	Clients     []*db.GetAppointmentClientsByAppointmentIDRow
}
