package api

import (
	"context"

	db "github.com/vention/booking_api/internal/repository"
)

// AdminsRepository describes the admin repository interface
type AdminsRepository interface {
	CreateProfessional(ctx context.Context, params *db.CreateProfessionalParams) (*db.User, error)
}
