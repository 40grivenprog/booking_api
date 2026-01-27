package users

import (
	"context"
	"database/sql"

	db "github.com/vention/booking_api/internal/repository"
)

// Service defines the business logic operations for clients
type Service interface {
	GetUserByChatID(ctx context.Context, chatID int64) (*db.GetUserByChatIDRow, error)
}

type service struct {
	usersRepo UsersRepository
}

// NewService creates a new users service
func NewService(repo UsersRepository) Service {
	return &service{
		usersRepo: repo,
	}
}

func (s *service) GetUserByChatID(ctx context.Context, chatID int64) (*db.GetUserByChatIDRow, error) {
	return s.usersRepo.GetUserByChatID(ctx, sql.NullInt64{Int64: chatID, Valid: true})
}
