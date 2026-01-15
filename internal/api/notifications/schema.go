package api

// SendAppointmentRequestRequest represents the request to send appointment notification
type SendAppointmentRequestRequest struct {
	ChatID      int64  `json:"chat_id" binding:"required"`
	ClientName  string `json:"client_name" binding:"required"` // "FirstName LastName"
	StartTime   string `json:"start_time" binding:"required"`  // RFC3339 format
	EndTime     string `json:"end_time" binding:"required"`    // RFC3339 format
	Description string `json:"description,omitempty"`
}

// SendAppointmentRequestResponse represents the response after sending notification
type SendAppointmentRequestResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}
