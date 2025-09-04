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
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// ConfirmAppointmentResponse represents the response for confirming an appointment
type ConfirmAppointmentResponse struct {
	Appointment AppointmentConfirm `json:"appointment"`
	Client      ClientConfirm      `json:"client"`
}

// AppointmentConfirm represents an appointment in the confirm response
type AppointmentConfirm struct {
	ID        string `json:"id"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// ClientConfirm represents a client in the confirm response
type ClientConfirm struct {
	ID        string `json:"id"`
	ChatID    int64  `json:"chat_id,omitempty"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// GetProfessionalAppointmentsResponse represents the response for getting professional appointments
type GetProfessionalAppointmentsResponse struct {
	Appointments []ProfessionalAppointment `json:"appointments"`
}

// ProfessionalAppointment represents an appointment with client details in professional context
type ProfessionalAppointment struct {
	ID        string                         `json:"id"`
	Type      string                         `json:"type"`
	StartTime string                         `json:"start_time"`
	EndTime   string                         `json:"end_time"`
	Status    string                         `json:"status"`
	CreatedAt string                         `json:"created_at"`
	UpdatedAt string                         `json:"updated_at"`
	Client    *ProfessionalAppointmentClient `json:"client,omitempty"`
}

// ProfessionalAppointmentClient represents client details in appointment context
type ProfessionalAppointmentClient struct {
	ID          string  `json:"id"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	PhoneNumber *string `json:"phone_number,omitempty"`
}
