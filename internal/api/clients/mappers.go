package api

import (
	common "github.com/vention/booking_api/internal/api/common"
	db "github.com/vention/booking_api/internal/repository"
)

// mapClientToClientRegisterResponse maps a client to a ClientRegisterResponse
func mapClientToClientRegisterResponse(client *db.Client, token string) ClientRegisterResponse {
	return ClientRegisterResponse{
		ID:        client.ID.String(),
		FirstName: client.FirstName,
		LastName:  client.LastName,
		Role:      common.RoleClient,
		ChatID:    common.FromNullInt64(client.ChatID),
		Token:     token,
		CreatedAt: common.FormatTimeWithTimezone(client.CreatedAt),
		UpdatedAt: common.FormatTimeWithTimezone(client.UpdatedAt),
	}
}

// mapAppointmentToGetClientAppointmentsResponse maps a list of appointments to a GetClientAppointmentsResponse
func mapAppointmentToGetClientAppointmentsResponse(appointments []*db.GetAppointmentsByClientWithStatusRow, total, page, pageSize int) GetClientAppointmentsResponse {
	var responseAppointments []ClientAppointment
	for _, appt := range appointments {
		appointment := ClientAppointment{
			ID:        appt.ID.String(),
			StartTime: common.FormatTimeRFC3339(appt.StartTime),
			EndTime:   common.FormatTimeRFC3339(appt.EndTime),
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

// // mapAppointmentToCancelClientAppointmentResponse maps an appointment to a CancelClientAppointmentResponse
// func mapAppointmentToCancelClientAppointmentResponse(appointment *db.CancelAppointmentByClientWithDetailsRow) CancelClientAppointmentResponse {
// 	return CancelClientAppointmentResponse{
// 		Appointment: CancelledAppointment{
// 			StartTime:          common.FormatTimeRFC3339(appointment.StartTime),
// 			EndTime:            common.FormatTimeRFC3339(appointment.EndTime),
// 			CancellationReason: appointment.CancellationReason.String,
// 		},
// 		Client: CancelledAppointmentClient{
// 			FirstName: appointment.ClientFirstName.String,
// 			LastName:  appointment.ClientLastName.String,
// 		},
// 		Professional: CancelledAppointmentProfessional{
// 			ChatID: common.FromNullInt64(appointment.ProfessionalChatID),
// 		},
// 	}
// }

// mapProfessionalsToGetSubscribedProfessionalsResponse maps a list of professionals to a GetSubscribedProfessionalsResponse
func mapProfessionalsToGetSubscribedProfessionalsResponse(professionals []*db.GetSubscriptionsByClientIDRow, total, page, pageSize int) GetSubscribedProfessionalsResponse {
	responseProfessionals := make([]GetSubscribedProfessionalsResponseItem, len(professionals))
	for i, professional := range professionals {
		responseProfessionals[i] = GetSubscribedProfessionalsResponseItem{
			ID:        professional.ID.String(),
			FirstName: professional.FirstName,
			LastName:  professional.LastName,
			ChatID:    common.FromNullInt64(professional.ChatID),
			Locale:    professional.Locale,
		}
	}

	offset := (page - 1) * pageSize
	hasNextPage := offset+pageSize < total

	return GetSubscribedProfessionalsResponse{
		Professionals: responseProfessionals,
		Pagination: common.PaginationResponse{
			HasNextPage: hasNextPage,
			Page:        page,
			PageSize:    pageSize,
		},
	}
}

func mapProfessionalsToGetProfessionalsResponse(rows []*db.GetProfessionalsRow, total, page, pageSize int) GetProfessionalsResponse {
	responseItems := make([]GetProfessionalsResponseItem, len(rows))
	for i, row := range rows {
		responseItems[i] = GetProfessionalsResponseItem{
			ID:        row.ID.String(),
			FirstName: row.FirstName,
			LastName:  row.LastName,
			ChatID:    common.FromNullInt64(row.ChatID),
			Locale:    row.Locale,
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

func mapInvitesToGetClientInvitesResponse(invites []*db.GetClientInvitesRow) GetClientInvitesResponse {
	responseInvites := make([]GetClientInvitesResponseItem, len(invites))
	for i, invite := range invites {
		responseInvites[i] = GetClientInvitesResponseItem{
			ID:               invite.ID.String(),
			AppointmentID:    invite.AppointmentID.String(),
			StartTime:        common.FormatTimeRFC3339(invite.StartTime),
			EndTime:          common.FormatTimeRFC3339(invite.EndTime),
			Description:      invite.Description,
			Type:             invite.Type,
			ProfessionalName: invite.ProfessionalName,
			ClientID:         invite.ClientID.String(),
		}
	}
	return GetClientInvitesResponse{
		Invites: responseInvites,
	}
}

func mapInviteToGetClientInviteResponse(invite *db.GetInviteByIDRow) GetClientInviteResponse {
	return GetClientInviteResponse{
		ID:               invite.ID.String(),
		AppointmentID:    invite.AppointmentID.String(),
		StartTime:        common.FormatTimeRFC3339(invite.StartTime),
		EndTime:          common.FormatTimeRFC3339(invite.EndTime),
		Description:      invite.Description,
		Type:             invite.Type,
		ProfessionalName: invite.ProfessionalName,
	}
}

// mapPreviousAppointmentsToGetClientPreviousAppointmentsResponse maps a list of previous appointments to a GetClientPreviousAppointmentsResponse
func mapPreviousAppointmentsToGetClientPreviousAppointmentsResponse(appointments []*db.GetPreviousAppointmentsByClientIDRow, total, page, pageSize int) GetClientPreviousAppointmentsResponse {
	responseAppointments := make([]GetClientPreviousAppointmentsResponseItem, len(appointments))
	for i, appt := range appointments {
		responseAppointments[i] = GetClientPreviousAppointmentsResponseItem{
			ID:        appt.ID.String(),
			StartTime: common.FormatTimeRFC3339(appt.StartTime),
			EndTime:   common.FormatTimeRFC3339(appt.EndTime),
			Type:      appt.Type,
			FirstName: appt.FirstName,
			LastName:  appt.LastName,
		}
	}

	offset := (page - 1) * pageSize
	hasNextPage := offset+pageSize < total

	return GetClientPreviousAppointmentsResponse{
		Appointments: responseAppointments,
		Pagination: common.PaginationResponse{
			HasNextPage: hasNextPage,
			Page:        page,
			PageSize:    pageSize,
		},
	}
}
