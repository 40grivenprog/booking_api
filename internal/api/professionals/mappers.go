package api

import (
	common "github.com/vention/booking_api/internal/api/common"
	db "github.com/vention/booking_api/internal/repository"
	"github.com/vention/booking_api/internal/services/professionals"
)

// import (
// 	"fmt"

// 	common "github.com/vention/booking_api/internal/api/common"
// 	db "github.com/vention/booking_api/internal/repository"
// )

func mapProfessionalToProfessionalSignInResponse(professional *db.Professional, token string) ProfessionalSignInResponse {
	responseUser := User{
		ID:        professional.ID.String(),
		Username:  professional.Username,
		FirstName: professional.FirstName,
		LastName:  professional.LastName,
		Role:      common.RoleProfessional,
		ChatID:    common.FromNullInt64(professional.ChatID),
		Token:     token,
		CreatedAt: common.FormatTimeWithTimezone(professional.CreatedAt),
		UpdatedAt: common.FormatTimeWithTimezone(professional.UpdatedAt),
	}

	return ProfessionalSignInResponse{
		User: responseUser,
	}
}

// mapAppointmentsToGetProfessionalAppointmentsResponse maps a list of appointments to a GetProfessionalAppointmentsResponse
func mapAppointmentsToGetProfessionalAppointmentsResponse(appointments []*db.GetAppointmentsByProfessionalWithStatusAndDateRow, total, page, pageSize int) GetProfessionalAppointmentsResponse {
	responseAppointments := make([]GetProfessionalAppointmentsResponseItem, len(appointments))
	for i, appt := range appointments {
		appointment := GetProfessionalAppointmentsResponseItem{
			ID:        appt.ID.String(),
			StartTime: common.FormatTimeRFC3339(appt.StartTime),
			EndTime:   common.FormatTimeRFC3339(appt.EndTime),
			Type:      appt.Type,
			Clients:   appt.Clients,
		}

		responseAppointments[i] = appointment
	}

	offset := (page - 1) * pageSize
	hasNextPage := offset+pageSize < total

	response := GetProfessionalAppointmentsResponse{
		Appointments: responseAppointments,
		Pagination: common.PaginationResponse{
			HasNextPage: hasNextPage,
			Page:        page,
			PageSize:    pageSize,
		},
	}

	return response
}

// mapTimetableAppointmentsToGetProfessionalTimetableResponse maps a list of appointments to a GetProfessionalTimetableResponse
func mapTimetableAppointmentsToGetProfessionalTimetableResponse(appointments []*db.GetProfessionalTimetableRow, dateStr string) GetProfessionalTimetableResponse {
	timetableAppointments := make([]TimetableAppointment, len(appointments))
	for i, apt := range appointments {
		description := ""
		if apt.Description.Valid {
			description = apt.Description.String
		}
		timetableAppointments[i] = TimetableAppointment{
			ID:          apt.ID.String(),
			StartTime:   common.FormatTimeRFC3339(apt.StartTime),
			EndTime:     common.FormatTimeRFC3339(apt.EndTime),
			Type:        apt.Type,
			Clients:     apt.Clients,
			Description: description,
		}
	}

	response := GetProfessionalTimetableResponse{
		Date:         dateStr,
		Appointments: timetableAppointments,
	}

	return response
}

func mapSubscriptionsToGetProfessionalSubscriptionsResponse(subscriptions []*db.GetSubscriptionsByProfessionalIDRow) GetProfessionalSubscriptionsResponse {
	responseSubscriptions := make([]ProfessionalSubscription, len(subscriptions))
	for i, sub := range subscriptions {
		responseSubscriptions[i] = ProfessionalSubscription{
			ID:        sub.ID.String(),
			FirstName: sub.FirstName,
			LastName:  sub.LastName,
			ChatID:    common.FromNullInt64(sub.ChatID),
			Locale:    sub.Locale,
		}
	}

	response := GetProfessionalSubscriptionsResponse{
		Subscriptions: responseSubscriptions,
	}

	return response
}

