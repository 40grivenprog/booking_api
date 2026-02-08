package clients

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	db "github.com/vention/booking_api/internal/repository"
)

type CancelAppointmentOutput struct {
	StartTime          time.Time
	EndTime            time.Time
	ProfessionalChatID sql.NullInt64
	ProfessionalLocale string
}

type GetClientPackageDetailsOutput struct {
	ID                  uuid.UUID
	IssuedAt            time.Time
	ExpiresAt           time.Time
	ApppointmentsNumber int32
	FirstName           string
	LastName            string
	Appointments        []*db.GetAppointmentsForThePackageRow
}
