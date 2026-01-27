package users

import (
	"context"
	"database/sql"

	db "github.com/vention/booking_api/internal/repository"
)

// UsersRepository defines the database operations needed by the users service
type UsersRepository interface {
	GetUserByChatID(ctx context.Context, chatID sql.NullInt64) (*db.GetUserByChatIDRow, error)
}
