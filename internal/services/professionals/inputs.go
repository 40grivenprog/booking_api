package professionals

import (
	"time"

	"github.com/google/uuid"
)

// SignInInput represents the input for professional sign-in
type SignInInput struct {
	Username string
	Password string
	ChatID   int64
}

// ConfirmAppointmentInput represents the input for confirming an appointment
type ConfirmAppointmentInput struct {
	ProfessionalID uuid.UUID
	AppointmentID  uuid.UUID
}

// CancelAppointmentInput represents the input for canceling an appointment
type CancelAppointmentInput struct {
	ProfessionalID     uuid.UUID
	AppointmentID      uuid.UUID
	CancellationReason string
}

// CreateUnavailableAppointmentInput represents the input for creating unavailable appointment
type CreateUnavailableAppointmentInput struct {
	ProfessionalID uuid.UUID
	StartTime      time.Time
	EndTime        time.Time
	Description    string
}
