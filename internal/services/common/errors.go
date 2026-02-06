package common

import "errors"

// Domain-level errors that are independent of HTTP
var (
	// Time validation errors
	ErrInvalidTimeRange = errors.New("invalid time range")
	ErrPastTime         = errors.New("time must be in the future")
	ErrPastAppointment  = errors.New("cannot modify appointments from the past")

	// Authentication errors
	ErrInvalidCredentials = errors.New("invalid credentials")

	// Authorization errors
	ErrForbidden = errors.New("access forbidden")

	// Appointment validation errors
	ErrAppointmentNotPending            = errors.New("appointment is not pending")
	ErrAppointmentNotPendingOrConfirmed = errors.New("appointment is not pending or confirmed")
	ErrAppointmentTimeConflict          = errors.New("you already have an appointment for this time. please cancel it first or wait for approval")
	ErrAppointmentFullyBooked           = errors.New("sorry. this appointment is fully booked")

	ErrProfessionalNotFound = errors.New("professional with such username not found")
)
