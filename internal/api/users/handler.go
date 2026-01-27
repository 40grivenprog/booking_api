package api

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/vention/booking_api/internal/services/users"
	"github.com/vention/booking_api/internal/token"
)

// UsersHandlerParams contains the parameters needed to register users handlers
type UsersHandlerParams struct {
	Router       *gin.RouterGroup
	TokenMaker   token.Maker
	UsersService users.Service
}

// UsersRegister registers all users-related routes
func UsersRegister(params UsersHandlerParams) error {
	if params.Router == nil {
		return errors.New("missing router")
	}

	if params.TokenMaker == nil {
		return errors.New("missing token maker")
	}

	if params.UsersService == nil {
		return errors.New("missing users service")
	}

	// Create controller
	controller := NewUsersController(params.TokenMaker, params.UsersService)

	// Create users group
	users := params.Router.Group("/users")
	{
		users.GET("/:chat_id", controller.GetUserByChatID)
	}

	return nil
}
