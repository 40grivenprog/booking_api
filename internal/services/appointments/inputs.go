package appointments

import (
	"time"

	"github.com/google/uuid"
)

// CreateAppointmentInput represents the input for creating an appointment
type CreateAppointmentInput struct {
	ClientID       uuid.UUID
	ProfessionalID uuid.UUID
	StartTime      time.Time
	EndTime        time.Time
	Description    string
}
