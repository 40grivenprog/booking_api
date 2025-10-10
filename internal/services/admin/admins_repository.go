package admin

import (
	"context"

	db "github.com/vention/booking_api/internal/repository"
)

// Repository defines the database operations needed by the admin service
type AdminsRepository interface {
	CreateProfessional(ctx context.Context, arg *db.CreateProfessionalParams) (*db.Professional, error)
}
