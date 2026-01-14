package api

import (
	"fmt"

	common "github.com/vention/booking_api/internal/api/common"
	db "github.com/vention/booking_api/internal/repository"
)

func mapProfessionalsToGetProfessionalsResponse(rows []*db.GetProfessionalsRow, total, page, pageSize int) GetProfessionalsResponse {
	responseItems := make([]GetProfessionalsResponseItem, len(rows))
	for i, row := range rows {
		responseItems[i] = GetProfessionalsResponseItem{
			ID:        row.ID.String(),
			FirstName: row.FirstName,
			LastName:  row.LastName,
		}
	}

	offset := (page - 1) * pageSize
	hasNextPage := offset+pageSize < total

	return GetProfessionalsResponse{
		Professionals: responseItems,
		Pagination: common.PaginationResponse{
			HasNextPage: hasNextPage,
			Page:        page,
			PageSize:    pageSize,
		},
	}
}

func mapProfessionalToProfessionalSignInResponse(professional *db.Professional) ProfessionalSignInResponse {
	responseUser := User{
		ID:          professional.ID.String(),
		Username:    professional.Username,
		FirstName:   professional.FirstName,
		LastName:    professional.LastName,
		Role:        common.RoleProfessional,
		PhoneNumber: common.FromNullString(professional.PhoneNumber),
		ChatID:      common.FromNullInt64(professional.ChatID),
		CreatedAt:   common.FormatTimeWithTimezone(professional.CreatedAt),
		UpdatedAt:   common.FormatTimeWithTimezone(professional.UpdatedAt),
	}

	return ProfessionalSignInResponse{
		User: responseUser,
	}
}

func mapAppointmentToConfirmAppointmentResponse(appointment *db.ConfirmAppointmentWithDetailsRow) ConfirmAppointmentResponse {
	return ConfirmAppointmentResponse{
		Appointment: AppointmentConfirm{
			ID:        appointment.ID.String(),
			Status:    string(appointment.Status.AppointmentStatus),
			StartTime: common.FormatTimeRFC3339(appointment.StartTime),
			EndTime:   common.FormatTimeRFC3339(appointment.EndTime),
			CreatedAt: common.FormatTimeRFC3339(appointment.CreatedAt),
			UpdatedAt: common.FormatTimeRFC3339(appointment.UpdatedAt),
		},
		Client: ClientConfirm{
			ID:        appointment.ClientID.UUID.String(),
			FirstName: appointment.ClientFirstName.String,
			LastName:  appointment.ClientLastName.String,
			ChatID:    appointment.ClientChatID.Int64,
		},
		Professional: ProfessionalConfirm{
			ID:        appointment.ProfessionalIDFull.String(),
			Username:  appointment.ProfessionalUsername.String,
			FirstName: appointment.ProfessionalFirstName.String,
			LastName:  appointment.ProfessionalLastName.String,
		},
	}
}

// mapAppointmentsToGetProfessionalAppointmentsResponse maps a list of appointments to a GetProfessionalAppointmentsResponse
func mapAppointmentsToGetProfessionalAppointmentsResponse(appointments []*db.GetAppointmentsByProfessionalWithStatusAndDateRow) GetProfessionalAppointmentsResponse {
	responseAppointments := make([]ProfessionalAppointment, len(appointments))
	for i, appt := range appointments {
		appointment := ProfessionalAppointment{
			ID:          appt.ID.String(),
			Type:        string(appt.Type),
			StartTime:   common.FormatTimeRFC3339(appt.StartTime),
			EndTime:     common.FormatTimeRFC3339(appt.EndTime),
			Description: appt.Description.String,
			Status:      string(appt.Status.AppointmentStatus),
			CreatedAt:   common.FormatTimeRFC3339(appt.CreatedAt),
			UpdatedAt:   common.FormatTimeRFC3339(appt.UpdatedAt),
		}
		appointment.Client = &ProfessionalAppointmentClient{
			ID:          appt.ClientID.UUID.String(),
			FirstName:   appt.ClientFirstName.String,
			LastName:    appt.ClientLastName.String,
			PhoneNumber: &appt.ClientPhoneNumber.String,
		}
		responseAppointments[i] = appointment
	}

	response := GetProfessionalAppointmentsResponse{
		Appointments: responseAppointments,
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
			ID:          appointment.ID.String(),
			Type:        string(appointment.Type),
			StartTime:   common.FormatTimeRFC3339(appointment.StartTime),
			EndTime:     common.FormatTimeRFC3339(appointment.EndTime),
			Status:      string(appointment.Status.AppointmentStatus),
			Description: appointment.Description.String,
			CreatedAt:   common.FormatTimeRFC3339(appointment.CreatedAt),
			UpdatedAt:   common.FormatTimeRFC3339(appointment.UpdatedAt),
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
