package api

// CreateAppointmentRequest represents the request to create an appointment
type CreateAppointmentRequest struct {
	ClientID       string `json:"client_id" binding:"required"`
	ProfessionalID string `json:"professional_id" binding:"required"`
	StartTime      string `json:"start_time" binding:"required"`
	EndTime        string `json:"end_time" binding:"required"`
}

// CreateAppointmentResponse represents the response after creating an appointment
type CreateAppointmentResponse struct {
	Appointment  Appointment  `json:"appointment"`
	Client       Client       `json:"client"`
	Professional Professional `json:"professional"`
}

// Appointment represents an appointment in the response
type Appointment struct {
	ID                 string  `json:"id"`
	StartTime          string  `json:"start_time"`
	EndTime            string  `json:"end_time"`
	Status             string  `json:"status"`
	CreatedAt          string  `json:"created_at"`
	UpdatedAt          string  `json:"updated_at"`
}

// Client represents a client in the response
type Client struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	ChatID      int64  `json:"chat_id,omitempty"`
}

// Professional represents a professional in the response
type Professional struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	ChatID    int64  `json:"chat_id,omitempty"`
}
