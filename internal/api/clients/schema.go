package api

// ClientRegisterRequest represents the request body for client registration
type ClientRegisterRequest struct {
	FirstName   string  `json:"first_name" binding:"required"`
	LastName    string  `json:"last_name" binding:"required"`
	ChatID      int64   `json:"chat_id" binding:"required"`
	PhoneNumber *string `json:"phone_number,omitempty"`
}

// ClientRegisterResponse represents the response for client registration
type ClientRegisterResponse struct {
	User User `json:"user"`
}

// User represents a user in API responses (using SQLC generated model)
type User struct {
	ID          string  `json:"id"`
	ChatID      *int64  `json:"chat_id,omitempty"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	PhoneNumber *string `json:"phone_number,omitempty"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
	Role        string  `json:"role"`
}

// GetClientAppointmentsResponse represents the response for getting client appointments
type GetClientAppointmentsResponse struct {
	Appointments []ClientAppointment `json:"appointments"`
}

// ClientAppointment represents an appointment with professional details in client context
type ClientAppointment struct {
	ID           string                         `json:"id"`
	Type         string                         `json:"type"`
	StartTime    string                         `json:"start_time"`
	EndTime      string                         `json:"end_time"`
	Status       string                         `json:"status"`
	Description  string                         `json:"description,omitempty"`
	CreatedAt    string                         `json:"created_at"`
	UpdatedAt    string                         `json:"updated_at"`
	Professional *ClientAppointmentProfessional `json:"professional,omitempty"`
}

// ClientAppointmentProfessional represents professional details in appointment context
type ClientAppointmentProfessional struct {
	ID          string  `json:"id"`
	Username    string  `json:"username"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	PhoneNumber *string `json:"phone_number,omitempty"`
	ChatID      *int64  `json:"chat_id,omitempty"`
}

// CancelClientAppointmentRequest represents the request to cancel an appointment by client
type CancelClientAppointmentRequest struct {
	CancellationReason string `json:"cancellation_reason" binding:"required"`
}

// CancelClientAppointmentResponse represents the response after cancelling an appointment by client
type CancelClientAppointmentResponse struct {
	Appointment  CancelledAppointment          `json:"appointment"`
	Client       ClientAppointmentClient       `json:"client"`
	Professional ClientAppointmentProfessional `json:"professional"`
}

// ClientAppointmentClient represents client details in appointment context
type ClientAppointmentClient struct {
	ID          string  `json:"id"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	PhoneNumber *string `json:"phone_number,omitempty"`
	ChatID      *int64  `json:"chat_id,omitempty"`
}

// CancelledAppointment represents a cancelled appointment
type CancelledAppointment struct {
	ID                 string `json:"id"`
	Type               string `json:"type"`
	StartTime          string `json:"start_time"`
	EndTime            string `json:"end_time"`
	Status             string `json:"status"`
	Description        string `json:"description,omitempty"`
	CancellationReason string `json:"cancellation_reason"`
	CancelledBy        string `json:"cancelled_by"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
}
