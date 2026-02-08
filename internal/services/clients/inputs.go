package clients

import (
	"time"

	"github.com/google/uuid"
	"github.com/vention/booking_api/internal/services/notifications"
)

// RegisterClientInput represents the input for registering a client
type RegisterClientInput struct {
	FirstName string
	LastName  string
	ChatID    int64
	Locale    string
}

// CancelAppointmentInput represents the input for canceling an appointment
type CancelAppointmentInput struct {
	ClientID           uuid.UUID
	AppointmentID      uuid.UUID
	CancellationReason string
}

// CreateAppointmentInput represents the input for creating an appointment
type CreateAppointmentInput struct {
	ClientID       uuid.UUID
	ProfessionalID uuid.UUID
	StartTime      time.Time
	EndTime        time.Time
	Description    string
}

// UpdateLocaleInput represents the input for changing the locale
type UpdateLocaleInput struct {
	ClientID uuid.UUID
	Locale   string
}

// SubscribeToProfessionalInput represents the input for subscribing to a professional
type SubscribeToProfessionalInput struct {
	ClientID       uuid.UUID
	ProfessionalID uuid.UUID
}

// UnsubscribeFromProfessionalInput represents the input for unsubscribing from a professional
type UnsubscribeFromProfessionalInput struct {
	ClientID       uuid.UUID
	ProfessionalID uuid.UUID
}

// AcceptClientInviteInput represents the input for accepting a client invite
type AcceptClientInviteInput struct {
	ClientID            uuid.UUID
	InviteID            uuid.UUID
	AppointmentID       uuid.UUID
	ClientName          string
	NotificationService notifications.Service
	Type                string
}
