package api

import (
	"context"
	db "github.com/vention/booking_api/internal/repository"
)

type ClientsRepository interface {
	CreateClient(ctx context.Context, arg *db.CreateClientParams) (*db.User, error)
	GetUserByUsername(ctx context.Context, username string) (*db.User, error)
}
