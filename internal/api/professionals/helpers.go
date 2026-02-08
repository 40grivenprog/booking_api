package api

import (
	"github.com/google/uuid"
	"github.com/vention/booking_api/internal/services/professionals"
)

// ConvertCreateGroupVisitAppointmentClientToInput converts a CreateGroupVisitAppointmentClient to a CreateGroupVisitAppointmentClientInput
func convertCreateGroupVisitAppointmentClientToInput(clients []CreateGroupVisitAppointmentClient) []professionals.InvitePartiallySelectedClientsClientInput {
	inputClients := make([]professionals.InvitePartiallySelectedClientsClientInput, len(clients))
	for i, client := range clients {
		inputClients[i] = professionals.InvitePartiallySelectedClientsClientInput{
			ID:     uuid.MustParse(client.ID),
			ChatID: client.ChatID,
			Locale: client.Locale,
		}
	}
	return inputClients
}

// convertCreateInviteForAppointmentClientToInput converts a CreateInviteForAppointmentClient to a CreateInviteForAppointmentClientInput
func convertCreateInviteForAppointmentClientToInput(clients []CreateInvitesForAppointmentRequestClientItem) []professionals.InvitePartiallySelectedClientsClientInput {
	inputClients := make([]professionals.InvitePartiallySelectedClientsClientInput, len(clients))
	for i, client := range clients {
		inputClients[i] = professionals.InvitePartiallySelectedClientsClientInput{
			ID:     uuid.MustParse(client.ID),
			ChatID: client.ChatID,
			Locale: client.Locale,
		}
	}
	return inputClients
}
