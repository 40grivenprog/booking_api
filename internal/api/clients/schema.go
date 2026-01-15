package api

import (
	common "github.com/vention/booking_api/internal/api/common"
)

// CreateAppointmentRequest represents the request to create an appointment
type CreateAppointmentRequest struct {
	ProfessionalID string `json:"professional_id" binding:"required"`
	StartTime      string `json:"start_time" binding:"required"`
	EndTime        string `json:"end_time" binding:"required"`
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
	FirstName   string  `json:"first_name" binding:"required"`
	LastName    string  `json:"last_name" binding:"required"`
	ChatID      int64   `json:"chat_id" binding:"required"`
	PhoneNumber *string `json:"phone_number,omitempty"`
}

// ClientRegisterResponse represents the response for client registration
type ClientRegisterResponse struct {
	ID          string  `json:"id"`
	ChatID      *int64  `json:"chat_id,omitempty"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	PhoneNumber *string `json:"phone_number,omitempty"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
	Role        string  `json:"role"`
	Token       string  `json:"token"`
}

// GetClientAppointmentsResponse represents the response for getting client appointments
type GetClientAppointmentsResponse struct {
	Appointments []ClientAppointment       `json:"appointments"`
	Pagination   common.PaginationResponse `json:"pagination"`
}

// ClientAppointment represents an appointment with professional details in client context
type ClientAppointment struct {
	ID           string                         `json:"id"`
	StartTime    string                         `json:"start_time"`
	EndTime      string                         `json:"end_time"`
	Description  string                         `json:"description,omitempty"`
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
