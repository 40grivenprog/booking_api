package api

import (
	"fmt"

	common "github.com/vention/booking_api/internal/api/common"
	db "github.com/vention/booking_api/internal/repository"
)

func mapProfessionalToProfessionalSignInResponse(professional *db.Professional, token string) ProfessionalSignInResponse {
	responseUser := User{
		ID:          professional.ID.String(),
		Username:    professional.Username,
		FirstName:   professional.FirstName,
		LastName:    professional.LastName,
		Role:        common.RoleProfessional,
		PhoneNumber: common.FromNullString(professional.PhoneNumber),
		ChatID:      common.FromNullInt64(professional.ChatID),
		Token:       token,
		CreatedAt:   common.FormatTimeWithTimezone(professional.CreatedAt),
		UpdatedAt:   common.FormatTimeWithTimezone(professional.UpdatedAt),
	}

	return ProfessionalSignInResponse{
		User: responseUser,
	}
}

func mapAppointmentToConfirmAppointmentResponse(appointment *db.ConfirmAppointmentWithDetailsRow) ConfirmAppointmentResponse {
	return ConfirmAppointmentResponse{
		Appointment: ConfirmAppointmentResponseAppointmentItem{
			StartTime: common.FormatTimeRFC3339(appointment.StartTime),
			EndTime:   common.FormatTimeRFC3339(appointment.EndTime),
		},
		Client: ConfirmAppointmentResponseClientItem{
			ChatID: appointment.ClientChatID.Int64,
		},
		Professional: ConfirmAppointmentResponseProfessionalItem{
			FirstName: appointment.ProfessionalFirstName.String,
			LastName:  appointment.ProfessionalLastName.String,
		},
	}
}

// mapAppointmentsToGetProfessionalAppointmentsResponse maps a list of appointments to a GetProfessionalAppointmentsResponse
func mapAppointmentsToGetProfessionalAppointmentsResponse(appointments []*db.GetAppointmentsByProfessionalWithStatusAndDateRow, total, page, pageSize int) GetProfessionalAppointmentsResponse {
	responseAppointments := make([]GetProfessionalAppointmentsResponseItem, len(appointments))
	for i, appt := range appointments {
		appointment := GetProfessionalAppointmentsResponseItem{
			ID:          appt.ID.String(),
			StartTime:   common.FormatTimeRFC3339(appt.StartTime),
			EndTime:     common.FormatTimeRFC3339(appt.EndTime),
			Description: appt.Description.String,
		}

		// Only include client if we have client data
		if appt.ClientFirstName.Valid && appt.ClientLastName.Valid {
			appointment.Client = &GetProfessionalAppointmentsResponseClient{
				FirstName: appt.ClientFirstName.String,
				LastName:  appt.ClientLastName.String,
			}
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

func mapAppointmentToCancelAppointmentResponse(appointment *db.CancelAppointmentByProfessionalWithDetailsRow) CancelAppointmentResponse {
	return CancelAppointmentResponse{
		Appointment: CancelledAppointment{
			ID:                 appointment.ID.String(),
			Type:               string(appointment.Type),
			StartTime:          common.FormatTimeRFC3339(appointment.StartTime),
			EndTime:            common.FormatTimeRFC3339(appointment.EndTime),
			Status:             string(appointment.Status.AppointmentStatus),
			CancellationReason: appointment.CancellationReason.String,
			CancelledBy:        common.CancelledByProfessional,
			CreatedAt:          common.FormatTimeRFC3339(appointment.CreatedAt),
			UpdatedAt:          common.FormatTimeRFC3339(appointment.UpdatedAt),
		},
		Client: ProfessionalAppointmentClient{
			ID:        appointment.ClientIDFull.String(),
			FirstName: appointment.ClientFirstName.String,
			LastName:  appointment.ClientLastName.String,
			ChatID:    common.FromNullInt64(appointment.ClientChatID),
		},
		Professional: ProfessionalInfo{
			ID:        appointment.ProfessionalIDFull.String(),
			Username:  appointment.ProfessionalUsername.String,
			FirstName: appointment.ProfessionalFirstName.String,
			LastName:  appointment.ProfessionalLastName.String,
		},
	}
}

func mapAppointmentToCreateUnavailableAppointmentResponse(appointment *db.Appointment) CreateUnavailableAppointmentResponse {
	return CreateUnavailableAppointmentResponse{
		Appointment: UnavailableAppointment{
			Type:        string(appointment.Type),
			StartTime:   common.FormatTimeRFC3339(appointment.StartTime),
			EndTime:     common.FormatTimeRFC3339(appointment.EndTime),
			Status:      string(appointment.Status.AppointmentStatus),
			Description: appointment.Description.String,
		},
	}
}

func mapTimetableAppointmentsToGetProfessionalTimetableResponse(appointments []*db.GetProfessionalTimetableRow, dateStr string) GetProfessionalTimetableResponse {
	timetableAppointments := make([]TimetableAppointment, len(appointments))
	for i, apt := range appointments {
		// Format description with client name if available
		description := apt.Description.String
		if apt.FirstName.Valid && apt.LastName.Valid {
			description = fmt.Sprintf("%s %s - %s",
				apt.FirstName.String,
				apt.LastName.String,
				apt.Description.String)
		}

		timetableAppointments[i] = TimetableAppointment{
			ID:          apt.ID.String(),
			StartTime:   common.FormatTimeRFC3339(apt.StartTime),
			EndTime:     common.FormatTimeRFC3339(apt.EndTime),
			Description: description,
		}
	}

	response := GetProfessionalTimetableResponse{
		Date:         dateStr,
		Appointments: timetableAppointments,
	}

	return response
}

func mapClientsToGetProfessionalClientsResponse(clients []*db.GetProfessionalClientsRow) GetProfessionalClientsResponse {
	responseClients := make([]ProfessionalClient, len(clients))
	for i, client := range clients {
		responseClients[i] = ProfessionalClient{
			ID:        client.ID.String(),
			FirstName: client.FirstName,
			LastName:  client.LastName,
		}
	}

	response := GetProfessionalClientsResponse{
		Clients: responseClients,
	}

	return response
}

func mapPreviousAppointmentsToGetPreviousAppointmentsByClientResponse(appointments []*db.GetPreviousProfessionalAppointmentsByClientRow) GetPreviousAppointmentsByClientResponse {
	responseAppointments := make([]PreviousAppointment, len(appointments))
	for i, apt := range appointments {
		responseAppointments[i] = PreviousAppointment{
			ID:          apt.ID.String(),
			StartTime:   common.FormatTimeRFC3339(apt.StartTime),
			EndTime:     common.FormatTimeRFC3339(apt.EndTime),
			Description: *common.FromNullString(apt.Description),
		}
	}

	response := GetPreviousAppointmentsByClientResponse{
		Appointments: responseAppointments,
	}

	return response
}
