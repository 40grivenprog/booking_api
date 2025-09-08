package api

import (
	"context"
	"database/sql"

	db "github.com/vention/booking_api/internal/repository"
)

// UsersRepository interface defines the contract for user-related database operations
type UsersRepository interface {
	GetUserByChatID(context.Context, sql.NullInt64) (*db.GetUserByChatIDRow, error)
}
