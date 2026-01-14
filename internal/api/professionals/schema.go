package api

import (
	common "github.com/vention/booking_api/internal/api/common"
)

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
	Professionals []GetProfessionalsResponseItem `json:"professionals"`
	Pagination    common.PaginationResponse      `json:"pagination"`
}

// GetProfessionalsResponseItem represents a professional in the response
type GetProfessionalsResponseItem struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
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
	ID          string `json:"id"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	Description string `json:"description"`
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
