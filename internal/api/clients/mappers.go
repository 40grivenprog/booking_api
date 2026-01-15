package api

import (
	common "github.com/vention/booking_api/internal/api/common"
	db "github.com/vention/booking_api/internal/repository"
)

// mapClientToClientRegisterResponse maps a client to a ClientRegisterResponse
func mapClientToClientRegisterResponse(client *db.Client, token string) ClientRegisterResponse {
	return ClientRegisterResponse{
		ID:          client.ID.String(),
		FirstName:   client.FirstName,
		LastName:    client.LastName,
		Role:        common.RoleClient,
		PhoneNumber: common.FromNullString(client.PhoneNumber),
		ChatID:      common.FromNullInt64(client.ChatID),
		Token:       token,
		CreatedAt:   common.FormatTimeWithTimezone(client.CreatedAt),
		UpdatedAt:   common.FormatTimeWithTimezone(client.UpdatedAt),
	}
}

// mapAppointmentToGetClientAppointmentsResponse maps a list of appointments to a GetClientAppointmentsResponse
func mapAppointmentToGetClientAppointmentsResponse(appointments []*db.GetAppointmentsByClientWithStatusRow, total, page, pageSize int) GetClientAppointmentsResponse {
	var responseAppointments []ClientAppointment
	for _, appt := range appointments {
		appointment := ClientAppointment{
			ID:          appt.ID.String(),
			StartTime:   common.FormatTimeRFC3339(appt.StartTime),
			EndTime:     common.FormatTimeRFC3339(appt.EndTime),
			Description: appt.Description.String,
		}
		professional := &ClientAppointmentProfessional{
			FirstName: appt.ProfessionalFirstName.String,
			LastName:  appt.ProfessionalLastName.String,
		}
		appointment.Professional = professional

		responseAppointments = append(responseAppointments, appointment)
	}

	offset := (page - 1) * pageSize
	hasNextPage := offset+pageSize < total

	response := GetClientAppointmentsResponse{
		Appointments: responseAppointments,
		Pagination: common.PaginationResponse{
			HasNextPage: hasNextPage,
			Page:        page,
			PageSize:    pageSize,
		},
	}

	return response
}

// mapAppointmentToCancelClientAppointmentResponse maps an appointment to a CancelClientAppointmentResponse
func mapAppointmentToCancelClientAppointmentResponse(appointment *db.CancelAppointmentByClientWithDetailsRow) CancelClientAppointmentResponse {
	return CancelClientAppointmentResponse{
		Appointment: CancelledAppointment{
			ID:                 appointment.ID.String(),
			Type:               string(appointment.Type),
			StartTime:          common.FormatTimeRFC3339(appointment.StartTime),
			EndTime:            common.FormatTimeRFC3339(appointment.EndTime),
			Status:             string(appointment.Status.AppointmentStatus),
			CancellationReason: appointment.CancellationReason.String,
			CancelledBy:        common.CancelledByClient,
			CreatedAt:          common.FormatTimeRFC3339(appointment.CreatedAt),
			UpdatedAt:          common.FormatTimeRFC3339(appointment.UpdatedAt),
		},
		Client: ClientAppointmentClient{
			ID:          appointment.ClientIDFull.String(),
			FirstName:   appointment.ClientFirstName.String,
			LastName:    appointment.ClientLastName.String,
			PhoneNumber: common.FromNullString(appointment.ClientPhoneNumber),
			ChatID:      common.FromNullInt64(appointment.ClientChatID),
		},
		Professional: ClientAppointmentProfessional{
			FirstName: appointment.ProfessionalFirstName.String,
			LastName:  appointment.ProfessionalLastName.String,
			ChatID:    common.FromNullInt64(appointment.ProfessionalChatID),
		},
	}
}

// mapAppointmentToCreateAppointmentResponse maps database result to API response
func mapAppointmentToCreateAppointmentResponse(appointment *db.CreateAppointmentWithDetailsRow) CreateAppointmentResponse {
	return CreateAppointmentResponse{
		Appointment: Appointment{
			StartTime:   common.FormatTimeRFC3339(appointment.StartTime),
			EndTime:     common.FormatTimeRFC3339(appointment.EndTime),
			Description: appointment.Description.String,
		},
		Client: Client{
			FirstName: appointment.ClientFirstName.String,
			LastName:  appointment.ClientLastName.String,
		},
		Professional: Professional{
			ChatID: common.Int64Value(common.FromNullInt64(appointment.ProfessionalChatID)),
		},
	}
}