// mapPreviousAppointmentsToGetPreviousAppointmentsResponse maps a list of previous appointments to a GetPreviousAppointmentsResponse
func mapPreviousAppointmentsToGetPreviousAppointmentsResponse(appointments []*db.GetPreviousAppointmentsByProfessionalIDRow, total, page, pageSize int) GetPreviousAppointmentsResponse {
	responseAppointments := make([]GetPreviousAppointmentsResponseItem, len(appointments))
	for i, apt := range appointments {
		description := ""
		if apt.Description.Valid {
			description = apt.Description.String
		}
		responseAppointments[i] = GetPreviousAppointmentsResponseItem{
			ID:          apt.ID.String(),
			StartTime:   common.FormatTimeRFC3339(apt.StartTime),
			EndTime:     common.FormatTimeRFC3339(apt.EndTime),
			Description: description,
			Type:        apt.Type,
			UsersCount:  apt.UsersCount,
		}
	}

	offset := (page - 1) * pageSize
	hasNextPage := offset+pageSize < total

	response := GetPreviousAppointmentsResponse{
		Appointments: responseAppointments,
		Pagination: common.PaginationResponse{
			HasNextPage: hasNextPage,
			Page:        page,
			PageSize:    pageSize,
		},
	}

	return response
}

func mapPreviousAppointmentsByClientIDToGetPreviousAppointmentsResponse(appointments []*db.GetPreviousAppointmentsByProfessionalIDAndClientIDRow, total, page, pageSize int) GetPreviousAppointmentsResponse {
	responseAppointments := make([]GetPreviousAppointmentsResponseItem, len(appointments))
	for i, apt := range appointments {
		description := ""
		if apt.Description.Valid {
			description = apt.Description.String
		}
		responseAppointments[i] = GetPreviousAppointmentsResponseItem{
			ID:          apt.ID.String(),
			StartTime:   common.FormatTimeRFC3339(apt.StartTime),
			EndTime:     common.FormatTimeRFC3339(apt.EndTime),
			Description: description,
			Type:        apt.Type,
			UsersCount:  apt.UsersCount,
		}
	}

	offset := (page - 1) * pageSize
	hasNextPage := offset+pageSize < total

	response := GetPreviousAppointmentsResponse{
		Appointments: responseAppointments,
		Pagination: common.PaginationResponse{
			HasNextPage: hasNextPage,
			Page:        page,
			PageSize:    pageSize,
		},
	}

	return response
}

// mapAppointmentDetailsToGetAppointmentDetailsResponse maps an appointment details to a GetAppointmentDetailsResponse
func mapAppointmentDetailsToGetAppointmentDetailsResponse(appointmentDetails *professionals.GetAppointmentDetailsOutput) GetAppointmentDetailsResponse {
	clients := make([]GetAppointmentDetailsResponseClient, len(appointmentDetails.Clients))
	for i, client := range appointmentDetails.Clients {
		clients[i] = GetAppointmentDetailsResponseClient{
			ID:        client.ID.String(),
			FirstName: client.FirstName,
			LastName:  client.LastName,
		}
	}

	response := GetAppointmentDetailsResponse{
		ID:          appointmentDetails.ID.String(),
		StartTime:   common.FormatTimeRFC3339(appointmentDetails.StartTime),
		EndTime:     common.FormatTimeRFC3339(appointmentDetails.EndTime),
		Description: appointmentDetails.Description.String,
		Type:        appointmentDetails.Type,
		Clients:     clients,
	}

	return response
}

// func mapClientsToGetProfessionalClientsResponse(clients []*db.GetProfessionalClientsRow) GetProfessionalClientsResponse {
// 	responseClients := make([]ProfessionalClient, len(clients))
// 	for i, client := range clients {
// 		responseClients[i] = ProfessionalClient{
// 			ID:        client.ID.String(),
// 			FirstName: client.FirstName,
// 			LastName:  client.LastName,
// 		}
// 	}

// 	response := GetProfessionalClientsResponse{
// 		Clients: responseClients,
// 	}

// 	return response
// }

// func mapPreviousAppointmentsToGetPreviousAppointmentsByClientResponse(appointments []*db.GetPreviousProfessionalAppointmentsByClientRow) GetPreviousAppointmentsByClientResponse {
// 	responseAppointments := make([]PreviousAppointment, len(appointments))
// 	for i, apt := range appointments {
// 		responseAppointments[i] = PreviousAppointment{
// 			ID:          apt.ID.String(),
// 			StartTime:   common.FormatTimeRFC3339(apt.StartTime),
// 			EndTime:     common.FormatTimeRFC3339(apt.EndTime),
// 			Description: *common.FromNullString(apt.Description),
// 		}
// 	}

// 	response := GetPreviousAppointmentsByClientResponse{
// 		Appointments: responseAppointments,
// 	}

// 	return response
// }
