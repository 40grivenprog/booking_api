package api

import (
	common "github.com/vention/booking_api/internal/api/common"
)

// ProfessionalSignInRequest represents the request body for professional sign in
type ProfessionalSignInRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	ChatID   int64  `json:"chat_id" binding:"required"`
	Locale   string `json:"locale" binding:"required"`
}

// ProfessionalSignInResponse represents the response for professional sign in
type ProfessionalSignInResponse struct {
	User User `json:"user"`
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
	Token       string  `json:"token"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// GetProfessionalAppointmentsResponse represents the response for getting professional appointments
type GetProfessionalAppointmentsResponse struct {
	Appointments []GetProfessionalAppointmentsResponseItem `json:"appointments"`
	Pagination   common.PaginationResponse                 `json:"pagination"`
}

// GetProfessionalAppointmentsResponseItem represents an appointment with client details in professional context
type GetProfessionalAppointmentsResponseItem struct {
	ID        string   `json:"id"`
	StartTime string   `json:"start_time"`
	EndTime   string   `json:"end_time"`
	Type      string   `json:"type"`
	Clients   []string `json:"clients"`
}

// CancelAppointmentRequest represents the request to cancel an appointment
type CancelAppointmentRequest struct {
	CancellationReason string `json:"cancellation_reason" binding:"required"`
}

// CreateUnavailableAppointmentRequest represents the request to create an unavailable appointment
type CreateUnavailableAppointmentRequest struct {
	Description string `json:"description,omitempty" binding:"required"`
	StartAt     string `json:"start_at" binding:"required"`
	EndAt       string `json:"end_at" binding:"required"`
}

// GetProfessionalAvailabilityResponse represents the response for professional availability
type GetProfessionalAvailabilityResponse struct {
	Date  string     `json:"date"`
	Slots []TimeSlot `json:"slots"`
}

// TimeSlot represents a one-hour time slot
type TimeSlot struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Available bool   `json:"available"`
}

// GetProfessionalAppointmentDatesResponse represents the response for getting appointment dates
type GetProfessionalAppointmentDatesResponse struct {
	Month string   `json:"month"`
	Dates []string `json:"dates"`
}

// TimetableAppointment represents an appointment in the timetable
type TimetableAppointment struct {
	ID          string   `json:"id"`
	StartTime   string   `json:"start_time"`
	EndTime     string   `json:"end_time"`
	Description string   `json:"description"`
	Type        string   `json:"type"`
	Clients     []string `json:"clients"`
}

// GetProfessionalTimetableResponse represents the response for getting professional timetable
type GetProfessionalTimetableResponse struct {
	Date         string                 `json:"date"`
	Appointments []TimetableAppointment `json:"appointments"`
}

// Client represents a client in API responses
type ProfessionalClient struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// GetProfessionalClientsResponse represents the response for getting professional clients
type GetProfessionalClientsResponse struct {
	Clients []ProfessionalClient `json:"clients"`
}

// PreviousAppointment represents a previous appointment
type PreviousAppointment struct {
	ID          string `json:"id"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	Description string `json:"description"`
}

// GetPreviousAppointmentsByClientResponse represents the response for getting previous appointments by client
type GetPreviousAppointmentsByClientResponse struct {
	Appointments []PreviousAppointment `json:"appointments"`
}

// UpdateLocaleRequest represents the request to change the locale
type UpdateLocaleRequest struct {
	Locale string `json:"locale" binding:"required"`
}

// CreateGroupVisitAppointmentRequest represents the request to create a group visit appointment
type CreateGroupVisitAppointmentRequest struct {
	Description     string                              `json:"description,omitempty" binding:"required"`
	StartAt         string                              `json:"start_at" binding:"required"`
	EndAt           string                              `json:"end_at" binding:"required"`
	Clients         []CreateGroupVisitAppointmentClient `json:"clients"`
	Type            string                              `json:"type" binding:"required"`
	ClientsSelected string                              `json:"clients_selected" binding:"required,oneof=all partially_selected"`
}

type CreateGroupVisitAppointmentClient struct {
	ID     string `json:"id" binding:"required"`
	ChatID int64  `json:"chat_id" binding:"required"`
	Locale string `json:"locale" binding:"required"`
}

// GetProfessionalSubscriptionsResponse represents the response for getting professional subscriptions
type GetProfessionalSubscriptionsResponse struct {
	Subscriptions []ProfessionalSubscription `json:"subscriptions"`
}

// ProfessionalSubscription represents a subscription in the response
type ProfessionalSubscription struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	ChatID    *int64 `json:"chat_id,omitempty"`
	Locale    string `json:"locale"`
}

// GetPreviousAppointmentsResponse represents the response for getting previous appointments
type GetPreviousAppointmentsResponse struct {
	Appointments []GetPreviousAppointmentsResponseItem `json:"appointments"`
	Pagination   common.PaginationResponse             `json:"pagination"`
}

// GetPreviousAppointmentsResponseItem represents a previous appointment
type GetPreviousAppointmentsResponseItem struct {
	ID          string `json:"id"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	Description string `json:"description"`
	Type        string `json:"type"`
	UsersCount  int64  `json:"users_count"`
}

// GetAppointmentDetailsResponse represents the response for getting appointment details
type GetAppointmentDetailsResponse struct {
	ID          string                                `json:"id"`
	StartTime   string                                `json:"start_time"`
	EndTime     string                                `json:"end_time"`
	Description string                                `json:"description"`
	Type        string                                `json:"type"`
	Clients     []GetAppointmentDetailsResponseClient `json:"clients"`
}

// GetAppointmentDetailsResponseClient represents a client in the response
type GetAppointmentDetailsResponseClient struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
