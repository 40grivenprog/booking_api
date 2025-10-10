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
		common.HandleErrorResponse(ctx, http.StatusBadRequest, common.ErrorTypeValidation, common.ErrorMsgInvalidClientID, err)
		return
	}

	// Get user from repository
	user, err := c.usersRepo.GetUserByChatID(ctx.Request.Context(), sql.NullInt64{Int64: chatID, Valid: true})
	if err != nil {
		common.HandleErrorResponse(ctx, http.StatusNotFound, common.ErrorTypeNotFound, common.ErrorMsgUserNotFound, err)
		return
	}

	// Return success response
	ctx.JSON(http.StatusOK, GetUserByChatIDResponse{
		User: User{
			ID:          user.ID.String(),
			ChatID:      common.FromNullInt64(user.ChatID),
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Role:        user.Role,
			PhoneNumber: common.FromNullString(user.PhoneNumber),
			CreatedAt:   common.FormatTimeWithTimezone(user.CreatedAt),
			UpdatedAt:   common.FormatTimeWithTimezone(user.UpdatedAt),
		},
	})
}
