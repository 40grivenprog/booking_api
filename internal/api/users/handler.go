package api

import (
	"errors"

	"github.com/gin-gonic/gin"
)

// UsersHandlerParams contains the parameters needed to register users handlers
type UsersHandlerParams struct {
	Router    *gin.RouterGroup
	UsersRepo UsersRepository
}

// UsersRegister registers all users-related routes
func UsersRegister(params UsersHandlerParams) error {
	if params.Router == nil {
		return errors.New("missing router")
	}

	if params.UsersRepo == nil {
		return errors.New("missing users repository")
	}

	// Create controller
	controller := NewUsersController(params.UsersRepo)

	// Create users group
	users := params.Router.Group("/users")
	{
		users.GET("/:chat_id", controller.GetUserByChatID)
	}

	return nil
}
