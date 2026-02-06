package professionals

import (
	"github.com/google/uuid"
	db "github.com/vention/booking_api/internal/repository"
)

func buildMapClientIdInviteID(invites []*db.GetInvitesByAppointmentIDAndClientIsRow) map[uuid.UUID]uuid.UUID {
	mapClientIdInviteID := make(map[uuid.UUID]uuid.UUID)
	for _, invite := range invites {
		mapClientIdInviteID[invite.ClientID] = invite.ID
	}
	return mapClientIdInviteID
}
