package api

import (
	common "github.com/vention/booking_api/internal/api/common"
)

// UpdateLocaleRequest represents the request to change the locale
type UpdateLocaleRequest struct {
	Locale string `json:"locale" binding:"required"`
}

// CreateAppointmentRequest represents the request to create an appointment
type CreateAppointmentRequest struct {
	ProfessionalID     string `json:"professional_id" binding:"required"`
	ProfessionalChatID int64  `json:"professional_chat_id" binding:"required"`
	ProfessionalLocale string `json:"professional_locale" binding:"required"`
	StartTime          string `json:"start_time" binding:"required"`
	EndTime            string `json:"end_time" binding:"required"`
}

// CreateAppointmentResponse represents the response after creating an appointment
type CreateAppointmentResponse struct {
	Appointment  Appointment  `json:"appointment"`
	Professional Professional `json:"professional"`
	Client       Client       `json:"client"`
}

// Appointment represents an appointment in the response
type Appointment struct {
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	Description string `json:"description,omitempty"`
}

// Client represents a client in the response
type Client struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// Professional represents a professional in the response
type Professional struct {
	ChatID int64 `json:"chat_id,omitempty"`
}

// ClientRegisterRequest represents the request body for client registration
type ClientRegisterRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	ChatID    int64  `json:"chat_id" binding:"required"`
	Locale    string `json:"locale" binding:"required"`
}

// ClientRegisterResponse represents the response for client registration
type ClientRegisterResponse struct {
	ID        string `json:"id"`
	ChatID    *int64 `json:"chat_id,omitempty"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Role      string `json:"role"`
	Token     string `json:"token"`
}

// GetClientAppointmentsResponse represents the response for getting client appointments
type GetClientAppointmentsResponse struct {
	Appointments []ClientAppointment       `json:"appointments"`
	Pagination   common.PaginationResponse `json:"pagination"`
}

// ClientAppointment represents an appointment with professional details in client context
type ClientAppointment struct {
	ID           string                         `json:"id"`
	Type         string                         `json:"type"`
	StartTime    string                         `json:"start_time"`
	EndTime      string                         `json:"end_time"`
	Professional *ClientAppointmentProfessional `json:"professional,omitempty"`
}

// ClientAppointmentProfessional represents professional details in appointment context
type ClientAppointmentProfessional struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	ChatID    *int64 `json:"chat_id,omitempty"`
}

// CancelClientAppointmentRequest represents the request to cancel an appointment by client
type CancelClientAppointmentRequest struct {
	CancellationReason string `json:"cancellation_reason" binding:"required"`
}

// GetSubscribedProfessionalsResponse represents the response for getting subscribed professionals
type GetSubscribedProfessionalsResponse struct {
	Professionals []GetSubscribedProfessionalsResponseItem `json:"professionals"`
	Pagination    common.PaginationResponse                `json:"pagination"`
}

// SubscribedProfessional represents a subscribed professional in the response
type GetSubscribedProfessionalsResponseItem struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	ChatID    *int64 `json:"chat_id,omitempty"`
	Locale    string `json:"locale"`
}

// GetProfessionalsResponse represents the response for getting all professionals
type GetProfessionalsResponse struct {
	Professionals []GetProfessionalsResponseItem `json:"professionals"`
	Pagination    common.PaginationResponse      `json:"pagination"`
}

// GetProfessionalsResponseItem represents a professional in the response
type GetProfessionalsResponseItem struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	ChatID    *int64 `json:"chat_id,omitempty"`
	Locale    string `json:"locale"`
}

// SubscribedProfessionalRequest represents the request to subscribe to a professional
type SubscribedProfessionalRequest struct {
	ProfessionalID string `json:"professional_id" binding:"required"`
	ChatID         int64  `json:"chat_id" binding:"required"`
	Locale         string `json:"locale" binding:"required"`
}

// GetClientInvitesResponse represents the response for getting client invites
type GetClientInvitesResponse struct {
	Invites []GetClientInvitesResponseItem `json:"invites"`
}

// GetClientInvitesResponseItem represents an invite in the response
type GetClientInvitesResponseItem struct {
	ID               string `json:"id"`
	AppointmentID    string `json:"appointment_id"`
	StartTime        string `json:"start_time"`
	EndTime          string `json:"end_time"`
	Description      string `json:"description"`
	Type             string `json:"type"`
	ProfessionalName string `json:"professional_name"`
	ClientID         string `json:"client_id"`
}

// GetClientInviteResponse represents the response for getting a client invite
type GetClientInviteResponse struct {
	ID               string `json:"id"`
	AppointmentID    string `json:"appointment_id"`
	StartTime        string `json:"start_time"`
	EndTime          string `json:"end_time"`
	Description      string `json:"description"`
	Type             string `json:"type"`
	ProfessionalName string `json:"professional_name"`
}

// AcceptClientInviteRequest represents the request to accept a client invite
type AcceptClientInviteRequest struct {
	AppointmentID string `json:"appointment_id" binding:"required"`
	Type          string `json:"type" binding:"required"`
}

// GetClientPreviousAppointmentsResponse represents the response for getting client previous appointments
type GetClientPreviousAppointmentsResponse struct {
	Appointments []GetClientPreviousAppointmentsResponseItem `json:"appointments"`
	Pagination   common.PaginationResponse                   `json:"pagination"`
}

// GetClientPreviousAppointmentsResponseItem represents a previous appointment
type GetClientPreviousAppointmentsResponseItem struct {
	ID        string `json:"id"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Type      string `json:"type"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// GetProfessionalsTimetableResponse represents the response for getting professionals timetable
type GetProfessionalsTimetableResponse struct {
	Date         string                                  `json:"date"`
	Appointments []GetProfessionalsTimetableResponseItem `json:"appointments"`
}

// GetProfessionalsTimetableResponseItem represents an appointment in the response
type GetProfessionalsTimetableResponseItem struct {
	ID        string `json:"id"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Type      string `json:"type"`
}

// GetClientPackagesResponse represents the response for getting client packages
type GetClientPackagesResponse struct {
	Packages []GetClientPackagesResponseItem `json:"packages"`
}

// GetClientPackagesResponseItem represents a package in the response
type GetClientPackagesResponseItem struct {
	ID                  string `json:"id"`
	IssuedAt            string `json:"issued_at"`
	ExpiresAt           string `json:"expires_at"`
	ApppointmentsNumber int64  `json:"apppointments_number"`
	FirstName           string `json:"first_name"`
	LastName            string `json:"last_name"`
}

// GetClientPackageDetailsResponse represents the response for getting client package details
type GetClientPackageDetailsResponse struct {
	ID                  string                                       `json:"id"`
	IssuedAt            string                                       `json:"issued_at"`
	ExpiresAt           string                                       `json:"expires_at"`
	ApppointmentsNumber int64                                        `json:"apppointments_number"`
	FirstName           string                                       `json:"first_name"`
	LastName            string                                       `json:"last_name"`
	Appointments        []GetClientPackageDetailsResponseAppointment `json:"appointments"`
}

// GetClientPackageDetailsResponseAppointment represents an appointment in the response
type GetClientPackageDetailsResponseAppointment struct {
	ID        string `json:"id"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Type      string `json:"type"`
}
