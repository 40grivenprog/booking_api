package common

import "errors"

// Domain-level errors that are independent of HTTP
var (
	// Time validation errors
	ErrInvalidTimeRange = errors.New("invalid time range")
	ErrPastTime         = errors.New("time must be in the future")

	// Authentication errors
	ErrInvalidCredentials = errors.New("invalid credentials")

	// Authorization errors
	ErrForbidden = errors.New("access forbidden")

	// Appointment validation errors
	ErrAppointmentNotPending            = errors.New("appointment is not pending")
	ErrAppointmentNotPendingOrConfirmed = errors.New("appointment is not pending or confirmed")
)
