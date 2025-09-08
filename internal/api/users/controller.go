package api

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vention/booking_api/internal/api/common"
)

// UsersController handles user-related HTTP requests
type UsersController struct {
	usersRepo UsersRepository
}

// NewUsersController creates a new users controller
func NewUsersController(usersRepo UsersRepository) *UsersController {
	return &UsersController{
		usersRepo: usersRepo,
	}
}

// GetUserByChatID handles GET /api/users/{chat_id}
func (c *UsersController) GetUserByChatID(ctx *gin.Context) {
	// Parse chat_id from URL parameter
	chatIDStr := ctx.Param("chat_id")
	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		common.HandleErrorResponse(ctx, http.StatusBadRequest, "validation_error", "Invalid chat_id format", err)
		return
	}

	// Get user from repository
	user, err := c.usersRepo.GetUserByChatID(ctx.Request.Context(), sql.NullInt64{Int64: chatID, Valid: true})
	if err != nil {
		common.HandleErrorResponse(ctx, http.StatusNotFound, "not_found", "User not found", err)
		return
	}

	// Return success response
	ctx.JSON(http.StatusOK, GetUserByChatIDResponse{
		User: User{
			ID:          user.ID.String(),
			ChatID:      &user.ChatID.Int64,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Role:        user.Role,
			CreatedAt:   user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
			PhoneNumber: &user.PhoneNumber.String,
		},
	})
}
