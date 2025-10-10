package api

import (
	common "github.com/vention/booking_api/internal/api/common"
	db "github.com/vention/booking_api/internal/repository"
)

func mapClientToClientRegisterResponse(client *db.Client) ClientRegisterResponse {
	return ClientRegisterResponse{
		ID:          client.ID.String(),
		FirstName:   client.FirstName,
		LastName:    client.LastName,
		Role:        common.RoleClient,
		PhoneNumber: common.FromNullString(client.PhoneNumber),
		ChatID:      common.FromNullInt64(client.ChatID),
		CreatedAt:   common.FormatTimeWithTimezone(client.CreatedAt),
		UpdatedAt:   common.FormatTimeWithTimezone(client.UpdatedAt),
	}
}

func mapAppointmentToGetClientAppointmentsResponse(appointments []*db.GetAppointmentsByClientWithStatusRow) GetClientAppointmentsResponse {
	var responseAppointments []ClientAppointment
	for _, appt := range appointments {
		appointment := ClientAppointment{
			ID:          appt.ID.String(),
			Type:        string(appt.Type),
			StartTime:   common.FormatTimeRFC3339(appt.StartTime),
			EndTime:     common.FormatTimeRFC3339(appt.EndTime),
			Description: appt.Description.String,
			Status:      string(appt.Status.AppointmentStatus),
			CreatedAt:   common.FormatTimeRFC3339(appt.CreatedAt),
			UpdatedAt:   common.FormatTimeRFC3339(appt.UpdatedAt),
		}
		professional := &ClientAppointmentProfessional{
			ID:        appt.ProfessionalIDFull.String(),
			Username:  appt.ProfessionalUsername.String,
			FirstName: appt.ProfessionalFirstName.String,
			LastName:  appt.ProfessionalLastName.String,
		}
		appointment.Professional = professional

		responseAppointments = append(responseAppointments, appointment)
	}

	response := GetClientAppointmentsResponse{
		Appointments: responseAppointments,
	}

	return response
}

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
			ID:          appointment.ProfessionalIDFull.String(),
			Username:    appointment.ProfessionalUsername.String,
			FirstName:   appointment.ProfessionalFirstName.String,
			LastName:    appointment.ProfessionalLastName.String,
			PhoneNumber: common.FromNullString(appointment.ProfessionalPhoneNumber),
			ChatID:      common.FromNullInt64(appointment.ProfessionalChatID),
		},
	}
}
