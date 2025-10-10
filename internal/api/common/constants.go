package common

// Error types
const (
	ErrorTypeValidation = "validation_error"
	ErrorTypeDatabase   = "database_error"
	ErrorTypeNotFound   = "not_found"
	ErrorTypeForbidden  = "forbidden"
	ErrorTypeConflict   = "conflict"
	ErrorTypeInternal   = "internal_error"
)

// Error messages
const (
	// Validation errors
	ErrorMsgInvalidRequestBody               = "Invalid request body"
	ErrorMsgInvalidAppointmentID             = "Invalid appointment_id format"
	ErrorMsgInvalidProfessionalID            = "Invalid professional_id format"
	ErrorMsgInvalidClientID                  = "Invalid client_id format"
	ErrorMsgInvalidDate                      = "Invalid date format. Use YYYY-MM-DD format (e.g., 2024-01-15)"
	ErrorMsgInvalidStatus                    = "Invalid status. Must be one of: pending, confirmed, cancelled, completed"
	ErrorMsgInvalidTime                      = "Invalid time format"
	ErrorMsgInvalidCredentials               = "Invalid username or password"
	ErrorMsgMissingRequiredField             = "Missing required field"
	ErrorMsgFutureTimeRequired               = "Appointment time must be in the future"
	ErrorMsgAppointmentNotPending            = "Appointment is not pending"
	ErrorMsgAppointmentNotPendingOrConfirmed = "Appointment is not pending or confirmed. Please check the status of the appointment."

	// Database errors
	ErrorMsgFailedToCreateAppointment     = "Failed to create appointment"
	ErrorMsgFailedToGetAppointment        = "Failed to get appointment"
	ErrorMsgFailedToUpdateAppointment     = "Failed to update appointment"
	ErrorMsgFailedToCreateClient          = "Failed to create client"
	ErrorMsgFailedToCreateProfessional    = "Failed to create professional"
	ErrorMsgFailedToUpdateProfessional    = "Failed to update professional"
	ErrorMsgFailedToRetrieveAppointments  = "Failed to retrieve appointments"
	ErrorMsgFailedToRetrieveProfessionals = "Failed to retrieve professionals"
	ErrorMsgFailedToGetTimetable          = "Failed to get professional timetable"

	// Not found errors
	ErrorMsgUserNotFound = "User not found"

	// Forbidden errors
	ErrorMsgNotAllowedToConfirmAppointment = "You are not allowed to confirm this appointment"
	ErrorMsgNotAllowedToCancelAppointment  = "You are not allowed to cancel this appointment"
	ErrorMsgNotAllowedToAccessResource     = "You are not allowed to access this resource"

	// Conflict errors
	ErrorMsgUsernameAlreadyExists = "Username already exists"

	// Internal errors
	ErrorMsgInternalServerError = "Internal server error"
)

// User roles
const (
	RoleProfessional = "professional"
	RoleClient       = "client"
	RoleAdmin        = "admin"
)

// User types (same as roles but used in different contexts)
const (
	UserTypeProfessional = "professional"
	UserTypeClient       = "client"
)

// Cancellation sources
const (
	CancelledByProfessional = "professional"
	CancelledByClient       = "client"
)

// Working hours configuration
const (
	WorkingHoursStart = 5  // 5:00 AM
	WorkingHoursEnd   = 23 // 11:00 PM (exclusive, so last slot is 22:00-23:00)
)

// Time slot configuration
const (
	SlotDurationMinutes = 60 // 1 hour slots
	SlotsPerDay         = WorkingHoursEnd - WorkingHoursStart
)

// Appointment type strings (for responses)
const (
	AppointmentTypeBooking     = "appointment"
	AppointmentTypeUnavailable = "unavailable"
)
