package clients

import (
	"database/sql"
	"time"
)

type CancelAppointmentOutput struct {
	StartTime          time.Time
	EndTime            time.Time
	ProfessionalChatID sql.NullInt64
	ProfessionalLocale string
}
