package professionals

import (
	"time"

	"github.com/google/uuid"
	"github.com/vention/booking_api/internal/services/notifications"
)

// SignInInput represents the input for professional sign-in
type SignInInput struct {
	Username string
	Password string
	ChatID   int64
	Locale   string
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

// UpdateLocaleInput represents the input for changing the locale
type UpdateLocaleInput struct {
	ProfessionalID uuid.UUID
	Locale         string
}

// CreateGroupVisitAppointmentInput represents the input for creating a group visit appointment
type CreateGroupVisitAppointmentInput struct {
	ProfessionalID uuid.UUID
	StartTime      time.Time
	EndTime        time.Time
	Description    string
	Type           string
}

// InviteAllClientsInput represents the input for inviting all clients
type InviteAllClientsInput struct {
	ProfessionalID       uuid.UUID
	AppointmentID        uuid.UUID
	StartTime            time.Time
	EndTime              time.Time
	Description          string
	Type                 string
	ProfessionalName     string
	ClientIDs            []uuid.UUID
	NotificationsService notifications.Service
}

// InvitePartiallySelectedClientsInput represents the input for inviting partially selected clients
type InvitePartiallySelectedClientsInput struct {
	AppointmentID        uuid.UUID
	StartTime            time.Time
	EndTime              time.Time
	Description          string
	Type                 string
	ProfessionalName     string
	Clients              []InvitePartiallySelectedClientsClientInput
	NotificationsService notifications.Service
}

// InvitePartiallySelectedClientsClientInput represents the input for inviting a partially selected client
type InvitePartiallySelectedClientsClientInput struct {
	ID     uuid.UUID
	ChatID int64
	Locale string
}

// UpdateAppointmentInput represents the input for updating an appointment
type UpdateAppointmentInput struct {
	ProfessionalID       uuid.UUID
	ProfessionalName     string
	AppointmentID        uuid.UUID
	Description          *string
	Type                 *string
	NotificationsService notifications.Service
}

// UpdatePreviousAppointmentInput represents the input for updating a previous appointment
type UpdatePreviousAppointmentInput struct {
	ProfessionalID uuid.UUID
	AppointmentID  uuid.UUID
	Type           string
	ClientsAdded   []uuid.UUID
	ClientsRemoved []uuid.UUID
}

// CreatePackageInput represents the input for creating a package
type CreatePackageInput struct {
	ProfessionalID      uuid.UUID
	ClientID            uuid.UUID
	IssuedAt            time.Time
	ExpiresAt           time.Time
	ApppointmentsNumber int64
}

// CreateInviteForAppointmentInput represents the input for creating an invite for an appointment
type CreateInviteForAppointmentInput struct {
	ProfessionalID uuid.UUID
	AppointmentID  uuid.UUID
	ClientID       uuid.UUID
}
