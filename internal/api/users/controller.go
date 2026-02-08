package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vention/booking_api/internal/api/common"
	"github.com/vention/booking_api/internal/services/users"
	"github.com/vention/booking_api/internal/token"
)

// UsersController handles user-related HTTP requests
type UsersController struct {
	tokenMaker   token.Maker
	usersService users.Service
}

// NewUsersController creates a new users controller
func NewUsersController(tokenMaker token.Maker, usersService users.Service) *UsersController {
	return &UsersController{
		tokenMaker:   tokenMaker,
		usersService: usersService,
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

	user, err := c.usersService.GetUserByChatID(ctx.Request.Context(), chatID)
	if err != nil {
		common.HandleErrorResponse(ctx, http.StatusNotFound, common.ErrorTypeNotFound, common.ErrorMsgUserNotFound, err)
		return
	}

	token, err := c.tokenMaker.CreateToken(user.ID, fmt.Sprintf("%s %s", user.FirstName, user.LastName))
	if err != nil {
		common.HandleErrorResponse(ctx, http.StatusInternalServerError, common.ErrorTypeInternal, common.ErrorMsgFailedToCreateToken, err)
		return
	}

	// Return success response
	ctx.JSON(http.StatusOK, GetUserByChatIDResponse{
		User: GetUserByChatIDResponseItem{
			ID:        user.ID.String(),
			ChatID:    common.FromNullInt64(user.ChatID),
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Role:      user.Role,
			Locale:    user.Locale,
			Token:     token,
		},
	})
}
