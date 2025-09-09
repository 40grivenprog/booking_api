package common

// Error types
const (
	ErrorTypeValidation   = "validation_error"
	ErrorTypeDatabase     = "database_error"
	ErrorTypeNotFound     = "not_found"
	ErrorTypeForbidden    = "forbidden"
	ErrorTypeUnauthorized = "unauthorized"
	ErrorTypeConflict     = "conflict"
	ErrorTypeInternal     = "internal_error"
	ErrorTypeBadRequest   = "bad_request"
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
	ErrorMsgInvalidEmail                     = "Invalid email format"
	ErrorMsgInvalidPhone                     = "Invalid phone number format"
	ErrorMsgFutureTimeRequired               = "Appointment time must be in the future"
	ErrorMsgAppointmentNotPending            = "Appointment is not pending"
	ErrorMsgAppointmentNotPendingOrConfirmed = "Appointment is not pending or confirmed. Please check the status of the appointment."

	// Database errors
	ErrorMsgFailedToCreateAppointment     = "Failed to create appointment"
	ErrorMsgFailedToGetAppointment        = "Failed to get appointment"
	ErrorMsgFailedToUpdateAppointment     = "Failed to update appointment"
	ErrorMsgFailedToDeleteAppointment     = "Failed to delete appointment"
	ErrorMsgFailedToCreateClient          = "Failed to create client"
	ErrorMsgFailedToGetClient             = "Failed to get client"
	ErrorMsgFailedToUpdateClient          = "Failed to update client"
	ErrorMsgFailedToCreateProfessional    = "Failed to create professional"
	ErrorMsgFailedToGetProfessional       = "Failed to get professional"
	ErrorMsgFailedToUpdateProfessional    = "Failed to update professional"
	ErrorMsgFailedToRetrieveAppointments  = "Failed to retrieve appointments"
	ErrorMsgFailedToRetrieveProfessionals = "Failed to retrieve professionals"
	ErrorMsgFailedToRetrieveClients       = "Failed to retrieve clients"

	// Not found errors
	ErrorMsgAppointmentNotFound  = "Appointment not found"
	ErrorMsgClientNotFound       = "Client not found"
	ErrorMsgProfessionalNotFound = "Professional not found"
	ErrorMsgUserNotFound         = "User not found"

	// Forbidden errors
	ErrorMsgNotAllowedToConfirmAppointment = "You are not allowed to confirm this appointment"
	ErrorMsgNotAllowedToCancelAppointment  = "You are not allowed to cancel this appointment"
	ErrorMsgNotAllowedToAccessResource     = "You are not allowed to access this resource"

	// Unauthorized errors
	ErrorMsgInvalidToken   = "Invalid or expired token"
	ErrorMsgTokenRequired  = "Authentication token required"
	ErrorMsgInvalidSession = "Invalid session"

	// Conflict errors
	ErrorMsgAppointmentConflict   = "Appointment time conflicts with existing appointment"
	ErrorMsgUsernameAlreadyExists = "Username already exists"
	ErrorMsgEmailAlreadyExists    = "Email already exists"
	ErrorMsgPhoneAlreadyExists    = "Phone number already exists"
	ErrorMsgChatIDAlreadyExists   = "Chat ID already exists"

	// Internal errors
	ErrorMsgInternalServerError = "Internal server error"
	ErrorMsgServiceUnavailable  = "Service temporarily unavailable"
	ErrorMsgTimeout             = "Request timeout"
)
