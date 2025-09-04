package api

import (
	"context"

	db "github.com/vention/booking_api/internal/repository"
)

// AdminRepository implements the admin repository interface
type AdminRepository interface {
	CreateProfessional(context.Context, *db.CreateProfessionalParams) (*db.Professional, error)
}
