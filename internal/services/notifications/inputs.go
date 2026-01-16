package notifications

type SendAppointmentRequestNotificationInput struct {
	ChatID      int64  `json:"chat_id" binding:"required"`
	ClientName  string `json:"client_name" binding:"required"` // "FirstName LastName"
	StartTime   string `json:"start_time" binding:"required"`  // RFC3339 format
	EndTime     string `json:"end_time" binding:"required"`    // RFC3339 format
	Description string `json:"description,omitempty"`
}

type SendAppointmentCancellationNotificationInput struct {
	ChatID             int64  `json:"chat_id" binding:"required"`
	StartTime          string `json:"start_time" binding:"required"`
	EndTime            string `json:"end_time" binding:"required"`
	RespondentName     string `json:"respondent_name" binding:"required"` // "FirstName LastName" - name of person who cancelled
	CancellationReason string `json:"cancellation_reason"`
	Type               string `json:"type" binding:"required,oneof=professional client"` // "professional" or "client"
}

type SendAppointmentConfirmationNotificationInput struct {
	ChatID           int64  `json:"chat_id" binding:"required"`
	StartTime        string `json:"start_time" binding:"required"`
	EndTime          string `json:"end_time" binding:"required"`
	ProfessionalName string `json:"professional_name" binding:"required"`
}
