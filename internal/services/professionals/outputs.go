package professionals

import (
	"time"

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
