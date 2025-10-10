package clients

import "github.com/google/uuid"

// RegisterClientInput represents the input for registering a client
type RegisterClientInput struct {
	FirstName   string
	LastName    string
	PhoneNumber string
	ChatID      int64
}

// CancelAppointmentInput represents the input for canceling an appointment
type CancelAppointmentInput struct {
	ClientID           uuid.UUID
	AppointmentID      uuid.UUID
	CancellationReason string
}
