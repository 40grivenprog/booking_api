package api

import (
	common "github.com/vention/booking_api/internal/api/common"
	db "github.com/vention/booking_api/internal/repository"
)

// mapAppointmentToCreateAppointmentResponse maps database result to API response
func mapAppointmentToCreateAppointmentResponse(appointment *db.CreateAppointmentWithDetailsRow) CreateAppointmentResponse {
	return CreateAppointmentResponse{
		Appointment: Appointment{
			ID:          appointment.ID.String(),
			StartTime:   common.FormatTimeRFC3339(appointment.StartTime),
			EndTime:     common.FormatTimeRFC3339(appointment.EndTime),
			Status:      string(appointment.Status.AppointmentStatus),
			Description: appointment.Description.String,
			CreatedAt:   common.FormatTimeRFC3339(appointment.CreatedAt),
			UpdatedAt:   common.FormatTimeRFC3339(appointment.UpdatedAt),
		},
		Client: Client{
			ID:          appointment.ClientIDFull.String(),
			FirstName:   appointment.ClientFirstName.String,
			LastName:    appointment.ClientLastName.String,
			PhoneNumber: common.StringValue(common.FromNullString(appointment.ClientPhoneNumber)),
			ChatID:      common.Int64Value(common.FromNullInt64(appointment.ClientChatID)),
		},
		Professional: Professional{
			ID:          appointment.ProfessionalIDFull.String(),
			Username:    appointment.ProfessionalUsername.String,
			FirstName:   appointment.ProfessionalFirstName.String,
			LastName:    appointment.ProfessionalLastName.String,
			PhoneNumber: common.StringValue(common.FromNullString(appointment.ProfessionalPhoneNumber)),
			ChatID:      common.Int64Value(common.FromNullInt64(appointment.ProfessionalChatID)),
		},
	}
}
