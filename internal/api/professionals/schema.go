package api

// ProfessionalSignInRequest represents the request body for professional sign in
type ProfessionalSignInRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	ChatID   int64  `json:"chat_id" binding:"required"`
}

// ProfessionalSignInResponse represents the response for professional sign in
type ProfessionalSignInResponse struct {
	User User `json:"user"`
}

// GetProfessionalsResponse represents the response for getting all professionals
type GetProfessionalsResponse struct {
	Professionals []User `json:"professionals"`
}

// User represents a user in API responses (using SQLC generated model)
type User struct {
	ID          string  `json:"id"`
	ChatID      *int64  `json:"chat_id,omitempty"`
	Username    string  `json:"username"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	UserType    string  `json:"user_type"`
	PhoneNumber *string `json:"phone_number,omitempty"`
	Role        string  `json:"role"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// ConfirmAppointmentResponse represents the response for confirming an appointment
type ConfirmAppointmentResponse struct {
	Appointment  AppointmentConfirm  `json:"appointment"`
	Client       ClientConfirm       `json:"client"`
	Professional ProfessionalConfirm `json:"professional"`
}

// AppointmentConfirm represents an appointment in the confirm response
type AppointmentConfirm struct {
	ID          string `json:"id"`
	Status      string `json:"status"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	Description string `json:"description,omitempty"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// ClientConfirm represents a client in the confirm response
type ClientConfirm struct {
	ID        string `json:"id"`
	ChatID    int64  `json:"chat_id,omitempty"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// ProfessionalConfirm represents a professional in the confirm response
type ProfessionalConfirm struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// GetProfessionalAppointmentsResponse represents the response for getting professional appointments
type GetProfessionalAppointmentsResponse struct {
	Appointments []ProfessionalAppointment `json:"appointments"`
}

// ProfessionalAppointment represents an appointment with client details in professional context
type ProfessionalAppointment struct {
	ID          string                         `json:"id"`
	Type        string                         `json:"type"`
	StartTime   string                         `json:"start_time"`
	EndTime     string                         `json:"end_time"`
	Status      string                         `json:"status"`
	Description string                         `json:"description,omitempty"`
	CreatedAt   string                         `json:"created_at"`
	UpdatedAt   string                         `json:"updated_at"`
	Client      *ProfessionalAppointmentClient `json:"client,omitempty"`
}

// ProfessionalAppointmentClient represents client details in appointment context
type ProfessionalAppointmentClient struct {
	ID          string  `json:"id"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	PhoneNumber *string `json:"phone_number,omitempty"`
	ChatID      *int64  `json:"chat_id,omitempty"`
}

// CancelAppointmentRequest represents the request to cancel an appointment
type CancelAppointmentRequest struct {
	CancellationReason string `json:"cancellation_reason" binding:"required"`
}

// CancelAppointmentResponse represents the response after cancelling an appointment
type CancelAppointmentResponse struct {
	Appointment  CancelledAppointment          `json:"appointment"`
	Client       ProfessionalAppointmentClient `json:"client"`
	Professional ProfessionalInfo              `json:"professional"`
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

// ProfessionalInfo represents professional details in appointment context
type ProfessionalInfo struct {
	ID          string  `json:"id"`
	Username    string  `json:"username"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	PhoneNumber *string `json:"phone_number,omitempty"`
	ChatID      *int64  `json:"chat_id,omitempty"`
}

// CreateUnavailableAppointmentRequest represents the request to create an unavailable appointment
type CreateUnavailableAppointmentRequest struct {
	Description string `json:"description,omitempty" binding:"required"`
	StartAt     string `json:"start_at" binding:"required"`
	EndAt       string `json:"end_at" binding:"required"`
}

// CreateUnavailableAppointmentResponse represents the response after creating an unavailable appointment
type CreateUnavailableAppointmentResponse struct {
	Appointment UnavailableAppointment `json:"appointment"`
}

// UnavailableAppointment represents an unavailable appointment
type UnavailableAppointment struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	Status      string `json:"status"`
	Description string `json:"description,omitempty"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// GetProfessionalAvailabilityResponse represents the response for professional availability
type GetProfessionalAvailabilityResponse struct {
	Date  string     `json:"date"`
	Slots []TimeSlot `json:"slots"`
}

// TimeSlot represents a one-hour time slot
type TimeSlot struct {
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	Available   bool   `json:"available"`
	Type        string `json:"type,omitempty"`        // "appointment", "unavailable", or empty if available
	Description string `json:"description,omitempty"` // Description with client info if available
}
