package api

import (
	"context"
	db "github.com/vention/booking_api/internal/repository"
)

type ProfessionalsRepository interface {
	GetProfessionals(ctx context.Context) ([]*db.Professional, error)
	GetProfessionalByUsername(ctx context.Context, username string) (*db.Professional, error)
	UpdateProfessionalChatID(ctx context.Context, arg *db.UpdateProfessionalChatIDParams) (*db.Professional, error)
}
