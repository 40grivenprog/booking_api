package api

import (
	"context"
	db "github.com/vention/booking_api/internal/repository"
)

type ProfessionalsRepository interface {
	GetProfessionals(ctx context.Context) ([]*db.User, error)
	GetUserByUsername(ctx context.Context, username string) (*db.User, error)
	UpdateUserChatID(ctx context.Context, arg *db.UpdateUserChatIDParams) (*db.User, error)
}
